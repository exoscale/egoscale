package global

import (
	"context"

	"github.com/exoscale/egoscale/v3/oapi"
	"github.com/exoscale/egoscale/v3/utils"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
)

type Operation struct {
	oapiClient *oapi.ClientWithResponses
}

func NewOperation(c *oapi.ClientWithResponses) *Operation {
	return &Operation{c}
}
func (a *Operation) Get(ctx context.Context, id openapi_types.UUID) (*oapi.Operation, error) {
	resp, err := a.oapiClient.GetOperationWithResponse(ctx, id)
	if err != nil {
		return nil, err
	}

	err = utils.ParseResponseError(resp.StatusCode(), resp.Body)
	if err != nil {
		return nil, err
	}

	return resp.JSON200, nil
}

