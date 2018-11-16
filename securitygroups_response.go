// Code generated; DO NOT EDIT.

package egoscale

import "fmt"

func (ListSecurityGroups) response() interface{} {
	return new(ListSecurityGroupsResponse)
}

// SetPage sets the current apge
func (ls *ListSecurityGroups) SetPage(page int) {
	ls.Page = page
}

// SetPageSize sets the page size
func (ls *ListSecurityGroups) SetPageSize(pageSize int) {
	ls.PageSize = pageSize
}

// each triggers the callback for each, valid answer or any non 404 issue
func (ListSecurityGroups) each(resp interface{}, callback IterateItemFunc) {
	items, ok := resp.(*ListSecurityGroupsResponse)
	if !ok {
		callback(nil, fmt.Errorf("wrong type, ListSecurityGroupsResponse was expected, got %T", resp))
		return
	}

	for i := range items.SecurityGroup {
		if !callback(&items.SecurityGroup[i], nil) {
			break
		}
	}
}
