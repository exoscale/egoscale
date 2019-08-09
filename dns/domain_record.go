package dns

import (
	"fmt"

	"github.com/exoscale/egoscale/api"
	egoapi "github.com/exoscale/egoscale/internal/egoscale"
)

// DomainRecordCreateOpts represents the DNS domain record creation options.
type DomainRecordCreateOpts struct {
	// Name represents the DNS domain record name.
	Name string
	// Type represents the DNS domain record type.
	Type string
	// Content represents the DNS domain record content.
	Content string
	// Priority represents the DNS domain record priority.
	Priority int
	// TTL represents the DNS domain record TTL.
	TTL int
}

// DomainRecordUpdateOpts represents the DNS domain record update options.
type DomainRecordUpdateOpts struct {
	// Name represents the DNS domain record name.
	Name string
	// Content represents the DNS domain record content.
	Content string
	// Priority represents the DNS domain record priority.
	Priority int
	// TTL represents the DNS domain record TTL.
	TTL int
}

// DomainRecord represents a DNS domain record.
type DomainRecord struct {
	api.Resource

	ID       int64
	Type     string
	Name     string
	Content  string
	Priority int
	TTL      int
	Domain   *Domain

	c *Client
}

// Update updates the DNS domain record.
func (r *DomainRecord) Update(opts *DomainRecordUpdateOpts) error {
	var req = egoapi.UpdateDNSRecord{
		DomainID: r.Domain.ID,
		ID:       r.ID,
		Type:     r.Type,
		Content:  r.Content,
		TTL:      r.TTL,
	}

	if opts == nil {
		return fmt.Errorf("no update options specified")
	}

	if opts.Name != "" {
		req.Name = opts.Name
	}
	if opts.Content != "" {
		req.Content = opts.Content
	}
	if opts.Priority > 0 {
		req.Priority = opts.Priority
	}
	if opts.TTL > 0 {
		req.TTL = opts.TTL
	}

	res, err := r.c.c.Request(&req)
	if err != nil {
		return err
	}
	record := res.(*egoapi.DNSRecord)

	r.Name = record.Name
	r.Content = record.Content
	r.Priority = record.Priority
	r.TTL = record.TTL

	return nil
}

// Delete deletes the DNS domain record.
func (r *DomainRecord) Delete() error {
	if err := r.c.csError(r.c.c.BooleanRequestWithContext(r.c.ctx, &egoapi.DeleteDNSRecord{
		DomainID: r.Domain.ID,
		ID:       r.ID,
	})); err != nil {
		return err
	}

	r.ID = 0
	r.Name = ""
	r.Type = ""
	r.Content = ""
	r.Priority = 0
	r.TTL = 0
	r.Domain = nil

	return nil
}

func (c *Client) domainRecordFromAPI(record *egoapi.DNSRecord, domain *Domain) *DomainRecord {
	return &DomainRecord{
		Resource: api.MarshalResource(record),
		ID:       record.ID,
		Type:     record.RecordType,
		Name:     record.Name,
		Content:  record.Content,
		Priority: record.Priority,
		TTL:      record.TTL,
		Domain:   domain,
		c:        c,
	}
}
