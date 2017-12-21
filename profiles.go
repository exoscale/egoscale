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
