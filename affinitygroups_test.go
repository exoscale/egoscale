package egoscale

import (
	"testing"
)

func TestAffinityGroups(t *testing.T) {
	var _ AsyncCommand = (*CreateAffinityGroup)(nil)
	var _ AsyncCommand = (*DeleteAffinityGroup)(nil)
	var _ Command = (*ListAffinityGroupTypes)(nil)
	var _ Command = (*ListAffinityGroups)(nil)
	var _ AsyncCommand = (*UpdateVMAffinityGroup)(nil)
}
