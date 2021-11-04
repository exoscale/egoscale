package v2

import (
	"context"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/exoscale/egoscale/v2/api"
	"github.com/exoscale/egoscale/v2/oapi"
)

var testSeededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

type dummyResource struct {
	id string
}

func (d dummyResource) get(_ context.Context, _ *Client, _, id string) (interface{}, error) {
	return &dummyResource{id: id}, nil
}

type testSuite struct {
	suite.Suite

	client *Client
}

func (ts *testSuite) SetupTest() {
	ts.client = &Client{
		oapiClient:   new(oapiClientMock),
		pollInterval: 10 * time.Millisecond,
	}
}

func (ts *testSuite) TearDownTest() {
	ts.client = nil
}

func (ts *testSuite) mock() *oapiClientMock {
	return ts.client.oapiClient.(*oapiClientMock)
}

func (ts *testSuite) mockGetOperation(o *oapi.Operation) {
	ts.mock().
		On("GetOperationWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(
			&oapi.GetOperationResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200:      o,
			},
			nil,
		)
}

func (ts *testSuite) randomID() string {
	id, err := uuid.NewV4()
	if err != nil {
		ts.T().Fatalf("unable to generate a new UUID: %s", err)
	}
	return id.String()
}

func (ts *testSuite) randomStringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[testSeededRand.Intn(len(charset))]
	}
	return string(b)
}

func (ts *testSuite) randomString(length int) string {
	const defaultCharset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	return ts.randomStringWithCharset(length, defaultCharset)
}

func (ts *testSuite) TestClient_SetHTTPClient() {
	testHTTPClient := http.DefaultClient

	client := new(Client)
	client.SetHTTPClient(testHTTPClient)

	ts.Require().Equal(testHTTPClient, client.httpClient)
}

func (ts *testSuite) TestClient_SetTimeout() {
	testTimeout := 5 * time.Minute

	client := new(Client)
	client.SetTimeout(testTimeout)

	ts.Require().Equal(testTimeout, client.timeout)
}

func (ts *testSuite) TestClient_SetTrace() {
	client := new(Client)
	client.SetTrace(true)

	ts.Require().Equal(true, client.trace)
}

func (ts *testSuite) TestClient_fetchfromIDs() {
	type args struct {
		ctx  context.Context
		zone string
		ids  []string
		rt   interface{}
	}

	tests := []struct {
		name     string
		args     args
		expected interface{}
		wantErr  bool
	}{
		{
			name: "with nil resource type",
			args: args{
				ids: nil,
				rt:  nil,
			},
			wantErr: true,
		},
		{
			name: "with concrete resource type",
			args: args{
				ids: []string{"id1", "id2"},
				rt:  dummyResource{},
			},
			wantErr: true,
		},
		{
			name: "with empty ids",
			args: args{
				ctx:  context.Background(),
				zone: testZone,
				ids:  nil,
				rt:   new(dummyResource),
			},
			expected: []*dummyResource{},
		},
		{
			name: "ok",
			args: args{
				ctx:  context.Background(),
				zone: testZone,
				ids:  []string{"id1", "id2"},
				rt:   new(dummyResource),
			},
			expected: []*dummyResource{{id: "id1"}, {id: "id2"}},
		},
	}

	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			actual, err := ts.client.fetchFromIDs(tt.args.ctx, tt.args.zone, tt.args.ids, tt.args.rt)
			if err != nil != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			ts.Require().Equal(tt.expected, actual)
		})
	}
}

func (ts *testSuite) TestDefaultTransport_RoundTrip() {
	testServer := httptest.NewServer(nil)
	defer testServer.Close()

	testClient := testServer.Client()
	testClient.Transport = &defaultTransport{next: testClient.Transport}

	req, err := http.NewRequest(http.MethodGet, testServer.URL, nil)
	ts.Require().NoError(err)

	_, err = testClient.Do(req)
	ts.Require().NoError(err)
	ts.Require().Equal(UserAgent, req.Header.Get("User-Agent"))
}

func TestSetEndpointFromContext(t *testing.T) {
	var (
		ctx                = context.Background()
		testReqEndpointEnv = "api"
		testURL            = "https://www.example.net/test.txt"
		req, _             = http.NewRequest("GET", testURL, nil)
	)

	// With empty context
	err := setEndpointFromContext(ctx, req)
	require.NoError(t, err)
	require.Equal(t, testURL, req.URL.String())

	// With augmented context
	reqEndpoint := api.NewReqEndpoint(testReqEndpointEnv, "")
	err = setEndpointFromContext(api.WithEndpoint(ctx, reqEndpoint), req)
	require.NoError(t, err)
	require.Equal(t, reqEndpoint.Host(), req.URL.Host)
}

func TestNewClient(t *testing.T) {
	var (
		testAPIKey        = "EXOxxxxxxxxxxxxxxxxxxxxxxxx"
		testAPISecret     = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
		testAPIEndpoint   = "https://example.net/"
		testHTTPTransport = http.Transport{}
		testHTTPClient    = &http.Client{Transport: &testHTTPTransport}
		testTimeout       = 5 * time.Second
		testPollInterval  = 10 * time.Second
	)

	client, err := NewClient(
		testAPIKey,
		testAPISecret,
		ClientOptCond(func() bool { return true }, ClientOptWithTrace()),
		ClientOptWithAPIEndpoint(testAPIEndpoint),
		ClientOptWithHTTPClient(testHTTPClient),
		ClientOptWithPollInterval(testPollInterval),
		ClientOptWithTimeout(testTimeout),
	)

	require.NoError(t, err)
	require.Equal(t, testAPIKey, client.apiKey)
	require.Equal(t, testAPISecret, client.apiSecret)
	require.Equal(t, testAPIEndpoint+api.Prefix, client.apiEndpoint)
	require.Equal(t, testHTTPClient, client.httpClient)
	require.Equal(t, testTimeout, client.timeout)
	require.Equal(t, testPollInterval, client.pollInterval)
	require.True(t, client.trace)
	require.IsType(t, &api.ErrorHandlerMiddleware{}, client.httpClient.Transport)
}

func TestSuiteClientTestSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}
