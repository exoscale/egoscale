package egoscale

import (
	"fmt"
)

// Volume represents a volume linked to a VM
type Volume struct {
	ID                         string        `json:"id"`
	Account                    string        `json:"account,omitempty"`
	AttachedAt                 string        `json:"attached,omitempty"`
	ChainInfo                  string        `json:"chaininfo,omitempty"`
	CreatedAt                  string        `json:"created,omitempty"`
	Destroyed                  bool          `json:"destroyed,omitempty"`
	DisplayVolume              bool          `json:"displayvolume,omitempty"`
	Domain                     string        `json:"domain,omitempty"`
	DomainID                   string        `json:"domainid,omitempty"`
	Name                       string        `json:"name,omitempty"`
	QuiesceVM                  bool          `json:"quiescevm,omitempty"`
	ServiceOfferingDisplayText string        `json:"serviceofferingdisplaytext,omitempty"`
	ServiceOfferingID          string        `json:"serviceofferingid,omitempty"`
	ServiceOfferingName        string        `json:"serviceofferingname,omitempty"`
	Size                       uint64        `json:"size,omitempty"`
	State                      string        `json:"state,omitempty"`
	Type                       string        `json:"type,omitempty"`
	VirtualMachineID           string        `json:"virtualmachineid,omitempty"`
	VMName                     string        `json:"vmname,omitempty"`
	VMState                    string        `json:"vmstate,omitempty"`
	ZoneID                     string        `json:"zoneid,omitempty"`
	ZoneName                   string        `json:"zonename,omitempty"`
	Tags                       []*VolumeTag  `json:"tags,omitempty"`
	JobID                      string        `json:"jobid,omitempty"`
	JobStatus                  JobStatusType `json:"jobstatus,omitempty"`
}

// VolumeTag represents a tag associated with a Volume
type VolumeTag struct {
	Account    string `json:"account,omitempty"`
	Customer   string `json:"customer,omitempty"`
	Domain     string `json:"domain,omitempty"`
	DomainID   string `json:"domainid,omitempty"`
	Key        string `json:"key,omitempty"`
	Project    string `json:"project,omitempty"`
	ProjectID  string `json:"projectid,omitempty"`
	Resource   string `json:"resource,omitempty"`
	ResourceID string `json:"resourceid,omitempty"`
	Value      string `json:"value,omitempty"`
}

// ListVolumesRequest represents a query listing volumes
type ListVolumesRequest struct {
	Account          string         `json:"account,omitempty"`
	DiskOfferingID   string         `json:"diskoffering,omitempty"`
	DisplayVolume    string         `json:"displayvolume,omitempty"` // root only
	DomainID         string         `json:"domainid,omitempty"`
	HostID           string         `json:"hostid,omitempty"`
	ID               string         `json:"id,omitempty"`
	IsRecursive      bool           `json:"isrecursive,omitempty"`
	Keyword          string         `json:"keyword,omitempty"`
	ListAll          bool           `json:"listall,omitempty"`
	Name             string         `json:"name,omitempty"`
	Page             int            `json:"page,omitempty"`
	PageSize         int            `json:"pagesize,omitempty"`
	PodID            string         `json:"podid,omitempty"`
	ProjectID        string         `json:"projectid,omitempty"`
	StorageID        string         `json:"storageid,omitempty"`
	Tags             []*ResourceTag `json:"tags,omitempty"`
	Type             string         `json:"type,omitempty"`
	VirtualMachineID string         `json:"virtualmachineid,omitempty"`
	ZoneID           string         `json:"zoneid,omitempty"`
}

// Command returns the CloudStack API command
func (req *ListVolumesRequest) Command() string {
	return "listVolumes"
}

// ListVolumesResponse represents a list of volumes
type ListVolumesResponse struct {
	Count  int       `json:"count"`
	Volume []*Volume `json:"volume"`
}

// GetRootVolumeForVirtualMachine returns the root volume of a VM
//
// Deprecated: helper function shouldn't be used
func (exo *Client) GetRootVolumeForVirtualMachine(virtualMachineID string) (*Volume, error) {
	r := new(ListVolumesResponse)
	err := exo.Request(&ListVolumesRequest{
		VirtualMachineID: virtualMachineID,
		Type:             "ROOT",
	}, r)
	if err != nil {
		return nil, err
	}

	if r.Count != 1 {
		return nil, fmt.Errorf("Expected exactly one volume for %v, got %d", virtualMachineID, r.Count)
	}

	return r.Volume[0], nil
}
