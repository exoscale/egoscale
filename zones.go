package egoscale

import (
	"strings"
)

type Zone struct {
	Allocationstate       string            `json:"allocationstate,omitempty"`
	Description           string            `json:"description,omitempty"`
	Displaytext           string            `json:"displaytext,omitempty"`
	Domain                string            `json:"domain,omitempty"`
	Domainid              string            `json:"domainid,omitempty"`
	Domainname            string            `json:"domainname,omitempty"`
	Id                    string            `json:"id,omitempty"`
	Internaldns1          string            `json:"internaldns1,omitempty"`
	Internaldns2          string            `json:"internaldns2,omitempty"`
	Ip6dns1               string            `json:"ip6dns1,omitempty"`
	Ip6dns2               string            `json:"ip6dns2,omitempty"`
	Localstorageenabled   bool              `json:"localstorageenabled,omitempty"`
	Name                  string            `json:"name,omitempty"`
	Networktype           string            `json:"networktype,omitempty"`
	Resourcedetails       map[string]string `json:"resourcedetails,omitempty"`
	Securitygroupsenabled bool              `json:"securitygroupsenabled,omitempty"`
	Vlan                  string            `json:"vlan,omitempty"`
	Zonetoken             string            `json:"zonetoken,omitempty"`
}

// ListZonesRequest represents a query for zones
type ListZonesRequest struct {
	Available      bool           `json:"available,omitempty"`
	DomainId       string         `json:"domainid,omitempty"`
	Id             string         `json:"id,omitempty"`
	Keyword        string         `json:"keyword,omitempty"`
	Name           string         `json:"name,omitempty"`
	Page           int            `json:"page,omitempty"`
	PageSize       int            `json:"pagesize,omitempty"`
	ShowCapacities bool           `json:"showcapacities,omitempty"`
	Tags           []*ResourceTag `json:"tags,omitempty"`
}

// Command returns the CloudStack API command
func (req *ListZonesRequest) Command() string {
	return "listZones"
}

// ListZonesResponse represents a list of zones
type ListZonesResponse struct {
	Count int     `json:"count"`
	Zone  []*Zone `json:"zone"`
}

// ListZones lists the zones
func (exo *Client) ListZones(req *ListZonesRequest) ([]*Zone, error) {
	var r ListZonesResponse
	err := exo.Request(req, &r)
	if err != nil {
		return nil, err
	}

	return r.Zone, nil
}

func (exo *Client) GetAllZones() (map[string]*Zone, error) {
	var zones map[string]*Zone
	response, err := exo.ListZones(&ListZonesRequest{})
	if err != nil {
		return zones, err
	}

	zones = make(map[string]*Zone)
	for _, zone := range response {
		zones[strings.ToLower(zone.Name)] = zone
	}
	return zones, nil
}
