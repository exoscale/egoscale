package egoscale

import (
	"testing"
)

func TestListInstanceGroups(t *testing.T) {
	req := &ListInstanceGroups{}
	_ = req.Response().(*ListInstanceGroupsResponse)
}

func TestCreateInstanceGroup(t *testing.T) {
	req := &CreateInstanceGroup{}
	_ = req.Response().(*InstanceGroup)
}

func TestUpdateInstanceGroup(t *testing.T) {
	req := &UpdateInstanceGroup{}
	_ = req.Response().(*InstanceGroup)
}

func TestDeleteInstanceGroup(t *testing.T) {
	req := &DeleteInstanceGroup{}
	_ = req.Response().(*booleanResponse)
}
