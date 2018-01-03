/*
Affinity and Anti-Affinity groups

Affinity and Anti-Affinity groups provide a way to influence where VMs should run. See: http://docs.cloudstack.apache.org/projects/cloudstack-administration/en/stable/virtual_machines.html#affinity-groups
*/

package egoscale

// AffinityGroup represents an (anti-)affinity group
type AffinityGroup struct {
	ID                string   `json:"id,omitempty"`
	Account           string   `json:"account,omitempty"`
	Description       string   `json:"description,omitempty"`
	Domain            string   `json:"domain,omitempty"`
	DomainID          string   `json:"domainid,omitempty"`
	Name              string   `json:"name,omitempty"`
	Type              string   `json:"type,omitempty"`
	VirtualMachineIDs []string `json:"virtualmachineIDs,omitempty"` // *I*ds is not a typo
}

// AffinityGroupType represent an affinity group type
type AffinityGroupType struct {
	Type string `json:"type"`
}

// CreateAffinityGroupRequest represents a new (anti-)affinity group
type CreateAffinityGroupRequest struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Account     string `json:"account,omitempty"`
	Description string `json:"description,omitempty"`
	DomainID    string `json:"domainid,omitempty"`
}

func (req *CreateAffinityGroupRequest) name() string {
	return "createAffinityGroup"
}

func (req *CreateAffinityGroupRequest) asyncResponse() interface{} {
	return new(CreateAffinityGroupResponse)
}

// CreateAffinityGroupResponse represents the response of the creation of an (anti-)affinity group
type CreateAffinityGroupResponse struct {
	AffinityGroup *AffinityGroup `json:"affinitygroup"`
}

// UpdateVMAffinityGroupRequest represents a modification of a (anti-)affinity group
type UpdateVMAffinityGroupRequest struct {
	ID                 string `json:"id"`
	AffinityGroupIDs   string `json:"affinitygroupids,omitempty"`   // mutually exclusive with names
	AffinityGroupNames string `json:"affinitygroupnames,omitempty"` // mutually exclusive with ids
}

func (req *UpdateVMAffinityGroupRequest) name() string {
	return "updateVMAffinityGroup"
}

func (req *UpdateVMAffinityGroupRequest) asyncResponse() interface{} {
	return new(UpdateVMAffinityGroupResponse)
}

// UpdateVMAffinityGroupResponse represents the new VM
type UpdateVMAffinityGroupResponse DeployVirtualMachineResponse

// DeleteAffinityGroupRequest represents an (anti-)affinity group to be deleted
type DeleteAffinityGroupRequest struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Type        string `json:"type,omitempty"`
	Account     string `json:"account,omitempty"`
	Description string `json:"description,omitempty"`
	DomainID    string `json:"domainid,omitempty"`
}

func (req *DeleteAffinityGroupRequest) name() string {
	return "deleteAffinityGroup"
}

func (req *DeleteAffinityGroupRequest) asyncResponse() interface{} {
	return new(BooleanResponse)
}

// ListAffinityGroupsRequest represents an (anti-)affinity groups search
type ListAffinityGroupsRequest struct {
	Account          string `json:"account,omitempty"`
	DomainID         string `json:"domainid,omitempty"`
	ID               string `json:"id,omitempty"`
	IsRecursive      bool   `json:"isrecursive,omitempty"`
	Keyword          string `json:"keyword,omitempty"`
	ListAll          bool   `json:"listall,omitempty"`
	Name             string `json:"name,omitempty"`
	Page             int    `json:"page,omitempty"`
	PageSize         int    `json:"pagesize,omitempty"`
	Type             string `json:"type,omitempty"`
	VirtualMachineID string `json:"virtualmachineid,omitempty"`
}

func (req *ListAffinityGroupsRequest) name() string {
	return "listAffinityGroups"
}

func (req *ListAffinityGroupsRequest) response() interface{} {
	return new(ListAffinityGroupsResponse)
}

// ListAffinityGroupTypesRequest represents an (anti-)affinity groups search
type ListAffinityGroupTypesRequest struct {
	Keyword  string `json:"keyword,omitempty"`
	Page     int    `json:"page,omitempty"`
	PageSize int    `json:"pagesize,omitempty"`
}

func (req *ListAffinityGroupTypesRequest) name() string {
	return "listAffinityGroupTypes"
}

func (req *ListAffinityGroupTypesRequest) response() interface{} {
	return new(ListAffinityGroupTypesResponse)
}

// ListAffinityGroupsResponse represents a list of (anti-)affinity groups
type ListAffinityGroupsResponse struct {
	Count         int              `json:"count"`
	AffinityGroup []*AffinityGroup `json:"affinitygroup"`
}

// ListAffinityGroupTypesResponse represents a list of (anti-)affinity group types
type ListAffinityGroupTypesResponse struct {
	Count             int                  `json:"count"`
	AffinityGroupType []*AffinityGroupType `json:"affinitygrouptype"`
}

// XXX UpdateVmAffinityGroup
//
// http://cloudstack.apache.org/api/apidocs-4.10/apis/updateVMAffinityGroup.html

// Legacy methods

// CreateAffinityGroup creates a group
//
// Deprecated: Use the API directly
func (exo *Client) CreateAffinityGroup(name string, async AsyncInfo) (*AffinityGroup, error) {
	req := &CreateAffinityGroupRequest{
		Name: name,
	}
	resp, err := exo.AsyncRequest(req, async)
	if err != nil {
		return nil, err
	}

	return resp.(CreateAffinityGroupResponse).AffinityGroup, nil
}

// DeleteAffinityGroup deletes a group
//
// Deprecated: Use the API directly
func (exo *Client) DeleteAffinityGroup(name string, async AsyncInfo) error {
	req := &DeleteAffinityGroupRequest{
		Name: name,
	}
	return exo.BooleanAsyncRequest(req, async)
}
