// code generated; DO NOT EDIT.

package egoscale

import "fmt"

// Response returns the struct to unmarshal
func (ListIsos) Response() interface{} {
	return new(ListIsosResponse)
}

// ListRequest returns itself
func (ls *ListIsos) ListRequest() (ListCommand, error) {
	if ls == nil {
		return nil, fmt.Errorf("%T cannot be nil", ls)
	}
	return ls, nil
}

// SetPage sets the current apge
func (ls *ListIsos) SetPage(page int) {
	ls.Page = page
}

// SetPageSize sets the page size
func (ls *ListIsos) SetPageSize(pageSize int) {
	ls.PageSize = pageSize
}

// Each triggers the callback for each, valid answer or any non 404 issue
func (ListIsos) Each(resp interface{}, callback IterateItemFunc) {
	items, ok := resp.(*ListIsosResponse)
	if !ok {
		callback(nil, fmt.Errorf("wrong type, ListIsosResponse was expected, got %T", resp))
		return
	}

	for i := range items.Iso {
		if !callback(&items.Iso[i], nil) {
			break
		}
	}
}
