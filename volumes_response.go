// Code generated; DO NOT EDIT.

package egoscale

import "fmt"

func (ListVolumes) response() interface{} {
	return new(ListVolumesResponse)
}

// SetPage sets the current apge
func (ls *ListVolumes) SetPage(page int) {
	ls.Page = page
}

// SetPageSize sets the page size
func (ls *ListVolumes) SetPageSize(pageSize int) {
	ls.PageSize = pageSize
}

// each triggers the callback for each, valid answer or any non 404 issue
func (ListVolumes) each(resp interface{}, callback IterateItemFunc) {
	items, ok := resp.(*ListVolumesResponse)
	if !ok {
		callback(nil, fmt.Errorf("wrong type, ListVolumesResponse was expected, got %T", resp))
		return
	}

	for i := range items.Volume {
		if !callback(&items.Volume[i], nil) {
			break
		}
	}
}
