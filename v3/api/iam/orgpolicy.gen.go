package iam

import (
	"context"

	"github.com/exoscale/egoscale/v3/oapi"
	"github.com/exoscale/egoscale/v3/utils"
)

type OrgPolicy struct {
	oapiClient *oapi.ClientWithResponses
}

func NewOrgPolicy(c *oapi.ClientWithResponses) *OrgPolicy {
	return &OrgPolicy{c}
}
func (a *OrgPolicy) Get(ctx context.Context) ([]oapi.IamPolicy, error) {
	resp, err := a.oapiClient.GetIamOrganizationPolicyWithResponse(ctx)
	if err != nil {
		return nil, err
	}

	err = utils.ParseResponseError(resp.StatusCode(), resp.Body)
	if err != nil {
		return nil, err
	}

	return *resp.JSON200, nil
}

func (a *OrgPolicy) Update(ctx context.Context, body oapi.UpdateIamOrganizationPolicyJSONRequestBody) (*oapi.Operation, error) {
	resp, err := a.oapiClient.UpdateIamOrganizationPolicyWithResponse(ctx, body)
	if err != nil {
		return nil, err
	}

	err = utils.ParseResponseError(resp.StatusCode(), resp.Body)
	if err != nil {
		return nil, err
	}

	return resp.JSON200, nil
}

