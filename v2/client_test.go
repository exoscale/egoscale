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

func TestDefaultClient_Retry(t *testing.T) {
	n := 0
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if n < 2 {
			n++
			w.WriteHeader(500)
		}
	}))
	defer testServer.Close()

	testClient, err := NewClient(
		"EXOxxxxxxxxxxxxxxxxxxxxxxxx",
		"XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
		ClientOptWithAPIEndpoint(testServer.URL),
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = testClient.httpClient.Get(testServer.URL)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetEndpointFromContext(t *testing.T) {
	ctx := context.Background()
	zonedCtx := api.WithEndpoint(ctx, api.NewReqEndpoint("api", "ch-gva-2"))
	environmentCtx := api.WithEndpoint(ctx, api.NewReqEndpoint("test", "ch-gva-2"))

	for _, tt := range []struct {
		name     string
		ctx      context.Context
		url      string
		wantHost string
	}{
		{name: "without endpoint context", ctx: ctx, url: "https://api.exoscale.com/v2/ssh-key", wantHost: "api.exoscale.com"},
		{name: "default endpoint", ctx: zonedCtx, url: "https://api.exoscale.com/v2/ssh-key", wantHost: "api-ch-gva-2.exoscale.com"},
		{name: "default endpoint with environment", ctx: environmentCtx, url: "https://api.exoscale.com/v2/ssh-key", wantHost: "test-ch-gva-2.exoscale.com"},
		{name: "environment endpoint", ctx: environmentCtx, url: "https://test.exoscale.com/v2/ssh-key", wantHost: "test-ch-gva-2.exoscale.com"},
		{name: "custom hostname", ctx: zonedCtx, url: "https://gateway.internal:8443/v2/ssh-key", wantHost: "gateway.internal:8443"},
		{name: "custom IPv4", ctx: zonedCtx, url: "http://127.0.0.1:8080/v2/ssh-key", wantHost: "127.0.0.1:8080"},
		{name: "custom IPv6", ctx: zonedCtx, url: "http://[::1]/v2/ssh-key", wantHost: "[::1]"},
		{name: "custom scheme", ctx: zonedCtx, url: "http://api.exoscale.com/v2/ssh-key", wantHost: "api.exoscale.com"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, tt.url, nil)
			require.NoError(t, err)
			require.NoError(t, setEndpointFromContext(tt.ctx, req))
			require.Equal(t, tt.wantHost, req.URL.Host)
			require.Equal(t, tt.wantHost, req.Host)
		})
	}
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
