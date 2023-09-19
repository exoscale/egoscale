// Code generated by v3/generator; DO NOT EDIT.
package replayer

import (
	"context"

v3 "github.com/exoscale/egoscale/v3"
	"github.com/exoscale/egoscale/v3/testing/recorder"
)

type OrgPolicyAPI struct {
    Replayer *Replayer

    GetHook func(ctx context.Context) error

    UpdateHook func(ctx context.Context, body v3.UpdateIamOrganizationPolicyJSONRequestBody) error

}


func (a *OrgPolicyAPI) Get(ctx context.Context) (*v3.IamPolicy, error) {
    resp := InitializeReturnType[*v3.IamPolicy](a.Get)

    expectedArgs := make(recorder.CallParameters)
    var returnErr error
    err := a.Replayer.GetTestCall(&resp, &expectedArgs, &returnErr)
    if err != nil {
        panic(err)
    }

    if a.GetHook == nil {
        
    } else {
        if err := a.GetHook(ctx); err != nil {
            panic(err)
        }
    }

    return resp, returnErr
}

func (a *OrgPolicyAPI) Update(ctx context.Context, body v3.UpdateIamOrganizationPolicyJSONRequestBody) (*v3.Operation, error) {
    resp := InitializeReturnType[*v3.Operation](a.Update)

    expectedArgs := make(recorder.CallParameters)
    var returnErr error
    err := a.Replayer.GetTestCall(&resp, &expectedArgs, &returnErr)
    if err != nil {
        panic(err)
    }

    if a.UpdateHook == nil {
        
             a.Replayer.AssertArgs(expectedArgs, body)
        
    } else {
        if err := a.UpdateHook(ctx, body); err != nil {
            panic(err)
        }
    }

    return resp, returnErr
}

