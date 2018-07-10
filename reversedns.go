package egoscale

import (
	"net"
)

// ReverseDNS represents the PTR record linked with an IPAddress or IP6Address belonging to a Virtual Machine or a Public IP Address (Elastic IP) instance
type ReverseDNS struct {
	DomainName       string `json:"domainname,omitempty" doc:"the domain name of the PTR record"`
	IP6Address       net.IP `json:"ip6address,omitempty" doc:"the IPv6 address linked with the PTR record (mutually exclusive with ipaddress)"`
	IPAddress        net.IP `json:"ipaddress,omitempty" doc:"the IPv4 address linked with the PTR record (mutually exclusive with ip6address)"`
	NicID            string `json:"nicid,omitempty" doc:"the virtual machine default NIC ID"`
	PublicIPID       string `json:"publicipid,omitempty" doc:"the public IP address ID"`
	VirtualMachineID string `json:"virtualmachineid,omitempty" doc:"the virtual machine ID"`
}

// UpdateReverseDNSForVirtualMachine is a command to create/update the PTR record(s) of a virtual machine
type UpdateReverseDNSForVirtualMachine struct {
	DomainName string `json:"domainname,omitempty" doc:"the domain name for the PTR record(s). It must have a valid TLD"`
	ID         string `json:"id,omitempty" doc:"the ID of the virtual machine"`
	_          bool   `name:"updateReverseDnsForVirtualMachine" description:"Update/create the PTR DNS record(s) for the Virtual Machine"`
}

func (*UpdateReverseDNSForVirtualMachine) response() interface{} {
	return new(VirtualMachine)
}

// UpdateReverseDNSForPublicIPAddress is a command to create/update the PTR record of a public IP address
type UpdateReverseDNSForPublicIPAddress struct {
	DomainName string `json:"domainname,omitempty" doc:"the domain name for the PTR record. It must have a valid TLD"`
	ID         string `json:"id,omitempty" doc:"the ID of the Public IP Address"`
	_          bool   `name:"updateReverseDnsForPublicIpAddress" description:"Update/create the PTR DNS record for the Public IP Address"`
}

func (*UpdateReverseDNSForPublicIPAddress) response() interface{} {
	return new(IPAddress)
}
