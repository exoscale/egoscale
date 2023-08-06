package dbaas

import (
	"context"

	"github.com/exoscale/egoscale/v3/oapi"
	"github.com/exoscale/egoscale/v3/utils"
)

type Integrations struct {
	oapiClient *oapi.ClientWithResponses
}

func NewIntegrations(c *oapi.ClientWithResponses) *Integrations {
	return &Integrations{c}
}
func (a *Integrations) ListSettings(ctx context.Context, integrationType string, sourceType string, destType string) (*oapi.DBaaSIntegrationSettings, error) {
	resp, err := a.oapiClient.ListDbaasIntegrationSettingsWithResponse(ctx, integrationType, sourceType, destType)
	if err != nil {
		return nil, err
	}

	err = utils.ParseResponseError(resp.StatusCode(), resp.Body)
	if err != nil {
		return nil, err
	}

	return oapi.FromListDbaasIntegrationSettingsResponse(resp), nil
}

