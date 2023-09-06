package v3

// ComputeAPI provides access to [Exoscale Compute] API resources.
//
// [Exoscale Compute]: https://community.exoscale.com/documentation/compute/
type ComputeAPI struct {
	client *Client
}

func (a *ComputeAPI) InstanceTypes() *InstanceTypesAPI {
	return &InstanceTypesAPI{a.client}
}

func (a *ComputeAPI) PrivateNetworks() *PrivateNetworksAPI {
	return &PrivateNetworksAPI{a.client}
}

func (a *ComputeAPI) SSHKeys() *SSHKeysAPI {
	return &SSHKeysAPI{a.client}
}

func (a *ComputeAPI) LoadBalancers() *LoadBalancersAPI {
	return &LoadBalancersAPI{a.client}
}
