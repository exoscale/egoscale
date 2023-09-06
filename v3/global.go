package v3

// GlobalAPI provides access to global API resources.
type GlobalAPI struct {
	client *Client
}

func (a *GlobalAPI) OrgQuotas() *OrgQuotasAPI {
	return &OrgQuotasAPI{a.client}
}

func (a *GlobalAPI) Operations() *OperationAPI {
	return &OperationAPI{a.client}
}
