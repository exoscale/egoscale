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
    return a.Recordee.List(ctx, )
}

func (a *SSHKeysAPI) Register(ctx context.Context, body v3.RegisterSshKeyJSONRequestBody) (*v3.Operation, error) {
    return a.Recordee.Register(ctx, body)
}

func (a *SSHKeysAPI) Delete(ctx context.Context, name string) (*v3.Operation, error) {
    return a.Recordee.Delete(ctx, name)
}

func (a *SSHKeysAPI) Get(ctx context.Context, name string) (*v3.SshKey, error) {
    return a.Recordee.Get(ctx, name)
}

