package egoscale

/*
Virtual Machines

... todo ...

See: http://docs.cloudstack.apache.org/projects/cloudstack-administration/en/stable/virtual_machines.html
*/

// VirtualMachine reprents a virtual machine
type VirtualMachine struct {
	ID                    string            `json:"id,omitempty"`
	Account               string            `json:"account,omitempty"`
	ClusterID             string            `json:"clusterid,omitempty"`
	ClusterName           string            `json:"clustername,omitempty"`
	CPUNumber             int64             `json:"cpunumber,omitempty"`
	CPUSpeed              int64             `json:"cpuspeed,omitempty"`
	CPUUsed               string            `json:"cpuused,omitempty"`
	Created               string            `json:"created,omitempty"`
	Details               map[string]string `json:"details,omitempty"`
	DiskIoRead            int64             `json:"diskioread,omitempty"`
	DiskIoWrite           int64             `json:"diskiowrite,omitempty"`
	DiskKbsRead           int64             `json:"diskkbsread,omitempty"`
	DiskKbsWrite          int64             `json:"diskkbswrite,omitempty"`
	DiskOfferingID        string            `json:"diskofferingid,omitempty"`
	DiskOfferingName      string            `json:"diskofferingname,omitempty"`
	DisplayName           string            `json:"displayname,omitempty"`
	DisplayVM             bool              `json:"displayvm,omitempty"`
	Domain                string            `json:"domain,omitempty"`
	DomainID              string            `json:"domainid,omitempty"`
	ForVirtualNetwork     bool              `json:"forvirtualnetwork,omitempty"`
	Group                 string            `json:"group,omitempty"`
	GroupID               string            `json:"groupid,omitempty"`
	GuestOsID             string            `json:"guestosid,omitempty"`
	HaEnable              bool              `json:"haenable,omitempty"`
	HostID                string            `json:"hostid,omitempty"`
	HostName              string            `json:"hostname,omitempty"`
	Hypervisor            string            `json:"hypervisor,omitempty"`
	InstanceName          string            `json:"instancename,omitempty"` // root only
	IsDynamicallyScalable bool              `json:"isdynamicallyscalable,omitempty"`
	IsoDisplayText        string            `json:"isodisplaytext,omitempty"`
	IsoID                 string            `json:"isoid,omitempty"`
	IsoName               string            `json:"isoname,omitempty"`
	KeyPair               string            `json:"keypair,omitempty"`
	Memory                int64             `json:"memory,omitempty"`
	MemoryIntFreeKbs      int64             `json:"memoryintfreekbs,omitempty"`
	MemoryKbs             int64             `json:"memorykbs,omitempty"`
	MemoryTargetKbs       int64             `json:"memorytargetkbs,omitempty"`
	Name                  string            `json:"name,omitempty"`
	NetworkKbsRead        int64             `json:"networkkbsread,omitempty"`
	NetworkKbsWrite       int64             `json:"networkkbswrite,omitempty"`
	OsCategoryID          string            `json:"oscategoryid,omitempty"`
	OsTypeID              string            `json:"ostypeid,omitempty"`
	Password              string            `json:"password,omitempty"`
	PasswordEnabled       bool              `json:"passwordenabled,omitempty"`
	PciDevices            string            `json:"pcidevices,omitempty"` // not in the doc
	PodID                 string            `json:"podid,omitempty"`
	PodName               string            `json:"podname,omitempty"`
	Project               string            `json:"project,omitempty"`
	ProjectID             string            `json:"projectid,omitempty"`
	PublicIP              string            `json:"publicip,omitempty"`
	PublicIPID            string            `json:"publicipid,omitempty"`
	RootDeviceTd          int64             `json:"rootdeviceid,omitempty"`
	RootDeviceType        string            `json:"rootdevicetype,omitempty"`
	ServiceOfferingID     string            `json:"serviceofferingid,omitempty"`
	ServiceOfferingName   string            `json:"serviceofferingname,omitempty"`
	ServiceState          string            `json:"servicestate,omitempty"`
	State                 string            `json:"state,omitempty"`
	TemplateDisplayText   string            `json:"templatedisplaytext,omitempty"`
	TemplateID            string            `json:"templateid,omitempty"`
	TemplateName          string            `json:"templatename,omitempty"`
	UserID                string            `json:"userid,omitempty"`   // not in the doc
	UserName              string            `json:"username,omitempty"` // not in the doc
	Vgpu                  string            `json:"vgpu,omitempty"`     // not in the doc
	ZoneID                string            `json:"zoneid,omitempty"`
	ZoneName              string            `json:"zonename,omitempty"`
	AffinityGroup         []*AffinityGroup  `json:"affinitygroup,omitempty"`
	Nic                   []*Nic            `json:"nic,omitempty"`
	SecurityGroup         []*SecurityGroup  `json:"securitygroup,omitempty"`
	Tags                  []*ResourceTag    `json:"tags,omitempty"`
	JobID                 string            `json:"jobid,omitempty"`
	JobStatus             JobStatusType     `json:"jobstatus,omitempty"`
}

// DeployVirtualMachineRequest represents the machine creation
type DeployVirtualMachineRequest struct {
	ServiceOfferingID  string         `json:"serviceofferingid"`
	TemplateID         string         `json:"templateid"`
	ZoneID             string         `json:"zoneid"`
	Account            string         `json:"account,omitempty"`
	AffinityGroupIDs   string         `json:"affinitygroupids,omitempty"`   // comma separated list, mutually exclusive with names
	AffinityGroupNames string         `json:"affinitygroupnames,omitempty"` // comma separated list, mutually exclusive with ids
	CustomID           string         `json:"customid,omitempty"`           // root only
	DeploymentPlanner  string         `json:"deploymentplanner,omitempty"`  // root only
	Details            string         `json:"details,omitempty"`
	DiskOfferingID     string         `json:"diskofferingid,omitempty"`
	DisplayName        string         `json:"displayname,omitempty"`
	DisplayVM          bool           `json:"displayvm,omitempty"`
	DomainID           string         `json:"domainid,omitempty"`
	Group              string         `json:"group,omitempty"`
	HostID             string         `json:"hostid,omitempty"`
	Hypervisor         string         `json:"hypervisor,omitempty"`
	IP6Address         string         `json:"ip6address,omitempty"`
	IPAddress          string         `json:"ipaddress,omitempty"`
	IPToNetworkList    []*IPToNetwork `json:"iptonetworklist,omitempty"`
	Keyboard           string         `json:"keyboard,omitempty"`
	KeyPair            string         `json:"keypair,omitempty"`
	Name               string         `json:"name,omitempty"`
	NetworkIDs         []string       `json:"networkids,omitempty"` // mutually exclusive with iptonetworklist
	ProjectID          string         `json:"projectid,omitempty"`
	RootDiskSize       int64          `json:"rootdisksize,omitempty"`       // in GiB
	SecurityGroupIDs   string         `json:"securitygroupids,omitempty"`   // comma separated list, exclusive with names
	SecurityGroupNames string         `json:"securitygroupnames,omitempty"` // comma separated list, exclusive with ids
	Size               string         `json:"size,omitempty"`               // mutually exclusive with diskofferingid
	StartVM            bool           `json:"startvm,omitempty"`
	UserData           []byte         `json:"userdata,omitempty"`
}

// Command returns the command name for the Cloud Stack API
func (req *DeployVirtualMachineRequest) Command() string {
	return "deployVirtualMachine"
}

// DeployVirtualMachineResponse represents a deployed VM instance
type DeployVirtualMachineResponse struct {
	VirtualMachine *VirtualMachine `json:"virtualmachine"`
}

// StartVirtualMachineRequest represents the creation of the virtual machine
type StartVirtualMachineRequest struct {
	ID               string `json:"id"`
	DeploymentPlaner string `json:"deploymentplanner,omitempty"` // root only
	HostID           string `json:"hostid,omitempty"`            // root only
}

// Command returns the command name for the Cloud Stack API
func (req *StartVirtualMachineRequest) Command() string {
	return "startVirtualMachine"
}

// StartVirtualMachineResponse represents a started VM instance
type StartVirtualMachineResponse DeployVirtualMachineResponse

// StopVirtualMachineRequest represents the stopping of the virtual machine
type StopVirtualMachineRequest struct {
	ID     string `json:"id"`
	Forced bool   `json:"forced,omitempty"`
}

// Command returns the command name for the Cloud Stack API
func (req *StopVirtualMachineRequest) Command() string {
	return "stopVirtualMachine"
}

// StopVirtualMachineResponse represents a stopped VM instance
type StopVirtualMachineResponse DeployVirtualMachineResponse

// RebootVirtualMachineRequest represents the rebooting of the virtual machine
type RebootVirtualMachineRequest struct {
	ID string `json:"id"`
}

// Command returns the command name for the Cloud Stack API
func (req *RebootVirtualMachineRequest) Command() string {
	return "rebootVirtualMachine"
}

// RebootVirtualMachineResponse represents a rebooted VM instance
type RebootVirtualMachineResponse DeployVirtualMachineResponse

// DestroyVirtualMachineRequest represents the destruction of the virtual machine
type DestroyVirtualMachineRequest struct {
	ID      string `json:"id"`
	Expunge bool   `json:"expunge,omitempty"`
}

// Command returns the command name for the Cloud Stack API
func (req *DestroyVirtualMachineRequest) Command() string {
	return "destroyVirtualMachine"
}

// DestroyVirtualMachineResponse represents a destroyed VM instance
type DestroyVirtualMachineResponse DeployVirtualMachineResponse

// ListVirtualMachinesRequest represents a search for a VM
type ListVirtualMachinesRequest struct {
	Account           string         `json:"account,omitempty"`
	AffinityGroupID   string         `json:"affinitygroupid,omitempty"`
	Details           string         `json:"details,omitempty"`   // comma separated list, all, group, nics, stats, ...
	DisplayVM         bool           `json:"displayvm,omitempty"` // root only
	DomainID          string         `json:"domainin,omitempty"`
	ForVirtualNetwork bool           `json:"forvirtualnetwork,omitempty"`
	GroupID           string         `json:"groupid,omitempty"`
	HostID            string         `json:"hostid,omitempty"`
	Hypervisor        string         `json:"hypervisor,omitempty"`
	ID                string         `json:"id,omitempty"`
	IDs               string         `json:"ids,omitempty"` // mutually exclusive with id
	IsoID             string         `json:"isoid,omitempty"`
	IsRecursive       bool           `json:"isrecursive,omitempty"`
	KeyPair           string         `json:"keypair,omitempty"`
	Keyword           string         `json:"keyword,omitempty"`
	ListAll           bool           `json:"listall,omitempty"`
	Name              string         `json:"name,omitempty"`
	NetworkID         string         `json:"networkid,omitempty"`
	Page              int            `json:"page,omitempty"`
	PageSize          int            `json:"pagesize,omitempty"`
	PodID             string         `json:"podid,omitempty"`
	ProjectID         string         `json:"projectid,omitempty"`
	ServiceOfferindID string         `json:"serviceofferingid,omitempty"`
	State             string         `json:"state,omitempty"` // Running, Stopped, Present, ...
	StorageID         string         `json:"storageid,omitempty"`
	Tags              []*ResourceTag `json:"tags,omitempty"`
	TemplateID        string         `json:"templateid,omitempty"`
	UserID            string         `json:"userid,omitempty"`
	VpcID             string         `json:"vpcid,omitempty"`
	ZoneID            string         `json:"zoneid,omitempty"`
}

// Command returns the command name for the Cloud Stack API
func (req *ListVirtualMachinesRequest) Command() string {
	return "listVirtualMachines"
}

// ListVirtualMachinesResponse represents a list of virtual machines
type ListVirtualMachinesResponse struct {
	Count          int               `json:"count"`
	VirtualMachine []*VirtualMachine `json:"virtualmachine"`
}

// IPToNetwork represents a mapping between ip and networks
type IPToNetwork struct {
	IP        string `json:"ip,omitempty"`
	IPV6      string `json:"ipv6,omitempty"`
	NetworkID string `json:"networkid,omitempty"`
}
