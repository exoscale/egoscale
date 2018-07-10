package egoscale

import (
	"testing"
)

func TestUpdateReverseDnsForVirtualMachine(t *testing.T) {
	req := &UpdateReverseDnsForVirtualMachine{}
	_ = req.response().(*VirtualMachine)
}

func TestUpdateReverseDnsForPublicIPAddress(t *testing.T) {
	req := &UpdateReverseDnsForPublicIPAddress{}
	_ = req.response().(*IPAddress)
}
