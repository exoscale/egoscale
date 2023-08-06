package iam

import (
	"context"

	"github.com/exoscale/egoscale/v3/oapi"
	"github.com/exoscale/egoscale/v3/utils"
)

type AccessKey struct {
	oapiClient *oapi.ClientWithResponses
}

func NewAccessKey(c *oapi.ClientWithResponses) *AccessKey {
	return &AccessKey{c}
}
func (a *AccessKey) List(ctx context.Context) ([]oapi.AccessKey, error) {
	resp, err := a.oapiClient.ListAccessKeysWithResponse(ctx)
	if err != nil {
		return nil, err
	}

	err = utils.ParseResponseError(resp.StatusCode(), resp.Body)
	if err != nil {
		return nil, err
	}

	return *resp.JSON200.AccessKeys, nil
}

func (a *AccessKey) ListKnownOperations(ctx context.Context) ([]oapi.AccessKeyOperation, error) {
	resp, err := a.oapiClient.ListAccessKeyKnownOperationsWithResponse(ctx)
	if err != nil {
		return nil, err
	}

	err = utils.ParseResponseError(resp.StatusCode(), resp.Body)
	if err != nil {
		return nil, err
	}

	return *resp.JSON200.AccessKeyOperations, nil
}

func (a *AccessKey) ListOperations(ctx context.Context) ([]oapi.AccessKeyOperation, error) {
	resp, err := a.oapiClient.ListAccessKeyOperationsWithResponse(ctx)
	if err != nil {
		return nil, err
	}

	err = utils.ParseResponseError(resp.StatusCode(), resp.Body)
	if err != nil {
		return nil, err
	}

	return *resp.JSON200.AccessKeyOperations, nil
}

func (a *AccessKey) Get(ctx context.Context, key string) (*oapi.AccessKey, error) {
	resp, err := a.oapiClient.GetAccessKeyWithResponse(ctx, key)
	if err != nil {
		return nil, err
	}

	err = utils.ParseResponseError(resp.StatusCode(), resp.Body)
	if err != nil {
		return nil, err
	}

	return resp.JSON200, nil
}

func (a *AccessKey) Create(ctx context.Context, body oapi.CreateAccessKeyJSONRequestBody) (*oapi.AccessKey, error) {
	resp, err := a.oapiClient.CreateAccessKeyWithResponse(ctx, body)
	if err != nil {
		return nil, err
	}

	err = utils.ParseResponseError(resp.StatusCode(), resp.Body)
	if err != nil {
		return nil, err
	}

	return resp.JSON200, nil
}

func (a *AccessKey) Revoke(ctx context.Context, key string) (*oapi.Operation, error) {
	resp, err := a.oapiClient.RevokeAccessKeyWithResponse(ctx, key)
	if err != nil {
		return nil, err
	}

	err = utils.ParseResponseError(resp.StatusCode(), resp.Body)
	if err != nil {
		return nil, err
	}

	return resp.JSON200, nil
}

