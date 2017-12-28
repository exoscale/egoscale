/*
Security Groups

Security Groups provide a way to isolate traffic to VMs.

	securityGroup, err := client.CreateSecurityGroup(&CreateSecurityGroupRequest{
		Name: "Load balancer",
		Description: "Opens HTTP/HTTPS ports from the outside world",
	})
	// ...
	err := client.DeleteSecurityGroup(&DeleteSecurityGroupRequest{
		Id: securityGroup.Id,
	})

See: http://docs.cloudstack.apache.org/projects/cloudstack-administration/en/stable/networking_and_traffic.html#security-groups

*/
package egoscale

import (
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
	Id                    string               `json:"ruleid,omitempty"`
	Account               string               `json:"account,omitempty"`
	Cidr                  string               `json:"cidr,omitempty"`
	IcmpType              int                  `json:"icmptype,omitempty"`
	IcmpCode              int                  `json:"icmpcode,omitempty"`
	StartPort             int                  `json:"startport,omitempty"`
	EndPort               int                  `json:"endport,omitempty"`
	Protocol              string               `json:"protocol,omitempty"`
	Tags                  []*ResourceTag       `json:"tags,omitempty"`
	SecurityGroupId       string               `json:"securitygroupid,omitempty"`
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

// Command returns the CloudStack API command
func (req *CreateSecurityGroupRequest) Command() string {
	return "createSecurityGroupRequest"
}

// DeleteSecurityGroupRequest represents a security group deletion
type DeleteSecurityGroupRequest struct {
	Account   string `json:"account,omitempty"`
	DomainId  string `json:"domainid,omitempty"`
	Id        string `json:"id,omitempty"`   // Mutually exclusive with name
	Name      string `json:"name,omitempty"` // Mutually exclusive with id
	ProjectId string `json:"project,omitempty"`
}

// Command returns the CloudStack API command
func (req *DeleteSecurityGroupRequest) Command() string {
	return "deleteSecurityGroupRequest"
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

// Command returns the CloudStack API command
func (req *AuthorizeSecurityGroupRequest) Command() string {
	return "authorizeSecurityGroupRequest"
}

// RevokeSecurityGroup represents the ingress/egress rule deletion
type RevokeSecurityGroupRequest struct {
	Id string `json:"id"`
}

// Command returns the CloudStack API command
func (req *RevokeSecurityGroupRequest) Command() string {
	return "revokeSecurityGroupRequest"
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
	Page              int            `json:"page,omitempty"`
	PageSize          int            `json:"pagesize,omitempty"`
	ProjectId         string         `json:"projectid,omitempty"`
	Type              string         `json:"type,omitempty"`
	SecurityGroupName string         `json:"securitygroupname,omitempty"`
	Tags              []*ResourceTag `json:"tags,omitempty"`
	VirtualMachineId  string         `json:"virtualmachineid,omitempty"`
}

// Command returns the CloudStack API command
func (req *ListSecurityGroupsRequest) Command() string {
	return "listSecurityGroups"
}

// ListSecurityGroupsResponse represents a list of security groups
type ListSecurityGroupsResponse struct {
	Count         int              `json:"count"`
	SecurityGroup []*SecurityGroup `json:"securitygroup"`
}

/*
Methods

All the methods requiring an AsyncInfo value are asynchronous and must be handled with care.
*/

// AuthorizeSecurityGroupEgress authorizes a particular egress rule for this security group
//
// http://cloudstack.apache.org/api/apidocs-4.10/apis/authorizeSecurityGroupEgress.html
func (exo *Client) AuthorizeSecurityGroupEgress(req *AuthorizeSecurityGroupRequest, async AsyncInfo) (*SecurityGroupRule, error) {
	return exo.addSecurityGroupRule(req, async)
}

// AuthorizeSecurityGroupIngress authorizes a particular ingress rule for this security group
//
// http://cloudstack.apache.org/api/apidocs-4.10/apis/authorizeSecurityGroupIngress.html
func (exo *Client) AuthorizeSecurityGroupIngress(req *AuthorizeSecurityGroupRequest, async AsyncInfo) (*SecurityGroupRule, error) {
	return exo.addSecurityGroupRule(req, async)
}

// RevokeSecurityGroupEgress revokes a particular egress rule for this security group
//
// http://cloudstack.apache.org/api/apidocs-4.10/apis/revokeSecurityGroupEgress.html
func (exo *Client) RevokeSecurityGroupEgress(req *RevokeSecurityGroupRequest, async AsyncInfo) error {
	return exo.BooleanAsyncRequest(req, async)
}

// RevokeSecurityGroupIngress revokes a particular ingress rule for this security group
//
// http://cloudstack.apache.org/api/apidocs-4.10/apis/revokeSecurityGroupIngress.html
func (exo *Client) RevokeSecurityGroupIngress(req *RevokeSecurityGroupRequest, async AsyncInfo) error {
	return exo.BooleanAsyncRequest(req, async)
}

func (exo *Client) addSecurityGroupRule(req *AuthorizeSecurityGroupRequest, async AsyncInfo) (*SecurityGroupRule, error) {
	// XXX has to be fixed
	/*
		if len(req.UserSecurityGroupList) > 0 {
			for i, usg := range req.UserSecurityGroupList {
				key := fmt.Sprintf("usersecuritygrouplist[%d]", i)
				params.Set(key+".account", usg.Account)
				params.Set(key+".group", usg.Group)
			}
		}
	*/

	var r AuthorizeSecurityGroupResponse
	err := exo.AsyncRequest(req, &r, async)
	if err != nil {
		return nil, err
	}

	// XXX THIS TOO
	kind := "Egress"
	if kind == "Egress" {
		return r.SecurityGroup.EgressRules[0], nil
	}
	return r.SecurityGroup.IngressRules[0], nil
}

// CreateSecurityGroup creates a SG
//
// http://cloudstack.apache.org/api/apidocs-4.10/apis/createSecurityGroup.html
func (exo *Client) CreateSecurityGroup(req *CreateSecurityGroupRequest) (*SecurityGroup, error) {
	var r AuthorizeSecurityGroupResponse
	err := exo.Request(req, &r)
	if err != nil {
		return nil, err
	}

	return r.SecurityGroup, err
}

// DeleteSecurityGroup deletes a Security Group by name
//
// http://cloudstack.apache.org/api/apidocs-4.10/apis/deleteSecurityGroup.html
func (exo *Client) DeleteSecurityGroup(req *DeleteSecurityGroupRequest) error {
	return exo.BooleanRequest(req)
}

// ListSecurityGroups lists the security groups.
//
// http://cloudstack.apache.org/api/apidocs-4.10/apis/listSecurityGroups.html
func (exo *Client) ListSecurityGroups(req *ListSecurityGroupsRequest) ([]*SecurityGroup, error) {
	var r ListSecurityGroupsResponse
	err := exo.Request(req, &r)
	if err != nil {
		return nil, err
	}

	return r.SecurityGroup, nil
}

// String formats the AuthorizeSecurityGroupRequest to a human version
func (p *AuthorizeSecurityGroupRequest) String() string {
	return fmt.Sprintf("%s: %s %s", p.SecurityGroupName, p.Protocol, p.Cidr)
}
