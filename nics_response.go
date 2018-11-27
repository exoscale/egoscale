// code generated; DO NOT EDIT.

package egoscale

import "fmt"

func (ListNics) response() interface{} {
	return new(ListNicsResponse)
}

// ListRequest returns itself
func (ls *ListNics) ListRequest() (ListCommand, error) {
	if ls == nil {
		return nil, fmt.Errorf("%T cannot be nil", ls)
	}
	return ls, nil
}

// SetPage sets the current apge
func (ls *ListNics) SetPage(page int) {
	ls.Page = page
}

// SetPageSize sets the page size
func (ls *ListNics) SetPageSize(pageSize int) {
	ls.PageSize = pageSize
}

// each triggers the callback for each, valid answer or any non 404 issue
func (ListNics) each(resp interface{}, callback IterateItemFunc) {
	items, ok := resp.(*ListNicsResponse)
	if !ok {
		callback(nil, fmt.Errorf("wrong type, ListNicsResponse was expected, got %T", resp))
		return
	}

	for i := range items.Nic {
		if !callback(&items.Nic[i], nil) {
			break
		}
	}
}
