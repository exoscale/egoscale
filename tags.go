package egoscale

// ResourceTag is a tag associated with a resource
//
// http://docs.cloudstack.apache.org/projects/cloudstack-administration/en/4.9/management.html
type ResourceTag struct {
	Account      string `json:"account,omitempty"`
	Customer     string `json:"customer,omitempty"`
	Domain       string `json:"domain,omitempty"`
	DomainID     string `json:"domainid,omitempty"`
	Key          string `json:"key"`
	Project      string `json:"project,omitempty"`
	ProjectID    string `json:"projectid,omitempty"`
	ResourceID   string `json:"resourceid,omitempty"`
	ResourceType string `json:"resourcetype,omitempty"`
	Value        string `json:"value"`
}

// CreateTagsRequest (Async) creates resource tag(s)
//
// CloudStack API: http://cloudstack.apache.org/api/apidocs-4.10/apis/createTags.html
type CreateTagsRequest struct {
	ResourceIds  []*string      `json:"resourceids"`
	ResourceType string         `json:"resourcetype"`
	Tags         []*ResourceTag `json:"tags"`
	Customer     string         `json:"customer,omitempty"`
}

func (req *CreateTagsRequest) name() string {
	return "createTags"
}

func (req *CreateTagsRequest) asyncResponse() interface{} {
	return new(BooleanResponse)
}

// DeleteTagsRequest (Async) deletes the resource tag(s)
//
// CloudStack API: http://cloudstack.apache.org/api/apidocs-4.10/apis/deleteTags.html
type DeleteTagsRequest struct {
	ResourceIds  []*string      `json:"resourceids"`
	ResourceType string         `json:"resourcetype"`
	Tags         []*ResourceTag `json:"tags,omitempty"`
}

func (req *DeleteTagsRequest) name() string {
	return "deleteTags"
}

func (req *DeleteTagsRequest) asyncResponse() interface{} {
	return new(BooleanResponse)
}

// ListTagsRequest list resource tag(s)
//
// CloudStack API: http://cloudstack.apache.org/api/apidocs-4.10/apis/listTags.html
type ListTagsRequest struct {
	Account      string `json:"account,omitempty"`
	Customer     string `json:"customer,omitempty"`
	DomainID     string `json:"domainid,omitempty"`
	IsRecursive  bool   `json:"isrecursive,omitempty"`
	Key          string `json:"key,omitempty"`
	Keyword      string `json:"keyword,omitempty"`
	ListAll      bool   `json:"listall,omitempty"`
	Page         int    `json:"page,omitempty"`
	PageSize     int    `json:"pagesize,omitempty"`
	ProjectID    string `json:"projectid,omitempty"`
	ResourceID   string `json:"resourceid,omitempty"`
	ResourceType string `json:"resourcetype,omitempty"`
	Value        string `json:"value,omitempty"`
}

func (req *ListTagsRequest) name() string {
	return "listTags"
}

func (req *ListTagsRequest) response() interface{} {
	return new(ListTagsResponse)
}

// ListTagsResponse represents a list of resource tags
type ListTagsResponse struct {
	Count int            `json:"count"`
	Tag   []*ResourceTag `json:"tag"`
}
