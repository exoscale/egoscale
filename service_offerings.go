package egoscale

/*
Service Offerings

A service offering correspond to some hardware features (CPU, RAM).

See: http://docs.cloudstack.apache.org/projects/cloudstack-administration/en/latest/service_offerings.html
*/

// ServiceOffering corresponds to the Compute Offerings
type ServiceOffering struct {
	CPUNumber              int               `json:"cpunumber,omitempty"`
	CPUSpeed               int               `json:"cpuspeed,omitempty"`
	DisplayText            string            `json:"displaytext,omitempty"`
	Domain                 string            `json:"domain,omitempty"`
	DomainID               string            `json:"domainid,omitempty"`
	HostTags               string            `json:"hosttags,omitempty"`
	ID                     string            `json:"id,omitempty"`
	IsCustomized           bool              `json:"iscustomized,omitempty"`
	IsSystem               bool              `json:"issystem,omitempty"`
	IsVolatile             bool              `json:"isvolatile,omitempty"`
	Memory                 int               `json:"memory,omitempty"`
	Name                   string            `json:"name,omitempty"`
	NetworkRate            int               `json:"networkrate,omitempty"`
	ServiceOfferingDetails map[string]string `json:"serviceofferingdetails,omitempty"`
}

// ListServiceOfferingsRequest represents a query for service offerings
type ListServiceOfferingsRequest struct {
	DomainID         string `json:"domainid,omitempty"`
	ID               string `json:"id,omitempty"`
	IsSystem         bool   `json:"issystem,omitempty"`
	Keyword          string `json:"keyword,omitempty"`
	Name             string `json:"name,omitempty"`
	Page             int    `json:"page,omitempty"`
	PageSize         int    `json:"pagesize,omitempty"`
	SystemVMType     string `json:"systemvmtype"`
	VirtualMachineID string `json:"virtualmachineid"`
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
