package egoscale

import (
	"context"
	"fmt"
)

// Delete deletes the given DNS domain by Name
func (record DNSRecord) Delete(ctx context.Context, client *Client) error {
	if record.ID == 0 {
		return fmt.Errorf("a DNS domain record may only be deleted using ID and DomainID")
	}

	return client.BooleanRequestWithContext(ctx, &DeleteDNSRecord{
		DomainID: record.DomainID,
		ID:       record.ID,
	})
}

// ListRequest builds the ListDNSRecords request
func (record DNSRecord) ListRequest(domainID int64) (ListCommand, error) {
	req := &ListDNSRecords{
		DomainID: domainID,
	}

	return req, nil
}

// CreateDNSRecord represents a new DNS domain record to be created
type CreateDNSRecord struct {
	Domain   string `json:"name" doc:"Name of the DNS domain the record belongs to"`
	Name     string `json:"record_name" doc:"Name of the DNS record"`
	Type     string `json:"record_type" doc:"Type of the DNS record"`
	Content  string `json:"content" doc:"Content of the DNS record"`
	Priority int    `json:"priority" doc:"Priority of the DNS record"`
	TTL      int    `json:"ttl" doc:"TTL of the DNS record"`
	_        bool   `name:"createDnsDomainRecord" description:"Create a new DNS domain"`
}

// Response returns the struct to unmarshal
func (CreateDNSRecord) Response() interface{} {
	return new(DNSRecord)
}

// UpdateDNSRecord represents a DNS record to be updated
type UpdateDNSRecord struct {
	DomainID int64  `json:"id,omitempty" description:"ID of the DNS domain the record belongs to"`
	ID       int64  `json:"record_id,omitempty" description:"ID of the DNS record"`
	Name     string `json:"record_name,omitempty" description:"Updated name of the DNS record"`
	Type     string `json:"record_type"`
	Content  string `json:"content,omitempty" description:"Updated content of the DNS record"`
	Priority int    `json:"priority,omitempty" description:"Updated priority of the DNS record"`
	TTL      int    `json:"ttl,omitempty" description:"Updated TTL of the DNS record"`
	_        bool   `name:"updateDnsDomainRecord" description:"Update a DNS domain record by ID and domain ID"`
}

// UpdateDNSRecordResponse represents the update of a DNS record
type UpdateDNSRecordResponse struct {
	Record UpdateDNSRecord `json:"record"`
}

// Response returns the struct to unmarshal
func (UpdateDNSRecord) Response() interface{} {
	return new(DNSRecord)
}

// DeleteDNSRecord represents a DNS domain record to be deleted
type DeleteDNSRecord struct {
	DomainID int64 `json:"id" doc:"ID of the DNS domain the record belongs to"`
	ID       int64 `json:"record_id" doc:"ID of the DNS domain record"`
	_        bool  `name:"deleteDnsDomainRecord" description:"Deletes a DNS domain by name"`
}

// Response returns the struct to unmarshal
func (DeleteDNSRecord) Response() interface{} {
	return new(BooleanResponse)
}

//go:generate go run generate/main.go -interface=Listable ListDNSRecords

// ListDNSDomains represents a query for a list of DNS domain records
type ListDNSRecords struct {
	DomainID   int64  `json:"id,omitempty" description:"List by domain ID"`
	DomainName string `json:"name,omitempty" description:"List by domain name"`
	Page       int    `json:"page,omitempty"`
	PageSize   int    `json:"pagesize,omitempty"`
	_          bool   `name:"listDnsDomainRecords" description:"List DNS domain records"`
}

// ListDNSRecordsResponse represents a list of DNS domain records
type ListDNSRecordsResponse struct {
	Count      int         `json:"count"`
	DNSRecords []DNSRecord `json:"records"`
}
