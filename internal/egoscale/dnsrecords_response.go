// code generated; DO NOT EDIT.

package egoscale

import "fmt"

// Response returns the struct to unmarshal
func (ListDNSRecords) Response() interface{} {
	return new(ListDNSRecordsResponse)
}

// ListRequest returns itself
func (ls *ListDNSRecords) ListRequest() (ListCommand, error) {
	if ls == nil {
		return nil, fmt.Errorf("%T cannot be nil", ls)
	}
	return ls, nil
}

// SetPage sets the current apge
func (ls *ListDNSRecords) SetPage(page int) {
	ls.Page = page
}

// SetPageSize sets the page size
func (ls *ListDNSRecords) SetPageSize(pageSize int) {
	ls.PageSize = pageSize
}

// Each triggers the callback for each, valid answer or any non 404 issue
func (ListDNSRecords) Each(resp interface{}, callback IterateItemFunc) {
	items, ok := resp.(*ListDNSRecordsResponse)
	if !ok {
		callback(nil, fmt.Errorf("wrong type, ListDNSRecordsResponse was expected, got %T", resp))
		return
	}

	for i := range items.DNSRecords {
		if !callback(&items.DNSRecords[i], nil) {
			break
		}
	}
}
