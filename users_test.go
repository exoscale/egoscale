package egoscale

import (
	"testing"
)

func TestRegisterUserKeys(t *testing.T) {
	req := &RegisterUserKeys{}
	_ = req.response().(*User)
}

func TestListUsers(t *testing.T) {
	req := &ListUsers{}
	_ = req.response().(*ListUsersResponse)
}
