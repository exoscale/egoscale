// code generated; DO NOT EDIT.

package egoscale

import "fmt"

func (ListUsers) response() interface{} {
	return new(ListUsersResponse)
}

// ListRequest returns itself
func (ls *ListUsers) ListRequest() (ListCommand, error) {
	if ls == nil {
		return nil, fmt.Errorf("%T cannot be nil", ls)
	}
	return ls, nil
}

// SetPage sets the current apge
func (ls *ListUsers) SetPage(page int) {
	ls.Page = page
}

// SetPageSize sets the page size
func (ls *ListUsers) SetPageSize(pageSize int) {
	ls.PageSize = pageSize
}

// each triggers the callback for each, valid answer or any non 404 issue
func (ListUsers) each(resp interface{}, callback IterateItemFunc) {
	items, ok := resp.(*ListUsersResponse)
	if !ok {
		callback(nil, fmt.Errorf("wrong type, ListUsersResponse was expected, got %T", resp))
		return
	}

	for i := range items.User {
		if !callback(&items.User[i], nil) {
			break
		}
	}
}
