package consumerapi

import (
	"context"

	"github.com/sauterp/egoscale-rewrite-experiments/oapi-codegen/gen/oapi"
)

// THIS CODE IS GENERATED

type ExoPlatform9 struct {
	*oapi.ClientWithResponses
}
type InstanceType struct {
	*oapi.ClientWithResponses
}

func (instancetype *InstanceType) Get(ctx context.Context, id string, reqEditors ...oapi.RequestEditorFn) *oapi.InstanceType {
	resp, err2 := instancetype.ClientWithResponses.GetInstanceTypeWithResponse(ctx, id, reqEditors...)
	if err2 != nil {
		panic(err2)
	}

	return resp.JSON200
}

func (exoplatform9 *ExoPlatform9) InstanceType() *InstanceType {
	return &InstanceType{
		exoplatform9.ClientWithResponses,
	}
}

func (client *Client) ExoPlatform9() *ExoPlatform9 {
	return &ExoPlatform9{
		client.ClientWithResponses,
	}
}
