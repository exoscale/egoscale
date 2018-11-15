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
	_ = req.response().(*AsyncJobResult)
	_ = req.asyncResponse().(*IPAddress)
}

func TestDisassociateIPAddress(t *testing.T) {
	req := &DisassociateIPAddress{}
	_ = req.response().(*AsyncJobResult)
	_ = req.asyncResponse().(*booleanResponse)
}

func TestListPublicIPAddresses(t *testing.T) {
	req := &ListPublicIPAddresses{}
	_ = req.response().(*ListPublicIPAddressesResponse)
}

func TestUpdateIPAddress(t *testing.T) {
	req := &UpdateIPAddress{}
	_ = req.response().(*AsyncJobResult)
	_ = req.asyncResponse().(*IPAddress)
}
