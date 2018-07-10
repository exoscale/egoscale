package egoscale

import (
	"testing"
)

func TestUpdateReverseDNSForVirtualMachine(t *testing.T) {
	req := &UpdateReverseDNSForVirtualMachine{}
	_ = req.response().(*VirtualMachine)
}

func TestUpdateReverseDNSForPublicIPAddress(t *testing.T) {
	req := &UpdateReverseDNSForPublicIPAddress{}
	_ = req.response().(*IPAddress)
}
