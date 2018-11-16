// Code generated; DO NOT EDIT.

package egoscale

import "fmt"

func (ListServiceOfferings) response() interface{} {
	return new(ListServiceOfferingsResponse)
}

// SetPage sets the current apge
func (ls *ListServiceOfferings) SetPage(page int) {
	ls.Page = page
}

// SetPageSize sets the page size
func (ls *ListServiceOfferings) SetPageSize(pageSize int) {
	ls.PageSize = pageSize
}

// each triggers the callback for each, valid answer or any non 404 issue
func (ListServiceOfferings) each(resp interface{}, callback IterateItemFunc) {
	items, ok := resp.(*ListServiceOfferingsResponse)
	if !ok {
		callback(nil, fmt.Errorf("wrong type, ListServiceOfferingsResponse was expected, got %T", resp))
		return
	}

	for i := range items.ServiceOffering {
		if !callback(&items.ServiceOffering[i], nil) {
			break
		}
	}
}
