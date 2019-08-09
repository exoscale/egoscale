package storage

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	egoerr "github.com/exoscale/egoscale/error"
)

// DefaultZone is the default Storage API zone.
const DefaultZone = "ch-gva-2"

// ClientOpts represents the Exoscale Storage API client options.
type ClientOpts struct {
	// APIEndpoint represents the Exoscale Storage API client endpoint to use.
	APIEndpoint string
	// Zone represents the Exoscale Storage API client default zone.
	Zone string
	// Tracing enables outgoing API calls and received responses display on the process standard error output.
	Tracing bool
}

// Client represents an Exoscale Storage API client.
type Client struct {
	s   *session.Session
	c   *s3.S3
	ctx context.Context

	apiEndpoint string
	zone        string
	tracing     bool
}

// NewClient returns a new Exoscale Storage API client.
func NewClient(ctx context.Context, apiKey, apiSecret string, opts *ClientOpts) (*Client, error) {
	var (
		client = Client{zone: DefaultZone}
		err    error
	)

	if apiKey == "" || apiSecret == "" {
		return nil, egoerr.ErrMissingAPICredentials
	}

	if opts.Zone != "" {
		client.zone = opts.Zone
	}

	if opts.APIEndpoint != "" {
		client.apiEndpoint = opts.APIEndpoint
	} else {
		client.apiEndpoint = fmt.Sprintf("https://sos-%s.exo.io", client.zone)
	}

	if client.s, err = session.NewSessionWithOptions(session.Options{Config: aws.Config{
		Region:      aws.String(client.zone),
		Endpoint:    aws.String(client.apiEndpoint),
		Credentials: credentials.NewStaticCredentials(apiKey, apiSecret, ""),
		// TODO: support tracing using https://godoc.org/github.com/aws/aws-sdk-go/aws#Config
	}}); err != nil {
		return nil, err
	}

	client.tracing = opts.Tracing
	client.c = s3.New(client.s)
	client.ctx = ctx

	return &client, nil
}
