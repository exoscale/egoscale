package egoscale

import (
	"testing"
)

func TestListNetworksIsACommand(t *testing.T) {
	var _ Command = (*ListNetworks)(nil)
}
