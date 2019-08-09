package egoscale

import (
	"testing"
)

func TestListAPIs(t *testing.T) {
	req := &ListAPIs{}
	_ = req.Response().(*ListAPIsResponse)
}
