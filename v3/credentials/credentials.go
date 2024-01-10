package credentials

import (
	"errors"
	"sync"
	"time"
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

type Credentials struct {
	credentials Value
	provider    Provider
	expireAt    time.Time

	sync.RWMutex
}

func NewCredentials(provider Provider) *Credentials {
	return &Credentials{
		provider: provider,
	}
}

func (c *Credentials) Expire() {
	c.Lock()
	defer c.Unlock()

	c.credentials = Value{}

}

func (c *Credentials) ExpiresAt() (time.Time, error) {
	c.RLock()
	defer c.RUnlock()

	return c.expireAt, nil
}

func (c *Credentials) Get() (Value, error) {
	c.RLock()
	defer c.RUnlock()
	if c.provider.IsExpired() {
		if err := c.retrieve(); err != nil {
			return Value{}, err
		}
	}

	return c.credentials, nil
}

func (c *Credentials) IsExpired() bool {
	return c.provider.IsExpired()
}

func (c *Credentials) retrieve() error {
	c.Lock()
	defer c.Unlock()

	v, err := c.provider.Retrieve()
	if err != nil {
		return err
	}

	c.credentials = v

	return nil
}
