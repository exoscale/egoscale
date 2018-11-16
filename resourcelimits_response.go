// Code generated; DO NOT EDIT.

package egoscale

import "fmt"

func (ListResourceLimits) response() interface{} {
	return new(ListResourceLimitsResponse)
}

// SetPage sets the current apge
func (ls *ListResourceLimits) SetPage(page int) {
	ls.Page = page
}

// SetPageSize sets the page size
func (ls *ListResourceLimits) SetPageSize(pageSize int) {
	ls.PageSize = pageSize
}

// each triggers the callback for each, valid answer or any non 404 issue
func (ListResourceLimits) each(resp interface{}, callback IterateItemFunc) {
	items, ok := resp.(*ListResourceLimitsResponse)
	if !ok {
		callback(nil, fmt.Errorf("wrong type, ListResourceLimitsResponse was expected, got %T", resp))
		return
	}

	for i := range items.ResourceLimit {
		if !callback(&items.ResourceLimit[i], nil) {
			break
		}
	}
}
