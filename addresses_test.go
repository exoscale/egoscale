package egoscale

import (
	"testing"
)

func TestAddressesRequests(t *testing.T) {
	var _ AsyncCommand = (*AssociateIPAddressRequest)(nil)
	var _ AsyncCommand = (*DisassociateIPAddressRequest)(nil)
	var _ Command = (*ListPublicIPAddressesRequest)(nil)
	var _ AsyncCommand = (*UpdateIPAddressRequest)(nil)
}
