package egoscale

import (
	"testing"
)

func TestIsoResourceType(t *testing.T) {
	instance := &Iso{}
	if instance.ResourceType() != "ISO" {
		t.Errorf("ResourceType doesn't match")
	}
}

func TestListIsos(t *testing.T) {
	req := &ListIsos{}
	_ = req.Response().(*ListIsosResponse)
}

func TestAttachIso(t *testing.T) {
	req := &AttachIso{}
	_ = req.Response().(*AsyncJobResult)
	_ = req.AsyncResponse().(*VirtualMachine)
}

func TestDetachIso(t *testing.T) {
	req := &DetachIso{}
	_ = req.Response().(*AsyncJobResult)
	_ = req.AsyncResponse().(*VirtualMachine)
}
