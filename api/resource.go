// Package api represents the low-level Exoscale API.
package api

import (
	"encoding/json"
)

// Resource represents a resource returned by the Exoscale API.
type Resource struct {
	raw []byte
}

// Raw returns the resource raw plaintext representation. Neither the  availability of the representation nor its
// format or structure are guaranteed: callers should not rely on the return value of this function in a stable
// implementation.
func (r Resource) Raw() []byte {
	return r.raw
}

// MarshalResource serializes the raw API response into a Resource struct.
func MarshalResource(r interface{}) Resource {
	jr, _ := json.Marshal(r)

	return Resource{raw: jr}
}
