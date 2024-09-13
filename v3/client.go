// Package v3 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/egoscale/v3/generator version v0.0.1 DO NOT EDIT.
package v3

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/exoscale/egoscale/v3/credentials"
	"github.com/go-playground/validator/v10"
)

// Endpoint represents a zone endpoint.
type Endpoint string

const (
	CHGva2 Endpoint = "https://api-ch-gva-2.exoscale.com/v2"
	CHDk2  Endpoint = "https://api-ch-dk-2.exoscale.com/v2"
	DEFra1 Endpoint = "https://api-de-fra-1.exoscale.com/v2"
	DEMuc1 Endpoint = "https://api-de-muc-1.exoscale.com/v2"
	ATVie1 Endpoint = "https://api-at-vie-1.exoscale.com/v2"
	ATVie2 Endpoint = "https://api-at-vie-2.exoscale.com/v2"
	BGSof1 Endpoint = "https://api-bg-sof-1.exoscale.com/v2"
)

func (c Client) GetZoneName(ctx context.Context, endpoint Endpoint) (ZoneName, error) {
	resp, err := c.ListZones(ctx)
	if err != nil {
		return "", fmt.Errorf("get zone name: list zones: %w", err)
	}

	zone, err := resp.FindZone(string(endpoint))
	if err != nil {
		return "", fmt.Errorf("get zone name: find zone: %w", err)
	}

	return zone.Name, nil
}

func (c Client) GetZoneAPIEndpoint(ctx context.Context, zoneName ZoneName) (Endpoint, error) {
	resp, err := c.ListZones(ctx)
	if err != nil {
		return "", fmt.Errorf("get zone api endpoint: list zones: %w", err)
	}

	zone, err := resp.FindZone(string(zoneName))
	if err != nil {
		return "", fmt.Errorf("get zone api endpoint: find zone: %w", err)
	}

	return zone.APIEndpoint, nil
}

// Client represents an Exoscale API client.
type Client struct {
	apiKey          string
	apiSecret       string
	userAgent       string
	serverEndpoint  string
	httpClient      *http.Client
	pollingInterval time.Duration
	validate        *validator.Validate
	trace           bool

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	requestInterceptors []RequestInterceptorFn
}

// RequestInterceptorFn is the function signature for the RequestInterceptor callback function
type RequestInterceptorFn func(ctx context.Context, req *http.Request) error

// Deprecated: use ClientUserAgent instead.
var UserAgent = getDefaultUserAgent()

const pollingInterval = 3 * time.Second

// ClientOpt represents a function setting Exoscale API client option.
type ClientOpt func(*Client) error

// ClientOptWithTrace returns a ClientOpt enabling HTTP request/response tracing.
func ClientOptWithTrace() ClientOpt {
	return func(c *Client) error {
		c.trace = true
		return nil
	}
}

// ClientOptWithUserAgent returns a ClientOpt setting the user agent header.
func ClientOptWithUserAgent(ua string) ClientOpt {
	return func(c *Client) error {
		c.userAgent = ua + " " + getDefaultUserAgent()
		return nil
	}
}

// ClientOptWithValidator returns a ClientOpt with a given validator.
func ClientOptWithValidator(validate *validator.Validate) ClientOpt {
	return func(c *Client) error {
		c.validate = validate
		return nil
	}
}

// ClientOptWithEndpoint returns a ClientOpt With a given zone Endpoint.
func ClientOptWithEndpoint(endpoint Endpoint) ClientOpt {
	return func(c *Client) error {
		c.serverEndpoint = string(endpoint)
		return nil
	}
}

// ClientOptWithRequestInterceptors returns a ClientOpt With given RequestInterceptors.
func ClientOptWithRequestInterceptors(f ...RequestInterceptorFn) ClientOpt {
	return func(c *Client) error {
		c.requestInterceptors = append(c.requestInterceptors, f...)
		return nil
	}
}

// ClientOptWithHTTPClient returns a ClientOpt overriding the default http.Client.
// Note: the Exoscale API client will chain additional middleware
// (http.RoundTripper) on the HTTP client internally, which can alter the HTTP
// requests and responses. If you don't want any other middleware than the ones
// currently set to your HTTP client, you should duplicate it and pass a copy
// instead.
func ClientOptWithHTTPClient(v *http.Client) ClientOpt {
	return func(c *Client) error {
		c.httpClient = v

		return nil
	}
}

// getDefaultUserAgent returns the "User-Agent" HTTP request header added to outgoing HTTP requests.
func getDefaultUserAgent() string {
	return fmt.Sprintf("egoscale/%s (%s; %s/%s)",
		Version,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH)
}

// NewClient returns a new Exoscale API client.
func NewClient(credentials *credentials.Credentials, opts ...ClientOpt) (*Client, error) {
	values, err := credentials.Get()
	if err != nil {
		return nil, err
	}

	client := &Client{
		apiKey:          values.APIKey,
		apiSecret:       values.APISecret,
		serverEndpoint:  string(CHGva2),
		httpClient:      http.DefaultClient,
		pollingInterval: pollingInterval,
		validate:        validator.New(),
		userAgent:       getDefaultUserAgent(),
	}

	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, fmt.Errorf("client configuration error: %s", err)
		}
	}

	return client, nil
}

// getUserAgent only for compatibility with UserAgent.
func (c *Client) getUserAgent() string {
	defaultUA := getDefaultUserAgent()

	if c.userAgent != defaultUA {
		return c.userAgent
	}

	if UserAgent != defaultUA {
		return UserAgent
	}

	return c.userAgent
}

// WithEndpoint returns a copy of Client with new zone Endpoint.
func (c *Client) WithEndpoint(endpoint Endpoint) *Client {
	clone := cloneClient(c)

	clone.serverEndpoint = string(endpoint)

	return clone
}

// WithUserAgent returns a copy of Client with new User-Agent.
func (c *Client) WithUserAgent(ua string) *Client {
	clone := cloneClient(c)

	clone.userAgent = ua + " " + getDefaultUserAgent()

	return clone
}

// WithTrace returns a copy of Client with tracing enabled.
func (c *Client) WithTrace() *Client {
	clone := cloneClient(c)

	clone.trace = true

	return clone
}

// WithHttpClient returns a copy of Client with new http.Client.
func (c *Client) WithHttpClient(client *http.Client) *Client {
	clone := cloneClient(c)

	clone.httpClient = client

	return clone
}

// WithRequestInterceptor returns a copy of Client with new RequestInterceptors.
func (c *Client) WithRequestInterceptor(f ...RequestInterceptorFn) *Client {
	clone := cloneClient(c)

	clone.requestInterceptors = append(clone.requestInterceptors, f...)

	return clone
}

func (c *Client) executeRequestInterceptors(ctx context.Context, req *http.Request) error {
	for _, fn := range c.requestInterceptors {
		if err := fn(ctx, req); err != nil {
			return err
		}
	}

	return nil
}

func cloneClient(c *Client) *Client {
	return &Client{
		apiKey:              c.apiKey,
		apiSecret:           c.apiSecret,
		userAgent:           c.userAgent,
		serverEndpoint:      c.serverEndpoint,
		httpClient:          c.httpClient,
		requestInterceptors: c.requestInterceptors,
		pollingInterval:     c.pollingInterval,
		trace:               c.trace,
		validate:            c.validate,
	}
}
