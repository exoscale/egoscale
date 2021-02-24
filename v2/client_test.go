package v2

import (
	"context"
	"encoding/json"
	"io/ioutil"
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

type clientTestSuite struct {
	client *Client

	suite.Suite
}

func (ts *clientTestSuite) SetupTest() {
	client, err := NewClient("x", "x")
	if err != nil {
		ts.T().Fatal(err)
	}

	// Overriding the internal public API client with mocked HTTP client.
	client.ClientWithResponses, err = papi.NewClientWithResponses("", papi.WithHTTPClient(client.httpClient))
	if err != nil {
		ts.T().Fatal(err)
	}

	httpmock.ActivateNonDefault(client.httpClient)

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

func (ts *clientTestSuite) unmarshalJSONRequestBody(req *http.Request, v interface{}) {
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		ts.T().Fatalf("error reading request body: %s", err)
	}
	if err = json.Unmarshal(data, v); err != nil {
		ts.T().Fatalf("error while unmarshalling JSON body: %s", err)
	}
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
	)

	client, err := NewClient(
		testAPIKey,
		testAPISecret,
		ClientOptWithAPIEndpoint(testAPIEndpoint),
		ClientOptWithHTTPClient(testHTTPClient),
		ClientOptWithTimeout(testTimeout),
	)

	require.NoError(t, err)
	require.Equal(t, testAPIKey, client.apiKey)
	require.Equal(t, testAPISecret, client.apiSecret)
	require.Equal(t, testAPIEndpoint+api.Prefix, client.apiEndpoint)
	require.Equal(t, testHTTPClient, client.httpClient)
	require.Equal(t, testTimeout, client.timeout)
	require.IsType(t, &api.ErrorHandlerMiddleware{}, client.httpClient.Transport)
}

func TestSuiteClientTestSuite(t *testing.T) {
	suite.Run(t, new(clientTestSuite))
}
