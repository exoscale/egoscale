package egoscale

import (
	"testing"
)

func TestISOResourceType(t *testing.T) {
	instance := &ISO{}
	if instance.ResourceType() != "ISO" {
		t.Errorf("ResourceType doesn't match")
	}
}

func TestListISOs(t *testing.T) {
	req := &ListISOs{}
	_ = req.Response().(*ListISOsResponse)
}

func TestAttachISO(t *testing.T) {
	req := &AttachISO{}
	_ = req.Response().(*AsyncJobResult)
	_ = req.AsyncResponse().(*VirtualMachine)
}

func TestDetachISO(t *testing.T) {
	req := &DetachISO{}
	_ = req.Response().(*AsyncJobResult)
	_ = req.AsyncResponse().(*VirtualMachine)
}
