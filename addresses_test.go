package egoscale

import (
	"testing"
)

func TestIPAddress(t *testing.T) {
	instance := &IPAddress{}
	if instance.ResourceType() != "PublicIpAddress" {
		t.Errorf("ResourceType doesn't match")
	}
}

func TestAssociateIPAddress(t *testing.T) {
	req := &AssociateIPAddress{}
	_ = req.Response().(*AsyncJobResult)
	_ = req.AsyncResponse().(*IPAddress)
}

func TestDisassociateIPAddress(t *testing.T) {
	req := &DisassociateIPAddress{}
	_ = req.Response().(*AsyncJobResult)
	_ = req.AsyncResponse().(*booleanResponse)
}

func TestListPublicIPAddresses(t *testing.T) {
	req := &ListPublicIPAddresses{}
	_ = req.Response().(*ListPublicIPAddressesResponse)
}

func TestUpdateIPAddress(t *testing.T) {
	req := &UpdateIPAddress{}
	_ = req.Response().(*AsyncJobResult)
	_ = req.AsyncResponse().(*IPAddress)
}
