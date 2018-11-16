// Code generated; DO NOT EDIT.

package egoscale

import "fmt"

func (ListAffinityGroups) response() interface{} {
	return new(ListAffinityGroupsResponse)
}

// SetPage sets the current apge
func (ls *ListAffinityGroups) SetPage(page int) {
	ls.Page = page
}

// SetPageSize sets the page size
func (ls *ListAffinityGroups) SetPageSize(pageSize int) {
	ls.PageSize = pageSize
}

// each triggers the callback for each, valid answer or any non 404 issue
func (ListAffinityGroups) each(resp interface{}, callback IterateItemFunc) {
	items, ok := resp.(*ListAffinityGroupsResponse)
	if !ok {
		callback(nil, fmt.Errorf("wrong type, ListAffinityGroupsResponse was expected, got %T", resp))
		return
	}

	for i := range items.AffinityGroup {
		if !callback(&items.AffinityGroup[i], nil) {
			break
		}
	}
}
