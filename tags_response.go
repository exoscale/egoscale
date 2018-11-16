// Code generated; DO NOT EDIT.

package egoscale

import "fmt"

func (ListTags) response() interface{} {
	return new(ListTagsResponse)
}

// SetPage sets the current apge
func (ls *ListTags) SetPage(page int) {
	ls.Page = page
}

// SetPageSize sets the page size
func (ls *ListTags) SetPageSize(pageSize int) {
	ls.PageSize = pageSize
}

// each triggers the callback for each, valid answer or any non 404 issue
func (ListTags) each(resp interface{}, callback IterateItemFunc) {
	items, ok := resp.(*ListTagsResponse)
	if !ok {
		callback(nil, fmt.Errorf("wrong type, ListTagsResponse was expected, got %T", resp))
		return
	}

	for i := range items.Tag {
		if !callback(&items.Tag[i], nil) {
			break
		}
	}
}
