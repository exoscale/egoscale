// Code generated; DO NOT EDIT.

package egoscale

import "fmt"

func (ListEvents) response() interface{} {
	return new(ListEventsResponse)
}

// SetPage sets the current apge
func (ls *ListEvents) SetPage(page int) {
	ls.Page = page
}

// SetPageSize sets the page size
func (ls *ListEvents) SetPageSize(pageSize int) {
	ls.PageSize = pageSize
}

// each triggers the callback for each, valid answer or any non 404 issue
func (ListEvents) each(resp interface{}, callback IterateItemFunc) {
	items, ok := resp.(*ListEventsResponse)
	if !ok {
		callback(nil, fmt.Errorf("wrong type, ListEventsResponse was expected, got %T", resp))
		return
	}

	for i := range items.Event {
		if !callback(&items.Event[i], nil) {
			break
		}
	}
}
