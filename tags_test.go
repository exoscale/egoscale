package egoscale

import (
	"testing"
)

func TestTagsRequests(t *testing.T) {
	var _ AsyncCommand = (*CreateTagsRequest)(nil)
	var _ AsyncCommand = (*DeleteTagsRequest)(nil)
	var _ Command = (*ListTagsRequest)(nil)
}
