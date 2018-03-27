package egoscale

// InstanceGroup represents a group of VM
type InstanceGroup struct {
	ID        string `json:"id"`
	Account   string `json:"account,omitempty"`
	Created   string `json:"created,omitempty"`
	Domain    string `json:"domain,omitempty"`
	DomainID  string `json:"domainid,omitempty"`
	Name      string `json:"name,omitempty"`
	Project   string `json:"project,omitempty"`
	ProjectID string `json:"projectid,omitempty"`
}

// InstanceGroupResponse represents a VM group
type InstanceGroupResponse struct {
	InstanceGroup InstanceGroup `json:"instancegroup"`
}

// CreateInstanceGroup creates a VM group
//
// CloudStack API: http://cloudstack.apache.org/api/apidocs-4.10/apis/createInstanceGroup.html
type CreateInstanceGroup struct {
	Name      string `json:"name"`
	Account   string `json:"account,omitempty"`
	DomainID  string `json:"domainid,omitempty"`
	ProjectID string `json:"projectid,omitempty"`
}

// APIName returns the CloudStack API command name
func (*CreateInstanceGroup) APIName() string {
	return "createInstanceGroup"
}

func (*CreateInstanceGroup) response() interface{} {
	return new(CreateInstanceGroupResponse)
}

// CreateInstanceGroupResponse represents a freshly created VM group
type CreateInstanceGroupResponse InstanceGroupResponse

// UpdateInstanceGroup creates a VM group
//
// CloudStack API: http://cloudstack.apache.org/api/apidocs-4.10/apis/updateInstanceGroup.html
type UpdateInstanceGroup struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

// APIName returns the CloudStack API command name
func (*UpdateInstanceGroup) APIName() string {
	return "updateInstanceGroup"
}

func (*UpdateInstanceGroup) response() interface{} {
	return new(UpdateInstanceGroupResponse)
}

// UpdateInstanceGroupResponse represents an updated VM group
type UpdateInstanceGroupResponse InstanceGroupResponse

// DeleteInstanceGroup creates a VM group
//
// CloudStack API: http://cloudstack.apache.org/api/apidocs-4.10/apis/deleteInstanceGroup.html
type DeleteInstanceGroup struct {
	Name      string `json:"name"`
	Account   string `json:"account,omitempty"`
	DomainID  string `json:"domainid,omitempty"`
	ProjectID string `json:"projectid,omitempty"`
}

// APIName returns the CloudStack API command name
func (*DeleteInstanceGroup) APIName() string {
	return "deleteInstanceGroup"
}

func (*DeleteInstanceGroup) response() interface{} {
	return new(booleanSyncResponse)
}

// ListInstanceGroups lists VM groups
//
// CloudStack API: http://cloudstack.apache.org/api/apidocs-4.10/apis/listInstanceGroups.html
type ListInstanceGroups struct {
	Account     string `json:"account,omitempty"`
	DomainID    string `json:"domainid,omitempty"`
	ID          string `json:"id,omitempty"`
	IsRecursive *bool  `json:"isrecursive,omitempty"`
	Keyword     string `json:"keyword,omitempty"`
	ListAll     *bool  `json:"listall,omitempty"`
	Page        int    `json:"page,omitempty"`
	PageSize    int    `json:"pagesize,omitempty"`
	State       string `json:"state,omitempty"`
	ProjectID   string `json:"projectid,omitempty"`
}

// APIName returns the CloudStack API command name
func (*ListInstanceGroups) APIName() string {
	return "listInstanceGroups"
}

func (*ListInstanceGroups) response() interface{} {
	return new(ListInstanceGroupsResponse)
}

func (req *ListInstanceGroups) onBeforeSend(params *url.Values) error {
	// When pagesize is set, the page must also be set
	if req.PageSize > 0 && req.Page == 0 {
		params.Set("page", "0")
	}
}

// ListInstanceGroupsResponse represents a list of instance groups
type ListInstanceGroupsResponse struct {
	Count         int             `json:"count"`
	InstanceGroup []InstanceGroup `json:"instancegroup"`
}
