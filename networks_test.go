package egoscale

import (
	"testing"
)

func TestListNetworksIsACommand(t *testing.T) {
	var _ Command = (*ListNetworks)(nil)
}

func TestListNetworks(t *testing.T) {
	req := &ListNetworks{}
	if req.name() != "listNetworks" {
		t.Errorf("API call doesn't match")
	}
	_ = req.response().(*ListNetworksResponse)
}
