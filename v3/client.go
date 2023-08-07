package v3

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/exoscale/egoscale/v3/oapi"
	"github.com/hashicorp/go-retryablehttp"
)

const (
	EnvKeyAPIEndpoint = "EXOSCALE_API_ENDPOINT"
)

// Client represents Exoscale V3 API Client.
type Client struct {
	creds      *Credentials
	oapiClient *oapi.ClientWithResponses
}

// NewClient returns a new Exoscale API V3 client, or an error if one couldn't be initialized.
// Client is generic (single EP) with no concept of zones/environments.
// For zone-aware client use ZonedClient.
// Default HTTP client is [go-retryablehttp] with static retry configuration.
// To change retry configuration, build new HTTP client and pass it using ClientOptWithHTTPClient.
// API credentials must be passed with ClientOptWithCredentials.
// If EXOSCALE_API_ENDPOINT environment variable is set, it replaces endpoint.
func NewClient(endpoint string, opts ...ClientOpt) (*Client, error) {
	// Env var override
	if h := os.Getenv(EnvKeyAPIEndpoint); h != "" {
		endpoint = h
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	config := ClientConfig{
		requestEditors: []oapi.RequestEditorFn{},
	}
	for _, opt := range opts {
		if err := opt(&config); err != nil {
			return nil, fmt.Errorf("client configuration error: %w", err)
		}
	}

	client := Client{
		creds: config.creds,
	}

	// Use retryablehttp client by default
	if config.httpClient == nil {
		rc := retryablehttp.NewClient()
		// TODO: attach to global logger when implemented
		rc.Logger = log.New(os.Stderr, "", 0)
		config.httpClient = rc.StandardClient()
	}

	// Mandatory oapi options.
	oapiOpts := []oapi.ClientOption{
		oapi.WithHTTPClient(config.httpClient),
		oapi.WithRequestEditorFn(NewUserAgentProvider(config.uaPrefix).Intercept),
	}

	// We are adding security middleware only if API credentials are specified
	// in order to allow generic usage and local testing.
	// TODO: add log line emphasizing the lack of credentials.
	if client.creds != nil {
		oapiOpts = append(
			oapiOpts,
			oapi.WithRequestEditorFn(NewSecurityProvider(client.creds).Intercept),
		)
	}

	// Attach any custom request editors
	for _, editor := range config.requestEditors {
		oapiOpts = append(
			oapiOpts,
			oapi.WithRequestEditorFn(editor),
		)
	}

	client.oapiClient, err = oapi.NewClientWithResponses(
		u.String(),
		oapiOpts...,
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
