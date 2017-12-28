/*
Elastic IPs

See: http://docs.cloudstack.apache.org/projects/cloudstack-administration/en/latest/networking_and_traffic.html#about-elastic-ips
*/

package egoscale

// IPAddress represents an IP Address
type IPAddress struct {
	ID                        string         `json:"id"`
	Account                   string         `json:"account,omitempty"`
	AllocatedAt               string         `json:"allocated,omitempty"`
	AssociatedNetworkID       string         `json:"associatednetworkid,omitempty"`
	AssociatedNetworkName     string         `json:"associatednetworkname,omitempty"`
	DomainID                  string         `json:"domainid,omitempty"`
	DomainName                string         `json:"domainname,omitempty"`
	ForDisplay                bool           `json:"fordisplay,omitempty"`
	ForVirtualNetwork         bool           `json:"forvirtualnetwork,omitempty"`
	IPAddress                 string         `json:"ipaddress"`
	IsElastic                 bool           `json:"iselastic,omitempty"`
	IsPortable                bool           `json:"isportable,omitempty"`
	IsSourceNat               bool           `json:"issourcenat,omitempty"`
	IsSystem                  bool           `json:"issystem,omitempty"`
	NetworkID                 string         `json:"networkid,omitempty"`
	PhysicalNetworkID         string         `json:"physicalnetworkid,omitempty"`
	Project                   string         `json:"project,omitempty"`
	ProjectID                 string         `json:"projectid,omitempty"`
	Purpose                   string         `json:"purpose,omitempty"`
	State                     string         `json:"state,omitempty"`
	VirtualMachineDisplayName string         `json:"virtualmachinedisplayname,omitempty"`
	VirtualMachineID          string         `json:"virtualmachineid,omitempty"`
	VirtualMachineName        string         `json:"virtualmachineName,omitempty"`
	VlanID                    string         `json:"vlanid,omitempty"`
	VlanName                  string         `json:"vlanname,omitempty"`
	VMIPAddress               string         `json:"vmipaddress,omitempty"`
	VpcID                     string         `json:"vpcid,omitempty"`
	ZoneID                    string         `json:"zoneid,omitempty"`
	ZoneName                  string         `json:"zonename,omitempty"`
	Tags                      []*ResourceTag `json:"tags,omitempty"`
	JobID                     string         `json:"jobid,omitempty"`
	JobStatus                 JobStatusType  `json:"jobstatus,omitempty"`
}

// AssociateIPAddressRequest represents the IP creation
type AssociateIPAddressRequest struct {
	Account    string `json:"account,omitempty"`
	DomainID   string `json:"domainid,omitempty"`
	ForDisplay bool   `json:"fordisplay,omitempty"`
	IsPortable bool   `json:"isportable,omitempty"`
	NetworkdID string `json:"networkid,omitempty"`
	ProjectID  string `json:"projectid,omitempty"`
	RegionID   string `json:"regionid,omitempty"`
	VpcID      string `json:"vpcid,omitempty"`
	ZoneID     string `json:"zoneid,omitempty"`
}

func (*AssociateIPAddressRequest) name() string {
	return "associateIPAddress"
}

func (*AssociateIPAddressRequest) response() interface{} {
	return new(AssociateIPAddressResponse)
}

// AssociateIPAddressResponse represents the response to the creation of an IPAddress
type AssociateIPAddressResponse struct {
	IPAddress *IPAddress `json:"ipaddress"`
}

// DisassociateIPAddressRequest represents the IP deletion
type DisassociateIPAddressRequest struct {
	ID string `json:"id"`
}

func (*DisassociateIPAddressRequest) name() string {
	return "disassociateIPAddress"
}
func (*DisassociateIPAddressRequest) response() interface{} {
	return new(BooleanResponse)
}

// ListPublicIPAddressesRequest represents a search for public IP addresses
type ListPublicIPAddressesRequest struct {
	Account            string         `json:"account,omitempty"`
	AllocatedOnly      bool           `json:"allocatedonly,omitempty"`
	AllocatedNetworkID string         `json:"allocatednetworkid,omitempty"`
	DomainID           string         `json:"domainid,omitempty"`
	ForDisplay         bool           `json:"fordisplay,omitempty"`
	ForLoadBalancing   bool           `json:"forloadbalancing,omitempty"`
	ForVirtualNetwork  string         `json:"forvirtualnetwork,omitempty"`
	ID                 string         `json:"id,omitempty"`
	IPAddress          string         `json:"ipaddress,omitempty"`
	IsElastic          bool           `json:"iselastic,omitempty"`
	IsRecursive        bool           `json:"isrecursive,omitempty"`
	IsSourceNat        bool           `json:"issourcenat,omitempty"`
	IsStaticNat        bool           `json:"isstaticnat,omitempty"`
	Keyword            string         `json:"keyword,omitempty"`
	ListAll            bool           `json:"listall,omiempty"`
	Page               int            `json:"page,omitempty"`
	PageSize           int            `json:"pagesize,omitempty"`
	PhysicalNetworkID  string         `json:"physicalnetworkid,omitempty"`
	ProjectID          string         `json:"projectid,omitempty"`
	Tags               []*ResourceTag `json:"tags,omitempty"`
	VlanID             string         `json:"vlanid,omitempty"`
	VpcID              string         `json:"vpcid,omitempty"`
	ZoneID             string         `json:"zoneid,omitempty"`
}

func (*ListPublicIPAddressesRequest) name() string {
	return "listPublicIPAddresses"
}

func (*ListPublicIPAddressesRequest) response() interface{} {
	return new(ListPublicIPAddressesResponse)
}

// ListPublicIPAddressesResponse represents a list of public IP addresses
type ListPublicIPAddressesResponse struct {
	Count           int          `json:"count"`
	PublicIPAddress []*IPAddress `json:"publicipaddress"`
}
