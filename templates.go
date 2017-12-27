package egoscale

import (
	"fmt"
	"regexp"
	"strings"
)

// Template represents a machine to be deployed
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

// ListTemplatesRequest represents a template query filter
type ListTemplatesRequest struct {
	TemplateFilter string         `json:"templatefilter"` // featured, etc.
	Account        string         `json:"account,omitempty"`
	DomainId       string         `json:"domainid,omitempty"`
	Hypervisor     string         `json:"hypervisor,omitempty"`
	Id             string         `json:"id,omitempty"`
	IsRecursive    bool           `json:"isrecursive,omitempty"`
	Keyword        string         `json:"keyword,omitempty"`
	ListAll        bool           `json:"listall,omitempty"`
	Name           string         `json:"name,omitempty"`
	Page           int            `json:"page,omitempty"`
	PageSize       int            `json:"pagesize,omitempty"`
	ProjectId      string         `json:"projectid,omitempty"`
	ShowRemoved    bool           `json:"showremoved,omitempty"`
	Tags           []*ResourceTag `json:"tags,omitempty"`
	ZoneId         string         `json:"zoneid,omitempty"`
}

// Command returns the CloudStack API command
func (req *ListTemplatesRequest) Command() string {
	return "listTemplates"
}

// ListTemplatesResponse represents a list of templates
type ListTemplatesResponse struct {
	Count    int         `json:"count"`
	Template []*Template `json:"template"`
}

// ListTemplates returns the templates
func (exo *Client) ListTemplates(req *ListTemplatesRequest) ([]*Template, error) {
	var r ListTemplatesResponse
	err := exo.Request(req, &r)
	if err != nil {
		return nil, err
	}

	return r.Template, nil
}

// GetImages list the available featured images and group them by name, then size.
func (exo *Client) GetImages() (map[string]map[int64]string, error) {
	var images map[string]map[int64]string
	images = make(map[string]map[int64]string)
	re := regexp.MustCompile(`^Linux (?P<name>.+?) (?P<version>[0-9.]+)\b`)

	templates, err := exo.ListTemplates(&ListTemplatesRequest{
		TemplateFilter: "featured",
	})
	if err != nil {
		return images, err
	}

	for _, template := range templates {
		size := int64(template.Size >> 30) // B to GiB

		fullname := strings.ToLower(template.Name)

		if _, present := images[fullname]; !present {
			images[fullname] = make(map[int64]string)
		}
		images[fullname][size] = template.Id

		submatch := re.FindStringSubmatch(template.Name)
		if len(submatch) > 0 {
			name := strings.Replace(strings.ToLower(submatch[1]), " ", "-", -1)
			version := submatch[2]
			image := fmt.Sprintf("%s-%s", name, version)

			if _, present := images[image]; !present {
				images[image] = make(map[int64]string)
			}
			images[image][size] = template.Id
		}
	}
	return images, nil
}
