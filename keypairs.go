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

func (req *CreateSSHKeyPairRequest) name() string {
	return "createSSHKeyPair"
}

func (req *CreateSSHKeyPairRequest) response() interface{} {
	return new(CreateSSHKeyPairResponse)
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

func (req *DeleteSSHKeyPairRequest) name() string {
	return "deleteSSHKeyPair"
}

func (req *DeleteSSHKeyPairRequest) response() interface{} {
	return new(BooleanResponse)
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

func (req *RegisterSSHKeyPairRequest) name() string {
	return "registerSSHKeyPair"
}

func (req *RegisterSSHKeyPairRequest) response() interface{} {
	return new(RegisterSSHKeyPairResponse)
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

func (req *ListSSHKeyPairsRequest) name() string {
	return "listSSHKeyPairs"
}

func (req *ListSSHKeyPairsRequest) response() interface{} {
	return new(ListSSHKeyPairsResponse)
}

// ListSSHKeyPairsResponse represents a list of SSH key pairs
type ListSSHKeyPairsResponse struct {
	Count      int           `json:"count"`
	SSHKeyPair []*SSHKeyPair `json:"sshkeypair"`
}

// ResetSSHKeyForVirtualMachineRequest represents a change for the key pairs
//
// http://cloudstack.apache.org/api/apidocs-4.10/apis/resetSSHKeyForVirtualMachine.html
type ResetSSHKeyForVirtualMachineRequest struct {
	ID        string `json:"id"`
	KeyPair   string `json:"keypair"`
	Account   string `json:"account,omitempty"`
	DomainID  string `json:"domainid,omitempty"`
	ProjectID string `json:"projectid,omitempty"`
}

func (req *ResetSSHKeyForVirtualMachineRequest) name() string {
	return "resetSSHKeyForVirtualMachine"
}

func (req *ResetSSHKeyForVirtualMachineRequest) asyncResponse() interface{} {
	return new(ResetSSHKeyForVirtualMachineResponse)
}

// ResetSSHKeyForVirtualMachineResponse represents the modified VirtualMachine
type ResetSSHKeyForVirtualMachineResponse DeployVirtualMachineResponse

// CreateKeypair create a new SSH Key Pair
//
// Deprecated: will go away, use the API directly
func (exo *Client) CreateKeypair(name string) (*SSHKeyPair, error) {
	req := &CreateSSHKeyPairRequest{
		Name: name,
	}
	resp, err := exo.Request(req)
	if err != nil {
		return nil, err
	}

	return resp.(*CreateSSHKeyPairResponse).KeyPair, nil
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
	resp, err := exo.Request(req)
	if err != nil {
		return nil, err
	}

	return resp.(*RegisterSSHKeyPairResponse).KeyPair, nil
}
