package egoscale

import (
	"testing"
)

func TestListAccounts(t *testing.T) {
	req := &ListAccounts{}
	_ = req.response().(*ListAccountsResponse)
}
