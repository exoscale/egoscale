// Code generated; DO NOT EDIT.

package egoscale

import "fmt"

func (ListAccounts) response() interface{} {
	return new(ListAccountsResponse)
}

// SetPage sets the current apge
func (ls *ListAccounts) SetPage(page int) {
	ls.Page = page
}

// SetPageSize sets the page size
func (ls *ListAccounts) SetPageSize(pageSize int) {
	ls.PageSize = pageSize
}

// each triggers the callback for each, valid answer or any non 404 issue
func (ListAccounts) each(resp interface{}, callback IterateItemFunc) {
	items, ok := resp.(*ListAccountsResponse)
	if !ok {
		callback(nil, fmt.Errorf("wrong type, ListAccountsResponse was expected, got %T", resp))
		return
	}

	for i := range items.Account {
		if !callback(&items.Account[i], nil) {
			break
		}
	}
}
