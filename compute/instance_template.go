package compute

import (
	"fmt"
	"time"

	"github.com/exoscale/egoscale/api"
	egoerr "github.com/exoscale/egoscale/error"
	egoapi "github.com/exoscale/egoscale/internal/egoscale"
)

// InstanceTemplate represents an Exoscale Compute instance template resource.
type InstanceTemplate struct {
	api.Resource

	ID                   string
	Name                 string
	Description          string
	Date                 time.Time
	Size                 int64
	SSHKeyEnabled        bool
	PasswordResetEnabled bool
	Username             string
	Zone                 *Zone

	c *Client
}

func (t *InstanceTemplate) String() string {
	return fmt.Sprintf("InstanceTemplate(ID=%q, Name=%q)", t.ID, t.Name)
}

// ListInstanceTemplates returns the list of Exoscale Compute instance templates optionally matching a name, template
// type (exoscale/mine).
func (c *Client) ListInstanceTemplates(zone *Zone, name, typ string) ([]*InstanceTemplate, error) {
	if zone == nil {
		return nil, egoerr.ErrMissingZone
	}

	if typ == "" {
		typ = "exoscale"
	}
	vtype, err := validateTemplateType(typ)
	if err != nil {
		return nil, err
	}

	res, err := c.c.ListWithContext(c.ctx, &egoapi.ListTemplates{
		ZoneID:         egoapi.MustParseUUID(zone.ID),
		Name:           name,
		TemplateFilter: vtype,
	})
	if err != nil {
		return nil, err
	}

	instanceTemplates := make([]*InstanceTemplate, 0)
	for _, i := range res {
		instanceTemplate, err := c.instanceTemplateFromAPI(i.(*egoapi.Template))
		if err != nil {
			return nil, err
		}
		instanceTemplates = append(instanceTemplates, instanceTemplate)
	}

	return instanceTemplates, nil
}

// GetInstanceTemplate returns an Exoscale Compute instance template by its unique identifier.
func (c *Client) GetInstanceTemplate(zone *Zone, id, typ string) (*InstanceTemplate, error) {
	if zone == nil {
		return nil, egoerr.ErrMissingZone
	}

	if typ == "" {
		typ = "exoscale"
	}
	vtype, err := validateTemplateType(typ)
	if err != nil {
		return nil, err
	}

	res, err := c.c.ListWithContext(c.ctx, &egoapi.ListTemplates{
		ZoneID:         egoapi.MustParseUUID(zone.ID),
		ID:             egoapi.MustParseUUID(id),
		TemplateFilter: vtype,
	})
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, egoerr.ErrResourceNotFound
	}

	return c.instanceTemplateFromAPI(res[0].(*egoapi.Template))
}

// TODO: RegisterInstanceTemplate()

func (c *Client) instanceTemplateFromAPI(t *egoapi.Template) (*InstanceTemplate, error) {
	zone, err := c.GetZoneByID(t.ZoneID.String())
	if err != nil {
		return nil, err
	}

	date, err := time.Parse("2006-01-02T15:04:05-0700", t.Created)
	if err != nil {
		return nil, err
	}

	username, ok := t.Details["username"]
	if !ok {
		username = ""
	}

	return &InstanceTemplate{
		Resource:             api.MarshalResource(t),
		ID:                   t.ID.String(),
		Name:                 t.Name,
		Description:          t.DisplayText,
		Zone:                 zone,
		Date:                 date,
		Size:                 t.Size,
		SSHKeyEnabled:        t.SSHKeyEnabled,
		PasswordResetEnabled: t.PasswordEnabled,
		Username:             username,
		c:                    c,
	}, nil
}

// validateTemplateType checks that the specified template instance type is valid, and returns the unaliased (CS) type
// to request to the Exoscale API.
func validateTemplateType(typ string) (string, error) {
	var templateTypes = map[string]string{"exoscale": "featured", "mine": "self"}

	if _, ok := templateTypes[typ]; typ != "" && !ok {
		return "", fmt.Errorf("invalid instance template type %q, supported values are: exoscale, mine", typ)
	}

	return templateTypes[typ], nil
}
