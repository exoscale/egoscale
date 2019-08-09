package runstatus

import (
	"context"

	egoerr "github.com/exoscale/egoscale/error"
	egoapi "github.com/exoscale/egoscale/internal/egoscale"
)

const DefaultAPIEndpoint = "https://api.runstatus.com"

// ClientOpts represents the Exoscale Runstatus API client option.
type ClientOpts struct {
	// APIEndpoint represents the Exoscale Runstatus API client endpoint to use.
	APIEndpoint string
	// Tracing enables outgoing API calls and received responses display on the process standard error output.
	Tracing bool
}

// Client represents an Exoscale Storage API client.
type Client struct {
	c   *egoapi.Client
	ctx context.Context

	apiEndpoint string
	tracing     bool
}

// NewClient returns a new Exoscale Runstatus API client.
func NewClient(ctx context.Context, apiKey, apiSecret string, opts *ClientOpts) (*Client, error) {
	var client = Client{apiEndpoint: DefaultAPIEndpoint}

	if apiKey == "" || apiSecret == "" {
		return nil, egoerr.ErrMissingAPICredentials
	}

	if opts.APIEndpoint != "" {
		client.apiEndpoint = opts.APIEndpoint
	}

	client.tracing = opts.Tracing
	client.c = egoapi.NewClient(client.apiEndpoint, apiKey, apiSecret)
	client.ctx = ctx

	return &client, nil
}
