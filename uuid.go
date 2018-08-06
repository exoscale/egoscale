package egoscale

import (
	"encoding/json"
	"fmt"

	"github.com/satori/go.uuid"
)

// UUID holds a UUID v4
type UUID struct {
	uuid.UUID
}

// UnmarshalJSON unmarshals the raw JSON into the MAC address
func (uuid *UUID) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	u, err := ParseUUID(s)
	if err == nil {
		uuid.UUID = u.UUID
	}
	return err
}

// MarshalJSON converts the UUID to a string representation
func (uuid UUID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", uuid.String())), nil
}

// ParseUUID parses a string into a UUID
func ParseUUID(s string) (*UUID, error) {
	u, err := uuid.FromString(s)
	if err != nil {
		return nil, err
	}
	return &UUID{u}, nil
}

// MustParseUUID acts like ParseUUID but panic in case of a failure
func MustParseUUID(s string) *UUID {
	u, e := ParseUUID(s)
	if e != nil {
		panic(e)
	}
	return u
}
