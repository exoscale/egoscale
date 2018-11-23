package admin

import "github.com/exoscale/egoscale"

type VirtualMachine struct {
	egoscale.VirtualMachine
	Account   string         `json:"account,omitempty"`
	AccountID *egoscale.UUID `json:"accountid,omitempty"`
	HostName  string         `json:"string,omitempty"`
}

type ListVirtualMachines struct {
	egoscale.ListVirtualMachines
	ListAll *bool `json:"listall,omitempty"`
}

type LisVirtualMachinesResponse struct {
	Count          int              `json:"count"`
	VirtualMachine []VirtualMachine `json:"virtualmachine"`
}
