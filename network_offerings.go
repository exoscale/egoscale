package egoscale

import (
	"net/url"
)

// NetworkOffering corresponds to the Compute Offerings
type NetworkOffering struct {
	ID                       string            `json:"id"`
	Availability             string            `json:"availability,omitempty"`
	ConserveMode             bool              `json:"conservemode,omitempty"`
	Created                  string            `json:"created"`
	Details                  map[string]string `json:"details,omitempty"`
	DisplayText              string            `json:"displaytext,omitempty"`
	EgressDefaultPolicy      bool              `json:"egressdefaultpolicy,omitempty"`
	ForVPC                   bool              `json:"forvpc,omitempty"`
	GuestIPType              string            `json:"guestiptype,omitempty"`
	IsDefault                bool              `json:"isdefault,omitempty"`
	IsPersistent             bool              `json:"ispersistent,omitempty"`
	MaxConnections           int               `json:"maxconnections,omitempty"`
	Name                     string            `json:"name,omitempty"`
	NetworkRate              int               `json:"networkrate,omitempty"`
	ServiceOfferingID        string            `json:"serviceofferingid,omitempty"`
	SpecifyIPRanges          bool              `json:"specifyipranges,omitempty"`
	SpecifyVlan              bool              `json:"specifyvlan,omitempty"`
	State                    string            `json:"state"` // Disabled/Enabled/Inactive
	SupportsPublicAccess     bool              `json:"supportspublicaccess,omitempty"`
	SupportsStrechedL2Subnet bool              `json:"supportsstrechedl2subnet,omitempty"`
	Tags                     []ResourceTag     `json:"tags,omitempty"`
	TrafficType              string            `json:"traffictype,omitempty"` // Public, Management, Control, ...
	Service                  []Service         `json:"service,omitempty"`
}

// ListNetworkOfferings represents a query for network offerings
//
// CloudStack API: https://cloudstack.apache.org/api/apidocs-4.10/apis/listNetworkOfferings.html
type ListNetworkOfferings struct {
	Availability       string        `json:"availability,omitempty"`
	DisplayText        string        `json:"displaytext,omitempty"`
	ForVPC             bool          `json:"forvpc,omitempty"`
	GuestIPType        string        `json:"guestiptype,omitempty"` // shared of isolated
	ID                 string        `json:"id,omitempty"`
	IsDefault          bool          `json:"isdefault,omitempty"`
	IsTagged           bool          `json:"istagged,omitempty"`
	Keyword            string        `json:"keyword,omitempty"`
	Name               string        `json:"name,omitempty"`
	NetworkID          string        `json:"networkid,omitempty"`
	Page               int           `json:"page,omitempty"`
	PageSize           int           `json:"pagesize,omitempty"`
	SourceNATSupported bool          `json:"sourcenatsupported,omitempty"`
	SpecifyIPRanges    bool          `json:"specifyipranges,omitempty"`
	SpecifyVlan        string        `json:"specifyvlan,omitempty"`
	State              string        `json:"state,omitempty"`
	SupportedServices  string        `json:"supportedservices,omitempty"`
	Tags               []ResourceTag `json:"tags,omitempty"`
	TrafficType        string        `json:"traffictype,omitempty"`
	ZoneID             string        `json:"zoneid,omitempty"`
}

// APIName returns the CloudStack API command name
func (*ListNetworkOfferings) APIName() string {
	return "listNetworkOfferings"
}

func (*ListNetworkOfferings) response() interface{} {
	return new(ListNetworkOfferingsResponse)
}

func (req *ListNetworkOfferings) onBeforeSend(params *url.Values) error {
	// When pagesize is set, the page must also be set
	if req.PageSize > 0 && req.Page == 0 {
		params.Set("page", "0")
	}

	return nil
}

// ListNetworkOfferingsResponse represents a list of service offerings
type ListNetworkOfferingsResponse struct {
	Count           int               `json:"count"`
	NetworkOffering []NetworkOffering `json:"networkoffering"`
}
