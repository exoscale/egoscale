package compute

import (
	"fmt"

	"github.com/exoscale/egoscale/api"
	egoerr "github.com/exoscale/egoscale/error"
	egoapi "github.com/exoscale/egoscale/internal/egoscale"
)

// Zone represents an Exoscale zone.
type Zone struct {
	api.Resource

	ID   string
	Name string

	c *Client
}

func (z *Zone) String() string {
	return fmt.Sprintf("Zone(ID=%q, Name=%q)", z.ID, z.Name)
}

// ListZones returns the list of available Exoscale zones.
func (c *Client) ListZones() ([]*Zone, error) {
	res, err := c.c.ListWithContext(c.ctx, &egoapi.Zone{})
	if err != nil {
		return nil, err
	}

	zones := make([]*Zone, 0)
	for _, i := range res {
		zones = append(zones, c.zoneFromAPI(i.(*egoapi.Zone)))
	}

	return zones, nil
}

// GetZone returns an Exoscale zone by its name.
func (c *Client) GetZoneByName(name string) (*Zone, error) {
	return c.getZone(nil, name)
}

// GetZone returns an Exoscale zone by its unique identifier.
func (c *Client) GetZoneByID(id string) (*Zone, error) {
	zoneID, err := egoapi.ParseUUID(id)
	if err != nil {
		return nil, err
	}

	return c.getZone(zoneID, "")
}

func (c *Client) getZone(id *egoapi.UUID, name string) (*Zone, error) {
	res, err := c.c.ListWithContext(c.ctx, &egoapi.Zone{
		ID:   id,
		Name: name,
	})
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, egoerr.ErrResourceNotFound
	}

	return c.zoneFromAPI(res[0].(*egoapi.Zone)), nil
}

func (c *Client) zoneFromAPI(zone *egoapi.Zone) *Zone {
	return &Zone{
		Resource: api.MarshalResource(zone),
		ID:       zone.ID.String(),
		Name:     zone.Name,
		c:        c,
	}
}
