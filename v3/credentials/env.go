package credentials

type EnvProvider struct {
}

func NewEnvCredentials() *Credentials {
	return NewCredentials(&EnvProvider{})
}

func (e *EnvProvider) Retrieve() (Value, error) {
	return Value{}, nil
}

func (e *EnvProvider) IsExpired() bool {
	return false
}
