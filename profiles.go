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

// AsyncJobResultProfile represents a query to fetch the status of async job
type AsyncJobResultProfile struct {
	JobId string `json:"jobid"`
}

// ListNic represents the NIC search
type ListNicsProfile struct {
	VirtualMachineId string `json:"virtualmachineid"`
	ForDisplay       bool   `json:"fordisplay,omitempty"`
	Keyword          string `json:"keyword,omitempty"`
	NetworkId        string `json:"networkid,omitempty"`
	NicId            string `json:"nicid,omitempty"`
	Page             string `json:"page,omitempty"`
	PageSize         string `json:"pagesize,omitempty"`
}
