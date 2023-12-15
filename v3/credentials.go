package v3

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
)

type InjectorFn func(ctx context.Context) (CredentialsValue, error)

type CredentialsValue struct {
	APIKey    string
	APISecret string
}

type Credentials struct {
	value     CredentialsValue
	expired   bool
	injectors []InjectorFn

	sync.RWMutex
}

func NewCredentialsFromFile(ctx context.Context) (*Credentials, error) {
	return NewCredentials(ctx, InjectorFile)
}

func NewCredentialsFromEnv(ctx context.Context) (*Credentials, error) {
	return NewCredentials(ctx, InjectorEnv)
}

func NewCredentials(ctx context.Context, injectors ...InjectorFn) (*Credentials, error) {
	credentials := &Credentials{}
	credentials.injectors = append(credentials.injectors, injectors...)

	if err := credentials.Retrieve(ctx); err != nil {
		return nil, err
	}

	if credentials.value.APIKey == "" || credentials.value.APISecret == "" {
		return nil, fmt.Errorf("missing or incomplete API credentials")
	}

	return credentials, nil
}

func (c *Credentials) Retrieve(ctx context.Context) error {
	c.Lock()
	defer c.Unlock()

	var err error
	var success bool
	for _, injector := range c.injectors {
		value, e := injector(ctx)
		if e != nil {
			err = errors.Join(err, e)
			continue
		}
		c.value = value
		success = true
	}
	if !success {
		return err
	}

	c.expired = false

	return nil
}

func (c *Credentials) Expire() {
	c.RLock()
	defer c.RUnlock()

	c.expired = true
}

func (c *Credentials) Get(ctx context.Context) (CredentialsValue, error) {
	c.RLock()
	defer c.RUnlock()

	if !c.expired {
		return CredentialsValue{
			APIKey:    c.value.APIKey,
			APISecret: c.value.APISecret,
		}, nil
	}

	if err := c.Retrieve(ctx); err != nil {
		return CredentialsValue{}, err
	}

	return CredentialsValue{
		APIKey:    c.value.APIKey,
		APISecret: c.value.APISecret,
	}, nil
}

// WithInjector add an
func (c *Credentials) WithInjectors(f ...InjectorFn) {
	c.Lock()
	defer c.Unlock()
	c.injectors = append(c.injectors, f...)
}

var InjectorFile InjectorFn = func(ctx context.Context) (CredentialsValue, error) {
	return CredentialsValue{}, nil
}

// InjectorEnv retrieve credentials from ENV variables.
// InjectorEnv credentials never expire.
var InjectorEnv InjectorFn = func(ctx context.Context) (CredentialsValue, error) {
	apiKey := os.Getenv("EXOSCALE_API_KEY")
	apiSecret := os.Getenv("EXOSCALE_API_SECRET")

	if apiKey == "" || apiSecret == "" {
		return CredentialsValue{}, fmt.Errorf("missing or incomplete API credentials")
	}

	return CredentialsValue{
		APIKey:    apiKey,
		APISecret: apiSecret,
	}, nil
}

func Test() {
	ctx := context.Background()

	// usable creds loaded if no error
	creds, _ := NewCredentialsFromFile(ctx)

	// add custom injector
	creds.WithInjectors(func(ctx context.Context) (CredentialsValue, error) {
		// custom
		return CredentialsValue{}, nil
	})
	// add a predifined injector
	creds.WithInjectors(InjectorEnv)

	// reload it after added injectors.
	if err := creds.Retrieve(ctx); err != nil {
		log.Fatal(err)
	}

	// create loaded creds from injector custom or not.
	_, _ = NewCredentials(ctx, func(ctx context.Context) (CredentialsValue, error) {
		// custom
		return CredentialsValue{}, nil
	}, InjectorFile)

	// create loaded creds from predifined injector.
	_, _ = NewCredentials(ctx, InjectorEnv, InjectorFile)
}
