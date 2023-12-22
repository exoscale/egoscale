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

	key := os.Getenv("AWS_ACCESS_KEY_ID")
	secret := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if key == "" || secret == "" {
		return Value{}, ErrMissingIncomplete
	}

	e.retrieved = true
	return Value{
		APIKey:    key,
		APISecret: secret,
	}, nil
}

// IsExpired returns if the credentials have been retrieved.
func (e *EnvProvider) IsExpired() bool {
	return !e.retrieved
}
