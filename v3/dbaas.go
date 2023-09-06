package v3

// DBaaSAPI provides access to [Exoscale DBaaS] API resources.
//
// [Exoscale DBaaS]: https://community.exoscale.com/documentation/dbaas/
type DBaaSAPI struct {
	client *Client
}

func (a *DBaaSAPI) Integrations() *IntegrationsAPI {
	return &IntegrationsAPI{a.client}
}
