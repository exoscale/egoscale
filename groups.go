package egoscale

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// CreateEgressRule is an alias of AuthorizeSecurityGroupEgress
func (exo *Client) CreateEgressRule(profile SecurityGroupRuleProfile, async AsyncInfo) (*SecurityGroupRule, error) {
	return exo.AuthorizeSecurityGroupEgress(profile, async)
}

// AuthorizeSecurityGroupEgress authorizes a particular egress rule for this security group
func (exo *Client) AuthorizeSecurityGroupEgress(profile SecurityGroupRuleProfile, async AsyncInfo) (*SecurityGroupRule, error) {
	return exo.addSecurityGroupRule("authorize", "Egress", profile, async)
}

// CreateIngressRule is an alias of AuthorizeSecurityGroupIngress
func (exo *Client) CreateIngressRule(profile SecurityGroupRuleProfile, async AsyncInfo) (*SecurityGroupRule, error) {
	return exo.AuthorizeSecurityGroupIngress(profile, async)
}

// AuthorizeSecurityGroupIngress authorizes a particular ingress rule for this security group
func (exo *Client) AuthorizeSecurityGroupIngress(profile SecurityGroupRuleProfile, async AsyncInfo) (*SecurityGroupRule, error) {
	return exo.addSecurityGroupRule("authorize", "Ingress", profile, async)
}

// DeleteEgressRule is an alias of RevokeSecurityGroupEgress
func (exo *Client) DeleteEgressRule(securityGroupRuleId string, async AsyncInfo) error {
	return exo.RevokeSecurityGroupEgress(securityGroupRuleId, async)
}

// RevokeSecurityGroupEgress revokes a particular egress rule for this security group
func (exo *Client) RevokeSecurityGroupEgress(securityGroupRuleId string, async AsyncInfo) error {
	return exo.delSecurityGroupRule("revoke", "Egress", securityGroupRuleId, async)
}

// DeleteIngressRule is an alias of RevokeSecurityGroupIngress
func (exo *Client) DeleteIngressRule(securityGroupRuleId string, async AsyncInfo) error {
	return exo.RevokeSecurityGroupIngress(securityGroupRuleId, async)
}

// RevokeSecurityGroupIngress revokes a particular ingress rule for this security group
func (exo *Client) RevokeSecurityGroupIngress(securityGroupRuleId string, async AsyncInfo) error {
	return exo.delSecurityGroupRule("revoke", "Ingress", securityGroupRuleId, async)
}

func (exo *Client) addSecurityGroupRule(action, kind string, profile SecurityGroupRuleProfile, async AsyncInfo) (*SecurityGroupRule, error) {
	params, err := prepareValues(profile)
	if err != nil {
		return nil, err
	}

	if len(profile.UserSecurityGroupList) > 0 {
		for i, usg := range profile.UserSecurityGroupList {
			key := fmt.Sprintf("usersecuritygrouplist[%d]", i)
			params.Set(key+".account", usg.Account)
			params.Set(key+".group", usg.Group)
		}
	}

	resp, err := exo.AsyncRequest(action+"SecurityGroup"+kind, *params, async)
	if err != nil {
		return nil, err
	}

	var r SecurityGroupRuleResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return nil, err
	}

	if kind == "Egress" {
		return r.SecurityGroup.EgressRules[0], nil
	}
	return r.SecurityGroup.IngressRules[0], nil
}

func (exo *Client) delSecurityGroupRule(action, kind, id string, async AsyncInfo) error {
	params := url.Values{}
	params.Set("id", id)

	resp, err := exo.AsyncRequest(action+"SecurityGroup"+kind, params, async)
	if err != nil {
		return err
	}

	var r BooleanResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return err
	}

	if !r.Success {
		return fmt.Errorf("Cannot %sSecurityGroup%s %s. %s", action, kind, id, r.DisplayText)
	}

	return nil
}

// CreateSecurityGroup creates a SG using the given profile
func (exo *Client) CreateSecurityGroup(profile SecurityGroupProfile) (*SecurityGroup, error) {
	params, err := prepareValues(profile)
	resp, err := exo.Request("createSecurityGroup", *params)
	if err != nil {
		return nil, err
	}

	var r CreateSecurityGroupResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return nil, err
	}

	return r.SecurityGroup, err
}

// CreateSecurityGroupWithRules creates a SG with the given set of rules
func (exo *Client) CreateSecurityGroupWithRules(name string, ingress, egress []SecurityGroupRuleProfile, async AsyncInfo) (*SecurityGroup, error) {

	params := url.Values{}
	params.Set("name", name)

	resp, err := exo.Request("createSecurityGroup", params)
	if err != nil {
		return nil, err
	}

	var r CreateSecurityGroupResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return nil, err
	}

	sgid := r.SecurityGroup.Id

	for _, erule := range egress {
		erule.SecurityGroupId = sgid
		_, err = exo.CreateEgressRule(erule, async)
		if err != nil {
			return nil, err
		}
	}

	for _, inrule := range ingress {
		inrule.SecurityGroupId = sgid
		_, err = exo.CreateIngressRule(inrule, async)
		if err != nil {
			return nil, err
		}
	}

	return r.SecurityGroup, nil
}

// DeleteSecurityGroup deletes a Security Group by name.
func (exo *Client) DeleteSecurityGroup(name string) error {
	params := url.Values{}
	params.Set("name", name)

	resp, err := exo.Request("deleteSecurityGroup", params)
	if err != nil {
		return err
	}

	fmt.Printf("## response: %+v\n", resp)
	return nil
}

// GetSecurityGroupById returns a security from its id
func (exo *Client) GetSecurityGroupById(securityGroupId string) (*SecurityGroup, error) {
	params := url.Values{}
	params.Set("id", securityGroupId)
	groups, err := exo.ListSecurityGroups(params)
	if err != nil {
		return nil, err
	}
	if len(groups) != 1 {
		return nil, fmt.Errorf("Expected exactly one security group for %v, got %d", securityGroupId, len(groups))
	}

	return groups[0], nil
}

// GetSecurityGroupByName returns a security from its name
func (exo *Client) GetSecurityGroupByName(securityGroupName string) (*SecurityGroup, error) {
	params := url.Values{}
	params.Set("securitygroupname", securityGroupName)
	groups, err := exo.ListSecurityGroups(params)
	if err != nil {
		return nil, err
	}
	if len(groups) != 1 {
		return nil, fmt.Errorf("Expected exactly one security group for %v, got %d", securityGroupName, len(groups))
	}

	return groups[0], nil
}

// ListSecurityGroups lists the security groups.
func (exo *Client) ListSecurityGroups(params url.Values) ([]*SecurityGroup, error) {
	resp, err := exo.Request("listSecurityGroups", params)
	if err != nil {
		return nil, err
	}

	var r ListSecurityGroupsResponse
	err = json.Unmarshal(resp, &r)
	if err != nil {
		return nil, err
	}

	return r.SecurityGroup, nil
}

// String formats the SecurityGroupRuleProfile to a human version
func (p *SecurityGroupRuleProfile) String() string {
	return fmt.Sprintf("%s: %s %s", p.SecurityGroupName, p.Protocol, p.Cidr)
}
