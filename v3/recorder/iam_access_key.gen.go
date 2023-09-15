// Code generated by v3/generator; DO NOT EDIT.
package recorder

import (
	"context"

v3 "github.com/exoscale/egoscale/v3"
)

type AccessKeyAPI struct {
    Recordee *v3.AccessKeyAPI
}


func (a *AccessKeyAPI) List(ctx context.Context) ([]v3.AccessKey, error) {
    req := argsToMap()

    resp, err := a.Recordee.List(ctx, )
    respMap := argsToMap(resp, err)

    writeErr := WriteTestdata(req, respMap, 0)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

func (a *AccessKeyAPI) ListKnownOperations(ctx context.Context) ([]v3.AccessKeyOperation, error) {
    req := argsToMap()

    resp, err := a.Recordee.ListKnownOperations(ctx, )
    respMap := argsToMap(resp, err)

    writeErr := WriteTestdata(req, respMap, 0)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

func (a *AccessKeyAPI) ListOperations(ctx context.Context) ([]v3.AccessKeyOperation, error) {
    req := argsToMap()

    resp, err := a.Recordee.ListOperations(ctx, )
    respMap := argsToMap(resp, err)

    writeErr := WriteTestdata(req, respMap, 0)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

func (a *AccessKeyAPI) Get(ctx context.Context, key string) (*v3.AccessKey, error) {
    req := argsToMap(key)

    resp, err := a.Recordee.Get(ctx, key)
    respMap := argsToMap(resp, err)

    writeErr := WriteTestdata(req, respMap, 0)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

func (a *AccessKeyAPI) Create(ctx context.Context, body v3.CreateAccessKeyJSONRequestBody) (*v3.AccessKey, error) {
    req := argsToMap(body)

    resp, err := a.Recordee.Create(ctx, body)
    respMap := argsToMap(resp, err)

    writeErr := WriteTestdata(req, respMap, 0)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

func (a *AccessKeyAPI) Revoke(ctx context.Context, key string) (*v3.Operation, error) {
    req := argsToMap(key)

    resp, err := a.Recordee.Revoke(ctx, key)
    respMap := argsToMap(resp, err)

    writeErr := WriteTestdata(req, respMap, 0)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

