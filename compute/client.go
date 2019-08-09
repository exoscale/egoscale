package compute

import (
	"context"
	"errors"
	"fmt"

	egoerr "github.com/exoscale/egoscale/error"
	egoapi "github.com/exoscale/egoscale/internal/egoscale"
)

const DefaultAPIEndpoint = "https://api.exoscale.com/v1"

// ClientOpts represents an Exoscale Compute API client options.
type ClientOpts struct {
	// APIEndpoint represents the Exoscale Compute API client endpoint to use.
	APIEndpoint string
	// Tracing enables outgoing API calls and received responses display on the process standard error output.
	Tracing bool
}

// Client represents an Exoscale Compute API client.
type Client struct {
	c   *egoapi.Client
	ctx context.Context

	apiEndpoint string
	tracing     bool
}

func (c *Client) String() string {
	return fmt.Sprintf("Client(Endpoint=%q, Key=%q)", c.c.Endpoint, c.c.APIKey)
}

func (c *Client) csError(err error) error {
	if _, ok := err.(*egoapi.ErrorResponse); ok {
		return errors.New(err.(*egoapi.ErrorResponse).ErrorText)
	}

	return err
}

// NewClient returns a new Exoscale Compute API client
func NewClient(ctx context.Context, apiKey, apiSecret string, opts *ClientOpts) (*Client, error) {
	var client = Client{apiEndpoint: DefaultAPIEndpoint}

	if apiKey == "" || apiSecret == "" {
		return nil, egoerr.ErrMissingAPICredentials
	}

	if opts == nil {
		opts = new(ClientOpts)
	}

	if opts.APIEndpoint != "" {
		client.apiEndpoint = opts.APIEndpoint
	}

	client.c = egoapi.NewClient(client.apiEndpoint, apiKey, apiSecret)
	client.ctx = ctx

	if opts.Tracing {
		client.c.TraceOn()
	}

	return &client, nil
}
