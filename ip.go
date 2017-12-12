package egoscale

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// AssociateIpAddress acquires and associates a public IP to a given zone.
func (exo *Client) AssociateIpAddress(zoneId string) (*IpAddress, error) {
	params := url.Values{}
	params.Set("zoneid", zoneId)
	resp, err := exo.Request("associateIpAddress", params)
	if err != nil {
		return nil, err
	}

	var jr JobResponse
	if err := json.Unmarshal(resp, &jr); err != nil {
		return nil, err
	}

	if jr.JobResultType == "object" {
		var r AssociateIpAddressResponse
		if err := json.Unmarshal(*jr.JobResult, &r); err != nil {
			return nil, err
		}

		return &r.IpAddress, nil
	}

	return nil, fmt.Errorf("Expected an object as a result, got %s", jr.JobResultType)
}

// DisassociateIpAddress disassociates a public IP from the account
func (exo *Client) DisassociateIpAddress(ipAddressId string) (bool, error) {
	params := url.Values{}
	params.Set("id", ipAddressId)
	resp, err := exo.Request("ipAddressId", params)
	if err != nil {
		return false, err
	}

	var jr JobResponse
	if err := json.Unmarshal(resp, &jr); err != nil {
		return false, err
	}

	if jr.JobResultType == "object" {
		var r DisassociateIpAddressResponse
		if err := json.Unmarshal(*jr.JobResult, &r); err != nil {
			return false, err
		}

		return r.Success, nil
	}

	return false, fmt.Errorf("Expected an object as a result, got %s", jr.JobResultType)
}

func (exo *Client) AddIpToNic(nic_id string, ip_address string) (string, error) {
	params := url.Values{}
	params.Set("nicid", nic_id)
	params.Set("ipaddress", ip_address)

	resp, err := exo.Request("addIpToNic", params)
	if err != nil {
		return "", err
	}

	var r AddIpToNicResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return "", err
	}

	return r.Id, nil
}

func (exo *Client) RemoveIpFromNic(nic_id string) (string, error) {
	params := url.Values{}
	params.Set("id", nic_id)

	resp, err := exo.Request("removeIpFromNic", params)
	if err != nil {
		return "", err
	}

	var r RemoveIpFromNicResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return "", err
	}
	return r.JobID, nil
}
