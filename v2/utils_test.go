package v2

import (
	"testing"
)

func TestIsValidUUID(t *testing.T) {
	uuid := "5361a11b-615c-42bf-9bdb-e2c3790ada14"
	if !IsValidUUID(uuid) {
		t.Error("uuid should be valid")
	}

	uuid = "invalid-uuid"
	if IsValidUUID(uuid) {
		t.Error("invalid uuid should be invalid")
	}
}
