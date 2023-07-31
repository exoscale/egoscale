package v3

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/exoscale/egoscale/v3/oapi"
	"github.com/hashicorp/go-retryablehttp"
)

// Client represents Exoscale V3 API Client.
type Client struct {
	endpoint   string
	creds      *Credentials
	httpClient *http.Client

	oapiClient oapi.ClientWithResponses
}

// NewClient returns a new Exoscale API V3 client, or an error if one couldn't be initialized.
// Default HTTP client is [go-retryablehttp] with static retry configuration.
// To change retry configuration, build new HTTP client and pass it using ClientOptWithHTTPClient.
//
// [go-retryablehttp]: https://github.com/hashicorp/go-retryablehttp
func NewClient(endpoint, apiKey, apiSecret string, opts ...ClientOpt) (*Client, error) {
	if apiKey == "" || apiSecret == "" {
		return nil, errors.New("missing or incomplete API credentials")
	}
	client := Client{
		creds: NewCredentials(apiKey, apiSecret),
	}

	for _, opt := range opts {
		if err := opt(&client); err != nil {
			return nil, fmt.Errorf("client configuration error: %w", err)
		}
	}

	// Use retryablehttp client by default
	if client.httpClient == nil {
		rc := retryablehttp.NewClient()
		// silence client
		// TODO: attach to global logger when implemented
		rc.Logger = log.New(io.Discard, "", 0)
		client.httpClient = rc.StandardClient()
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}

	client.oapiClient, err = oapi.NewClient(
		u.String(),
		oapi.WithHTTPClient(client.httpClient),
		oapi.WithRequestEditorFn(SetUserAgent),
		oapi.WithRequestEditorFn(NewSecurityProvider(client.crends).Intercept),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize API client: %w", err)
	}

	return &client, nil
}

// OAPIClient returns configured instance of OpenAPI generated (low-level) API client.
func (c *Client) OAPIClient() *oapi.Client {
	return c.oapiClient
}

// ClientOpt represents a function setting Exoscale API client option.
type ClientOpt func(*Client) error

// ClientOptWithCredentials returns a ClientOpt that sets credentials.
func ClientOptWithCredentials(apiKey, apiSecret string) ClientOpt {
	return func(c *Client) error {
		c.creds = NewCredentials(apiKey, apiSecret)

		return nil
	}
}

// ClientOptWithHTTPClient returns a ClientOpt overriding the default http.Client.
func ClientOptWithHTTPClient(v *http.Client) ClientOpt {
	return func(c *Client) error {
		c.httpClient = v

		return nil
	}
}
