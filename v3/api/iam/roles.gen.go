package iam

import (
	"context"

	"github.com/exoscale/egoscale/v3/oapi"
	"github.com/exoscale/egoscale/v3/utils"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
)

type Roles struct {
	oapiClient *oapi.ClientWithResponses
}

func NewRoles(c *oapi.ClientWithResponses) *Roles {
	return &Roles{c}
}
func (a *Roles) List(ctx context.Context) ([]oapi.IamRole, error) {
	resp, err := a.oapiClient.ListIamRolesWithResponse(ctx)
	if err != nil {
		return nil, err
	}

	err = utils.ParseResponseError(resp.StatusCode(), resp.Body)
	if err != nil {
		return nil, err
	}

	return *resp.JSON200.IamRoles, nil
}

func (a *Roles) Get(ctx context.Context, id openapi_types.UUID) (*oapi.IamRole, error) {
	resp, err := a.oapiClient.GetIamRoleWithResponse(ctx, id)
	if err != nil {
		return nil, err
	}

	err = utils.ParseResponseError(resp.StatusCode(), resp.Body)
	if err != nil {
		return nil, err
	}

	return resp.JSON200, nil
}

func (a *Roles) Create(ctx context.Context, body oapi.CreateIamRoleJSONRequestBody) (*oapi.Operation, error) {
	resp, err := a.oapiClient.CreateIamRoleWithResponse(ctx, body)
	if err != nil {
		return nil, err
	}

	err = utils.ParseResponseError(resp.StatusCode(), resp.Body)
	if err != nil {
		return nil, err
	}

	return resp.JSON200, nil
}

func (a *Roles) Delete(ctx context.Context, id openapi_types.UUID) (*oapi.Operation, error) {
	resp, err := a.oapiClient.DeleteIamRoleWithResponse(ctx, id)
	if err != nil {
		return nil, err
	}

	err = utils.ParseResponseError(resp.StatusCode(), resp.Body)
	if err != nil {
		return nil, err
	}

	return resp.JSON200, nil
}

func (a *Roles) Update(ctx context.Context, id openapi_types.UUID, body oapi.UpdateIamRoleJSONRequestBody) (*oapi.Operation, error) {
	resp, err := a.oapiClient.UpdateIamRoleWithResponse(ctx, id, body)
	if err != nil {
		return nil, err
	}

	err = utils.ParseResponseError(resp.StatusCode(), resp.Body)
	if err != nil {
		return nil, err
	}

	return resp.JSON200, nil
}

func (a *Roles) UpdatePolicy(ctx context.Context, id openapi_types.UUID, body oapi.UpdateIamRolePolicyJSONRequestBody) (*oapi.Operation, error) {
	resp, err := a.oapiClient.UpdateIamRolePolicyWithResponse(ctx, id, body)
	if err != nil {
		return nil, err
	}

	err = utils.ParseResponseError(resp.StatusCode(), resp.Body)
	if err != nil {
		return nil, err
	}

	return resp.JSON200, nil
}

