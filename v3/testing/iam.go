package testing

import (
	"context"

	v3 "github.com/exoscale/egoscale/v3"
	"github.com/exoscale/egoscale/v3/testing/recorder"
	"github.com/exoscale/egoscale/v3/testing/replayer"
)

type IAMAPIIface interface {
	Roles() *v3.RolesAPI
	OrgPolicy() *v3.OrgPolicyAPI
	AccessKey() AccessKeyAPIIface
}

type IAMAPIRecorder struct {
	client   *TestClient
	Recorder *recorder.Recorder
}

func (a *IAMAPIRecorder) Roles() *v3.RolesAPI {
	panic("not implemented")
}

func (a *IAMAPIRecorder) OrgPolicy() *v3.OrgPolicyAPI {
	panic("not implemented")
}

func (a *IAMAPIRecorder) AccessKey() AccessKeyAPIIface {
	return &recorder.AccessKeyAPI{
		Recordee: a.client.Client.IAM().AccessKey(),
		Recorder: a.Recorder,
	}
}

type IAMAPIReplayer struct {
	Replayer *replayer.Replayer
}

func (a *IAMAPIReplayer) Roles() *v3.RolesAPI {
	panic("not implemented")
}

func (a *IAMAPIReplayer) OrgPolicy() *v3.OrgPolicyAPI {
	panic("not implemented")
}

func (a *IAMAPIReplayer) AccessKey() AccessKeyAPIIface {
	return &replayer.AccessKeyAPI{
		Replayer: a.Replayer,
	}
}

type AccessKeyAPIIface interface {
	List(ctx context.Context) ([]v3.AccessKey, error)
	ListKnownOperations(ctx context.Context) ([]v3.AccessKeyOperation, error)
	ListOperations(ctx context.Context) ([]v3.AccessKeyOperation, error)
	Get(ctx context.Context, key string) (*v3.AccessKey, error)
	Create(ctx context.Context, body v3.CreateAccessKeyJSONRequestBody) (*v3.AccessKey, error)
	Revoke(ctx context.Context, key string) (*v3.Operation, error)
}
