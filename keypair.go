/*
SSH Key Pair

In addition to username and password (disabled on Exoscale), SSH keys are used to log into the infrastructure.

See: http://docs.cloudstack.apache.org/projects/cloudstack-administration/en/stable/virtual_machines.html#creating-the-ssh-keypair
*/
package egoscale

import (
	"encoding/json"
)

// SshKeyPair represents an SSH key pair
type SshKeyPair struct {
	Account     string `json:"account,omitempty"`
	DomainId    string `json:"domainid,omitempty"`
	ProjectId   string `json:"projectid,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
	Name        string `json:"name,omitempty"`
	PrivateKey  string `json:"privatekey,omitempty"`
}

// CreateSshKeyPairRequest represents a new keypair to be created
type CreateSshKeyPairRequest struct {
	Name      string `json:"name"`
	Account   string `json:"account,omitempty"`
	DomainId  string `json:"domainid,omitempty"`
	ProjectId string `json:"projectid,omitempty"`
}

// CreateSshKeyPairResponse represents the creation of an SSH Key Pair
type CreateSshKeyPairResponse struct {
	KeyPair SshKeyPair `json:"keypair"`
}

// DeleteSshKeyPairRequest represents a new keypair to be created
type DeleteSshKeyPairRequest struct {
	Name      string `json:"name"`
	Account   string `json:"account,omitempty"`
	DomainId  string `json:"domainid,omitempty"`
	ProjectId string `json:"projectid,omitempty"`
}

// SshKeyPairRequest represents a new registration of a public key in a keypair
type RegisterSshKeyPairRequest struct {
	Name      string `json:"name"`
	PublicKey string `json:"publickey"`
	Account   string `json:"account,omitempty"`
	DomainId  string `json:"domainid,omitempty"`
	ProjectId string `json:"projectid,omitempty"`
}

// RegisterSshKeyPairResponse represents the creation of an SSH Key Pair
type RegisterSshKeyPairResponse struct {
	KeyPair SshKeyPair `json:"keypair"`
}

// ListSshKeyPairsRequest represents a query for a list of SSH KeyPairs
type ListSshKeyPairsRequest struct {
	Account     string `json:"account,omitempty"`
	DomainId    string `json:"domainid,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
	IsRecursive bool   `json:"isrecursive,omitempty"`
	Keyword     string `json:"keyword,omitempty"`
	ListAll     bool   `json:"listall,omitempty"`
	Name        string `json:"name,omitempty"`
	Page        string `json:"page,omitempty"`
	PageSize    string `json:"pagesize,omitempty"`
	ProjectId   string `json:"projectid,omitempty"`
}

// ListSshKeyPairsResponse
type ListSshKeyPairsResponse struct {
	Count      int           `json:"count"`
	SshKeyPair []*SshKeyPair `json:"sshkeypair"`
}

// CreateSshKeyPair create a new SSH Key Pair
//
// http://cloudstack.apache.org/api/apidocs-4.10/apis/createSSHKeyPair.html
func (exo *Client) CreateSshKeyPair(req CreateSshKeyPairRequest) (*SshKeyPair, error) {
	params, err := prepareValues(req)
	if err != nil {
		return nil, err
	}
	resp, err := exo.Request("createSSHKeyPair", *params)
	if err != nil {
		return nil, err
	}

	var r CreateSshKeyPairResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return nil, err
	}

	return &r.KeyPair, nil
}

// DeleteSshKeyPair deletes an SSH key pair
//
// http://cloudstack.apache.org/api/apidocs-4.10/apis/deleteSSHKeyPair.html
func (exo *Client) DeleteSshKeyPair(req DeleteSshKeyPairRequest) error {
	params, err := prepareValues(req)
	if err != nil {
		return err
	}
	return exo.BooleanRequest("deleteSSHKeyPair", *params)
}

// RegisterSshKeyPair registers a public key in a keypair
//
// http://cloudstack.apache.org/api/apidocs-4.10/apis/registerSSHKeyPair.html
func (exo *Client) RegisterSshKeyPair(req RegisterSshKeyPairRequest) (*SshKeyPair, error) {
	params, err := prepareValues(req)
	if err != nil {
		return nil, nil
	}
	resp, err := exo.Request("registerSSHKeyPair", *params)
	if err != nil {
		return nil, err
	}

	var r RegisterSshKeyPairResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return nil, err
	}

	return &r.KeyPair, nil
}

// ListSshKeyPairs lists the key pairs
//
// http://cloudstack.apache.org/api/apidocs-4.10/apis/listSSHKeyPairs.html
func (exo *Client) ListSshKeyPairs(req ListSshKeyPairsRequest) ([]*SshKeyPair, error) {
	params, err := prepareValues(req)
	if err != nil {
		return nil, err
	}

	resp, err := exo.Request("listSSHKeyPairs", *params)
	if err != nil {
		return nil, err
	}

	var r ListSshKeyPairsResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return nil, err
	}

	return r.SshKeyPair, nil
}

// XXX ResetSshKeyForVirtualMachine
//
// http://cloudstack.apache.org/api/apidocs-4.10/apis/resetSSHKeyForVirtualMachine.html
