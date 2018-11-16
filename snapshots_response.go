// Code generated; DO NOT EDIT.

package egoscale

import "fmt"

func (ListSnapshots) response() interface{} {
	return new(ListSnapshotsResponse)
}

// SetPage sets the current apge
func (ls *ListSnapshots) SetPage(page int) {
	ls.Page = page
}

// SetPageSize sets the page size
func (ls *ListSnapshots) SetPageSize(pageSize int) {
	ls.PageSize = pageSize
}

// each triggers the callback for each, valid answer or any non 404 issue
func (ListSnapshots) each(resp interface{}, callback IterateItemFunc) {
	items, ok := resp.(*ListSnapshotsResponse)
	if !ok {
		callback(nil, fmt.Errorf("wrong type, ListSnapshotsResponse was expected, got %T", resp))
		return
	}

	for i := range items.Snapshot {
		if !callback(&items.Snapshot[i], nil) {
			break
		}
	}
}
