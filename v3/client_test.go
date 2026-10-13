package v3

import (
	"strings"
	"testing"
)

func TestZoneEndpoints(t *testing.T) {
	if got := ZoneEndpoints[ZoneNameCHGva2]; got != CHGva2 {
		t.Fatalf("ZoneEndpoints[%q] = %q, want %q", ZoneNameCHGva2, got, CHGva2)
	}

	for zone, endpoint := range ZoneEndpoints {
		if !strings.Contains(string(endpoint), "-"+string(zone)+".") {
			t.Errorf("ZoneEndpoints[%q] = %q, not the endpoint of that zone", zone, endpoint)
		}
	}
}
