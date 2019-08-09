package compute

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/exoscale/egoscale/api"
	egoerr "github.com/exoscale/egoscale/error"
	egoapi "github.com/exoscale/egoscale/internal/egoscale"
)

// ElasticIPCreateOpts represents the Elastic IP resource creation options.
type ElasticIPCreateOpts struct {
	// HealthcheckMode represents the Elastic IP health checking mode (for managed Elastic IPs only).
	HealthcheckMode string
	// HealthcheckPort represents the Elastic IP health checking port (for managed Elastic IPs only).
	HealthcheckPort uint16
	// HealthcheckPath represents the Elastic IP health checking path (for managed Elastic IPs only).
	HealthcheckPath string
	// HealthcheckInterval represents the Elastic IP health checking interval (for managed Elastic IPs only).
	HealthcheckInterval time.Duration
	// HealthcheckTimeout represents the Elastic IP health checking timeout (for managed Elastic IPs only).
	HealthcheckTimeout time.Duration
	// HealthcheckStrikesOK represents the Elastic IP health checking number of OK strikes (for managed Elastic IPs
	// only).
	HealthcheckStrikesOK int
	// HealthcheckStrikesFail represents the Elastic IP health checking number of failed strikes (for managed Elastic
	// IPs only).
	HealthcheckStrikesFail int
}

// ElasticIPUpdateOpts represents the Elastic IP resource update options.
type ElasticIPUpdateOpts struct {
	// HealthcheckMode represents the Elastic IP health checking mode (for managed Elastic IPs only).
	HealthcheckMode string
	// HealthcheckPort represents the Elastic IP health checking port (for managed Elastic IPs only).
	HealthcheckPort uint16
	// HealthcheckPath represents the Elastic IP health checking path (for managed Elastic IPs only).
	HealthcheckPath string
	// HealthcheckInterval represents the Elastic IP health checking interval (for managed Elastic IPs only).
	HealthcheckInterval time.Duration
	// HealthcheckTimeout represents the Elastic IP health checking timeout (for managed Elastic IPs only).
	HealthcheckTimeout time.Duration
	// HealthcheckStrikesOK represents the Elastic IP health checking number of OK strikes (for managed Elastic IPsi
	// only).
	HealthcheckStrikesOK int
	// HealthcheckStrikesFail represents the Elastic IP health checking number of failed strikes (for managed Elastic
	// IPs only).
	HealthcheckStrikesFail int
}

// ElasticIP represents an Exoscale Elastic IP resource.
type ElasticIP struct {
	api.Resource

	ID                     string
	Address                net.IP
	HealthcheckMode        string
	HealthcheckPort        uint16
	HealthcheckPath        string
	HealthcheckInterval    time.Duration
	HealthcheckTimeout     time.Duration
	HealthcheckStrikesOK   int
	HealthcheckStrikesFail int
	Zone                   *Zone

	c *Client
}

func (e *ElasticIP) String() string {
	return fmt.Sprintf("ElasticIP(ID=%q, Address=%q)", e.ID, e.Address.String())
}

// ReverseDNS retrieves the current Elastic IP address reverse DNS record.
func (e *ElasticIP) ReverseDNS() (string, error) {
	res, err := e.c.c.RequestWithContext(e.c.ctx, &egoapi.QueryReverseDNSForPublicIPAddress{
		ID: egoapi.MustParseUUID(e.ID),
	})
	if err != nil {
		return "", e.c.csError(err)
	}

	ip := res.(*egoapi.IPAddress)
	if len(ip.ReverseDNS) == 0 {
		return "", nil
	}

	return ip.ReverseDNS[0].DomainName, nil
}

// SetReverseDNS sets the Elastic IP address reverse DNS record.
func (e *ElasticIP) SetReverseDNS(record string) error {
	if _, err := e.c.c.RequestWithContext(e.c.ctx, &egoapi.UpdateReverseDNSForPublicIPAddress{
		ID:         egoapi.MustParseUUID(e.ID),
		DomainName: record,
	}); err != nil {
		return e.c.csError(err)
	}

	return nil
}

// UnsetReverseDNS unsets the Elastic IP address reverse DNS record.
func (e *ElasticIP) UnsetReverseDNS() error {
	if err := e.c.c.BooleanRequestWithContext(e.c.ctx, &egoapi.DeleteReverseDNSFromPublicIPAddress{
		ID: egoapi.MustParseUUID(e.ID),
	}); err != nil {
		return e.c.csError(err)
	}

	return nil
}

// Instances returns the list of Compute instances attached to the Elastic IP.
func (e *ElasticIP) Instances() ([]*Instance, error) {
	res, err := e.c.c.ListWithContext(e.c.ctx, &egoapi.VirtualMachine{ZoneID: egoapi.MustParseUUID(e.Zone.ID)})
	if err != nil {
		return nil, err
	}

	instances := make([]*Instance, 0)
	for _, i := range res {
		instance, err := e.c.instanceFromAPI(i.(*egoapi.VirtualMachine))
		if err != nil {
			return nil, err
		}

		nic, err := instance.defaultNIC()
		if err != nil {
			return nil, err
		}

		for _, ip := range nic.SecondaryIP {
			if ip.IPAddress.Equal(e.Address) {
				instances = append(instances, instance)
			}
		}
	}

	return instances, nil
}

// AttachInstance attaches the Elastic IP to a Compute instance.
func (e *ElasticIP) AttachInstance(instance *Instance) error {
	if instance == nil {
		return errors.New("missing Compute instance")
	}

	return instance.AttachElasticIP(e)
}

// DetachInstance detaches the Elastic IP from a Compute instance.
func (e *ElasticIP) DetachInstance(instance *Instance) error {
	if instance == nil {
		return errors.New("missing Compute instance")
	}

	return instance.DetachElasticIP(e)
}

// Update updates the Elastic IP properties.
func (e *ElasticIP) Update(opts *ElasticIPUpdateOpts) error {
	if opts == nil {
		return fmt.Errorf("no update options specified")
	}

	if _, err := e.c.c.RequestWithContext(e.c.ctx, &egoapi.UpdateIPAddress{
		ID:                     egoapi.MustParseUUID(e.ID),
		HealthcheckMode:        opts.HealthcheckMode,
		HealthcheckPort:        int64(opts.HealthcheckPort),
		HealthcheckPath:        opts.HealthcheckPath,
		HealthcheckInterval:    int64(opts.HealthcheckInterval.Seconds()),
		HealthcheckTimeout:     int64(opts.HealthcheckTimeout.Seconds()),
		HealthcheckStrikesOk:   int64(opts.HealthcheckStrikesOK),
		HealthcheckStrikesFail: int64(opts.HealthcheckStrikesFail),
	}); err != nil {
		return e.c.csError(err)
	}

	if opts != nil {
		if opts.HealthcheckMode != "" {
			e.HealthcheckMode = opts.HealthcheckMode
		}
		if opts.HealthcheckPort > 0 {
			e.HealthcheckPort = opts.HealthcheckPort
		}
		if opts.HealthcheckPath != "" {
			e.HealthcheckPath = opts.HealthcheckPath
		}
		if opts.HealthcheckInterval > 0 {
			e.HealthcheckInterval = opts.HealthcheckInterval
		}
		if opts.HealthcheckTimeout > 0 {
			e.HealthcheckTimeout = opts.HealthcheckTimeout
		}
		if opts.HealthcheckStrikesOK > 0 {
			e.HealthcheckStrikesOK = opts.HealthcheckStrikesOK
		}
		if opts.HealthcheckStrikesFail > 0 {
			e.HealthcheckStrikesFail = opts.HealthcheckStrikesFail
		}
	}

	return nil
}

// Delete deletes the Elastic IP.
func (e *ElasticIP) Delete() error {
	if err := e.c.csError(e.c.c.BooleanRequestWithContext(e.c.ctx, &egoapi.DisassociateIPAddress{
		ID: egoapi.MustParseUUID(e.ID)},
	)); err != nil {
		return e.c.csError(err)
	}

	e.ID = ""
	e.Zone = nil
	e.Address = nil
	e.HealthcheckMode = ""
	e.HealthcheckPort = 0
	e.HealthcheckPath = ""
	e.HealthcheckInterval = 0
	e.HealthcheckTimeout = 0
	e.HealthcheckStrikesOK = 0
	e.HealthcheckStrikesFail = 0

	return nil
}

// CreateElasticIP creates a new Elastic IP resource.
func (c *Client) CreateElasticIP(zone *Zone, opts *ElasticIPCreateOpts) (*ElasticIP, error) {
	if zone == nil {
		return nil, egoerr.ErrMissingZone
	}

	if opts == nil {
		opts = new(ElasticIPCreateOpts)
	}

	res, err := c.c.RequestWithContext(c.ctx, &egoapi.AssociateIPAddress{
		ZoneID:                 egoapi.MustParseUUID(zone.ID),
		HealthcheckMode:        opts.HealthcheckMode,
		HealthcheckPort:        int64(opts.HealthcheckPort),
		HealthcheckPath:        opts.HealthcheckPath,
		HealthcheckInterval:    int64(opts.HealthcheckInterval.Seconds()),
		HealthcheckTimeout:     int64(opts.HealthcheckTimeout.Seconds()),
		HealthcheckStrikesOk:   int64(opts.HealthcheckStrikesOK),
		HealthcheckStrikesFail: int64(opts.HealthcheckStrikesFail),
	})
	if err != nil {
		return nil, c.csError(err)
	}

	return c.elasticIPFromAPI(res.(*egoapi.IPAddress))
}

// ListElasticIPs returns the list of Elastic IPs.
func (c *Client) ListElasticIPs(zone *Zone) ([]*ElasticIP, error) {
	if zone == nil {
		return nil, egoerr.ErrMissingZone
	}

	res, err := c.c.ListWithContext(c.ctx, &egoapi.IPAddress{ZoneID: egoapi.MustParseUUID(zone.ID)})
	if err != nil {
		return nil, c.csError(err)
	}

	elasticIPs := make([]*ElasticIP, 0)
	for _, i := range res {
		elasticIP, err := c.elasticIPFromAPI(i.(*egoapi.IPAddress))
		if err != nil {
			return nil, err
		}

		elasticIPs = append(elasticIPs, elasticIP)
	}

	return elasticIPs, nil
}

// GetElasticIPByID returns an Elastic IP by its unique identifier.
func (c *Client) GetElasticIPByID(zone *Zone, id string) (*ElasticIP, error) {
	if zone == nil {
		return nil, egoerr.ErrMissingZone
	}

	eipID, err := egoapi.ParseUUID(id)
	if err != nil {
		return nil, err
	}

	return c.getElasticIP(egoapi.MustParseUUID(zone.ID), eipID, nil)
}

// GetElasticIPByAddress returns an Elastic IP by its IP address.
func (c *Client) GetElasticIPByAddress(zone *Zone, address string) (*ElasticIP, error) {
	if zone == nil {
		return nil, egoerr.ErrMissingZone
	}

	return c.getElasticIP(egoapi.MustParseUUID(zone.ID), nil, net.ParseIP(address))
}

func (c *Client) getElasticIP(zoneID, id *egoapi.UUID, address net.IP) (*ElasticIP, error) {
	res, err := c.c.ListWithContext(c.ctx, &egoapi.IPAddress{
		ZoneID:    zoneID,
		ID:        id,
		IPAddress: address,
	})
	if err != nil {
		return nil, c.csError(err)
	}

	if len(res) == 0 {
		return nil, egoerr.ErrResourceNotFound
	}

	return c.elasticIPFromAPI(res[0].(*egoapi.IPAddress))
}

func (c *Client) elasticIPFromAPI(eip *egoapi.IPAddress) (*ElasticIP, error) {
	var elasticIP = ElasticIP{

		Resource: api.MarshalResource(eip),
		ID:       eip.ID.String(),
		Address:  eip.IPAddress,
		c:        c,
	}

	zone, err := c.GetZoneByID(eip.ZoneID.String())
	if err != nil {
		return nil, err
	}
	elasticIP.Zone = zone

	if healthcheck := eip.Healthcheck; healthcheck != nil {
		elasticIP.HealthcheckMode = healthcheck.Mode
		elasticIP.HealthcheckPort = uint16(healthcheck.Port)
		elasticIP.HealthcheckPath = healthcheck.Path
		elasticIP.HealthcheckInterval = time.Second * time.Duration(healthcheck.Interval)
		elasticIP.HealthcheckTimeout = time.Second * time.Duration(healthcheck.Timeout)
		elasticIP.HealthcheckStrikesOK = int(healthcheck.StrikesOk)
		elasticIP.HealthcheckStrikesFail = int(healthcheck.StrikesFail)
	}

	return &elasticIP, nil
}
