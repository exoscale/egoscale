package egoscale

/*
Security Groups

Security Groups provide a way to isolate traffic to VMs.

	resp = new(CreateSecurityGroupResponse)
	err := client.Request(&CreateSecurityGroupRequest{
		Name: "Load balancer",
		Description: "Opens HTTP/HTTPS ports from the outside world",
	}, resp)
	// ...
	err := client.BooleanRequest(&DeleteSecurityGroupRequest{
		ID: resp.SecurityGroup.ID,
	})
	// ...

See: http://docs.cloudstack.apache.org/projects/cloudstack-administration/en/stable/networking_and_traffic.html#security-groups

*/

// SecurityGroup represent a firewalling set of rules
type SecurityGroup struct {
	ID                  string         `json:"id"`
	Account             string         `json:"account,omitempty"`
	Description         string         `json:"description,omitempty"`
	Domain              string         `json:"domain,omitempty"`
	Domainid            string         `json:"domainid,omitempty"`
	Name                string         `json:"name"`
	Project             string         `json:"project,omitempty"`
	Projectid           string         `json:"projectid,omitempty"`
	VirtualMachineCount int            `json:"virtualmachinecount,omitempty"`
	VirtualMachineIDs   []string       `json:"virtualmachineids,omitempty"`
	IngressRules        []*IngressRule `json:"ingressrule"`
	EgressRules         []*EgressRule  `json:"egressrule"`
	Tags                []*ResourceTag `json:"tags,omitempty"`
	JobID               string         `json:"jobid,omitempty"`
	JobStatus           JobStatusType  `json:"jobstatus,omitempty"`
}

// IngressRule represents the ingress rule
type IngressRule struct {
	RuleID                string               `json:"ruleid"`
	Account               string               `json:"account,omitempty"`
	Cidr                  string               `json:"cidr,omitempty"`
	IcmpType              int                  `json:"icmptype,omitempty"`
	IcmpCode              int                  `json:"icmpcode,omitempty"`
	StartPort             int                  `json:"startport,omitempty"`
	EndPort               int                  `json:"endport,omitempty"`
	Protocol              string               `json:"protocol,omitempty"`
	Tags                  []*ResourceTag       `json:"tags,omitempty"`
	SecurityGroupID       string               `json:"securitygroupid,omitempty"`
	SecurityGroupName     string               `json:"securitygroupname,omitempty"`
	UserSecurityGroupList []*UserSecurityGroup `json:"usersecuritygrouplist,omitempty"`
	JobID                 string               `json:"jobid,omitempty"`
	JobStatus             JobStatusType        `json:"jobstatus,omitempty"`
}

// EgressRule represents the ingress rule
type EgressRule IngressRule

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

// CreateSecurityGroupResponse represents a new security group
type CreateSecurityGroupResponse struct {
	SecurityGroup *SecurityGroup `json:"securitygroup"`
}

// DeleteSecurityGroupRequest represents a security group deletion
type DeleteSecurityGroupRequest struct {
	Account   string `json:"account,omitempty"`
	DomainID  string `json:"domainid,omitempty"`
	ID        string `json:"id,omitempty"`   // Mutually exclusive with name
	Name      string `json:"name,omitempty"` // Mutually exclusive with id
	ProjectID string `json:"project,omitempty"`
}

// Command returns the CloudStack API command
func (req *DeleteSecurityGroupRequest) Command() string {
	return "deleteSecurityGroupRequest"
}

// AuthorizeSecurityGroupIngressRequest represents the ingress rule creation
type AuthorizeSecurityGroupIngressRequest struct {
	Account               string               `json:"account,omitempty"`
	Cidr                  string               `json:"cidrlist,omitempty"`
	IcmpType              int                  `json:"icmptype,omitempty"`
	IcmpCode              int                  `json:"icmpcode,omitempty"`
	StartPort             int                  `json:"startport,omitempty"`
	EndPort               int                  `json:"endport,omitempty"`
	Protocol              string               `json:"protocol,omitempty"`
	SecurityGroupID       string               `json:"securitygroupid,omitempty"`
	SecurityGroupName     string               `json:"securitygroupname,omitempty"`
	UserSecurityGroupList []*UserSecurityGroup // manually done... `json:"usersecuritygrouplist,omitempty"`
}

// Command returns the CloudStack API command
func (req *AuthorizeSecurityGroupIngressRequest) Command() string {
	return "authorizeSecurityGroupIngress"
}

// AuthorizeSecurityGroupIngressResponse represents the new egress rule
// /!\ the Cloud Stack API document is not fully accurate. /!\
type AuthorizeSecurityGroupIngressResponse CreateSecurityGroupResponse

// AuthorizeSecurityGroupEgressRequest represents the egress rule creation
type AuthorizeSecurityGroupEgressRequest AuthorizeSecurityGroupIngressRequest

// Command returns the CloudStack API command
func (req *AuthorizeSecurityGroupEgressRequest) Command() string {
	return "authorizeSecurityGroupEgress"
}

// AuthorizeSecurityGroupEgressResponse represents the new egress rule
// /!\ the Cloud Stack API document is not fully accurate. /!\
type AuthorizeSecurityGroupEgressResponse CreateSecurityGroupResponse

// RevokeSecurityGroupRequest represents the ingress/egress rule deletion
type RevokeSecurityGroupRequest struct {
	ID string `json:"id"`
}

// Command returns the CloudStack API command
func (req *RevokeSecurityGroupRequest) Command() string {
	return "revokeSecurityGroupRequest"
}

// ListSecurityGroupsRequest represents a search for security groups
type ListSecurityGroupsRequest struct {
	Account           string         `json:"account,omitempty"`
	DomainID          string         `json:"domainid,omitempty"`
	ID                string         `json:"id,omitempty"`
	IsRecursive       bool           `json:"isrecursive,omitempty"`
	Keyword           string         `json:"keyword,omitempty"`
	ListAll           bool           `json:"listall,omitempty"`
	Page              int            `json:"page,omitempty"`
	PageSize          int            `json:"pagesize,omitempty"`
	ProjectID         string         `json:"projectid,omitempty"`
	Type              string         `json:"type,omitempty"`
	SecurityGroupName string         `json:"securitygroupname,omitempty"`
	Tags              []*ResourceTag `json:"tags,omitempty"`
	VirtualMachineID  string         `json:"virtualmachineid,omitempty"`
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

// CreateIngressRule creates a set of ingress rules
//
// Deprecated: use the API directly
func (exo *Client) CreateIngressRule(req *AuthorizeSecurityGroupIngressRequest, async AsyncInfo) ([]*IngressRule, error) {
	resp := new(AuthorizeSecurityGroupIngressResponse)
	err := exo.AsyncRequest(req, resp, async)
	if err != nil {
		return nil, err
	}
	return resp.SecurityGroup.IngressRules, nil
}

// CreateEgressRule creates a set of egress rules
//
// Deprecated: use the API directly
func (exo *Client) CreateEgressRule(req *AuthorizeSecurityGroupEgressRequest, async AsyncInfo) ([]*EgressRule, error) {
	resp := new(AuthorizeSecurityGroupEgressResponse)
	err := exo.AsyncRequest(req, resp, async)
	if err != nil {
		return nil, err
	}
	return resp.SecurityGroup.EgressRules, nil
}

// CreateSecurityGroupWithRules create a security group with its rules
// Warning: it doesn't rollback in case of a failure!
//
// Deprecated: use the API directly
func (exo *Client) CreateSecurityGroupWithRules(name string, ingress []*AuthorizeSecurityGroupIngressRequest, egress []*AuthorizeSecurityGroupEgressRequest, async AsyncInfo) (*SecurityGroup, error) {
	req := &CreateSecurityGroupRequest{
		Name: name,
	}
	resp := new(CreateSecurityGroupResponse)
	err := exo.Request(req, resp)
	if err != nil {
		return nil, err
	}

	for _, ereq := range egress {
		ereq.SecurityGroupID = resp.SecurityGroup.ID

		resp := new(AuthorizeSecurityGroupEgressResponse)
		err := exo.AsyncRequest(ereq, resp, async)
		if err != nil {
			return nil, err
		}
	}
	for _, ireq := range ingress {
		ireq.SecurityGroupID = resp.SecurityGroup.ID

		resp := new(AuthorizeSecurityGroupIngressResponse)
		err := exo.AsyncRequest(ireq, resp, async)
		if err != nil {
			return nil, err
		}
	}

	r := new(ListSecurityGroupsResponse)
	err = exo.Request(&ListSecurityGroupsRequest{ID: resp.SecurityGroup.ID}, r)
	if err != nil {
		return nil, err
	}

	return r.SecurityGroup[0], nil
}

// DeleteSecurityGroup deletes a security group
//
// Deprecated: use the API directly
func (exo *Client) DeleteSecurityGroup(name string) error {
	req := &DeleteSecurityGroupRequest{
		Name: name,
	}
	return exo.BooleanRequest(req)
}
