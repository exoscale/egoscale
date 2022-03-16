package oapi

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

// newTestMockPollFunc returns a mocked polling function that sleeps for a specified duration,
// then returns the provided completion flag, resource and error.
func newTestMockPollFunc(duration time.Duration, done bool, res interface{}, err error) PollFunc {
	return func(_ context.Context) (bool, interface{}, error) {
		time.Sleep(duration)
		return done, res, err
	}
}

func (ts *testSuite) TestNewPoller() {
	ts.Require().Equal(
		&Poller{interval: DefaultPollingInterval},
		NewPoller())
}

func (ts *testSuite) TestPoller_WithInterval() {
	poller := NewPoller()
	ts.Require().Equal(&Poller{interval: time.Second}, poller.WithInterval(time.Second))
}

func (ts *testSuite) TestPoller_WithTimeout() {
	poller := NewPoller()
	ts.Require().Equal(&Poller{
		interval: DefaultPollingInterval,
		timeout:  time.Second,
	},
		poller.WithTimeout(time.Second),
	)
}

func (ts *testSuite) TestPoller_Poll() {
	tests := []struct {
		name     string
		testFunc func(ctx context.Context, ts *testSuite) func() bool
	}{
		{
			name: "failure by timeout",
			testFunc: func(ctx context.Context, ts *testSuite) func() bool {
				return func() bool {
					_, err := NewPoller().
						Poll(
							ctx,
							newTestMockPollFunc(10*time.Second, true, nil, nil),
						)
					ts.Require().ErrorIs(err, context.DeadlineExceeded)

					return true
				}
			},
		},
		{
			name: "failure by error during polling",
			testFunc: func(ctx context.Context, ts *testSuite) func() bool {
				return func() bool {
					_, err := NewPoller().
						WithInterval(time.Second).
						Poll(
							ctx,
							newTestMockPollFunc(time.Second, true, nil, errors.New("o noes")),
						)
					return err != nil
				}
			},
		},
		{
			name: "ok",
			testFunc: func(ctx context.Context, ts *testSuite) func() bool {
				return func() bool {
					res, err := NewPoller().
						WithInterval(time.Second).
						Poll(
							ctx,
							newTestMockPollFunc(time.Second, true, "yay", nil),
						)
					return res.(string) == "yay" && err == nil
				}
			},
		},
	}

	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			ts.Require().Eventually(
				tt.testFunc(ctx, ts),
				5*time.Second,
				100*time.Millisecond,
			)
		})
	}
}

func (ts *testSuite) TestOperationPoller() {
	var (
		operationID           = ts.randomID()
		operationReferenceID  = ts.randomID()
		operationStatePending = OperationStatePending
		operationStateSuccess = OperationStateSuccess
		operationStateFailure = OperationStateFailure
		operationStateTimeout = OperationStateTimeout
	)

	tests := []struct {
		name               string
		httpResponseStatus int
		operationResponse  *Operation
		setupFunc          func(ts *testSuite)
		testFunc           func(ts *testSuite, done bool, res interface{}, err error)
	}{
		{
			name:               "pending",
			httpResponseStatus: http.StatusOK,
			operationResponse: &Operation{
				Id:        &operationID,
				State:     &operationStatePending,
				Reference: NewReference(nil, &operationReferenceID, nil),
			},
			testFunc: func(ts *testSuite, done bool, res interface{}, err error) {
				ts.Require().NoError(err)
				ts.Require().False(done)
			},
		},
		{
			name:               "success",
			httpResponseStatus: http.StatusOK,
			operationResponse: &Operation{
				Id:        &operationID,
				State:     &operationStateSuccess,
				Reference: NewReference(nil, &operationReferenceID, nil),
			},
			testFunc: func(ts *testSuite, done bool, res interface{}, err error) {
				ts.Require().NoError(err)
				ts.Require().Equal(NewReference(nil, &operationReferenceID, nil), res)
				ts.Require().True(done)
			},
		},
		{
			name:               "failure",
			httpResponseStatus: http.StatusOK,
			operationResponse: &Operation{
				Id:        &operationID,
				State:     &operationStateFailure,
				Reference: NewReference(nil, &operationReferenceID, nil),
			},
			testFunc: func(ts *testSuite, done bool, res interface{}, err error) {
				ts.Require().Error(err)
				ts.Require().True(done)
			},
		},
		{
			name:               "timeout",
			httpResponseStatus: http.StatusOK,
			operationResponse: &Operation{
				Id:        &operationID,
				State:     &operationStateTimeout,
				Reference: NewReference(nil, &operationReferenceID, nil),
			},
			testFunc: func(ts *testSuite, done bool, res interface{}, err error) {
				ts.Require().Error(err)
				ts.Require().True(done)
			},
		},
		{
			name:               "API error",
			httpResponseStatus: http.StatusInternalServerError,
			testFunc: func(ts *testSuite, done bool, res interface{}, err error) {
				ts.Require().Error(err)
			},
		},
	}

	for _, tt := range tests {
		// Reset the OAPI client mock calls stack between test cases
		ts.client.(*oapiClientMock).ExpectedCalls = nil

		ts.T().Run(tt.name, func(t *testing.T) {
			ts.client.(*oapiClientMock).
				On("GetOperationWithResponse", mock.Anything, mock.Anything, ([]RequestEditorFn)(nil)).
				Return(&GetOperationResponse{
					HTTPResponse: &http.Response{StatusCode: tt.httpResponseStatus},
					JSON200:      tt.operationResponse,
				}, nil)

			done, res, err := OperationPoller(ts.client, "", ts.randomID())(context.Background())
			tt.testFunc(ts, done, res, err)
		})
	}
}
