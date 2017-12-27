package egoscale

type Template struct {
	Account               string            `json:"account,omitempty"`
	AccountId             string            `json:"accountid,omitempty"`
	Bootable              bool              `json:"bootable,omitempty"`
	Checksum              string            `json:"checksum,omitempty"`
	CreatedAt             string            `json:"created,omitempty"`
	CrossZones            bool              `json:"crossZones,omitempty"`
	Details               map[string]string `json:"details,omitempty"`
	DisplayText           string            `json:"displaytext,omitempty"`
	Domain                string            `json:"domain,omitempty"`
	DomainId              string            `json:"domainid,omitempty"`
	Format                string            `json:"format,omitempty"`
	HostId                string            `json:"hostid,omitempty"`
	HostName              string            `json:"hostname,omitempty"`
	Hypervisor            string            `json:"hypervisor,omitempty"`
	Id                    string            `json:"id,omitempty"`
	IsDynamicallyScalable bool              `json:"isdynamicallyscalable,omitempty"`
	IsExtractable         bool              `json:"isextractable,omitempty"`
	IsFeatured            bool              `json:"isfeatured,omitempty"`
	IsPublic              bool              `json:"ispublic,omitempty"`
	IsReady               bool              `json:"isready,omitempty"`
	Name                  string            `json:"name,omitempty"`
	OsTypeId              string            `json:"ostypeid,omitempty"`
	OsTypeName            string            `json:"ostypename,omitempty"`
	PasswordEnabled       bool              `json:"passwordenabled,omitempty"`
	Project               string            `json:"project,omitempty"`
	ProjectId             string            `json:"projectid,omitempty"`
	RemovedAt             string            `json:"removed,omitempty"`
	Size                  int64             `json:"size,omitempty"`
	SourceTemplateId      string            `json:"sourcetemplateid,omitempty"`
	SshKeyEnabled         bool              `json:"sshkeyenabled,omitempty"`
	Status                string            `json:"status,omitempty"`
	Zoneid                string            `json:"zoneid,omitempty"`
	Zonename              string            `json:"zonename,omitempty"`
}

// ListTemplatesResponse represents a list of templates
type ListTemplatesResponse struct {
	Count    int         `json:"count"`
	Template []*Template `json:"template"`
}
