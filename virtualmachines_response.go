// Code generated; DO NOT EDIT.

package egoscale

import "fmt"

func (ListVirtualMachines) response() interface{} {
	return new(ListVirtualMachinesResponse)
}

// SetPage sets the current apge
func (ls *ListVirtualMachines) SetPage(page int) {
	ls.Page = page
}

// SetPageSize sets the page size
func (ls *ListVirtualMachines) SetPageSize(pageSize int) {
	ls.PageSize = pageSize
}

// each triggers the callback for each, valid answer or any non 404 issue
func (ListVirtualMachines) each(resp interface{}, callback IterateItemFunc) {
	items, ok := resp.(*ListVirtualMachinesResponse)
	if !ok {
		callback(nil, fmt.Errorf("wrong type, ListVirtualMachinesResponse was expected, got %T", resp))
		return
	}

	for i := range items.VirtualMachine {
		if !callback(&items.VirtualMachine[i], nil) {
			break
		}
	}
}
