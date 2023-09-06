package v3

// IAMAPI provides access to [Exoscale IAM] API resources.
//
// [Exoscale IAM]: https://community.exoscale.com/documentation/iam/
type IAMAPI struct {
	client *Client
}

func (a *IAMAPI) Roles() *RolesAPI {
	return &RolesAPI{a.client}
}

func (a *IAMAPI) OrgPolicy() *OrgPolicyAPI {
	return &OrgPolicyAPI{a.client}
}

func (a *IAMAPI) AccessKey() *AccessKeyAPI {
	return &AccessKeyAPI{a.client}
}
