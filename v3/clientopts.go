package v3

import (
	"fmt"
	"net/http"
	"os"

	"github.com/exoscale/egoscale/v3/oapi"
)

const (
	EnvKeyAPIKey    = "EXOSCALE_API_KEY"
	EnvKeyAPISecret = "EXOSCALE_API_SECRET"
)

// ClientOpt represents a function setting Exoscale API client option.
type ClientOpt func(*Client) error

// ClientOptWithCredentials returns a ClientOpt that sets credentials.
func ClientOptWithCredentials(apiKey, apiSecret string) ClientOpt {
	return func(c *Client) error {
		c.creds = NewCredentials(apiKey, apiSecret)

		return nil
	}
}

// ClientOptWithCredentialsFromEnv returns a ClientOpt that reads credentials from environment.
// Returns error of any value is missing in environment.
func ClientOptWithCredentialsFromEnv() ClientOpt {
	return func(c *Client) error {
		key := os.Getenv(EnvKeyAPIKey)
		secret := os.Getenv(EnvKeyAPISecret)
		if key == "" || secret == "" {
			return fmt.Errorf("API credentials not found in environment: %s %s", EnvKeyAPIKey, EnvKeyAPISecret)
		}

		c.creds = NewCredentials(key, secret)

		return nil
	}
}

// ClientOptWithHTTPClient returns a ClientOpt overriding the default http.Client.
// Default HTTP client is [go-retryablehttp] with static retry configuration.
// If you want to keep it your custom client should extend it.
//
// [go-retryablehttp]: https://github.com/hashicorp/go-retryablehttp
func ClientOptWithHTTPClient(v *http.Client) ClientOpt {
	return func(c *Client) error {
		c.httpClient = v

		return nil
	}
}

// ClientOptWithRequestEdotor returns a ClientOpt that adds oapi.RequestEditorFn to oapi client.
// Editors run sequentialy and this function appends provided editor funtion to the end of the list.
func ClientOptWithRequestEdotor(e oapi.RequestEditorFn) ClientOpt {
	return func(c *Client) error {
		c.requestEditors = append(c.requestEditors, e)

		return nil
	}
}
