// code generated; DO NOT EDIT.

package egoscale

import "fmt"

// Response returns the struct to unmarshal
func (ListDNSDomains) Response() interface{} {
	return new(ListDNSDomainsResponse)
}

// ListRequest returns itself
func (ls *ListDNSDomains) ListRequest() (ListCommand, error) {
	if ls == nil {
		return nil, fmt.Errorf("%T cannot be nil", ls)
	}
	return ls, nil
}

// SetPage sets the current apge
func (ls *ListDNSDomains) SetPage(page int) {
	ls.Page = page
}

// SetPageSize sets the page size
func (ls *ListDNSDomains) SetPageSize(pageSize int) {
	ls.PageSize = pageSize
}

// Each triggers the callback for each, valid answer or any non 404 issue
func (ListDNSDomains) Each(resp interface{}, callback IterateItemFunc) {
	items, ok := resp.(*ListDNSDomainsResponse)
	if !ok {
		callback(nil, fmt.Errorf("wrong type, ListDNSDomainsResponse was expected, got %T", resp))
		return
	}

	for i := range items.DNSDomains {
		if !callback(&items.DNSDomains[i], nil) {
			break
		}
	}
}
