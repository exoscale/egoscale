// Code generated; DO NOT EDIT.

package egoscale

import "fmt"

func (ListSSHKeyPairs) response() interface{} {
	return new(ListSSHKeyPairsResponse)
}

// SetPage sets the current apge
func (ls *ListSSHKeyPairs) SetPage(page int) {
	ls.Page = page
}

// SetPageSize sets the page size
func (ls *ListSSHKeyPairs) SetPageSize(pageSize int) {
	ls.PageSize = pageSize
}

// each triggers the callback for each, valid answer or any non 404 issue
func (ListSSHKeyPairs) each(resp interface{}, callback IterateItemFunc) {
	items, ok := resp.(*ListSSHKeyPairsResponse)
	if !ok {
		callback(nil, fmt.Errorf("wrong type, ListSSHKeyPairsResponse was expected, got %T", resp))
		return
	}

	for i := range items.SSHKeyPair {
		if !callback(&items.SSHKeyPair[i], nil) {
			break
		}
	}
}
