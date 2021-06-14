package v2

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/exoscale/egoscale/v2/api"
	papi "github.com/exoscale/egoscale/v2/internal/public-api"
)

var testSeededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

type dummyResource struct {
	id string
}

func (d dummyResource) get(_ context.Context, _ *Client, _, id string) (interface{}, error) {
	return &dummyResource{id: id}, nil
}

type clientTestSuite struct {
	client *Client

	suite.Suite
}

func (ts *clientTestSuite) SetupTest() {
	httpmock.Activate()

	client, err := NewClient("x", "x", ClientOptWithPollInterval(10*time.Millisecond))
	if err != nil {
		ts.T().Fatal(err)
	}

	// Overriding the internal public API client with mocked HTTP client.
	client.ClientWithResponses, err = papi.NewClientWithResponses("", papi.WithHTTPClient(client.httpClient))
	if err != nil {
		ts.T().Fatal(err)
	}

	ts.client = client
}

func (ts *clientTestSuite) TearDownTest() {
	ts.client = nil

	httpmock.DeactivateAndReset()
}

func (ts *clientTestSuite) mockAPIRequest(method, url string, body interface{}) {
	httpmock.RegisterResponder(method, url, func(_ *http.Request) (*http.Response, error) {
		resp, err := httpmock.NewJsonResponse(http.StatusOK, body)
		if err != nil {
			ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
		}
		return resp, nil
	})
}

func (ts *clientTestSuite) randomID() string {
	id, err := uuid.NewV4()
	if err != nil {
		ts.T().Fatalf("unable to generate a new UUID: %s", err)
	}
	return id.String()
}

func (ts *clientTestSuite) randomStringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[testSeededRand.Intn(len(charset))]
	}
	return string(b)
}

func (ts *clientTestSuite) randomString(length int) string {
	const defaultCharset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	return ts.randomStringWithCharset(length, defaultCharset)
}

func (ts *clientTestSuite) unmarshalJSONRequestBody(req *http.Request, v interface{}) {
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		ts.T().Fatalf("error reading request body: %s", err)
	}
	if err = json.Unmarshal(data, v); err != nil {
		ts.T().Fatalf("error while unmarshalling JSON body: %s", err)
	}
}

func (ts *clientTestSuite) TestClient_SetHTTPClient() {
	testHTTPClient := http.DefaultClient

	client := new(Client)
	client.SetHTTPClient(testHTTPClient)

	ts.Require().Equal(testHTTPClient, client.httpClient)
}

func (ts *clientTestSuite) TestClient_SetTimeout() {
	testTimeout := 5 * time.Minute

	client := new(Client)
	client.SetTimeout(testTimeout)

	ts.Require().Equal(testTimeout, client.timeout)
}

func (ts *clientTestSuite) TestClient_SetTrace() {
	client := new(Client)
	client.SetTrace(true)

	ts.Require().Equal(true, client.trace)
}

func (ts *clientTestSuite) TestClient_fetchfromIDs() {
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

func (ts *clientTestSuite) TestDefaultTransport_RoundTrip() {
	var ok bool

	httpmock.RegisterResponder("GET", "/zone",
		func(req *http.Request) (*http.Response, error) {
			ts.Require().Equal(UserAgent, req.Header.Get("User-Agent"))
			ok = true

			resp, err := httpmock.NewJsonResponse(http.StatusOK, struct {
				Zones *[]papi.Zone `json:"zones,omitempty"`
			}{
				Zones: new([]papi.Zone),
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}
			return resp, nil
		})

	_, err := ts.client.ListZones(context.Background())
	ts.Require().NoError(err)
	ts.Require().True(ok)
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
		ClientOptWithAPIEndpoint(testAPIEndpoint),
		ClientOptWithHTTPClient(testHTTPClient),
		ClientOptWithTimeout(testTimeout),
		ClientOptWithPollInterval(testPollInterval),
	)

	require.NoError(t, err)
	require.Equal(t, testAPIKey, client.apiKey)
	require.Equal(t, testAPISecret, client.apiSecret)
	require.Equal(t, testAPIEndpoint+api.Prefix, client.apiEndpoint)
	require.Equal(t, testHTTPClient, client.httpClient)
	require.Equal(t, testTimeout, client.timeout)
	require.Equal(t, testPollInterval, client.pollInterval)
	require.IsType(t, &api.ErrorHandlerMiddleware{}, client.httpClient.Transport)
}

func TestSuiteClientTestSuite(t *testing.T) {
	suite.Run(t, new(clientTestSuite))
}
