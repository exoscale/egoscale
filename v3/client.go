// Package v3 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/egoscale/v3/generator version v0.0.1 DO NOT EDIT.
package v3

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/exoscale/egoscale/v3/credentials"
	"github.com/exoscale/egoscale/version"
)

// URL represents a zoned url endpoint.
type URL string

const (
	CHGva2 URL = "https://api-ch-gva-2.exoscale.com/v2"
	CHDk2  URL = "https://api-ch-dk-2.exoscale.com/v2"
	DEFra1 URL = "https://api-de-fra-1.exoscale.com/v2"
	DEMuc1 URL = "https://api-de-muc-1.exoscale.com/v2"
	ATVie1 URL = "https://api-at-vie-1.exoscale.com/v2"
	ATVie2 URL = "https://api-at-vie-2.exoscale.com/v2"
	BGSof1 URL = "https://api-bg-sof-1.exoscale.com/v2"
)

// Zones represents a list of all Exoscale zone.
var Zones map[string]URL = map[string]URL{
	"ch-gva-2": CHGva2,
	"ch-dk-2":  CHDk2,
	"de-fra-1": DEFra1,
	"de-muc-1": DEMuc1,
	"at-vie-1": ATVie1,
	"at-vie-2": ATVie2,
	"bg-sof-1": BGSof1,
}

// Client represents an Exoscale API client.
type Client struct {
	apiKey          string
	apiSecret       string
	serverURL       string
	httpClient      *http.Client
	timeout         time.Duration
	pollingInterval time.Duration
	trace           bool

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	requestInterceptors []RequestInterceptorFn
}

// RequestInterceptorFn is the function signature for the RequestInterceptor callback function
type RequestInterceptorFn func(ctx context.Context, req *http.Request) error

// UserAgent is the "User-Agent" HTTP request header added to outgoing HTTP requests.
var UserAgent = fmt.Sprintf("egoscale/%s (%s; %s/%s)",
	version.Version,
	runtime.Version(),
	runtime.GOOS,
	runtime.GOARCH)

const pollingInterval = 3 * time.Second

// ClientOpt represents a function setting Exoscale API client option.
type ClientOpt func(*Client) error

// ClientOptWithTimeout returns a ClientOpt overriding the default client timeout.
func ClientOptWithTimeout(v time.Duration) ClientOpt {
	return func(c *Client) error {
		if v <= 0 {
			return errors.New("timeout value must be greater than 0")
		}
		c.timeout = v

		return nil
	}
}

// ClientOptWithTrace returns a ClientOpt enabling HTTP request/response tracing.
func ClientOptWithTrace() ClientOpt {
	return func(c *Client) error {
		c.trace = true
		return nil
	}
}

// ClientOptWithURL returns a ClientOpt With a given zone URL.
func ClientOptWithURL(url URL) ClientOpt {
	return func(c *Client) error {
		c.serverURL = string(url)
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

func NewClient(credentials *credentials.Credentials, opts ...ClientOpt) (*Client, error) {
	values, err := credentials.Get()
	if err != nil {
		return nil, err
	}

	client := &Client{
		apiKey:          values.APIKey,
		apiSecret:       values.APISecret,
		serverURL:       string(CHGva2),
		httpClient:      http.DefaultClient,
		pollingInterval: pollingInterval,
	}

	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, fmt.Errorf("client configuration error: %s", err)
		}
	}

	return client, nil
}

// WithURL returns a copy of Client with new zone URL.
func (c *Client) WithURL(url URL) *Client {
	return &Client{
		apiKey:              c.apiKey,
		apiSecret:           c.apiSecret,
		serverURL:           string(url),
		httpClient:          c.httpClient,
		requestInterceptors: c.requestInterceptors,
		pollingInterval:     c.pollingInterval,
	}
}

// WithTrace returns a copy of Client with tracing enabled.
func (c *Client) WithTrace() *Client {
	return &Client{
		apiKey:              c.apiKey,
		apiSecret:           c.apiSecret,
		serverURL:           c.serverURL,
		httpClient:          c.httpClient,
		requestInterceptors: c.requestInterceptors,
		pollingInterval:     c.pollingInterval,
		trace:               true,
	}
}

// WithHttpClient returns a copy of Client with new http.Client.
func (c *Client) WithHttpClient(client *http.Client) *Client {
	return &Client{
		apiKey:              c.apiKey,
		apiSecret:           c.apiSecret,
		serverURL:           c.serverURL,
		httpClient:          client,
		requestInterceptors: c.requestInterceptors,
		pollingInterval:     c.pollingInterval,
	}
}

// WithRequestInterceptor returns a copy of Client with new RequestInterceptors.
func (c *Client) WithRequestInterceptor(f ...RequestInterceptorFn) *Client {
	return &Client{
		apiKey:              c.apiKey,
		apiSecret:           c.apiSecret,
		serverURL:           c.serverURL,
		httpClient:          c.httpClient,
		requestInterceptors: append(c.requestInterceptors, f...),
		pollingInterval:     c.pollingInterval,
	}
}

func (c *Client) executeRequestInterceptors(ctx context.Context, req *http.Request) error {
	for _, fn := range c.requestInterceptors {
		if err := fn(ctx, req); err != nil {
			return err
		}
	}

	return nil
}
