// affinitygroups contains the methods related to (anti-)affinity groups
package egoscale

import (
	"encoding/json"
	"fmt"
)

// AffinityGroup represents an (anti-)affinity group
type AffinityGroup struct {
	Id                string   `json:"id,omitempty"`
	Account           string   `json:"account,omitempty"`
	Description       string   `json:"description,omitempty"`
	Domain            string   `json:"domain,omitempty"`
	DomainId          string   `json:"domainid,omitempty"`
	Name              string   `json:"name,omitempty"`
	Type              string   `json:"type,omitempty"`
	VirtualMachineIds []string `json:"virtualmachineIds,omitempty"` // *I*ds is not a typo
}

// CreateAffinityGroupRequest represents a new (anti-)affinity group
type CreateAffinityGroupRequest struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Account     string `json:"account,omitempty"`
	Description string `json:"description,omitempty"`
	DomainId    string `json:"domainid,omitempty"`
}

// DeleteAffinityGroupRequest represents an (anti-)affinity group to be deleted
type DeleteAffinityGroupRequest struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Type        string `json:"type,omitempty"`
	Account     string `json:"account,omitempty"`
	Description string `json:"description,omitempty"`
	DomainId    string `json:"domainid,omitempty"`
}

// ListAffinityGroupsRequest represents an (anti-)affinity groups search
type ListAffinityGroupsRequest struct {
	Account          string `json:"account,omitempty"`
	DomainId         string `json:"domainid,omitempty"`
	Id               string `json:"id,omitempty"`
	IsRecursive      bool   `json:"isrecursive,omitempty"`
	Keyword          string `json:"keyword,omitempty"`
	ListAll          bool   `json:"listall,omitempty"`
	Name             string `json:"name,omitempty"`
	Page             string `json:"page,omitempty"`
	PageSize         string `json:"pagesize,omitempty"`
	Type             string `json:"type,omitempty"`
	VirtualMachineId string `json:"virtualmachineid,omitempty"`
}

// CreateAffinityGroupResponse represents the response of the creation of an (anti-)affinity group
type CreateAffinityGroupResponse struct {
	AffinityGroup AffinityGroup `json:"affinitygroup"`
}

// ListAffinityGroupTypesRequest represents an (anti-)affinity groups search
type ListAffinityGroupTypesRequest struct {
	Keyword  string `json:"keyword,omitempty"`
	Page     string `json:"page,omitempty"`
	PageSize string `json:"pagesize,omitempty"`
}

// ListAffinityGroupsResponse represents a list of (anti-)affinity groups
type ListAffinityGroupsResponse struct {
	Count         int              `json:"count"`
	AffinityGroup []*AffinityGroup `json:"affinitygroup"`
}

// ListAffinityGroupTypesResponse represents a list of (anti-)affinity group types
type ListAffinityGroupTypesResponse struct {
	Count             int                  `json:"count"`
	AffinityGroupType []*AffinityGroupType `json:"affinitygrouptype"`
}

// CreateAffinityGroup creates an (anti-)affinity group
func (exo *Client) CreateAffinityGroup(req CreateAffinityGroupRequest, async AsyncInfo) (*AffinityGroup, error) {
	params, err := prepareValues(req)
	if err != nil {
		return nil, err
	}

	resp, err := exo.AsyncRequest("createAffinityGroup", *params, async)
	if err != nil {
		return nil, err
	}

	var r CreateAffinityGroupResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return nil, err
	}

	return &r.AffinityGroup, nil
}

// DeleteAffinityGroup deletes an affinity group by name
func (exo *Client) DeleteAffinityGroup(req DeleteAffinityGroupRequest, async AsyncInfo) error {
	params, err := prepareValues(req)
	if err != nil {
		return err
	}

	resp, err := exo.AsyncRequest("deleteAffinityGroup", *params, async)
	if err != nil {
		return err
	}

	var r BooleanResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return err
	}

	if !r.Success {
		return fmt.Errorf("Cannot delete affinity group: %s", r.DisplayText)
	}

	return nil
}

// ListAffinityGroups lists the affinity groups
func (exo *Client) ListAffinityGroups(req ListAffinityGroupsRequest) ([]*AffinityGroup, error) {
	params, err := prepareValues(req)
	if err != nil {
		return nil, err
	}

	resp, err := exo.Request("listAffinityGroups", *params)
	if err != nil {
		return nil, err
	}

	var r ListAffinityGroupsResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return nil, err
	}

	return r.AffinityGroup, nil
}

// ListAffinityGroupTypes lists the affinity group type
func (exo *Client) ListAffinityGroupTypes(req ListAffinityGroupTypesRequest) ([]*AffinityGroupType, error) {
	params, err := prepareValues(req)
	if err != nil {
		return nil, err
	}

	resp, err := exo.Request("listAffinityGroupTypes", *params)
	if err != nil {
		return nil, err
	}

	var r ListAffinityGroupTypesResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return nil, err
	}

	return r.AffinityGroupType, nil
}
