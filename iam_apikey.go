package egoscale

// APIKeyType holds the type of the API key
type APIKeyType string

const (
	// APIKeyTypeUnrestricted is unrestricted
	APIKeyTypeUnrestricted APIKeyType = "unrestricted"
	// APIKeyTypeRestricted is restricted
	APIKeyTypeRestricted APIKeyType = "restricted"
)

// APIKey represents an API key
type APIKey struct {
	Name       string     `json:"name"`
	Key        string     `json:"key"`
	Secret     string     `json:"secret,omitempty"`
	Operations []string   `json:"operations,omitempty"`
	Type       APIKeyType `json:"type"`
}

// CreateAPIKey represents an API key creation
type CreateAPIKey struct {
	Name       string `json:"name"`
	Operations string `json:"operations,omitempty"`
	_          bool   `name:"createApiKey" description:"Create an apikey."`
}

// Response returns the struct to unmarshal
func (CreateAPIKey) Response() interface{} {
	return new(APIKey)
}

// ListAPIKeys represents a search for API keys
type ListAPIKeys struct {
	_ bool `name:"listApiKeys" description:"List apikeys."`
}

// ListAPIKeysResponse represents a list of API keys
type ListAPIKeysResponse struct {
	Count   int      `json:"count"`
	APIKeys []APIKey `json:"apikeys"`
}

// Response returns the struct to unmarshal
func (ListAPIKeys) Response() interface{} {
	return new(ListAPIKeysResponse)
}

// RevokeAPIKey represents a revocation of an API key
type RevokeAPIKey struct {
	Key string `json:"key"`
	_   bool   `name:"revokeApiKey" description:"Revoke an apikey."`
}

// RevokeAPIKeyResponse represents the response to an API key revocation
type RevokeAPIKeyResponse struct {
	Success bool `json:"success"`
}

// Response returns the struct to unmarshal
func (RevokeAPIKey) Response() interface{} {
	return new(RevokeAPIKeyResponse)
}
