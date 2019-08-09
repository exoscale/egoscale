package compute

import (
	"fmt"

	"github.com/exoscale/egoscale/api"
	egoerr "github.com/exoscale/egoscale/error"
	egoapi "github.com/exoscale/egoscale/internal/egoscale"
)

// AntiAffinityGroupOpts represents the Anti-Affinity Group resource creation options.
type AntiAffinityGroupCreateOpts struct {
	// Description represents the Anti-Affinity Group description.
	Description string
}

// AntiAffinityGroup represents an Anti-Affinity Group resource.
type AntiAffinityGroup struct {
	api.Resource

	ID          string
	Name        string
	Description string

	c *Client
}

func (a *AntiAffinityGroup) String() string {
	return fmt.Sprintf("AntiAffinityGroup(ID=%q, Name=%q)", a.ID, a.Name)
}

// Delete deletes the Anti-Affinity Group.
func (a *AntiAffinityGroup) Delete() error {
	if err := a.c.csError(a.c.c.BooleanRequestWithContext(a.c.ctx, &egoapi.DeleteAffinityGroup{
		ID: egoapi.MustParseUUID(a.ID),
	})); err != nil {
		return err
	}

	a.ID = ""
	a.Name = ""
	a.Description = ""

	return nil
}

// CreateAntiAffinityGroup creates a new Anti-Affinity Group resource.
func (c *Client) CreateAntiAffinityGroup(name string, opts *AntiAffinityGroupCreateOpts) (*AntiAffinityGroup, error) {
	if opts == nil {
		opts = new(AntiAffinityGroupCreateOpts)
	}

	res, err := c.c.RequestWithContext(c.ctx, &egoapi.CreateAffinityGroup{
		Type:        "host anti-affinity",
		Name:        name,
		Description: opts.Description,
	})
	if err != nil {
		return nil, err
	}

	return c.antiAffinityGroupFromAPI(res.(*egoapi.AffinityGroup)), nil
}

// ListAntiAffinityGroups returns the list of Anti-Affinity Groups.
func (c *Client) ListAntiAffinityGroups() ([]*AntiAffinityGroup, error) {
	res, err := c.c.ListWithContext(c.ctx, &egoapi.AffinityGroup{})
	if err != nil {
		return nil, err
	}

	antiAffinityGroups := make([]*AntiAffinityGroup, 0)
	for _, i := range res {
		antiAffinityGroups = append(antiAffinityGroups, c.antiAffinityGroupFromAPI(i.(*egoapi.AffinityGroup)))
	}

	return antiAffinityGroups, nil
}

// GetAntiAffinityGroupByName returns an Anti-Affinity Group by its name.
func (c *Client) GetAntiAffinityGroupByName(name string) (*AntiAffinityGroup, error) {
	return c.getAntiAffinityGroup(nil, name)
}

// GetAntiAffinityGroupByID returns an Anti-Affinity Group by its unique identifier.
func (c *Client) GetAntiAffinityGroupByID(id string) (*AntiAffinityGroup, error) {
	sgID, err := egoapi.ParseUUID(id)
	if err != nil {
		return nil, err
	}

	return c.getAntiAffinityGroup(sgID, "")
}

func (c *Client) getAntiAffinityGroup(id *egoapi.UUID, name string) (*AntiAffinityGroup, error) {
	res, err := c.c.ListWithContext(c.ctx, &egoapi.AffinityGroup{
		ID:   id,
		Name: name,
	})
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, egoerr.ErrResourceNotFound
	}

	return c.antiAffinityGroupFromAPI(res[0].(*egoapi.AffinityGroup)), nil
}

func (c *Client) antiAffinityGroupFromAPI(aag *egoapi.AffinityGroup) *AntiAffinityGroup {
	return &AntiAffinityGroup{
		Resource:    api.MarshalResource(aag),
		ID:          aag.ID.String(),
		Name:        aag.Name,
		Description: aag.Description,
		c:           c,
	}
}
