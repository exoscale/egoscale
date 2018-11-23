package egoscale

import (
	"testing"
)

func TestListResourceLimits(t *testing.T) {
	req := &ListResourceLimits{}
	_ = req.Response().(*ListResourceLimitsResponse)
}

func TestUpdateResourceLimit(t *testing.T) {
	req := &UpdateResourceLimit{}
	_ = req.Response().(*UpdateResourceLimitResponse)
}
