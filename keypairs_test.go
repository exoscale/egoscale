package egoscale

import (
	"testing"
)

func TestSSHKeyPairs(t *testing.T) {
	var _ AsyncCommand = (*ResetSSHKeyForVirtualMachine)(nil)
	var _ Command = (*RegisterSSHKeyPair)(nil)
	var _ Command = (*CreateSSHKeyPair)(nil)
	var _ Command = (*DeleteSSHKeyPair)(nil)
	var _ Command = (*ListSSHKeyPairs)(nil)
}
