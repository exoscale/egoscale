package global

import (
	"context"

	"github.com/exoscale/egoscale/v3/oapi"
	"github.com/exoscale/egoscale/v3/utils"
)

type OrgQuotas struct {
	oapiClient *oapi.ClientWithResponses
}

func NewOrgQuotas(c *oapi.ClientWithResponses) *OrgQuotas {
	return &OrgQuotas{c}
}
func (a *OrgQuotas) List(ctx context.Context) ([]oapi.Quota, error) {
	resp, err := a.oapiClient.ListQuotasWithResponse(ctx)
	if err != nil {
		return nil, err
	}

	err = utils.ParseResponseError(resp.StatusCode(), resp.Body)
	if err != nil {
		return nil, err
	}

	return *resp.JSON200.Quotas, nil
}

func (a *OrgQuotas) Get(ctx context.Context, entity string) (*oapi.Quota, error) {
	resp, err := a.oapiClient.GetQuotaWithResponse(ctx, entity)
	if err != nil {
		return nil, err
	}

	err = utils.ParseResponseError(resp.StatusCode(), resp.Body)
	if err != nil {
		return nil, err
	}

	return resp.JSON200, nil
}

