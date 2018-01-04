package egoscale

import (
	"testing"
)

func TestGroupsRequests(t *testing.T) {
	var _ AsyncCommand = (*AuthorizeSecurityGroupEgress)(nil)
	var _ AsyncCommand = (*AuthorizeSecurityGroupIngress)(nil)
	var _ Command = (*CreateSecurityGroup)(nil)
	var _ Command = (*DeleteSecurityGroup)(nil)
	var _ Command = (*ListSecurityGroups)(nil)
	var _ AsyncCommand = (*RevokeSecurityGroupEgress)(nil)
	var _ AsyncCommand = (*RevokeSecurityGroupIngress)(nil)
}
