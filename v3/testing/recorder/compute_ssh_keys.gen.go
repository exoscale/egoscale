// Code generated by v3/generator; DO NOT EDIT.
package recorder

import (
	"context"

v3 "github.com/exoscale/egoscale/v3"
)

type SSHKeysAPI struct {
    Recordee *v3.SSHKeysAPI
    Recorder *Recorder
}


func (a *SSHKeysAPI) List(ctx context.Context) ([]v3.SshKey, error) {
    req := ArgsToMap()

    resp, err := a.Recordee.List(ctx)

    writeErr := a.Recorder.RecordCall("SSHKeysAPI.List", req, resp, err)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

func (a *SSHKeysAPI) Register(ctx context.Context, body v3.RegisterSshKeyJSONRequestBody) (*v3.Operation, error) {
    req := ArgsToMap(body)

    resp, err := a.Recordee.Register(ctx, body)

    writeErr := a.Recorder.RecordCall("SSHKeysAPI.Register", req, resp, err)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

func (a *SSHKeysAPI) Delete(ctx context.Context, name string) (*v3.Operation, error) {
    req := ArgsToMap(name)

    resp, err := a.Recordee.Delete(ctx, name)

    writeErr := a.Recorder.RecordCall("SSHKeysAPI.Delete", req, resp, err)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

func (a *SSHKeysAPI) Get(ctx context.Context, name string) (*v3.SshKey, error) {
    req := ArgsToMap(name)

    resp, err := a.Recordee.Get(ctx, name)

    writeErr := a.Recorder.RecordCall("SSHKeysAPI.Get", req, resp, err)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

