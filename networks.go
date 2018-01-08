package egoscale

// Network represents a network
type Network struct {
	ID                          string         `json:"id"`
	Account                     string         `json:"account"`
	ACLID                       string         `json:"aclid,omitempty"`
	ACLType                     string         `json:"acltype,omitempty"`
	BroadcastDomainType         string         `json:"broadcastdomaintype,omitempty"`
	BroadcastURI                string         `json:"broadcasturi,omitempty"`
	CanUseForDeploy             bool           `json:"canusefordeploy,omitempty"`
	Cidr                        string         `json:"cidr,omitempty"`
	DisplayNetwork              bool           `json:"diplaynetwork,omitempty"`
	DisplayText                 string         `json:"displaytext"`
	DNS1                        string         `json:"dns1,omitempty"`
	DNS2                        string         `json:"dns2,omitempty"`
	Domain                      string         `json:"domain,omitempty"`
	DomainID                    string         `json:"domainid,omitempty"`
	Gateway                     string         `json:"gateway,omitempty"`
	IP6Cidr                     string         `json:"ip6cidr,omitempty"`
	IP6Gateway                  string         `json:"ip6gateway,omitempty"`
	IsDefault                   bool           `json:"isdefault,omitempty"`
	IsPersistent                bool           `json:"ispersistent,omitempty"`
	Name                        string         `json:"name"`
	Netmask                     string         `json:"netmask,omitempty"`
	NetworkCidr                 string         `json:"networkcidr,omitempty"`
	NetworkDomain               string         `json:"networkdomain,omitempty"`
	NetworkOfferingAvailability string         `json:"networkofferingavailability,omitempty"`
	NetworkOfferingConserveMode bool           `json:"networkofferingconservemode,omitempty"`
	NetworkOfferingDisplayText  string         `json:"networkofferingdisplaytext,omitempty"`
	NetworkOfferingID           string         `json:"networkofferingid,omitempty"`
	NetworkOfferingName         string         `json:"networkofferingname,omitempty"`
	PhysicalNetworkID           string         `json:"physicalnetworkid,omitempty"`
	Project                     string         `json:"project,omitempty"`
	ProjectID                   string         `json:"projectid,omitempty"`
	Related                     string         `json:"related,omitempty"`
	ReserveIPRange              string         `json:"reserveiprange,omitempty"`
	RestartRequired             bool           `json:"restartrequired,omitempty"`
	SpecifyIPRanges             bool           `json:"specifyipranges,omitempty"`
	State                       string         `json:"state"`
	StrechedL2Subnet            bool           `json:"strechedl2subnet,omitempty"`
	SubdomainAccess             bool           `json:"subdomainaccess,omitempty"`
	TrafficType                 string         `json:"traffictype"`
	Type                        string         `json:"type"`
	Vlan                        string         `json:"vlan,omitemtpy"` // root only
	VpcID                       string         `json:"vpcid,omitempty"`
	ZoneID                      string         `json:"zoneid,omitempty"`
	ZoneName                    string         `json:"zonename,omitempty"`
	ZonesNetworkSpans           string         `json:"zonesnetworkspans,omitempty"`
	Service                     []*Service     `json:"service"`
	Tags                        []*ResourceTag `json:"tags"`
}

// Service is a feature of a network
type Service struct {
	Name       string               `json:"name"`
	Capability []*ServiceCapability `json:"capability,omitempty"`
	Provider   []*ServiceProvider   `json:"provider,omitempty"`
}

// ServiceCapability represents optional capability of a service
type ServiceCapability struct {
	CanChooseServiceCapability bool   `json:"canchooseservicecapability"`
	Name                       string `json:"name"`
	Value                      string `json:"value"`
}

// ServiceProvider represents the provider of the service
type ServiceProvider struct {
	ID                           string   `json:"id"`
	CanEnableIndividualService   bool     `json:"canenableindividualservice"`
	DestinationPhysicalNetworkID string   `json:"destinationphysicalnetworkid"`
	Name                         string   `json:"name"`
	PhysicalNetworkID            string   `json:"physicalnetworkid"`
	ServiceList                  []string `json:"servicelist,omitempty"`
}

// ListNetworks represents a query to a network
type ListNetworks struct {
	Account           string         `json:"account,omitempty"`
	ACLType           string         `json:"acltype,omitempty"` // Account or Domain
	CanUseForDeploy   bool           `json:"canusefordeploy,omitempty"`
	DisplayNetwork    bool           `json:"displaynetwork,omitempty"` // root only
	DomainID          string         `json:"domainid,omitempty"`
	ForVpc            string         `json:"forvpc,omitempty"`
	ID                string         `json:"id,omitempty"`
	IsRecursive       bool           `json:"isrecursive,omitempty"`
	IsSystem          bool           `json:"issystem,omitempty"`
	Keyword           string         `json:"keyword,omitempty"`
	ListAll           bool           `json:"listall,omitempty"`
	Page              int            `json:"page,omitempty"`
	PageSize          int            `json:"pagesize,omitempty"`
	PhysicalNetworkID string         `json:"physicalnetworkid,omitempty"`
	ProjectID         string         `json:"projectid,omitempty"`
	RestartRequired   bool           `json:"restartrequired,omitempty"`
	SpecifyRanges     bool           `json:"specifyranges,omitempty"`
	SupportedServices []*Service     `json:"supportedservices,omitempty"`
	Tags              []*ResourceTag `json:"resourcetag,omitempty"`
	TrafficType       string         `json:"traffictype,omitempty"`
	Type              string         `json:"type,omitempty"`
	VpcID             string         `json:"vpcid,omitempty"`
	ZoneID            string         `json:"zoneid,omitempty"`
}

func (req *ListNetworks) name() string {
	return "listNetworks"
}

func (req *ListNetworks) response() interface{} {
	return new(ListNetworksResponse)
}

// ListNetworksResponse represents the list of networks
type ListNetworksResponse struct {
	Count   int        `json:"count"`
	Network []*Network `json:"network"`
}
