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

    writeErr := RecordCall("AccessKeyAPI.List", req, resp, err)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

func (a *AccessKeyAPI) ListKnownOperations(ctx context.Context) ([]v3.AccessKeyOperation, error) {
    req := argsToMap()

    resp, err := a.Recordee.ListKnownOperations(ctx, )

    writeErr := RecordCall("AccessKeyAPI.ListKnownOperations", req, resp, err)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

func (a *AccessKeyAPI) ListOperations(ctx context.Context) ([]v3.AccessKeyOperation, error) {
    req := argsToMap()

    resp, err := a.Recordee.ListOperations(ctx, )

    writeErr := RecordCall("AccessKeyAPI.ListOperations", req, resp, err)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

func (a *AccessKeyAPI) Get(ctx context.Context, key string) (*v3.AccessKey, error) {
    req := argsToMap(key)

    resp, err := a.Recordee.Get(ctx, key)

    writeErr := RecordCall("AccessKeyAPI.Get", req, resp, err)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

func (a *AccessKeyAPI) Create(ctx context.Context, body v3.CreateAccessKeyJSONRequestBody) (*v3.AccessKey, error) {
    req := argsToMap(body)

    resp, err := a.Recordee.Create(ctx, body)

    writeErr := RecordCall("AccessKeyAPI.Create", req, resp, err)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

func (a *AccessKeyAPI) Revoke(ctx context.Context, key string) (*v3.Operation, error) {
    req := argsToMap(key)

    resp, err := a.Recordee.Revoke(ctx, key)

    writeErr := RecordCall("AccessKeyAPI.Revoke", req, resp, err)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

