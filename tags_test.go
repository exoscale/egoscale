package egoscale

import (
	"testing"
)

func TestTagss(t *testing.T) {
	var _ AsyncCommand = (*CreateTags)(nil)
	var _ AsyncCommand = (*DeleteTags)(nil)
	var _ Command = (*ListTags)(nil)
}
