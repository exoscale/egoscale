package v2

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewPoller(t *testing.T) {
	require.Equal(t,
		&Poller{interval: defaultPollingInterval},
		NewPoller())
}

func TestPoller_WithInterval(t *testing.T) {
	testPoller := NewPoller()
	require.Equal(t,
		&Poller{interval: time.Second},
		testPoller.WithInterval(time.Second))
}

func TestPoller_WithTimeout(t *testing.T) {
	testPoller := NewPoller()
	require.Equal(t,
		&Poller{
			interval: defaultPollingInterval,
			timeout:  time.Second,
		},
		testPoller.WithTimeout(time.Second))
}

func TestPoller_Poll(t *testing.T) {
	{
		poller := NewPoller().WithInterval(time.Second)
		require.Eventually(t,
			func() bool {
				res, err := poller.Poll(context.Background(),
					newTestMockPollFunc(time.Second, true, "yay", nil))
				return res.(string) == "yay" && err == nil
			},
			5*time.Second,
			time.Second,
			"polling must complete successfully before the timeout")
	}

	{
		poller := NewPoller().WithInterval(time.Second)
		require.Eventually(t,
			func() bool {
				_, err := poller.Poll(context.Background(),
					newTestMockPollFunc(time.Second, true, nil, errors.New("o noes")))
				return err != nil
			},
			5*time.Second,
			time.Second,
			"polling must complete with error before the timeout")
	}

	{
		poller := NewPoller().WithInterval(time.Second)
		require.Never(t,
			func() bool {
				_, err := poller.Poll(context.Background(),
					newTestMockPollFunc(10*time.Second, true, nil, nil))
				return err == nil
			},
			5*time.Second,
			time.Second,
			"polling must NOT complete before the timeout")
	}

	{
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		poller := NewPoller()
		require.Eventually(t,
			func() bool {
				_, err := poller.Poll(ctx,
					newTestMockPollFunc(10*time.Second, true, nil, nil))
				return err != nil
			},
			5*time.Second,
			time.Second,
			"polling must abort on context cancellation")
	}
}

func newTestMockPollFunc(duration time.Duration, done bool, res interface{}, err error) PollFunc {
	return func(_ context.Context) (bool, interface{}, error) {
		time.Sleep(duration)
		return done, res, err
	}
}

func TestClientWithResponses_JobOperationPoller(t *testing.T) {
	// A pending job must return done=false and no error
	{
		mockAPIOperationPending := newTestMockAPIOperationServer(operationStatePending)
		defer mockAPIOperationPending.Close()

		c, err := NewClientWithResponses(mockAPIOperationPending.URL)
		require.NoError(t, err)
		pf := c.OperationPoller("", "")
		done, _, err := pf(context.Background())
		require.NoError(t, err)
		require.False(t, done)
	}

	// A successful job must return done=true and no error
	{
		mockAPIOperationSuccess := newTestMockAPIOperationServer(operationStateSuccess)
		defer mockAPIOperationSuccess.Close()

		c, err := NewClientWithResponses(mockAPIOperationSuccess.URL)
		require.NoError(t, err)
		pf := c.OperationPoller("", "")
		done, _, err := pf(context.Background())
		require.NoError(t, err)
		require.True(t, done)
	}

	// A failed job must return done=true and and an error
	{
		mockAPIOperationFail := newTestMockAPIOperationServer(operationStateFailure)
		defer mockAPIOperationFail.Close()

		c, err := NewClientWithResponses(mockAPIOperationFail.URL)
		require.NoError(t, err)
		pf := c.OperationPoller("", "")
		done, _, err := pf(context.Background())
		require.Error(t, err)
		require.True(t, done)
	}

	// A timed-out job must return done=true and and an error
	{
		mockAPIOperationTimeout := newTestMockAPIOperationServer(operationStateTimeout)
		defer mockAPIOperationTimeout.Close()

		c, err := NewClientWithResponses(mockAPIOperationTimeout.URL)
		require.NoError(t, err)
		pf := c.OperationPoller("", "")
		done, _, err := pf(context.Background())
		require.Error(t, err)
		require.True(t, done)
	}
}

type testMockAPIOperation struct {
	state string
}

func newTestMockAPIOperationServer(state string) *httptest.Server {
	return httptest.NewServer(&testMockAPIOperation{state: state})
}

func (t *testMockAPIOperation) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{` +
		`"id":"021ee8b0-a1a4-11ea-aed0-6329b72edcc5",` +
		`"state":"` + t.state + `",` +
		`"reference":{` +
		`"id":"31161e61-2354-47e6-9df0-36c855ef2a10",` +
		`"command":"some-command",` +
		`"link":"/v2.alpha/some-resource/31161e61-2354-47e6-9df0-36c855ef2a10"` +
		`}}`))
}
