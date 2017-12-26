// http://docs.cloudstack.apache.org/projects/cloudstack-administration/en/stable/networking_and_traffic.html#security-groups
package egoscale

import (
	"encoding/json"
	"fmt"
)

// SecurityGroup represent a firewalling set of rules
type SecurityGroup struct {
	Id                  string               `json:"id,omitempty"`
	Account             string               `json:"account,omitempty"`
	Description         string               `json:"description,omitempty"`
	Domain              string               `json:"domain,omitempty"`
	Domainid            string               `json:"domainid,omitempty"`
	Name                string               `json:"name,omitempty"`
	Project             string               `json:"project,omitempty"`
	Projectid           string               `json:"projectid,omitempty"`
	VirtualMachineCount int                  `json:"virtualmachinecount,omitempty"`
	VirtualMachineIds   []string             `json:"virtualmachineids,omitempty"`
	IngressRules        []*SecurityGroupRule `json:"ingressrule,omitempty"`
	EgressRules         []*SecurityGroupRule `json:"egressrule,omitempty"`
	Tags                []*ResourceTag       `json:"tags,omitempty"`
	JobId               string               `json:"jobid,omitempty"`
	JobStatus           JobStatusType        `json:"jobstatus,omitempty"`
}

// SecurityGroupRule represents the ingress/egress rule
type SecurityGroupRule struct {
	Id                    string         `json:"ruleid,omitempty"`
	Account               string         `json:"account,omitempty"`
	Cidr                  string         `json:"cidr,omitempty"`
	IcmpType              int            `json:"icmptype,omitempty"`
	IcmpCode              int            `json:"icmpcode,omitempty"`
	StartPort             int            `json:"startport,omitempty"`
	EndPort               int            `json:"endport,omitempty"`
	Protocol              string         `json:"protocol,omitempty"`
	Tags                  []*ResourceTag `json:"tags,omitempty"`
	SecurityGroupId       string
	SecurityGroupName     string               `json:"securitygroupname,omitempty"`
	UserSecurityGroupList []*UserSecurityGroup `json:"usersecuritygrouplist,omitempty"`
	JobId                 string               `json:"jobid,omitempty"`
	JobStatus             JobStatusType        `json:"jobstatus,omitempty"`
}

// UserSecurityGroup represents the traffic of another security group
type UserSecurityGroup struct {
	Group   string `json:"group,omitempty"`
	Account string `json:"account,omitempty"`
}

// CreateSecurityGroupRequest represents a security group creation
type CreateSecurityGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// DeleteSecurityGroupRequest represents a security group deletion
type DeleteSecurityGroupRequest struct {
	Account   string `json:"account,omitempty"`
	DomainId  string `json:"domainid,omitempty"`
	Id        string `json:"id,omitempty"`   // Mutually exclusive with name
	Name      string `json:"name,omitempty"` // Mutually exclusive with id
	ProjectId string `json:"project,omitempty"`
}

// AuthorizeSecurityGroupRequest represents the ingress/egress rule creation
type AuthorizeSecurityGroupRequest struct {
	Account               string               `json:"account,omitempty"`
	Cidr                  string               `json:"cidrlist,omitempty"`
	IcmpType              int                  `json:"icmptype,omitempty"`
	IcmpCode              int                  `json:"icmpcode,omitempty"`
	StartPort             int                  `json:"startport,omitempty"`
	EndPort               int                  `json:"endport,omitempty"`
	Protocol              string               `json:"protocol,omitempty"`
	SecurityGroupId       string               `json:"securitygroupid,omitempty"`
	SecurityGroupName     string               `json:"securitygroupname,omitempty"`
	UserSecurityGroupList []*UserSecurityGroup // manually done... `json:"usersecuritygrouplist,omitempty"`
}

// RevokeSecurityGroup represents the ingress/egress rule deletion
type RevokeSecurityGroupRequest struct {
	Id string `json:"id"`
}

// AuthorizeSecurityGroupResponse represents a new security group
// as well as a deployed security group (rule)
// /!\ the Cloud Stack API document is not fully accurate. /!\
type AuthorizeSecurityGroupResponse struct {
	SecurityGroup *SecurityGroup `json:"securitygroup"`
}

// ListSecurityGroupsRequest represents a search for security groups
type ListSecurityGroupsRequest struct {
	Account           string         `json:"account,omitempty"`
	DomainId          string         `json:"domainid,omitempty"`
	Id                string         `json:"id,omitempty"`
	IsRecursive       bool           `json:"isrecursive,omitempty"`
	Keyword           string         `json:"keyword,omitempty"`
	ListAll           bool           `json:"listall,omitempty"`
	Page              string         `json:"page,omitempty"`
	PageSize          string         `json:"pagesize,omitempty"`
	ProjectId         string         `json:"projectid,omitempty"`
	Type              string         `json:"type,omitempty"`
	SecurityGroupName string         `json:"securitygroupname,omitempty"`
	Tags              []*ResourceTag // XXX `json:"tags,omitempty"`
	VirtualMachineId  string         `json:"virtualmachineid,omitempty"`
}

// ListSecurityGroupsResponse represents a list of security groups
type ListSecurityGroupsResponse struct {
	Count         int              `json:"count"`
	SecurityGroup []*SecurityGroup `json:"securitygroup"`
}

// AuthorizeSecurityGroupEgress authorizes a particular egress rule for this security group
func (exo *Client) AuthorizeSecurityGroupEgress(profile AuthorizeSecurityGroupRequest, async AsyncInfo) (*SecurityGroupRule, error) {
	return exo.addSecurityGroupRule("authorize", "Egress", profile, async)
}

// AuthorizeSecurityGroupIngress authorizes a particular ingress rule for this security group
func (exo *Client) AuthorizeSecurityGroupIngress(profile AuthorizeSecurityGroupRequest, async AsyncInfo) (*SecurityGroupRule, error) {
	return exo.addSecurityGroupRule("authorize", "Ingress", profile, async)
}

// RevokeSecurityGroupEgress revokes a particular egress rule for this security group
func (exo *Client) RevokeSecurityGroupEgress(req RevokeSecurityGroupRequest, async AsyncInfo) error {
	return exo.delSecurityGroupRule("revoke", "Egress", req, async)
}

// RevokeSecurityGroupIngress revokes a particular ingress rule for this security group
func (exo *Client) RevokeSecurityGroupIngress(req RevokeSecurityGroupRequest, async AsyncInfo) error {
	return exo.delSecurityGroupRule("revoke", "Ingress", req, async)
}

func (exo *Client) addSecurityGroupRule(action, kind string, profile AuthorizeSecurityGroupRequest, async AsyncInfo) (*SecurityGroupRule, error) {
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

	var r AuthorizeSecurityGroupResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return nil, err
	}

	if kind == "Egress" {
		return r.SecurityGroup.EgressRules[0], nil
	}
	return r.SecurityGroup.IngressRules[0], nil
}

func (exo *Client) delSecurityGroupRule(action, kind string, req RevokeSecurityGroupRequest, async AsyncInfo) error {
	params, err := prepareValues(req)
	if err != nil {
		return err
	}
	resp, err := exo.AsyncRequest(action+"SecurityGroup"+kind, *params, async)
	if err != nil {
		return err
	}

	var r BooleanResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return err
	}

	if !r.Success {
		return fmt.Errorf("Cannot %sSecurityGroup%s %#v. %s", action, kind, req, r.DisplayText)
	}

	return nil
}

// CreateSecurityGroup creates a SG using the given profile
func (exo *Client) CreateSecurityGroup(req CreateSecurityGroupRequest) (*SecurityGroup, error) {
	params, err := prepareValues(req)
	resp, err := exo.Request("createSecurityGroup", *params)
	if err != nil {
		return nil, err
	}

	var r AuthorizeSecurityGroupResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return nil, err
	}

	return r.SecurityGroup, err
}

// DeleteSecurityGroup deletes a Security Group by name.
func (exo *Client) DeleteSecurityGroup(req DeleteSecurityGroupRequest) error {
	params, err := prepareValues(req)
	resp, err := exo.Request("deleteSecurityGroup", *params)
	if err != nil {
		return err
	}

	var r BooleanResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return err
	}

	if !r.Success {
		return fmt.Errorf("Cannot delete security group: %s", r.DisplayText)
	}

	return nil
}

// GetSecurityGroups is an alias for ListSecurityGroups
func (exo *Client) GetSecurityGroups(req ListSecurityGroupsRequest) ([]*SecurityGroup, error) {
	return exo.ListSecurityGroups(req)
}

// ListSecurityGroups lists the security groups.
func (exo *Client) ListSecurityGroups(req ListSecurityGroupsRequest) ([]*SecurityGroup, error) {
	params, err := prepareValues(req)
	if err != nil {
		return nil, err
	}

	resp, err := exo.Request("listSecurityGroups", *params)
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

// String formats the AuthorizeSecurityGroupRequest to a human version
func (p *AuthorizeSecurityGroupRequest) String() string {
	return fmt.Sprintf("%s: %s %s", p.SecurityGroupName, p.Protocol, p.Cidr)
}
