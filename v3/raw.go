package v3

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// RawRequest is an API request built by the caller rather than by a generated operation, for
// endpoints this version of the package does not know yet (generic API tools, MCP servers,
// scripts following the live OpenAPI spec).
type RawRequest struct {
	// Method is the HTTP method (GET, POST, PUT, PATCH, DELETE).
	Method string
	// Path is relative to the client endpoint and starts with "/", e.g. "/instance/<id>". Path
	// parameters must already be escaped (url.PathEscape); the query string goes in Query.
	Path string
	// Query is the query string. Like the generated operations, only single-valued parameters
	// are covered by the request signature.
	Query url.Values
	// Body is sent as application/json when not nil: a []byte or json.RawMessage is sent as
	// is, anything else is encoded with encoding/json.
	Body any
	// OperationID, if set, names the request in traces (ClientOptWithTrace).
	OperationID string
}

// Do signs and sends a raw API request on the client endpoint, with the same user agent, request
// interceptors, tracing and HTTP client as the generated operations.
//
// Like them, a 4xx or 5xx response is returned as an error wrapping an *APIError (use
// errors.Is with ErrNotFound… and errors.As to read the status, the message and the raw body
// with APIError.Body). Otherwise the response is returned and the caller must close its body.
func (c Client) Do(ctx context.Context, r RawRequest) (*http.Response, error) {
	name := r.OperationID
	if name == "" {
		name = "raw request"
	}
	request, err := c.newRawRequest(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("%s: new request: %w", name, err)
	}

	request.Header.Add("User-Agent", c.getUserAgent())
	request.Header.Add("Accept", "application/json")
	if r.Body != nil {
		request.Header.Add("Content-Type", "application/json")
	}

	if err := c.executeRequestInterceptors(ctx, request); err != nil {
		return nil, fmt.Errorf("%s: execute request editors: %w", name, err)
	}

	if err := c.signRequest(request); err != nil {
		return nil, fmt.Errorf("%s: sign request: %w", name, err)
	}

	if c.trace {
		dumpRequest(request, name)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("%s: http client do: %w", name, err)
	}

	if c.trace {
		dumpResponse(response)
	}

	if err := handleHTTPErrorResp(response); err != nil {
		return nil, fmt.Errorf("%s: http response: %w", name, err)
	}

	return response, nil
}

// DoJSON sends a raw API request with [Client.Do] and decodes its JSON response into v (for
// example a *map[string]any, or a generated type such as *Operation). A nil v discards the
// response body.
func (c Client) DoJSON(ctx context.Context, r RawRequest, v any) error {
	response, err := c.Do(ctx, r)
	if err != nil {
		return err
	}
	if v == nil {
		_, _ = io.Copy(io.Discard, response.Body)
		return response.Body.Close()
	}

	if err := prepareJSONResponse(response, v); err != nil {
		name := r.OperationID
		if name == "" {
			name = "raw request"
		}
		return fmt.Errorf("%s: prepare JSON response: %w", name, err)
	}

	return nil
}

func (c Client) newRawRequest(ctx context.Context, r RawRequest) (*http.Request, error) {
	if r.Method == "" {
		return nil, errors.New("missing method")
	}
	if !strings.HasPrefix(r.Path, "/") {
		return nil, fmt.Errorf("path %q must start with /", r.Path)
	}
	if strings.ContainsAny(r.Path, "?#") {
		return nil, fmt.Errorf("path %q must not contain a query string or fragment: use Query", r.Path)
	}

	u := c.serverEndpoint + r.Path
	if len(r.Query) > 0 {
		u += "?" + r.Query.Encode()
	}

	var body io.Reader
	switch b := r.Body.(type) {
	case nil:
	case []byte:
		body = bytes.NewReader(b)
	case json.RawMessage:
		body = bytes.NewReader(b)
	default:
		buf, err := prepareJSONBody(b)
		if err != nil {
			return nil, fmt.Errorf("prepare JSON body: %w", err)
		}
		body = buf
	}

	return http.NewRequestWithContext(ctx, strings.ToUpper(r.Method), u, body)
}

// Body returns the raw body of the error response.
func (e *APIError) Body() []byte { return e.body }
