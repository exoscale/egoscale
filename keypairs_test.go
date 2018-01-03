package egoscale

import (
	"testing"
)

func TestSSHKeyPairRequests(t *testing.T) {
	var _ AsyncCommand = (*ResetSSHKeyForVirtualMachineRequest)(nil)
	var _ Command = (*RegisterSSHKeyPairRequest)(nil)
	var _ Command = (*CreateSSHKeyPairRequest)(nil)
	var _ Command = (*DeleteSSHKeyPairRequest)(nil)
	var _ Command = (*ListSSHKeyPairsRequest)(nil)
}
