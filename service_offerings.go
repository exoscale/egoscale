package egoscale

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

// ListServiceOfferings represents a query for service offerings
//
// CloudStack API: https://cloudstack.apache.org/api/apidocs-4.10/apis/listServiceOfferings.html
type ListServiceOfferings struct {
	DomainID         string `json:"domainid,omitempty"`
	ID               string `json:"id,omitempty"`
	IsSystem         bool   `json:"issystem,omitempty"`
	Keyword          string `json:"keyword,omitempty"`
	Name             string `json:"name,omitempty"`
	Page             int    `json:"page,omitempty"`
	PageSize         int    `json:"pagesize,omitempty"`
	SystemVMType     string `json:"systemvmtype,omitempty"`
	VirtualMachineID string `json:"virtualmachineid,omitempty"`
}

func (req *ListServiceOfferings) name() string {
	return "listServiceOfferings"
}

func (req *ListServiceOfferings) response() interface{} {
	return new(ListServiceOfferingsResponse)
}

// ListServiceOfferingsResponse represents a list of service offerings
type ListServiceOfferingsResponse struct {
	Count           int                `json:"count"`
	ServiceOffering []*ServiceOffering `json:"serviceoffering"`
}
