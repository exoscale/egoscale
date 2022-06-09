package v2

import (
	"context"
	"time"

	apiv2 "github.com/exoscale/egoscale/v2/api"
	"github.com/exoscale/egoscale/v2/oapi"
)

// DnsDomain represents a DNS domain.
type DnsDomain struct {
	CreatedAt   *time.Time
	ID          *string `req-for:"delete"`
	UnicodeName *string `req-for:"create"`
}

// DnsDomainRecord represents a DNS record.
type DnsDomainRecord struct {
	Content   *string `req-for:"create"`
	CreatedAt *time.Time
	ID        *string `req-for:"delete,update"`
	Name      *string `req-for:"create"`
	Priority  *int64
	Ttl       *int64
	Type      *string `req-for:"create"`
	UpdatedAt *time.Time
}

func dnsDomainFromAPI(d *oapi.DnsDomain) *DnsDomain {
	return &DnsDomain{
		CreatedAt:   d.CreatedAt,
		ID:          d.Id,
		UnicodeName: d.UnicodeName,
	}
}

func dnsDomainRecordFromAPI(d *oapi.DnsDomainRecord) *DnsDomainRecord {
	var t *string
	if d.Type != nil {
		x := string(*d.Type)
		t = &x
	}
	return &DnsDomainRecord{
		Content:   d.Content,
		CreatedAt: d.CreatedAt,
		ID:        d.Id,
		Name:      d.Name,
		Priority:  d.Priority,
		Ttl:       d.Ttl,
		Type:      t,
		UpdatedAt: d.UpdatedAt,
	}
}

// ListDnsDomains returns the list of DNS domains.
func (c *Client) ListDnsDomains(ctx context.Context, zone string) ([]DnsDomain, error) {
	var list []DnsDomain

	resp, err := c.ListDnsDomainsWithResponse(apiv2.WithZone(ctx, zone))
	if err != nil {
		return nil, err
	}

	if resp.JSON200.DnsDomains != nil {
		for _, domain := range *resp.JSON200.DnsDomains {
			list = append(list, *dnsDomainFromAPI(&domain))
		}
	}

	return list, nil
}

// GetDnsDomain returns DNS domain details.
func (c *Client) GetDnsDomain(ctx context.Context, zone, id string) (*DnsDomain, error) {
	resp, err := c.GetDnsDomainWithResponse(apiv2.WithZone(ctx, zone), id)
	if err != nil {
		return nil, err
	}

	return dnsDomainFromAPI(resp.JSON200), nil
}

// DeleteDnsDomain deletes a DNS domain.
func (c *Client) DeleteDnsDomain(ctx context.Context, zone string, domain *DnsDomain) error {
	if err := validateOperationParams(domain, "delete"); err != nil {
		return err
	}

	resp, err := c.DeleteDnsDomainWithResponse(apiv2.WithZone(ctx, zone), *domain.ID)
	if err != nil {
		return err
	}

	_, err = oapi.NewPoller().
		WithTimeout(c.timeout).
		WithInterval(c.pollInterval).
		Poll(ctx, oapi.OperationPoller(c, zone, *resp.JSON200.Id))
	if err != nil {
		return err
	}

	return nil
}

// CreateDnsDomain adds a new DNS domain.
func (c *Client) CreateDnsDomain(
	ctx context.Context,
	zone string,
	domain *DnsDomain,
) (*DnsDomain, error) {
	if err := validateOperationParams(domain, "create"); err != nil {
		return nil, err
	}

	resp, err := c.CreateDnsDomainWithResponse(apiv2.WithZone(ctx, zone), oapi.CreateDnsDomainJSONRequestBody{
		UnicodeName: domain.UnicodeName,
	})
	if err != nil {
		return nil, err
	}

	r, err := oapi.NewPoller().
		WithTimeout(c.timeout).
		WithInterval(c.pollInterval).
		Poll(ctx, oapi.OperationPoller(c, zone, *resp.JSON200.Id))
	if err != nil {
		return nil, err
	}

	return c.GetDnsDomain(ctx, zone, *r.(*struct {
		Command *string `json:"command,omitempty"`
		Id      *string `json:"id,omitempty"` // revive:disable-line
		Link    *string `json:"link,omitempty"`
	}).Id)
}

// GetDnsDomainZoneFile returns zone file of a DNS domain.
func (c *Client) GetDnsDomainZoneFile(ctx context.Context, zone, id string) ([]byte, error) {
	resp, err := c.GetDnsDomainZoneFileWithResponse(apiv2.WithZone(ctx, zone), id)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

// ListDnsDomainRecords returns the list of records for DNS domain.
func (c *Client) ListDnsDomainRecords(ctx context.Context, zone, id string) ([]DnsDomainRecord, error) {
	var list []DnsDomainRecord

	resp, err := c.ListDnsDomainRecordsWithResponse(apiv2.WithZone(ctx, zone), id)
	if err != nil {
		return nil, err
	}

	if resp.JSON200.DnsDomainRecords != nil {
		for _, record := range *resp.JSON200.DnsDomainRecords {
			list = append(list, *dnsDomainRecordFromAPI(&record))
		}
	}

	return list, nil
}

// GetDnsDomainRecord returns a single DNS domain record.
func (c *Client) GetDnsDomainRecord(ctx context.Context, zone, domainID, recordID string) (*DnsDomainRecord, error) {
	resp, err := c.GetDnsDomainRecordWithResponse(apiv2.WithZone(ctx, zone), domainID, recordID)
	if err != nil {
		return nil, err
	}

	return dnsDomainRecordFromAPI(resp.JSON200), nil
}

// CreateDnsDomainRecord adds a new DNS record for domain.
func (c *Client) CreateDnsDomainRecord(
	ctx context.Context,
	zone string,
	domainID string,
	record *DnsDomainRecord,
) (*DnsDomainRecord, error) {
	if err := validateOperationParams(record, "create"); err != nil {
		return nil, err
	}

	resp, err := c.CreateDnsDomainRecordWithResponse(apiv2.WithZone(ctx, zone), domainID, oapi.CreateDnsDomainRecordJSONRequestBody{
		Content:  *record.Content,
		Name:     *record.Name,
		Priority: record.Priority,
		Ttl:      record.Ttl,
		Type:     oapi.CreateDnsDomainRecordJSONBodyType(*record.Type),
	})
	if err != nil {
		return nil, err
	}

	r, err := oapi.NewPoller().
		WithTimeout(c.timeout).
		WithInterval(c.pollInterval).
		Poll(ctx, oapi.OperationPoller(c, zone, *resp.JSON200.Id))
	if err != nil {
		return nil, err
	}

	return c.GetDnsDomainRecord(ctx, zone, domainID, *r.(*struct {
		Command *string `json:"command,omitempty"`
		Id      *string `json:"id,omitempty"` // revive:disable-line
		Link    *string `json:"link,omitempty"`
	}).Id)
}

// DeleteDnsDomainRecord deletes a DNS domain record.
func (c *Client) DeleteDnsDomainRecord(ctx context.Context, zone, domainID string, record *DnsDomainRecord) error {
	if err := validateOperationParams(record, "delete"); err != nil {
		return err
	}

	resp, err := c.DeleteDnsDomainRecordWithResponse(apiv2.WithZone(ctx, zone), domainID, *record.ID)
	if err != nil {
		return err
	}

	_, err = oapi.NewPoller().
		WithTimeout(c.timeout).
		WithInterval(c.pollInterval).
		Poll(ctx, oapi.OperationPoller(c, zone, *resp.JSON200.Id))
	if err != nil {
		return err
	}

	return nil
}

// UpdateDnsDomainRecord updates existing DNS domain record.
func (c *Client) UpdateDnsDomainRecord(ctx context.Context, zone, domainID string, record *DnsDomainRecord) error {
	if err := validateOperationParams(record, "update"); err != nil {
		return err
	}

	resp, err := c.UpdateDnsDomainRecordWithResponse(apiv2.WithZone(ctx, zone), domainID, *record.ID, oapi.UpdateDnsDomainRecordJSONRequestBody{
		Content:  record.Content,
		Name:     record.Name,
		Priority: record.Priority,
		Ttl:      record.Ttl,
	})
	if err != nil {
		return err
	}

	_, err = oapi.NewPoller().
		WithTimeout(c.timeout).
		WithInterval(c.pollInterval).
		Poll(ctx, oapi.OperationPoller(c, zone, *resp.JSON200.Id))
	if err != nil {
		return err
	}

	return nil
}
