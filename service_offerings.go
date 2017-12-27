package egoscale

import (
	"strings"
)

type ServiceOffering struct {
	CpuNumber              int               `json:"cpunumber,omitempty"`
	CpuSpeed               int               `json:"cpuspeed,omitempty"`
	DisplayText            string            `json:"displaytext,omitempty"`
	Domain                 string            `json:"domain,omitempty"`
	DomainId               string            `json:"domainid,omitempty"`
	HostTags               string            `json:"hosttags,omitempty"`
	Id                     string            `json:"id,omitempty"`
	IsCustomized           bool              `json:"iscustomized,omitempty"`
	IsSystem               bool              `json:"issystem,omitempty"`
	IsVolatile             bool              `json:"isvolatile,omitempty"`
	Memory                 int               `json:"memory,omitempty"`
	Name                   string            `json:"name,omitempty"`
	NetworkRate            int               `json:"networkrate,omitempty"`
	ServiceOfferingDetails map[string]string `json:"serviceofferingdetails,omitempty"`
}

// ListServiceOfferingRequest represents a query for service offerings
type ListServiceOfferingsRequest struct {
	DomainId         string `json:"domainid,omitempty"`
	Id               string `json:"id,omitempty"`
	IsSystem         bool   `json:"issystem,omitempty"`
	Keyword          string `json:"keyword,omitempty"`
	Name             string `json:"name,omitempty"`
	Page             int    `json:"page,omitempty"`
	PageSize         int    `json:"pagesize,omitempty"`
	SystemVmType     string `json:"systemvmtype"`
	VirtualMachineId string `json:"virtualmachineid"`
}

// Command returns the CloudStack API command
func (req *ListServiceOfferingsRequest) Command() string {
	return "listServiceOfferings"
}

// ListServiceOfferingsResponse represents a list of service offerings
type ListServiceOfferingsResponse struct {
	Count           int                `json:"count"`
	ServiceOffering []*ServiceOffering `json:"serviceoffering"`
}

func (exo *Client) ListServiceOfferings(req *ListServiceOfferingsRequest) ([]*ServiceOffering, error) {
	var r ListServiceOfferingsResponse
	err := exo.Request(req, &r)
	if err != nil {
		return nil, err
	}

	return r.ServiceOffering, nil
}

func (exo *Client) GetProfiles() (map[string]string, error) {
	profiles := make(map[string]string)
	serviceOfferings, err := exo.ListServiceOfferings(&ListServiceOfferingsRequest{})
	if err != nil {
		return profiles, nil
	}

	for _, offering := range serviceOfferings {
		profiles[strings.ToLower(offering.Name)] = offering.Id
	}

	return profiles, nil
}
