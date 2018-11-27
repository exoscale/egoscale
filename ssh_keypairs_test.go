package egoscale

import (
	"testing"
)

func TestResetSSHKeyForVirtualMachine(t *testing.T) {
	req := &ResetSSHKeyForVirtualMachine{}
	_ = req.Response().(*AsyncJobResult)
	_ = req.AsyncResponse().(*VirtualMachine)
}

func TestRegisterSSHKeyPair(t *testing.T) {
	req := &RegisterSSHKeyPair{}
	_ = req.Response().(*SSHKeyPair)
}

func TestCreateSSHKeyPair(t *testing.T) {
	req := &CreateSSHKeyPair{}
	_ = req.Response().(*SSHKeyPair)
}

func TestDeleteSSHKeyPair(t *testing.T) {
	req := &DeleteSSHKeyPair{}
	_ = req.Response().(*booleanResponse)
}

func TestListSSHKeyPairsResponse(t *testing.T) {
	req := &ListSSHKeyPairs{}
	_ = req.Response().(*ListSSHKeyPairsResponse)
}
