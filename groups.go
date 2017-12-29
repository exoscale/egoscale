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

func (req *CreateSecurityGroupRequest) name() string {
	return "createSecurityGroupRequest"
}

func (req *CreateSecurityGroupRequest) response() interface{} {
	return new(CreateSecurityGroupResponse)
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

func (req *DeleteSecurityGroupRequest) name() string {
	return "deleteSecurityGroupRequest"
}

func (req *DeleteSecurityGroupRequest) response() interface{} {
	return new(BooleanResponse)
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

func (req *AuthorizeSecurityGroupIngressRequest) name() string {
	return "authorizeSecurityGroupIngress"
}

func (req *AuthorizeSecurityGroupIngressRequest) asyncResponse() interface{} {
	return new(AuthorizeSecurityGroupIngressResponse)
}

// AuthorizeSecurityGroupIngressResponse represents the new egress rule
// /!\ the Cloud Stack API document is not fully accurate. /!\
type AuthorizeSecurityGroupIngressResponse CreateSecurityGroupResponse

// AuthorizeSecurityGroupEgressRequest represents the egress rule creation
type AuthorizeSecurityGroupEgressRequest AuthorizeSecurityGroupIngressRequest

func (req *AuthorizeSecurityGroupEgressRequest) name() string {
	return "authorizeSecurityGroupEgress"
}

func (req *AuthorizeSecurityGroupEgressRequest) asyncResponse() interface{} {
	return new(AuthorizeSecurityGroupEgressResponse)
}

// AuthorizeSecurityGroupEgressResponse represents the new egress rule
// /!\ the Cloud Stack API document is not fully accurate. /!\
type AuthorizeSecurityGroupEgressResponse CreateSecurityGroupResponse

// RevokeSecurityGroupRequest represents the ingress/egress rule deletion
type RevokeSecurityGroupRequest struct {
	ID string `json:"id"`
}

func (req *RevokeSecurityGroupRequest) name() string {
	return "revokeSecurityGroupRequest"
}

func (req *RevokeSecurityGroupRequest) asyncResponse() interface{} {
	return new(BooleanResponse)
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

func (req *ListSecurityGroupsRequest) name() string {
	return "listSecurityGroups"
}

func (req *ListSecurityGroupsRequest) response() interface{} {
	return new(ListSecurityGroupsResponse)
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
	resp, err := exo.AsyncRequest(req, async)
	if err != nil {
		return nil, err
	}
	return resp.(*AuthorizeSecurityGroupIngressResponse).SecurityGroup.IngressRules, nil
}

// CreateEgressRule creates a set of egress rules
//
// Deprecated: use the API directly
func (exo *Client) CreateEgressRule(req *AuthorizeSecurityGroupEgressRequest, async AsyncInfo) ([]*EgressRule, error) {
	resp, err := exo.AsyncRequest(req, async)
	if err != nil {
		return nil, err
	}
	return resp.(*AuthorizeSecurityGroupEgressResponse).SecurityGroup.EgressRules, nil
}

// CreateSecurityGroupWithRules create a security group with its rules
// Warning: it doesn't rollback in case of a failure!
//
// Deprecated: use the API directly
func (exo *Client) CreateSecurityGroupWithRules(name string, ingress []*AuthorizeSecurityGroupIngressRequest, egress []*AuthorizeSecurityGroupEgressRequest, async AsyncInfo) (*SecurityGroup, error) {
	req := &CreateSecurityGroupRequest{
		Name: name,
	}
	resp, err := exo.Request(req)
	if err != nil {
		return nil, err
	}

	sg := resp.(*CreateSecurityGroupResponse).SecurityGroup

	for _, ereq := range egress {
		ereq.SecurityGroupID = sg.ID

		_, err := exo.AsyncRequest(ereq, async)
		if err != nil {
			return nil, err
		}
	}
	for _, ireq := range ingress {
		ireq.SecurityGroupID = sg.ID

		_, err := exo.AsyncRequest(ireq, async)
		if err != nil {
			return nil, err
		}
	}

	r, err := exo.Request(&ListSecurityGroupsRequest{
		ID: sg.ID,
	})
	if err != nil {
		return nil, err
	}

	return r.(*ListSecurityGroupsResponse).SecurityGroup[0], nil
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
