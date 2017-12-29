package egoscale

/*
NICs

See: http://docs.cloudstack.apache.org/projects/cloudstack-administration/en/latest/networking_and_traffic.html#configuring-multiple-ip-addresses-on-a-single-nic
*/

// Nic represents a Network Interface Controller (NIC)
type Nic struct {
	ID               string            `json:"id,omitempty"`
	BroadcastURI     string            `json:"broadcasturi,omitempty"`
	Gateway          string            `json:"gateway,omitempty"`
	IP6Address       string            `json:"ip6address,omitempty"`
	IP6Cidr          string            `json:"ip6cidr,omitempty"`
	IP6Gateway       string            `json:"ip6gateway,omitempty"`
	IPAddress        string            `json:"ipaddress,omitempty"`
	IsDefault        bool              `json:"isdefault,omitempty"`
	IsolationURI     string            `json:"isolationuri,omitempty"`
	MacAddress       string            `json:"macaddress,omitempty"`
	Netmask          string            `json:"netmask,omitempty"`
	NetworkID        string            `json:"networkid,omitempty"`
	NetworkName      string            `json:"networkname,omitempty"`
	SecondaryIP      []*NicSecondaryIP `json:"secondaryip,omitempty"`
	Traffictype      string            `json:"traffictype,omitempty"`
	Type             string            `json:"type,omitempty"`
	VirtualMachineID string            `json:"virtualmachineid,omitempty"`
}

// NicSecondaryIP represents a link between NicID and IPAddress
type NicSecondaryIP struct {
	ID               string `json:"id"`
	IPAddress        string `json:"ipaddress"`
	NetworkID        string `json:"networkid"`
	NicID            string `json:"nicid"`
	VirtualMachineID string `json:"virtualmachineid,omitempty"`
}

// ListNicsRequest represents the NIC search
type ListNicsRequest struct {
	VirtualMachineID string `json:"virtualmachineid"`
	ForDisplay       bool   `json:"fordisplay,omitempty"`
	Keyword          string `json:"keyword,omitempty"`
	NetworkID        string `json:"networkid,omitempty"`
	NicID            string `json:"nicid,omitempty"`
	Page             string `json:"page,omitempty"`
	PageSize         string `json:"pagesize,omitempty"`
}

func (req *ListNicsRequest) name() string {
	return "listNics"
}

func (req *ListNicsRequest) response() interface{} {
	return new(ListNicsResponse)
}

// ListNicsResponse represents a list of templates
type ListNicsResponse struct {
	Count int    `json:"count"`
	Nic   []*Nic `json:"nic"`
}

// AddIPToNicRequest represents the assignation of a secondary IP
type AddIPToNicRequest struct {
	NicID     string `json:"nicid"`
	IPAddress string `json:"ipaddress"`
}

func (req *AddIPToNicRequest) name() string {
	return "addIPToNic"
}
func (req *AddIPToNicRequest) asyncResponse() interface{} {
	return new(AddIPToNicResponse)
}

// AddIPToNicResponse represents the addition of an IP to a NIC
type AddIPToNicResponse struct {
	NicSecondaryIP *NicSecondaryIP `json:"nicsecondaryip"`
}

// RemoveIPFromNicRequest represents a deletion request
type RemoveIPFromNicRequest struct {
	ID string `json:"id"`
}

func (req *RemoveIPFromNicRequest) name() string {
	return "removeIPFromNic"
}

func (req *RemoveIPFromNicRequest) asyncResponse() interface{} {
	return new(BooleanResponse)
}

// ListNics lists the NIC of a VM
//
// Deprecated: use the API directly
func (exo *Client) ListNics(req *ListNicsRequest) ([]*Nic, error) {
	resp, err := exo.Request(req)
	if err != nil {
		return nil, err
	}

	return resp.(*ListNicsResponse).Nic, nil
}

// AddIPToNic adds an IP to a NIC
//
// Deprecated: use the API directly
func (exo *Client) AddIPToNic(nicID string, ipAddress string, async AsyncInfo) (*NicSecondaryIP, error) {
	req := &AddIPToNicRequest{
		NicID:     nicID,
		IPAddress: ipAddress,
	}
	resp, err := exo.AsyncRequest(req, async)
	if err != nil {
		return nil, err
	}

	return resp.(AddIPToNicResponse).NicSecondaryIP, nil
}

// RemoveIPFromNic removes an IP from a NIC
//
// Deprecated: use the API directly
func (exo *Client) RemoveIPFromNic(secondaryNicID string, async AsyncInfo) error {
	req := &RemoveIPFromNicRequest{
		ID: secondaryNicID,
	}
	return exo.BooleanAsyncRequest(req, async)
}
