package credentials

import "os"

type EnvProvider struct {
	retrieved bool
}

func NewEnvCredentials() *Credentials {
	return NewCredentials(&EnvProvider{})
}

// Retrieve retrieves the keys from the environment.
func (e *EnvProvider) Retrieve() (Value, error) {
	e.retrieved = false

	v := Value{
		APIKey: readFromEnv(
			"EXOSCALE_API_KEY",
			"EXOSCALE_KEY",
			"CLOUDSTACK_KEY",
			"CLOUDSTACK_API_KEY",
		),
		APISecret: readFromEnv(
			"EXOSCALE_API_SECRET",
			"EXOSCALE_SECRET",
			"EXOSCALE_SECRET_KEY",
			"CLOUDSTACK_SECRET",
			"CLOUDSTACK_SECRET_KEY",
		),
	}

	if !v.IsSet() {
		return Value{}, ErrMissingIncomplete
	}

	e.retrieved = true

	return v, nil
}

// IsExpired returns if the credentials have been retrieved.
func (e *EnvProvider) IsExpired() bool {
	return !e.retrieved
}

func readFromEnv(keys ...string) string {
	for _, key := range keys {
		if value, ok := os.LookupEnv(key); ok {
			return value
		}
	}
	return ""
}
