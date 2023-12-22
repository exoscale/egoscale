package credentials

// A StaticProvider is a set of credentials which are set programmatically,
// and will never expire.
type StaticProvider struct {
	Value
}

func NewStaticCredentials(apiKey, apiSecret string) *Credentials {
	return NewCredentials(
		&StaticProvider{Value{APIKey: apiKey, APISecret: apiSecret}},
	)
}

// Retrieve returns the credentials or error if the credentials are invalid.
func (s *StaticProvider) Retrieve() (Value, error) {
	if !s.HasKeys() {
		return Value{}, ErrMissingIncomplete
	}

	return s.Value, nil
}

// IsExpired returns if the credentials are expired.
//
// For StaticProvider, the credentials never expired.
func (s *StaticProvider) IsExpired() bool {
	return false
}
