package egoscale

// Nic represents a Network Interface Controller (NIC)
type Nic struct {
	Id               string            `json:"id,omitempty"`
	BroadcastUri     string            `json:"broadcasturi,omitempty"`
	Gateway          string            `json:"gateway,omitempty"`
	Ip6Address       string            `json:"ip6address,omitempty"`
	Ip6Cidr          string            `json:"ip6cidr,omitempty"`
	Ip6Gateway       string            `json:"ip6gateway,omitempty"`
	IpAddress        string            `json:"ipaddress,omitempty"`
	IsDefault        bool              `json:"isdefault,omitempty"`
	IsolationUri     string            `json:"isolationuri,omitempty"`
	MacAddress       string            `json:"macaddress,omitempty"`
	Netmask          string            `json:"netmask,omitempty"`
	NetworkId        string            `json:"networkid,omitempty"`
	NetworkName      string            `json:"networkname,omitempty"`
	SecondaryIp      []*NicSecondaryIp `json:"secondaryip,omitempty"`
	Traffictype      string            `json:"traffictype,omitempty"`
	Type             string            `json:"type,omitempty"`
	VirtualMachineId string            `json:"virtualmachineid,omitempty"`
}

// NicSecondaryIp represents a link between NicId and IpAddress
type NicSecondaryIp struct {
	Id               string `json:"id"`
	IpAddress        string `json:"ipaddress"`
	NetworkId        string `json:"networkid"`
	NicId            string `json:"nicid"`
	VirtualMachineId string `json:"virtualmachineid,omitempty"`
}

// ListNic represents the NIC search
type ListNicsRequest struct {
	VirtualMachineId string `json:"virtualmachineid"`
	ForDisplay       bool   `json:"fordisplay,omitempty"`
	Keyword          string `json:"keyword,omitempty"`
	NetworkId        string `json:"networkid,omitempty"`
	NicId            string `json:"nicid,omitempty"`
	Page             string `json:"page,omitempty"`
	PageSize         string `json:"pagesize,omitempty"`
}

// Command() returns the CloudStack API command
func (req *ListNicsRequest) Command() string {
	return "listNics"
}

// ListNicsResponse represents a list of templates
type ListNicsResponse struct {
	Count int    `json:"count"`
	Nic   []*Nic `json:"nic"`
}

// AddIpToNicRequest represents the assignation of a secondary IP
type AddIpToNicRequest struct {
	NidId     string `json:"nicid"`
	IpAddress string `json:"ipaddress"`
}

// Command returns the CloudStack API command
func (req *AddIpToNicRequest) Command() string {
	return "addIpToNic"
}

// AddIpToNicResponse represents the addition of an IP to a NIC
type AddIpToNicResponse struct {
	NicSecondaryIp *NicSecondaryIp `json:"nicsecondaryip"`
}

// RemoveIpFromNicRequest
type RemoveIpFromNicRequest struct {
	Id string `json:"id"`
}

// Command returns the CloudStack API command
func (req *RemoveIpFromNicRequest) Command() string {
	return "removeIpFromNic"
}

// ListNics lists the NIC of a VM
func (exo *Client) ListNics(req *ListNicsRequest) ([]*Nic, error) {
	var r ListNicsResponse
	err := exo.Request(req, &r)
	if err != nil {
		return nil, err
	}

	return r.Nic, nil
}

// AddIpToNic adds the IP address to the given NIC
func (exo *Client) AddIpToNic(req *AddIpToNicRequest, async AsyncInfo) (*NicSecondaryIp, error) {
	var r AddIpToNicResponse
	err := exo.AsyncRequest(req, &r, async)
	if err != nil {
		return nil, err
	}

	return r.NicSecondaryIp, nil
}

// RemoveIpFromNic removes the IP address (by Id) from the NIC
func (exo *Client) RemoveIpFromNic(req *RemoveIpFromNicRequest, async AsyncInfo) error {
	return exo.BooleanAsyncRequest(req, async)
}
