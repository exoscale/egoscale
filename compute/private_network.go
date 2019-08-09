package compute

import (
	"fmt"
	"net"

	"github.com/exoscale/egoscale/api"
	egoerr "github.com/exoscale/egoscale/error"
	egoapi "github.com/exoscale/egoscale/internal/egoscale"
)

// PrivateNetworkCreateOpts represents the Private Network resource creation options.
type PrivateNetworkCreateOpts struct {
	// Name represents the Private Network name.
	Name string
	// Description represents the Private Network name.
	Description string
	// StartIP represents the Private Network start IP (for managed Private Networks only).
	StartIP net.IP
	// EndIP represents the Private Network end IP (for managed Private Networks only).
	EndIP net.IP
	// Netmask represents the Private Network netmask (for managed Private Networks only).
	Netmask net.IP
}

// PrivateNetworkUpdateOpts represents the Private Network resource update options.
type PrivateNetworkUpdateOpts struct {
	// Name represents the Private Network name.
	Name string
	// Description represents the Private Network name.
	Description string
	// StartIP represents the Private Network start IP (for managed Private Networks only).
	StartIP net.IP
	// EndIP represents the Private Network end IP (for managed Private Networks only).
	EndIP net.IP
	// Netmask represents the Private Network netmask (for managed Private Networks only).
	Netmask net.IP
}

// PrivateNetwork represents a Private Network resource.
type PrivateNetwork struct {
	api.Resource

	ID          string
	Name        string
	Description string
	StartIP     net.IP
	EndIP       net.IP
	Netmask     net.IP
	Zone        *Zone

	c *Client
}

func (p *PrivateNetwork) String() string {
	return fmt.Sprintf("PrivateNetwork(ID=%q, Name=%q)", p.ID, p.Name)
}

// TODO: PrivateNetwork.Instances()

// TODO: PrivateNetwork.AttachInstance()

// TODO: PrivateNetwork.DetachInstance()

// Update updates the Private Network properties.
func (p *PrivateNetwork) Update(opts *PrivateNetworkUpdateOpts) error {
	if opts == nil {
		return fmt.Errorf("no update options specified")
	}

	if _, err := p.c.c.RequestWithContext(p.c.ctx, &egoapi.UpdateNetwork{
		ID:          egoapi.MustParseUUID(p.ID),
		Name:        opts.Name,
		DisplayText: opts.Description,
		StartIP:     opts.StartIP,
		EndIP:       opts.EndIP,
		Netmask:     opts.Netmask,
	}); err != nil {
		return err
	}

	if opts != nil {
		if opts.Name != "" {
			p.Name = opts.Name
		}
		if opts.Description != "" {
			p.Description = opts.Description
		}
		if opts.StartIP != nil {
			p.StartIP = opts.StartIP
		}
		if opts.EndIP != nil {
			p.EndIP = opts.EndIP
		}
		if opts.Netmask != nil {
			p.Netmask = opts.Netmask
		}
	}

	return nil
}

// Delete deletes the Private Network.
func (p *PrivateNetwork) Delete() error {
	if err := p.c.csError(p.c.c.BooleanRequestWithContext(p.c.ctx,
		&egoapi.DeleteNetwork{ID: egoapi.MustParseUUID(p.ID)})); err != nil {
		return err
	}

	p.ID = ""
	p.Name = ""
	p.Description = ""
	p.StartIP = nil
	p.EndIP = nil
	p.Netmask = nil

	return nil
}

// CreatePrivateNetwork creates a new Private Network resource identified by name.
func (c *Client) CreatePrivateNetwork(zone *Zone, opts *PrivateNetworkCreateOpts) (*PrivateNetwork, error) {
	if zone == nil {
		return nil, egoerr.ErrMissingZone
	}

	if opts == nil {
		opts = new(PrivateNetworkCreateOpts)
	}

	res, err := c.c.RequestWithContext(c.ctx, &egoapi.CreateNetwork{
		ZoneID:      egoapi.MustParseUUID(zone.ID),
		Name:        opts.Name,
		DisplayText: opts.Description,
		StartIP:     opts.StartIP,
		EndIP:       opts.EndIP,
		Netmask:     opts.Netmask,
	})
	if err != nil {
		return nil, err
	}

	return c.privateNetworkFromAPI(res.(*egoapi.Network))
}

// ListPrivateNetworks returns the list of Private Networks.
func (c *Client) ListPrivateNetworks(zone *Zone) ([]*PrivateNetwork, error) {
	if zone == nil {
		return nil, egoerr.ErrMissingZone
	}

	res, err := c.c.ListWithContext(c.ctx, &egoapi.Network{ZoneID: egoapi.MustParseUUID(zone.ID)})
	if err != nil {
		return nil, err
	}

	privateNetworks := make([]*PrivateNetwork, 0)
	for _, i := range res {
		privateNetwork, err := c.privateNetworkFromAPI(i.(*egoapi.Network))
		if err != nil {
			return nil, err
		}

		privateNetworks = append(privateNetworks, privateNetwork)
	}

	return privateNetworks, nil
}

// GetPrivateNetwork returns a Private Network located in a zone by its unique identifier.
func (c *Client) GetPrivateNetwork(zone *Zone, id string) (*PrivateNetwork, error) {
	if zone == nil {
		return nil, egoerr.ErrMissingZone
	}

	res, err := c.c.ListWithContext(c.ctx, &egoapi.Network{
		ID:     egoapi.MustParseUUID(id),
		ZoneID: egoapi.MustParseUUID(zone.ID),
	})
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, egoerr.ErrResourceNotFound
	}

	return c.privateNetworkFromAPI(res[0].(*egoapi.Network))
}

func (c *Client) privateNetworkFromAPI(n *egoapi.Network) (*PrivateNetwork, error) {
	zone, err := c.GetZoneByID(n.ZoneID.String())
	if err != nil {
		return nil, err
	}

	return &PrivateNetwork{
		Resource:    api.MarshalResource(n),
		ID:          n.ID.String(),
		Name:        n.Name,
		Description: n.DisplayText,
		Zone:        zone,
		StartIP:     n.StartIP,
		EndIP:       n.EndIP,
		Netmask:     n.Netmask,
		c:           c,
	}, nil
}
