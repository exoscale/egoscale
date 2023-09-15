// Code generated by v3/generator; DO NOT EDIT.
package recorder

import (
	"context"


	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
v3 "github.com/exoscale/egoscale/v3"
)

type PrivateNetworksAPI struct {
    Recordee *v3.PrivateNetworksAPI
}


func (a *PrivateNetworksAPI) List(ctx context.Context) ([]v3.PrivateNetwork, error) {
    req := argsToMap()

    resp, err := a.Recordee.List(ctx, )
    respMap := argsToMap(resp, err)

    writeErr := WriteTestdata(req, respMap, 0)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

func (a *PrivateNetworksAPI) Get(ctx context.Context, id openapi_types.UUID) (*v3.PrivateNetwork, error) {
    req := argsToMap(id)

    resp, err := a.Recordee.Get(ctx, id)
    respMap := argsToMap(resp, err)

    writeErr := WriteTestdata(req, respMap, 0)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

func (a *PrivateNetworksAPI) Create(ctx context.Context, body v3.CreatePrivateNetworkJSONRequestBody) (*v3.Operation, error) {
    req := argsToMap(body)

    resp, err := a.Recordee.Create(ctx, body)
    respMap := argsToMap(resp, err)

    writeErr := WriteTestdata(req, respMap, 0)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

func (a *PrivateNetworksAPI) Update(ctx context.Context, id openapi_types.UUID, body v3.UpdatePrivateNetworkJSONRequestBody) (*v3.Operation, error) {
    req := argsToMap(id, body)

    resp, err := a.Recordee.Update(ctx, id, body)
    respMap := argsToMap(resp, err)

    writeErr := WriteTestdata(req, respMap, 0)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

func (a *PrivateNetworksAPI) Delete(ctx context.Context, id openapi_types.UUID) (*v3.Operation, error) {
    req := argsToMap(id)

    resp, err := a.Recordee.Delete(ctx, id)
    respMap := argsToMap(resp, err)

    writeErr := WriteTestdata(req, respMap, 0)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

