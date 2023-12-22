package credentials

import (
	"fmt"
	"sync"
)

type Value struct {
	APIKey    string
	APISecret string
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
	value    Value
	Provider Provider

	sync.RWMutex
}

func NewFileCredentials() *Credentials {
	return NewCredentials(&FileProvider{})
}

func NewEnvCredentials() *Credentials {
	return NewCredentials(&EnvProvider{})
}

func NewCredentials(provider Provider) *Credentials {
	return &Credentials{
		Provider: provider,
	}
}

func (c *Credentials) Get() (Value, error) {
	return Value{}, fmt.Errorf("not implemented")
}

type EnvProvider struct {
}

func (e *EnvProvider) Retrieve() (Value, error) {
	return Value{}, nil
}

func (e *EnvProvider) IsExpired() bool {
	return false
}

type FileProvider struct {
}

func (f *FileProvider) Retrieve() (Value, error) {
	return Value{}, nil
}

func (f *FileProvider) IsExpired() bool {
	return false
}
