package credentials

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var (
	ErrMissingIncomplete = errors.New("missing or incomplete API credentials")
)

type Value struct {
	APIKey    string
	APISecret string
}

// HasKeys returns true if the credentials Value has both APIKey and APISecret.
func (v Value) HasKeys() bool {
	return v.APIKey != "" && v.APISecret != ""
}

type Provider interface {
	// Retrieve returns nil if it successfully retrieved the value.
	// Error is returned if the value were not obtainable, or empty.
	Retrieve() (Value, error)

	// IsExpired returns if the credentials are no longer valid, and need
	// to be retrieved.
	IsExpired() bool
}

type ProviderWithContext interface {
	Provider

	RetrieveWithContext(context.Context) (Value, error)
}

type Credentials struct {
	credentials Value
	provider    Provider

	sync.RWMutex
}

func NewCredentials(provider Provider) *Credentials {
	return &Credentials{
		provider: provider,
	}
}

func (c *Credentials) Get() (Value, error) {
	return Value{}, fmt.Errorf("not implemented")
}

func (c *Credentials) retrieve(ctx context.Context) {
	c.Lock()
	defer c.Unlock()

}
