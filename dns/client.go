package dns

import (
	"context"
	"errors"

	egoerr "github.com/exoscale/egoscale/error"
	egoapi "github.com/exoscale/egoscale/internal/egoscale"
)

const DefaultAPIEndpoint = "https://api.exoscale.com/v1"

// ClientOpts represents an Exoscale DNS API client option.
type ClientOpts struct {
	// Endpoint represents the Exoscale DNS API client endpoint to use.
	APIEndpoint string
	// Tracing enables outgoing API calls and received responses display on the process standard error output.
	Tracing bool
}

// Client represents an Exoscale DNS API client.
type Client struct {
	c   *egoapi.Client
	ctx context.Context

	apiEndpoint string
	tracing     bool
}

func (c *Client) csError(err error) error {
	if _, ok := err.(*egoapi.ErrorResponse); ok {
		return errors.New(err.(*egoapi.ErrorResponse).ErrorText)
	}

	return err
}

// NewClient returns a new Exoscale DNS API client.
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

	client.tracing = opts.Tracing
	client.c = egoapi.NewClient(client.apiEndpoint, apiKey, apiSecret)
	client.ctx = ctx

	return &client, nil
}
