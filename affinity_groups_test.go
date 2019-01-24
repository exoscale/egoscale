package egoscale

import (
	"net/url"
	"testing"
)

func TestCreateAffinityGroup(t *testing.T) {
	req := &CreateAffinityGroup{}
	_ = req.Response().(*AsyncJobResult)
	_ = req.AsyncResponse().(*AffinityGroup)
}

func TestDeleteAffinityGroup(t *testing.T) {
	req := &DeleteAffinityGroup{}
	_ = req.Response().(*AsyncJobResult)
	_ = req.AsyncResponse().(*BooleanResponse)
}

func TestListAffinityGroups(t *testing.T) {
	req := &ListAffinityGroups{}
	_ = req.Response().(*ListAffinityGroupsResponse)
}

func TestListAffinityGroupTypes(t *testing.T) {
	req := &ListAffinityGroupTypes{}
	_ = req.Response().(*ListAffinityGroupTypesResponse)
}

func TestUpdateVMAffinityGroup(t *testing.T) {
	req := &UpdateVMAffinityGroup{}
	_ = req.Response().(*AsyncJobResult)
	_ = req.AsyncResponse().(*VirtualMachine)
}

func TestUpdateVMOnBeforeSend(t *testing.T) {
	req := &UpdateVMAffinityGroup{}
	params := url.Values{}

	if err := req.onBeforeSend(params); err != nil {
		t.Error(err)
	}

	if _, ok := params["affinitygroupids"]; !ok {
		t.Errorf("affinitygroupids should have been set")
	}
}
