package egoscale

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// APIKeyType holds the type of the apikey
type APIKeyType string

const (
	// APIKeyUnrestricted is unrestricted
	APIKeyUnrestricted APIKeyType = "restricted"
	// APIKeyRestricted is restricted
	APIKeyRestricted APIKeyType = "unrestricted"
)

// APIKey represents an apikey
type APIKey struct {
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
	Description string   `json:"description"`
	Key         string   `json:"key"`
	Secret      string   `json:"secret"`
	Operations  []string `json:"operations"`
	Type        string   `json:"type"`
}

// APIKeyErrorResponse represents an error in the API
type APIKeyErrorResponse struct {
	Message string              `json:"message,omitempty"`
	Errors  map[string][]string `json:"errors"`
}

// Error formats the APIkeyErrorResponse into a string
func (req *APIKeyErrorResponse) Error() string {
	if len(req.Errors) > 0 {
		errs := []string{}
		for name, ss := range req.Errors {
			if len(ss) > 0 {
				errs = append(errs, fmt.Sprintf("%s: %s", name, strings.Join(ss, ", ")))
			}
		}
		return fmt.Sprintf("apikey error: %s (%s)", req.Message, strings.Join(errs, "; "))
	}
	return fmt.Sprintf("apikey error: %s", req.Message)
}

// APIKeyResponse represents an apikey response
type APIKeyResponse struct {
	APIKey *APIKey `json:"api_key"`
}

// CreateAPIKey create an apikey
func (client *Client) CreateAPIKey(ctx context.Context, description string, operations []string) (*APIKey, error) {
	m, err := json.Marshal(APIKeyResponse{
		APIKey: &APIKey{
			Description: description,
			Operations:  operations,
		},
	})
	if err != nil {
		return nil, err
	}

	req, err := client.iamRequest(ctx, "/v1/", nil, string(m), "POST")
	if err != nil {
		return nil, err
	}

	var resp *APIKeyResponse
	if err := json.Unmarshal(req, &resp); err != nil {
		return nil, err
	}

	return resp.APIKey, nil
}

// ListAPIKeys represents a list of apikeys
type ListAPIKeys struct {
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
	Description string   `json:"description"`
	Key         string   `json:"key"`
	Operations  []string `json:"operations"`
	Type        string   `json:"type"`
}

// ListAPIKeysResponse represents a list of apikeys response
type ListAPIKeysResponse struct {
	Count  int64         `json:"count,omitempty"`
	APIKey []ListAPIKeys `json:"api_key,omitempty"`
}

// ListAPIKeys list the apikeys
func (client *Client) ListAPIKeys(ctx context.Context) ([]ListAPIKeysResponse, error) {
	req, err := client.iamRequest(ctx, "/v1/", nil, "", "GET")
	if err != nil {
		return nil, err
	}

	var resp []ListAPIKeysResponse
	if err := json.Unmarshal(req, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// DeleteAPIKey delete an apikey
func (client *Client) DeleteAPIKey(ctx context.Context, key string) error {
	_, err := client.iamRequest(ctx, "/v1/"+key, nil, "", "DELETE")
	return err
}

func (client *Client) iamRequest(ctx context.Context, uri string, urlValues url.Values, params, method string) (json.RawMessage, error) {
	rawURL := client.Endpoint + uri
	url, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	q := url.Query()
	for k, vs := range urlValues {
		for _, v := range vs {
			q.Add(k, v)
		}
	}
	url.RawQuery = q.Encode()

	req, err := http.NewRequest(method, url.String(), strings.NewReader(params))
	if err != nil {
		return nil, err
	}

	var hdr = make(http.Header)
	hdr.Add("", client.APIKey+":"+client.apiSecret)
	hdr.Add("User-Agent", UserAgent)
	hdr.Add("Accept", "application/json")
	if params != "" {
		hdr.Add("Content-Type", "application/json")
	}
	req.Header = hdr

	resp, err := client.HTTPClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("content-type")
	if !strings.Contains(contentType, "application/json") {
		return nil, fmt.Errorf(`response content-type expected to be "application/json", got %q`, contentType)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		e := new(APIKeyErrorResponse)
		if err := json.Unmarshal(b, e); err != nil {
			return nil, err
		}
		return nil, e
	}

	return b, nil
}
