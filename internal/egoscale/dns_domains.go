package egoscale

import (
	"context"
	"fmt"
)

// Delete deletes the given DNS domain by Name
func (domain DNSDomain) Delete(ctx context.Context, client *Client) error {
	if domain.Name == "" {
		return fmt.Errorf("a DNS domain may only be deleted using Name")
	}

	return client.BooleanRequestWithContext(ctx, &DeleteDNSDomain{
		Name: domain.Name,
	})
}

// ListRequest builds the ListDNSDomains request
func (domain DNSDomain) ListRequest() (ListCommand, error) {
	req := &ListDNSDomains{}

	return req, nil
}

// CreateDNSDomain represents a new DNS domain to be created
type CreateDNSDomain struct {
	Name string `json:"name" doc:"Name of the DNS domain"`
	_    bool   `name:"createDnsDomain" description:"Create a new DNS domain"`
}

// Response returns the struct to unmarshal
func (CreateDNSDomain) Response() interface{} {
	return new(DNSDomain)
}

// DeleteDNSDomain represents a DNS domain to be deleted
type DeleteDNSDomain struct {
	Name string `json:"name" doc:"Name of the DNS domain"`
	_    bool   `name:"deleteDnsDomain" description:"Deletes a DNS domain by name"`
}

// Response returns the struct to unmarshal
func (DeleteDNSDomain) Response() interface{} {
	return new(BooleanResponse)
}

//go:generate go run generate/main.go -interface=Listable ListDNSDomains

// ListDNSDomains represents a query for a list of DNS domains
type ListDNSDomains struct {
	Page     int  `json:"page,omitempty"`
	PageSize int  `json:"pagesize,omitempty"`
	_        bool `name:"listDnsDomains" description:"List DNS domains"`
}

// ListDNSDomainsResponse represents a list of DNS domains
type ListDNSDomainsResponse struct {
	Count      int         `json:"count"`
	DNSDomains []DNSDomain `json:"dnsdomain"`
}
