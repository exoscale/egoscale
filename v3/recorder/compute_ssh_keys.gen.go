// Code generated by v3/generator; DO NOT EDIT.
package recorder

import (
	"context"

v3 "github.com/exoscale/egoscale/v3"
)

type SSHKeysAPI struct {
    Recordee *v3.SSHKeysAPI
}


func (a *SSHKeysAPI) List(ctx context.Context) ([]v3.SshKey, error) {
    req := argsToMap()

    resp, err := a.Recordee.List(ctx, )

    writeErr := WriteTestdata("SSHKeysAPI.List", req, resp, err)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

func (a *SSHKeysAPI) Register(ctx context.Context, body v3.RegisterSshKeyJSONRequestBody) (*v3.Operation, error) {
    req := argsToMap(body)

    resp, err := a.Recordee.Register(ctx, body)

    writeErr := WriteTestdata("SSHKeysAPI.Register", req, resp, err)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

func (a *SSHKeysAPI) Delete(ctx context.Context, name string) (*v3.Operation, error) {
    req := argsToMap(name)

    resp, err := a.Recordee.Delete(ctx, name)

    writeErr := WriteTestdata("SSHKeysAPI.Delete", req, resp, err)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

func (a *SSHKeysAPI) Get(ctx context.Context, name string) (*v3.SshKey, error) {
    req := argsToMap(name)

    resp, err := a.Recordee.Get(ctx, name)

    writeErr := WriteTestdata("SSHKeysAPI.Get", req, resp, err)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

