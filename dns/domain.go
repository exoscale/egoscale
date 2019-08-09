package dns

import (
	"fmt"

	"github.com/exoscale/egoscale/api"
	egoerr "github.com/exoscale/egoscale/error"
	egoapi "github.com/exoscale/egoscale/internal/egoscale"
)

// Domain represents a DNS domain.
type Domain struct {
	api.Resource

	ID          int64
	Name        string
	UnicodeName string

	c *Client
}

// Records returns the DNS domain records.
func (d *Domain) Records() ([]*DomainRecord, error) {
	res, err := d.c.c.ListWithContext(d.c.ctx, &egoapi.ListDNSRecords{DomainID: d.ID})
	if err != nil {
		return nil, err
	}

	records := make([]*DomainRecord, 0)
	for _, i := range res {
		records = append(records, d.c.domainRecordFromAPI(i.(*egoapi.DNSRecord), d))
	}

	return records, nil
}

// AddRecord adds a new DNS domain record.
func (d *Domain) AddRecord(opts *DomainRecordCreateOpts) (*DomainRecord, error) {
	if opts == nil {
		return nil, fmt.Errorf("no creation options specified")
	}

	res, err := d.c.c.Request(&egoapi.CreateDNSRecord{
		Domain:   d.Name,
		Name:     opts.Name,
		Type:     opts.Type,
		Content:  opts.Content,
		Priority: opts.Priority,
		TTL:      opts.TTL,
	})
	if err != nil {
		return nil, err
	}

	return d.c.domainRecordFromAPI(res.(*egoapi.DNSRecord), d), nil
}

// Delete deletes the DNS domain.
func (d *Domain) Delete() error {
	if err := d.c.csError(d.c.c.BooleanRequestWithContext(d.c.ctx,
		&egoapi.DeleteDNSDomain{Name: d.Name})); err != nil {
		return err
	}

	d.ID = 0
	d.Name = ""
	d.UnicodeName = ""

	return nil
}

// CreateDomain creates a new DNS domain resource identified by name, and returns a Domain object if successful or an
// error.
func (c *Client) CreateDomain(name string) (*Domain, error) {
	res, err := c.c.Request(&egoapi.CreateDNSDomain{Name: name})
	if err != nil {
		return nil, err
	}

	return c.domainFromAPI(res.(*egoapi.DNSDomain)), nil
}

// ListDomains returns the list of DNS domains, or an error if the API query failed.
func (c *Client) ListDomains() ([]*Domain, error) {
	res, err := c.c.ListWithContext(c.ctx, &egoapi.DNSDomain{})
	if err != nil {
		return nil, err
	}

	domains := make([]*Domain, 0)
	for _, i := range res {
		domains = append(domains, c.domainFromAPI(i.(*egoapi.DNSDomain)))
	}

	return domains, nil
}

// GetDomainByName returns a DNS domain by its name.
func (c *Client) GetDomainByName(name string) (*Domain, error) {
	return c.getDomain(0, name)
}

// GetDomain returns a DNS domain by its unique identifier.
func (c *Client) GetDomainByID(id int64) (*Domain, error) {
	return c.getDomain(id, "")
}

func (c *Client) getDomain(id int64, name string) (*Domain, error) {
	domains, err := c.ListDomains()
	if err != nil {
		return nil, err
	}

	for _, domain := range domains {
		if (name != "" && domain.Name == name) || (id > 0 && domain.ID == id) {
			return domain, nil
		}
	}

	return nil, egoerr.ErrResourceNotFound
}

func (c *Client) domainFromAPI(domain *egoapi.DNSDomain) *Domain {
	return &Domain{
		Resource:    api.MarshalResource(domain),
		ID:          domain.ID,
		Name:        domain.Name,
		UnicodeName: domain.UnicodeName,
		c:           c,
	}
}
