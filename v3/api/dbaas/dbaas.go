package dbaas

import "github.com/exoscale/egoscale/v3/oapi"

type DBaaSIface interface {
	Integrations() *Integrations
}

//
// [Exoscale DBaaS]: https://community.exoscale.com/documentation/dbaas/
type DBaaS struct {
	oapiClient *oapi.ClientWithResponses
}

// NewDBaaS initializes DBaas with provided oapi Client.
func NewDBaaS(c *oapi.ClientWithResponses) *DBaaS {
	return &DBaaS{c}
}

func (a *DBaaS) Integrations() *Integrations {
	return NewIntegrations(a.oapiClient)
}

type MockDBaaS struct {
}

func NewMockDBaaS(c *oapi.ClientWithResponses) DBaaSIface {
	return &MockDBaaS{}
}

// func (a *MockDBaaS) Integrations() IntegrationsIface {
// 	return NewMockIntegrations()
// }
