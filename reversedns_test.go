package egoscale

import (
	"testing"
)

func TestDeleteReverseDNSFromVirtualMachine(t *testing.T) {
	req := &DeleteReverseDNSFromVirtualMachine{}
	_ = req.Response().(*BooleanResponse)
}

func TestDeleteReverseDNSFromPublicIPAddress(t *testing.T) {
	req := &DeleteReverseDNSFromPublicIPAddress{}
	_ = req.Response().(*BooleanResponse)
}

func TestQueryReverseDNSForVirtualMachine(t *testing.T) {
	req := &QueryReverseDNSForVirtualMachine{}
	_ = req.Response().(*VirtualMachine)
}

func TestQueryReverseDNSForPublicIPAddress(t *testing.T) {
	req := &QueryReverseDNSForPublicIPAddress{}
	_ = req.Response().(*IPAddress)
}

func TestUpdateReverseDNSForVirtualMachine(t *testing.T) {
	req := &UpdateReverseDNSForVirtualMachine{}
	_ = req.Response().(*VirtualMachine)
}

func TestUpdateReverseDNSForPublicIPAddress(t *testing.T) {
	req := &UpdateReverseDNSForPublicIPAddress{}
	_ = req.Response().(*IPAddress)
}
