package egoscale

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// DNSDomain represents a domain
type DNSDomain struct {
	ID             int64  `json:"id"`
	AccountID      int64  `json:"account_id,omitempty"`
	UserID         int64  `json:"user_id,omitempty"`
	RegistrantID   int64  `json:"registrant_id,omitempty"`
	Name           string `json:"name"`
	UnicodeName    string `json:"unicode_name"`
	Token          string `json:"token"`
	State          string `json:"state"`
	Language       string `json:"language,omitempty"`
	Lockable       bool   `json:"lockable"`
	AutoRenew      bool   `json:"auto_renew"`
	WhoisProtected bool   `json:"whois_protected"`
	RecordCount    int64  `json:"record_count"`
	ServiceCount   int64  `json:"service_count"`
	ExpiresOn      string `json:"expires_on,omitempty"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// DNSDomainResponse represents a domain creation response
type DNSDomainResponse struct {
	Domain *DNSDomain `json:"domain"`
}

// DNSRecord represents a DNS record
type DNSRecord struct {
	ID         int64  `json:"id,omitempty"`
	DomainID   int64  `json:"domain_id,omitempty"`
	Name       string `json:"name"`
	TTL        int    `json:"ttl,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
	Content    string `json:"content"`
	RecordType string `json:"record_type"`
	Prio       int    `json:"prio,omitempty"`
}

// DNSRecordResponse represents the creation of a DNS record
type DNSRecordResponse struct {
	Record DNSRecord `json:"record"`
}

// UpdateDNSRecord represents a DNS record
type UpdateDNSRecord struct {
	ID         int64  `json:"id,omitempty"`
	DomainID   int64  `json:"domain_id,omitempty"`
	Name       string `json:"name,omitempty"`
	TTL        int    `json:"ttl,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
	Content    string `json:"content,omitempty"`
	RecordType string `json:"record_type,omitempty"`
	Prio       int    `json:"prio,omitempty"`
}

// UpdateDNSRecordResponse represents the creation of a DNS record
type UpdateDNSRecordResponse struct {
	Record UpdateDNSRecord `json:"record"`
}

// DNSErrorResponse represents an error in the API
type DNSErrorResponse struct {
	Message string              `json:"message,omitempty"`
	Errors  map[string][]string `json:"errors"`
}

// Record represent record type
type Record int

//go:generate stringer -type=Record
const (
	// A record type
	A Record = iota
	// AAAA record type
	AAAA
	// ALIAS record type
	ALIAS
	// CNAME record type
	CNAME
	// HINFO record type
	HINFO
	// MX record type
	MX
	// NAPTR record type
	NAPTR
	// NS record type
	NS
	// POOL record type
	POOL
	// SPF record type
	SPF
	// SRV record type
	SRV
	// SSHFP record type
	SSHFP
	// TXT record type
	TXT
	// URL record type
	URL
)

// Error formats the DNSerror into a string
func (req *DNSErrorResponse) Error() string {
	if len(req.Errors) > 0 {
		errs := []string{}
		for name, ss := range req.Errors {
			if len(ss) > 0 {
				errs = append(errs, fmt.Sprintf("%s: %s", name, strings.Join(ss, ", ")))
			}
		}
		return fmt.Sprintf("dns error: %s (%s)", req.Message, strings.Join(errs, "; "))
	}
	return fmt.Sprintf("dns error: %s", req.Message)
}

// CreateDomain creates a DNS domain
func (client *Client) CreateDomain(name string) (*DNSDomain, error) {
	m, err := json.Marshal(DNSDomainResponse{
		Domain: &DNSDomain{
			Name: name,
		},
	})
	if err != nil {
		return nil, err
	}

	resp, err := client.dnsRequest("/v1/domains", nil, string(m), "POST")
	if err != nil {
		return nil, err
	}

	var d *DNSDomainResponse
	if err := json.Unmarshal(resp, &d); err != nil {
		return nil, err
	}

	return d.Domain, nil
}

// GetDomain gets a DNS domain
func (client *Client) GetDomain(name string) (*DNSDomain, error) {
	resp, err := client.dnsRequest("/v1/domains/"+name, nil, "", "GET")
	if err != nil {
		return nil, err
	}

	var d *DNSDomainResponse
	if err := json.Unmarshal(resp, &d); err != nil {
		return nil, err
	}

	return d.Domain, nil
}

// GetDomains gets DNS domains
func (client *Client) GetDomains() ([]DNSDomain, error) {
	resp, err := client.dnsRequest("/v1/domains", nil, "", "GET")
	if err != nil {
		return nil, err
	}

	var d []DNSDomainResponse
	if err := json.Unmarshal(resp, &d); err != nil {
		return nil, err
	}

	domains := make([]DNSDomain, len(d))
	for i := range d {
		domains[i] = *d[i].Domain
	}
	return domains, nil
}

// DeleteDomain delets a DNS domain
func (client *Client) DeleteDomain(name string) error {
	_, err := client.dnsRequest("/v1/domains/"+name, nil, "", "DELETE")
	return err
}

// GetRecord returns a DNS record
func (client *Client) GetRecord(domain string, recordID int64) (*DNSRecord, error) {
	id := strconv.FormatInt(recordID, 10)
	resp, err := client.dnsRequest("/v1/domains/"+domain+"/records/"+id, nil, "", "GET")
	if err != nil {
		return nil, err
	}

	var r DNSRecordResponse
	if err = json.Unmarshal(resp, &r); err != nil {
		return nil, err
	}

	return &(r.Record), nil
}

// GetRecords returns the DNS records
func (client *Client) GetRecords(domain string) ([]DNSRecord, error) {
	resp, err := client.dnsRequest("/v1/domains/"+domain+"/records", nil, "", "GET")
	if err != nil {
		return nil, err
	}

	var r []DNSRecordResponse
	if err = json.Unmarshal(resp, &r); err != nil {
		return nil, err
	}

	records := make([]DNSRecord, 0, len(r))
	for _, rec := range r {
		records = append(records, rec.Record)
	}

	return records, nil
}

// GetRecordsWithFilters returns the DNS records (filters can be empty)
func (client *Client) GetRecordsWithFilters(domain, name, content, recordType string) ([]DNSRecord, error) {

	filters := map[string]string{}
	if name != "" {
		filters["name"] = name
	}
	if content != "" {
		filters["content"] = content
	}
	if recordType != "" {
		filters["record_type"] = recordType
	}

	resp, err := client.dnsRequest("/v1/domains/"+domain+"/records", filters, "", "GET")
	if err != nil {
		return nil, err
	}

	var r []DNSRecordResponse
	if err = json.Unmarshal(resp, &r); err != nil {
		return nil, err
	}

	records := make([]DNSRecord, 0, len(r))
	for _, rec := range r {
		records = append(records, rec.Record)
	}

	return records, nil
}

// CreateRecord creates a DNS record
func (client *Client) CreateRecord(name string, rec DNSRecord) (*DNSRecord, error) {
	body, err := json.Marshal(DNSRecordResponse{
		Record: rec,
	})
	if err != nil {
		return nil, err
	}

	resp, err := client.dnsRequest("/v1/domains/"+name+"/records", nil, string(body), "POST")
	if err != nil {
		return nil, err
	}

	var r DNSRecordResponse
	if err = json.Unmarshal(resp, &r); err != nil {
		return nil, err
	}

	return &(r.Record), nil
}

// UpdateRecord updates a DNS record
func (client *Client) UpdateRecord(name string, rec UpdateDNSRecord) (*DNSRecord, error) {
	body, err := json.Marshal(UpdateDNSRecordResponse{
		Record: rec,
	})
	if err != nil {
		return nil, err
	}

	id := strconv.FormatInt(rec.ID, 10)
	resp, err := client.dnsRequest("/v1/domains/"+name+"/records/"+id, nil, string(body), "PUT")
	if err != nil {
		return nil, err
	}

	var r DNSRecordResponse
	if err = json.Unmarshal(resp, &r); err != nil {
		return nil, err
	}

	return &(r.Record), nil
}

// DeleteRecord deletes a record
func (client *Client) DeleteRecord(name string, recordID int64) error {
	id := strconv.FormatInt(recordID, 10)
	_, err := client.dnsRequest("/v1/domains/"+name+"/records/"+id, nil, "", "DELETE")

	return err
}

func (client *Client) dnsRequest(uri string, uriParams map[string]string, params, method string) (json.RawMessage, error) {
	url := client.Endpoint + uri
	req, err := http.NewRequest(method, url, strings.NewReader(params))
	if err != nil {
		return nil, err
	}

	var hdr = make(http.Header)
	hdr.Add("X-DNS-TOKEN", client.APIKey+":"+client.apiSecret)
	hdr.Add("User-Agent", fmt.Sprintf("exoscale/egoscale (%v)", Version))
	hdr.Add("Accept", "application/json")
	if params != "" {
		hdr.Add("Content-Type", "application/json")
	}
	req.Header = hdr

	q := req.URL.Query()

	for k, v := range uriParams {
		q.Add(k, v)
	}

	req.URL.RawQuery = q.Encode()

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // nolint: errcheck

	contentType := resp.Header.Get("content-type")
	if !strings.Contains(contentType, "application/json") {
		return nil, fmt.Errorf(`response content-type expected to be "application/json", got %q`, contentType)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		e := new(DNSErrorResponse)
		if err := json.Unmarshal(b, e); err != nil {
			return nil, err
		}
		return nil, e
	}

	return b, nil
}

// GetRecordIDByName get record ID by name
func (client *Client) GetRecordIDByName(domainName, recordName string) (int64, error) {
	records, err := client.GetRecords(domainName)
	if err != nil {
		return 0, err
	}

	resRecID := []int64{}

	for _, r := range records {
		id := fmt.Sprintf("%d", r.ID)
		if id == recordName {
			return r.ID, nil
		}
		if recordName == r.Name {
			resRecID = append(resRecID, r.ID)
		}
	}
	if len(resRecID) > 1 {
		return 0, fmt.Errorf("more than one records were found")
	}
	if len(resRecID) == 1 {
		return resRecID[0], nil
	}

	return 0, fmt.Errorf("no records were found")
}
