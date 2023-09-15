// Code generated by v3/generator; DO NOT EDIT.
package recorder

import (
	"context"


	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
v3 "github.com/exoscale/egoscale/v3"
)

type InstanceTypesAPI struct {
    Recordee *v3.InstanceTypesAPI
}


func (a *InstanceTypesAPI) List(ctx context.Context) ([]v3.InstanceType, error) {
    req := argsToMap()

    resp, err := a.Recordee.List(ctx, )
    respMap := argsToMap(resp, err)

    writeErr := WriteTestdata(req, respMap, 0)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

func (a *InstanceTypesAPI) Get(ctx context.Context, id openapi_types.UUID) (*v3.InstanceType, error) {
    req := argsToMap(id)

    resp, err := a.Recordee.Get(ctx, id)
    respMap := argsToMap(resp, err)

    writeErr := WriteTestdata(req, respMap, 0)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

