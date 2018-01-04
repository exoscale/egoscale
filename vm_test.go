package egoscale

import (
	"testing"
)

func TestVirtualMachines(t *testing.T) {
	var _ AsyncCommand = (*DeployVirtualMachine)(nil)
	var _ AsyncCommand = (*DestroyVirtualMachine)(nil)
	var _ AsyncCommand = (*RebootVirtualMachine)(nil)
	var _ AsyncCommand = (*StartVirtualMachine)(nil)
	var _ AsyncCommand = (*StopVirtualMachine)(nil)
	var _ AsyncCommand = (*ResetPasswordForVirtualMachine)(nil)
	var _ Command = (*UpdateVirtualMachine)(nil)
	var _ Command = (*ListVirtualMachines)(nil)
	var _ Command = (*GetVMPassword)(nil)
	var _ AsyncCommand = (*RestoreVirtualMachine)(nil)
	var _ Command = (*ChangeServiceForVirtualMachine)(nil)
	var _ AsyncCommand = (*ScaleVirtualMachine)(nil)
	var _ Command = (*RecoverVirtualMachine)(nil)
	var _ AsyncCommand = (*ExpungeVirtualMachine)(nil)
	// TODO implement
	//var _ AsyncCommand = (*RemoveNICFromVirtualMachine)(nil)
	//var _ AsyncCommand = (*UpdateDefaultNICFromVirtualMachine)(nil)
}
