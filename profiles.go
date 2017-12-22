// profiles contains the structs that can be created using Egoscale
package egoscale

// VirtualMachineProfile represents the machine creation
type VirtualMachineProfile struct {
	Name            string
	DiskSize        uint64
	SecurityGroups  []string
	Keypair         string
	UserData        string
	ServiceOffering string
	Template        string
	Zone            string
	AffinityGroups  []string
}

// IpProfile represents the IP creation
type IpAddressProfile struct {
	Zone string `json:"zoneid,omitempty"`
}

// SecurityGroupProfile represents a security group creation
type SecurityGroupProfile struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// SecurityGroupRuleProfile represents the ingress/egress rule creation
type SecurityGroupRuleProfile struct {
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

// AsyncJobResultProfile represents a query to fetch the status of async job
type AsyncJobResultProfile struct {
	JobId string `json:"jobid"`
}
