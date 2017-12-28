package egoscale

/*
SSH Key Pairs

In addition to username and password (disabled on Exoscale), SSH keys are used to log into the infrastructure.

See: http://docs.cloudstack.apache.org/projects/cloudstack-administration/en/stable/virtual_machines.html#creating-the-ssh-keypair
*/

// SSHKeyPair represents an SSH key pair
type SSHKeyPair struct {
	Account     string `json:"account,omitempty"`
	DomainID    string `json:"domainid,omitempty"`
	ProjectID   string `json:"projectid,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
	Name        string `json:"name,omitempty"`
	PrivateKey  string `json:"privatekey,omitempty"`
}

// CreateSSHKeyPairRequest represents a new keypair to be created
//
// http://cloudstack.apache.org/api/apidocs-4.10/apis/createSSHKeyPair.html
type CreateSSHKeyPairRequest struct {
	Name      string `json:"name"`
	Account   string `json:"account,omitempty"`
	DomainID  string `json:"domainid,omitempty"`
	ProjectID string `json:"projectid,omitempty"`
}

// Command returns the CloudStack API command
func (req *CreateSSHKeyPairRequest) Command() string {
	return "createSSHKeyPair"
}

// CreateSSHKeyPairResponse represents the creation of an SSH Key Pair
type CreateSSHKeyPairResponse struct {
	KeyPair *SSHKeyPair `json:"keypair"`
}

// DeleteSSHKeyPairRequest represents a new keypair to be created
//
// http://cloudstack.apache.org/api/apidocs-4.10/apis/deleteSSHKeyPair.html
type DeleteSSHKeyPairRequest struct {
	Name      string `json:"name"`
	Account   string `json:"account,omitempty"`
	DomainID  string `json:"domainid,omitempty"`
	ProjectID string `json:"projectid,omitempty"`
}

// Command returns the CloudStack API command
func (req *DeleteSSHKeyPairRequest) Command() string {
	return "deleteSSHKeyPair"
}

// RegisterSSHKeyPairRequest represents a new registration of a public key in a keypair
//
// http://cloudstack.apache.org/api/apidocs-4.10/apis/registerSSHKeyPair.html
type RegisterSSHKeyPairRequest struct {
	Name      string `json:"name"`
	PublicKey string `json:"publickey"`
	Account   string `json:"account,omitempty"`
	DomainID  string `json:"domainid,omitempty"`
	ProjectID string `json:"projectid,omitempty"`
}

// Command returns the CloudStack API command
func (req *RegisterSSHKeyPairRequest) Command() string {
	return "registerSSHKeyPair"
}

// RegisterSSHKeyPairResponse represents the creation of an SSH Key Pair
type RegisterSSHKeyPairResponse struct {
	KeyPair *SSHKeyPair `json:"keypair"`
}

// ListSSHKeyPairsRequest represents a query for a list of SSH KeyPairs
//
// http://cloudstack.apache.org/api/apidocs-4.10/apis/listSSHKeyPairs.html
type ListSSHKeyPairsRequest struct {
	Account     string `json:"account,omitempty"`
	DomainID    string `json:"domainid,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
	IsRecursive bool   `json:"isrecursive,omitempty"`
	Keyword     string `json:"keyword,omitempty"`
	ListAll     bool   `json:"listall,omitempty"`
	Name        string `json:"name,omitempty"`
	Page        string `json:"page,omitempty"`
	PageSize    string `json:"pagesize,omitempty"`
	ProjectID   string `json:"projectid,omitempty"`
}

// Command returns the CloudStack API command
func (req *ListSSHKeyPairsRequest) Command() string {
	return "listSSHKeyPairs"
}

// ListSSHKeyPairsResponse represents a list of SSH key pairs
type ListSSHKeyPairsResponse struct {
	Count      int           `json:"count"`
	SSHKeyPair []*SSHKeyPair `json:"sshkeypair"`
}

// XXX ResetSSHKeyForVirtualMachine
//
// http://cloudstack.apache.org/api/apidocs-4.10/apis/resetSSHKeyForVirtualMachine.html

// CreateKeypair create a new SSH Key Pair
//
// Deprecated: will go away, use the API directly
func (exo *Client) CreateKeypair(name string) (*SSHKeyPair, error) {
	req := &CreateSSHKeyPairRequest{
		Name: name,
	}
	r := new(CreateSSHKeyPairResponse)
	err := exo.Request(req, r)
	if err != nil {
		return nil, err
	}

	return r.KeyPair, nil
}

// DeleteKeypair deletes an SSH key pair
//
// Deprecated: will go away, use the API directly
func (exo *Client) DeleteKeypair(name string) error {
	req := &DeleteSSHKeyPairRequest{
		Name: name,
	}
	return exo.BooleanRequest(req)
}

// RegisterKeypair registers a public key in a keypair
//
// Deprecated: will go away, use the API directly
func (exo *Client) RegisterKeypair(name string, publicKey string) (*SSHKeyPair, error) {
	req := &RegisterSSHKeyPairRequest{
		Name:      name,
		PublicKey: publicKey,
	}
	r := new(RegisterSSHKeyPairResponse)
	err := exo.Request(req, r)
	if err != nil {
		return nil, err
	}

	return r.KeyPair, nil
}
