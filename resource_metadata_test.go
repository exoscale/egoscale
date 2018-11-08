package egoscale

import (
	"testing"
)

func TestResourceMetadata(t *testing.T) {
	var _ Command = (*ListResourceDetails)(nil)
}

func TestListResourceDetails(t *testing.T) {
	req := &ListResourceDetails{}
	_ = req.response().(*ListResourceDetailsResponse)
}
