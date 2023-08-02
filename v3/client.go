package v3

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/exoscale/egoscale/v3/oapi"
	"github.com/hashicorp/go-retryablehttp"
)

// Client represents Exoscale V3 API Client.
type Client struct {
	endpoint   string
	creds      *Credentials
	httpClient *http.Client

	requestEditors []oapi.RequestEditorFn
	// TODO: implement response editors (not available in oapi, should be embeded in consumer API.

	oapiClient *oapi.ClientWithResponses
}

// NewClient returns a new Exoscale API V3 client, or an error if one couldn't be initialized.
// Client is generic (single EP) with no concept of zones/environments.
// Default HTTP client is [go-retryablehttp] with static retry configuration.
// To change retry configuration, build new HTTP client and pass it using ClientOptWithHTTPClient.
// API credentials must be passed with ClientOptWithCredentials.
func NewClient(endpoint string, opts ...ClientOpt) (*Client, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	client := Client{
		endpoint:       endpoint,
		requestEditors: []oapi.RequestEditorFn{},
	}

	// Use retryablehttp client by default
	rc := retryablehttp.NewClient()
	// TODO: attach to global logger when implemented
	rc.Logger = log.New(os.Stderr, "", 0)
	client.httpClient = rc.StandardClient()

	for _, opt := range opts {
		if err := opt(&client); err != nil {
			return nil, fmt.Errorf("client configuration error: %w", err)
		}
	}

	oapiOpts := []oapi.ClientOption{
		oapi.WithHTTPClient(client.httpClient),
		oapi.WithRequestEditorFn(SetUserAgent),
	}

	// We are adding security middleware only if API credentials are specified
	// in order to allow generic usage and local testing.
	// In production consumers are expected to check that non-empty credentials are set
	// before initializing client.
	// TODO: add log line emphasizing the lack of credentials.
	if client.creds != nil {
		oapiOpts = append(
			oapiOpts,
			oapi.WithRequestEditorFn(NewSecurityProvider(client.creds).Intercept),
		)
	}

	// Attach any custom request editors
	for _, editor := range client.requestEditors {
		oapiOpts = append(
			oapiOpts,
			oapi.WithRequestEditorFn(editor),
		)
	}

	client.oapiClient, err = oapi.NewClientWithResponses(
		u.String(),
		oapi.WithHTTPClient(client.httpClient),
		oapi.WithRequestEditorFn(SetUserAgent),
		oapi.WithRequestEditorFn(NewSecurityProvider(client.creds).Intercept),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize API client: %w", err)
	}

	return &client, nil
}

// OAPIClient returns configured instance of OpenAPI generated (low-level) API client.
func (c *Client) OAPIClient() *oapi.ClientWithResponses {
	return c.oapiClient
}
