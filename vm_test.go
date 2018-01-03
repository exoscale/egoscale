package egoscale

import (
	"testing"
)

func TestVirtualMachineRequests(t *testing.T) {
	var _ AsyncCommand = (*DeployVirtualMachineRequest)(nil)
	var _ AsyncCommand = (*DestroyVirtualMachineRequest)(nil)
	var _ AsyncCommand = (*RebootVirtualMachineRequest)(nil)
	var _ AsyncCommand = (*StartVirtualMachineRequest)(nil)
	var _ AsyncCommand = (*StopVirtualMachineRequest)(nil)
	var _ AsyncCommand = (*ResetPasswordForVirtualMachineRequest)(nil)
	var _ Command = (*UpdateVirtualMachineRequest)(nil)
	var _ Command = (*ListVirtualMachinesRequest)(nil)
	var _ Command = (*GetVMPasswordRequest)(nil)
	var _ AsyncCommand = (*RestoreVirtualMachineRequest)(nil)
	var _ Command = (*ChangeServiceForVirtualMachineRequest)(nil)
	var _ AsyncCommand = (*ScaleVirtualMachineRequest)(nil)
	var _ Command = (*RecoverVirtualMachineRequest)(nil)
	var _ AsyncCommand = (*ExpungeVirtualMachineRequest)(nil)
	// TODO implement
	//var _ AsyncCommand = (*RemoveNICFromVirtualMachineRequest)(nil)
	//var _ AsyncCommand = (*UpdateDefaultNICFromVirtualMachineRequest)(nil)
}
