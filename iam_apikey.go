package egoscale

// APIKeyType holds the type of the apikey
type APIKeyType string

const (
	// APIKeyUnrestricted is unrestricted
	APIKeyUnrestricted APIKeyType = "restricted"
	// APIKeyRestricted is restricted
	APIKeyRestricted APIKeyType = "unrestricted"
)

// APIKey represents an apikey
type APIKey struct {
	Description string     `json:"description"`
	Key         string     `json:"key"`
	Secret      string     `json:"secret,omitempty"`
	Operations  string     `json:"operations,omitempty"`
	Type        APIKeyType `json:"type"`
}

// CreateAPIKey represents the apikey creation
type CreateAPIKey struct {
	Description string `json:"description"`
	Operations  string `json:"operations,omitempty"`
	_           bool   `name:"createApiKey" description:"Create an apikey."`
}

// Response returns the struct to unmarshal
func (CreateAPIKey) Response() interface{} {
	return new(APIKey)
}

// ListAPIKeys represents a search for an apikeys
type ListAPIKeys struct {
	_ bool `name:"listApiKeys" description:"List apikeys."`
}

// ListAPIKeysResponse represents a list of apikeys
type ListAPIKeysResponse struct {
	Count      int      `json:"count"`
	ListAPIKey []APIKey `json:"apikey"`
}

// Response returns the struct to unmarshal
func (ListAPIKeys) Response() interface{} {
	return new(ListAPIKeysResponse)
}

// RevokeAPIKey represents a revoke of apikey
type RevokeAPIKey struct {
	Key string `json:"key"`
	_   bool   `name:"revokeApiKey" description:"Revoke an apikey."`
}

// RevokeAPIKeyResponse represents a revoke of apikey
type RevokeAPIKeyResponse struct {
	Success bool `json:"success"`
}

// Response returns the struct to unmarshal
func (RevokeAPIKey) Response() interface{} {
	return new(RevokeAPIKeyResponse)
}
