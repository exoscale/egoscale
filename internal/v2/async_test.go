package v2

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
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
	var (
		operationID              = "021ee8b0-a1a4-11ea-aed0-6329b72edcc5"
		mockOperationReferenceID = "31161e61-2354-47e6-9df0-36c855ef2a10"

		newTestClient = func(state string) (*ClientWithResponses, error) {
			mockClient := NewMockClient()
			mockClient.RegisterResponder("GET", "/operation/"+operationID,
				func(req *http.Request) (*http.Response, error) {
					resp, err := httpmock.NewJsonResponse(http.StatusOK, Operation{
						Id:        &operationID,
						State:     &state,
						Reference: &Resource{Id: &mockOperationReferenceID},
					})
					if err != nil {
						t.Fatalf("error initializing mock HTTP responder: %s", err)
					}
					return resp, nil
				})

			return NewClientWithResponses("", WithHTTPClient(mockClient))
		}
	)

	// A pending job must return done=false and no error
	{
		c, err := newTestClient(operationStatePending)
		require.NoError(t, err)
		done, _, err := c.OperationPoller("", operationID)(context.Background())
		require.NoError(t, err)
		require.False(t, done)
	}

	// A successful job must return done=true and no error
	{
		c, err := newTestClient(operationStateSuccess)
		require.NoError(t, err)
		done, res, err := c.OperationPoller("", operationID)(context.Background())
		require.NoError(t, err)
		require.Equal(t, &Resource{Id: &mockOperationReferenceID}, res)
		require.True(t, done)
	}

	// A failed job must return done=true and and an error
	{
		c, err := newTestClient(operationStateFailure)
		require.NoError(t, err)
		done, _, err := c.OperationPoller("", operationID)(context.Background())
		require.Error(t, err)
		require.True(t, done)
	}

	// A timed-out job must return done=true and and an error
	{
		c, err := newTestClient(operationStateTimeout)
		require.NoError(t, err)
		done, _, err := c.OperationPoller("", operationID)(context.Background())
		require.Error(t, err)
		require.True(t, done)
	}
}
