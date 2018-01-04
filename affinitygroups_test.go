package egoscale

import (
	"testing"
)

func TestAffinityGroupRequests(t *testing.T) {
	var _ AsyncCommand = (*CreateAffinityGroupRequest)(nil)
	var _ AsyncCommand = (*DeleteAffinityGroupRequest)(nil)
	var _ Command = (*ListAffinityGroupTypesRequest)(nil)
	var _ Command = (*ListAffinityGroupsRequest)(nil)
	var _ AsyncCommand = (*UpdateVMAffinityGroupRequest)(nil)
}
