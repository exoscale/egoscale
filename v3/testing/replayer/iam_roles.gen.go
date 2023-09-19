// Code generated by v3/generator; DO NOT EDIT.
package replayer

import (
	"context"


	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
v3 "github.com/exoscale/egoscale/v3"
	"github.com/exoscale/egoscale/v3/testing/recorder"
)

type RolesAPI struct {
    Replayer *Replayer

    ListHook func(ctx context.Context) error

    GetHook func(ctx context.Context, id openapi_types.UUID) error

    CreateHook func(ctx context.Context, body v3.CreateIamRoleJSONRequestBody) error

    DeleteHook func(ctx context.Context, id openapi_types.UUID) error

    UpdateHook func(ctx context.Context, id openapi_types.UUID, body v3.UpdateIamRoleJSONRequestBody) error

    UpdatePolicyHook func(ctx context.Context, id openapi_types.UUID, body v3.UpdateIamRolePolicyJSONRequestBody) error

}


func (a *RolesAPI) List(ctx context.Context) ([]v3.IamRole, error) {
    resp := InitializeReturnType[[]v3.IamRole](a.List)

    expectedArgs := make(recorder.CallParameters)
    var returnErr error
    err := a.Replayer.GetTestCall(&resp, &expectedArgs, &returnErr)
    if err != nil {
        panic(err)
    }

    if a.ListHook == nil {
        a.Replayer.AssertArgs(expectedArgs, ctx)
    } else {
        if err := a.ListHook(ctx); err != nil {
            panic(err)
        }
    }

    return resp, returnErr
}

func (a *RolesAPI) Get(ctx context.Context, id openapi_types.UUID) (*v3.IamRole, error) {
    resp := InitializeReturnType[*v3.IamRole](a.Get)

    expectedArgs := make(recorder.CallParameters)
    var returnErr error
    err := a.Replayer.GetTestCall(&resp, &expectedArgs, &returnErr)
    if err != nil {
        panic(err)
    }

    if a.GetHook == nil {
        a.Replayer.AssertArgs(expectedArgs, ctx, id)
    } else {
        if err := a.GetHook(ctx, id); err != nil {
            panic(err)
        }
    }

    return resp, returnErr
}

func (a *RolesAPI) Create(ctx context.Context, body v3.CreateIamRoleJSONRequestBody) (*v3.Operation, error) {
    resp := InitializeReturnType[*v3.Operation](a.Create)

    expectedArgs := make(recorder.CallParameters)
    var returnErr error
    err := a.Replayer.GetTestCall(&resp, &expectedArgs, &returnErr)
    if err != nil {
        panic(err)
    }

    if a.CreateHook == nil {
        a.Replayer.AssertArgs(expectedArgs, ctx, body)
    } else {
        if err := a.CreateHook(ctx, body); err != nil {
            panic(err)
        }
    }

    return resp, returnErr
}

func (a *RolesAPI) Delete(ctx context.Context, id openapi_types.UUID) (*v3.Operation, error) {
    resp := InitializeReturnType[*v3.Operation](a.Delete)

    expectedArgs := make(recorder.CallParameters)
    var returnErr error
    err := a.Replayer.GetTestCall(&resp, &expectedArgs, &returnErr)
    if err != nil {
        panic(err)
    }

    if a.DeleteHook == nil {
        a.Replayer.AssertArgs(expectedArgs, ctx, id)
    } else {
        if err := a.DeleteHook(ctx, id); err != nil {
            panic(err)
        }
    }

    return resp, returnErr
}

func (a *RolesAPI) Update(ctx context.Context, id openapi_types.UUID, body v3.UpdateIamRoleJSONRequestBody) (*v3.Operation, error) {
    resp := InitializeReturnType[*v3.Operation](a.Update)

    expectedArgs := make(recorder.CallParameters)
    var returnErr error
    err := a.Replayer.GetTestCall(&resp, &expectedArgs, &returnErr)
    if err != nil {
        panic(err)
    }

    if a.UpdateHook == nil {
        a.Replayer.AssertArgs(expectedArgs, ctx, id, body)
    } else {
        if err := a.UpdateHook(ctx, id, body); err != nil {
            panic(err)
        }
    }

    return resp, returnErr
}

func (a *RolesAPI) UpdatePolicy(ctx context.Context, id openapi_types.UUID, body v3.UpdateIamRolePolicyJSONRequestBody) (*v3.Operation, error) {
    resp := InitializeReturnType[*v3.Operation](a.UpdatePolicy)

    expectedArgs := make(recorder.CallParameters)
    var returnErr error
    err := a.Replayer.GetTestCall(&resp, &expectedArgs, &returnErr)
    if err != nil {
        panic(err)
    }

    if a.UpdatePolicyHook == nil {
        a.Replayer.AssertArgs(expectedArgs, ctx, id, body)
    } else {
        if err := a.UpdatePolicyHook(ctx, id, body); err != nil {
            panic(err)
        }
    }

    return resp, returnErr
}

