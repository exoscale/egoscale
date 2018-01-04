package egoscale

import (
	"testing"
)

func TestAddressess(t *testing.T) {
	var _ AsyncCommand = (*AssociateIPAddress)(nil)
	var _ AsyncCommand = (*DisassociateIPAddress)(nil)
	var _ Command = (*ListPublicIPAddresses)(nil)
	var _ AsyncCommand = (*UpdateIPAddress)(nil)
}
