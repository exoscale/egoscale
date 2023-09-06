package v3

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

const (
	EnvKeyAPIEndpoint = "EXOSCALE_API_ENDPOINT"

	PollingInterval = 3 * time.Second
)

// Client represents Exoscale V3 API Client.
type Client struct {
	server     string
	httpClient *http.Client
	reqEditors []RequestEditorFn
	logger     Logger

	creds *Credentials
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

	// Validate that string is URL
	_, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	// Ensure the server URL always has a trailing slash
	if !strings.HasSuffix(endpoint, "/") {
		endpoint += "/"
	}

	config := ClientConfig{
		requestEditors: []RequestEditorFn{},
	}
	for _, opt := range opts {
		if err := opt(&config); err != nil {
			return nil, fmt.Errorf("client configuration error: %w", err)
		}
	}

	// Use retryablehttp client by default
	if config.httpClient == nil {
		rc := retryablehttp.NewClient()
		rc.Logger = log.New(io.Discard, "", 0)
		if config.logger != nil {
			rc.Logger = config.logger
		}
		config.httpClient = rc.StandardClient()
	}

	// Use dummy discard logger if none provided.
	if config.logger == nil {
		config.logger = NewStandardLogger(io.Discard, false)
	}

	client := Client{
		server:     endpoint,
		httpClient: config.httpClient,
		logger:     config.logger,
		creds:      config.creds,
	}

	// Mandatory request editors.
	fns := []RequestEditorFn{
		NewUserAgentProvider(config.uaPrefix).Intercept,
	}

	// We are adding security middleware only if API credentials are specified
	// in order to allow generic usage and local testing.
	// TODO: add log line emphasizing the lack of credentials
	if client.creds != nil {
		fns = append(fns, NewSecurityProvider(client.creds).Intercept)
	}

	// Attach any custom request editors
	for _, editor := range config.requestEditors {
		fns = append(fns, editor)
	}

	client.reqEditors = fns

	return &client, nil
}

// Wait is a helper that waits for async operation to reach the final state.
// Final states are one of: failure, success, timeout.
func (c *Client) Wait(
	ctx context.Context,
	f func(ctx context.Context) (*Operation, error),
) (*Operation, error) {
	ticker := time.NewTicker(PollingInterval)
	defer ticker.Stop()

	op, err := f(ctx)
	if err != nil {
		return nil, err
	}
	// Exit right away if operation is already done.
	if *op.State != OperationStatePending {
		return op, nil
	}

	for {
		select {
		case <-ticker.C:
			op, err := c.Global().Operations().Get(ctx, *op.Id)
			if err != nil {
				return nil, err
			}
			if *op.State != OperationStatePending {
				continue
			}

			return op, nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

// IAM provides access to IAM resources on Exoscale platform.
func (c *Client) IAM() *IAMAPI {
	return &IAMAPI{c}
}

// DBaaS provides access to DBaaS resources on Exoscale platform.
func (c *Client) DBaaS() *DBaaSAPI {
	return &DBaaSAPI{c}
}

// Compute provides access to Compute resources on Exoscale platform.
func (c *Client) Compute() *ComputeAPI {
	return &ComputeAPI{c}
}

// DNS provides access to DNS resources on Exoscale platform.
func (c *Client) DNS() *DNSAPI {
	return &DNSAPI{c}
}

// Global provides access to global resources on Exoscale platform.
func (c *Client) Global() *GlobalAPI {
	return &GlobalAPI{c}
}
