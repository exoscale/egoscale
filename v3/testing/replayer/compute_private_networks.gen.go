// Code generated by v3/generator; DO NOT EDIT.
package replayer

import (
	"context"


	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
v3 "github.com/exoscale/egoscale/v3"
	"github.com/exoscale/egoscale/v3/testing/recorder"
)

type PrivateNetworksAPI struct {
    Replayer *Replayer

    ListHook func(ctx context.Context) error

    GetHook func(ctx context.Context, id openapi_types.UUID) error

    CreateHook func(ctx context.Context, body v3.CreatePrivateNetworkJSONRequestBody) error

    UpdateHook func(ctx context.Context, id openapi_types.UUID, body v3.UpdatePrivateNetworkJSONRequestBody) error

    DeleteHook func(ctx context.Context, id openapi_types.UUID) error

}


func (a *PrivateNetworksAPI) List(ctx context.Context) ([]v3.PrivateNetwork, error) {
    resp := InitializeReturnType[[]v3.PrivateNetwork](a.List)

    expectedArgs := make(recorder.CallParameters)
    var returnErr error
    err := a.Replayer.GetTestCall(&resp, &expectedArgs, &returnErr)
    if err != nil {
        panic(err)
    }

    if a.ListHook == nil {
        
    } else {
        if err := a.ListHook(ctx); err != nil {
            panic(err)
        }
    }

    return resp, returnErr
}

func (a *PrivateNetworksAPI) Get(ctx context.Context, id openapi_types.UUID) (*v3.PrivateNetwork, error) {
    resp := InitializeReturnType[*v3.PrivateNetwork](a.Get)

    expectedArgs := make(recorder.CallParameters)
    var returnErr error
    err := a.Replayer.GetTestCall(&resp, &expectedArgs, &returnErr)
    if err != nil {
        panic(err)
    }

    if a.GetHook == nil {
        
             a.Replayer.AssertArgs(expectedArgs, id)
        
    } else {
        if err := a.GetHook(ctx, id); err != nil {
            panic(err)
        }
    }

    return resp, returnErr
}

func (a *PrivateNetworksAPI) Create(ctx context.Context, body v3.CreatePrivateNetworkJSONRequestBody) (*v3.Operation, error) {
    resp := InitializeReturnType[*v3.Operation](a.Create)

    expectedArgs := make(recorder.CallParameters)
    var returnErr error
    err := a.Replayer.GetTestCall(&resp, &expectedArgs, &returnErr)
    if err != nil {
        panic(err)
    }

    if a.CreateHook == nil {
        
             a.Replayer.AssertArgs(expectedArgs, body)
        
    } else {
        if err := a.CreateHook(ctx, body); err != nil {
            panic(err)
        }
    }

    return resp, returnErr
}

func (a *PrivateNetworksAPI) Update(ctx context.Context, id openapi_types.UUID, body v3.UpdatePrivateNetworkJSONRequestBody) (*v3.Operation, error) {
    resp := InitializeReturnType[*v3.Operation](a.Update)

    expectedArgs := make(recorder.CallParameters)
    var returnErr error
    err := a.Replayer.GetTestCall(&resp, &expectedArgs, &returnErr)
    if err != nil {
        panic(err)
    }

    if a.UpdateHook == nil {
        
             a.Replayer.AssertArgs(expectedArgs, id, body)
        
    } else {
        if err := a.UpdateHook(ctx, id, body); err != nil {
            panic(err)
        }
    }

    return resp, returnErr
}

func (a *PrivateNetworksAPI) Delete(ctx context.Context, id openapi_types.UUID) (*v3.Operation, error) {
    resp := InitializeReturnType[*v3.Operation](a.Delete)

    expectedArgs := make(recorder.CallParameters)
    var returnErr error
    err := a.Replayer.GetTestCall(&resp, &expectedArgs, &returnErr)
    if err != nil {
        panic(err)
    }

    if a.DeleteHook == nil {
        
             a.Replayer.AssertArgs(expectedArgs, id)
        
    } else {
        if err := a.DeleteHook(ctx, id); err != nil {
            panic(err)
        }
    }

    return resp, returnErr
}

