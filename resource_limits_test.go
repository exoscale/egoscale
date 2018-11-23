package egoscale

import (
	"testing"
)

func TestListResourceLimits(t *testing.T) {
	req := &ListResourceLimits{}
	_ = req.Response().(*ListResourceLimitsResponse)
}
