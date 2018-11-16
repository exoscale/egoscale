// Code generated; DO NOT EDIT.

package egoscale

import "fmt"

func (ListTemplates) response() interface{} {
	return new(ListTemplatesResponse)
}

// SetPage sets the current apge
func (ls *ListTemplates) SetPage(page int) {
	ls.Page = page
}

// SetPageSize sets the page size
func (ls *ListTemplates) SetPageSize(pageSize int) {
	ls.PageSize = pageSize
}

// each triggers the callback for each, valid answer or any non 404 issue
func (ListTemplates) each(resp interface{}, callback IterateItemFunc) {
	items, ok := resp.(*ListTemplatesResponse)
	if !ok {
		callback(nil, fmt.Errorf("wrong type, ListTemplatesResponse was expected, got %T", resp))
		return
	}

	for i := range items.Template {
		if !callback(&items.Template[i], nil) {
			break
		}
	}
}
