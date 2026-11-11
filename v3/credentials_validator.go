package v3

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/exoscale/egoscale/v3/credentials"
)

// credentialsValidationTimeout bounds the API call checking new credentials.
const credentialsValidationTimeout = 30 * time.Second

// CredentialsValidator returns a validator for credentials.ReloadableOptWithValidator.
// It makes one API call (ListQuotas) with the new credentials, using a client configured with
// opts (e.g. ClientOptWithEndpoint), and only rejects them when the API refuses the key
// (unknown or deleted key, wrong secret): a key lacking the IAM permission for that call, or a
// transient error, doesn't block a rotation.
func CredentialsValidator(opts ...ClientOpt) func(ctx context.Context, v credentials.Value) error {
	return func(ctx context.Context, v credentials.Value) error {
		client, err := NewClient(credentials.NewStaticCredentials(v.APIKey, v.APISecret), opts...)
		if err != nil {
			return fmt.Errorf("validate credentials: %w", err)
		}

		ctx, cancel := context.WithTimeout(ctx, credentialsValidationTimeout)
		defer cancel()

		if _, err := client.ListQuotas(ctx); isInvalidKeyError(err) {
			return fmt.Errorf("validate credentials: %w", err)
		}

		return nil
	}
}

// isInvalidKeyError reports whether err is the API refusing the API key of a request.
func isInvalidKeyError(err error) bool {
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		return false
	}

	return apiErr.StatusCode == http.StatusUnauthorized ||
		(apiErr.StatusCode == http.StatusForbidden && strings.Contains(strings.ToLower(apiErr.Message), "invalid key"))
}
