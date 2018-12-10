package egoscale

// Iso represents an attachable ISO disc
type Iso struct {
	Account           string            `json:"account,omitempty" doc:"the account name to which the template belongs"`
	AccountID         *UUID             `json:"accountid,omitempty" doc:"the account id to which the template belongs"`
	Bootable          bool              `json:"bootable,omitempty" doc:"true if the ISO is bootable, false otherwise"`
	Checksum          string            `json:"checksum,omitempty" doc:"checksum of the template"`
	Created           string            `json:"created,omitempty" doc:"the date this template was created"`
	CrossZones        bool              `json:"crossZones,omitempty" doc:"true if the template is managed across all Zones, false otherwise"`
	Details           map[string]string `json:"details,omitempty" doc:"additional key/value details tied with template"`
	DisplayText       string            `json:"displaytext,omitempty" doc:"the template display text"`
	Format            string            `json:"format,omitempty" doc:"the format of the template."`
	ID                *UUID             `json:"id,omitempty" doc:"the template ID"`
	IsExtractable     bool              `json:"isextractable,omitempty" doc:"true if the template is extractable, false otherwise"`
	IsFeatured        bool              `json:"isfeatured,omitempty" doc:"true if this template is a featured template, false otherwise"`
	IsPublic          bool              `json:"ispublic,omitempty" doc:"true if this template is a public template, false otherwise"`
	IsReady           bool              `json:"isready,omitempty" doc:"true if the template is ready to be deployed from, false otherwise."`
	Name              string            `json:"name,omitempty" doc:"the template name"`
	OSCategoryID      *UUID             `json:"oscategoryid,omitempty" doc:"the ID of the OS category for this template"`
	OSCategoryName    string            `json:"oscategoryname,omitempty" doc:"the name of the OS category for this template"`
	OSTypeID          *UUID             `json:"ostypeid,omitempty" doc:"the ID of the OS type for this template."`
	OSTypeName        string            `json:"ostypename,omitempty" doc:"the name of the OS type for this template."`
	PasswordEnabled   bool              `json:"passwordenabled,omitempty" doc:"true if the reset password feature is enabled, false otherwise"`
	Removed           string            `json:"removed,omitempty" doc:"the date this template was removed"`
	Size              int64             `json:"size,omitempty" doc:"the size of the template"`
	SourceTemplateID  *UUID             `json:"sourcetemplateid,omitempty" doc:"the template ID of the parent template if present"`
	SSHKeyEnabled     bool              `json:"sshkeyenabled,omitempty" doc:"true if template is sshkey enabled, false otherwise"`
	Status            string            `json:"status,omitempty" doc:"the status of the template"`
	Tags              []ResourceTag     `json:"tags,omitempty" doc:"the list of resource tags associated with tempate"`
	TemplateDirectory string            `json:"templatedirectory,omitempty" doc:"Template directory"`
	TemplateTag       string            `json:"templatetag,omitempty" doc:"the tag of this template"`
	TemplateType      string            `json:"templatetype,omitempty" doc:"the type of the template"`
	URL               string            `json:"url,omitempty" doc:"Original URL of the template where it was downloaded"`
	ZoneID            *UUID             `json:"zoneid,omitempty" doc:"the ID of the zone for this template"`
	ZoneName          string            `json:"zonename,omitempty" doc:"the name of the zone for this template"`
}

// ResourceType returns the type of the resource
func (Iso) ResourceType() string {
	return "ISO"
}

// ListRequest produces the ListIsos command.
func (iso Iso) ListRequest() (ListCommand, error) {
	req := &ListIsos{
		ID:     iso.ID,
		ZoneID: iso.ZoneID,
	}
	if iso.Bootable {
		req.Bootable = &iso.Bootable
	}
	if iso.IsFeatured {
		req.IsoFilter = "featured"
	}
	if iso.IsPublic {
		req.IsPublic = &iso.IsPublic
	}
	if iso.IsReady {
		req.IsReady = &iso.IsReady
	}

	return req, nil
}

//go:generate go run generate/main.go -interface=Listable ListIsos

// ListIsos represents the list all available ISO files request
type ListIsos struct {
	_           bool          `name:"listIsos" description:"Lists all available ISO files."`
	Bootable    *bool         `json:"bootable,omitempty" doc:"True if the ISO is bootable, false otherwise"`
	ID          *UUID         `json:"id,omitempty" doc:"List ISO by id"`
	IsoFilter   string        `json:"isofilter,omitempty" doc:"Possible values are \"featured\", \"self\", \"selfexecutable\",\"sharedexecutable\",\"executable\", and \"community\". * featured : templates that have been marked as featured and public. * self : templates that have been registered or created by the calling user. * selfexecutable : same as self, but only returns templates that can be used to deploy a new VM. * sharedexecutable : templates ready to be deployed that have been granted to the calling user by another user. * executable : templates that are owned by the calling user, or public templates, that can be used to deploy a VM. * community : templates that have been marked as public but not featured. * all : all templates (only usable by admins)."`
	IsPublic    *bool         `json:"ispublic,omitempty" doc:"True if the ISO is publicly available to all users, false otherwise."`
	IsReady     *bool         `json:"isready,omitempty" doc:"True if this ISO is ready to be deployed"`
	Keyword     string        `json:"keyword,omitempty" doc:"List by keyword"`
	Name        string        `json:"name,omitempty" doc:"List all isos by name"`
	Page        int           `json:"page,omitempty"`
	PageSize    int           `json:"pagesize,omitempty"`
	ShowRemoved *bool         `json:"showremoved,omitempty" doc:"Show removed ISOs as well"`
	Tags        []ResourceTag `json:"tags,omitempty" doc:"List resources by tags (key/value pairs)"`
	ZoneID      *UUID         `json:"zoneid,omitempty" doc:"The ID of the zone"`
}

// ListIsosResponse represents a list of ISO files
type ListIsosResponse struct {
	Count int   `json:"count"`
	Iso   []Iso `json:"iso"`
}

// AttachIso represents the request to attach an ISO to a virtual machine.
type AttachIso struct {
	_                bool  `name:"attachIso" description:"Attaches an ISO to a virtual machine."`
	ID               *UUID `json:"id" doc:"the ID of the ISO file"`
	VirtualNachineID *UUID `json:"virtualmachineid" doc:"the ID of the virtual machine"`
}

// Response returns the struct to unmarshal
func (AttachIso) Response() interface{} {
	return new(AsyncJobResult)
}

// AsyncResponse returns the struct to unmarshal the async job
func (AttachIso) AsyncResponse() interface{} {
	return new(VirtualMachine)
}

// DetachIso represents the request to detach an ISO to a virtual machine.
type DetachIso struct {
	_                bool  `name:"detachIso" description:"Detaches any ISO file (if any) currently attached to a virtual machine."`
	VirtualNachineID *UUID `json:"virtualmachineid" doc:"The ID of the virtual machine"`
}

// Response returns the struct to unmarshal
func (DetachIso) Response() interface{} {
	return new(AsyncJobResult)
}

// AsyncResponse returns the struct to unmarshal the async job
func (DetachIso) AsyncResponse() interface{} {
	return new(VirtualMachine)
}
