package egoscale

import (
	"testing"
)

func TestResetSSHKeyForVirtualMachine(t *testing.T) {
	req := &ResetSSHKeyForVirtualMachine{}
	_ = req.response().(*AsyncJobResult)
	_ = req.asyncResponse().(*VirtualMachine)
}

func TestRegisterSSHKeyPair(t *testing.T) {
	req := &RegisterSSHKeyPair{}
	_ = req.response().(*SSHKeyPair)
}

func TestCreateSSHKeyPair(t *testing.T) {
	req := &CreateSSHKeyPair{}
	_ = req.response().(*SSHKeyPair)
}

func TestDeleteSSHKeyPair(t *testing.T) {
	req := &DeleteSSHKeyPair{}
	_ = req.response().(*booleanResponse)
}

func TestListSSHKeyPairsResponse(t *testing.T) {
	req := &ListSSHKeyPairs{}
	_ = req.response().(*ListSSHKeyPairsResponse)
}
