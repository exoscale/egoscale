package consumerapi

import (
	"context"

	"github.com/sauterp/egoscale-rewrite-experiments/oapi-codegen/gen/oapi"
)

// THIS CODE IS GENERATED

type ExoPlatform8 struct {
	*oapi.ClientWithResponses
}
type InstancePool struct {
	*oapi.ClientWithResponses
}

func (instancepool *InstancePool) Get(ctx context.Context, id string, reqEditors ...oapi.RequestEditorFn) *oapi.InstancePool {
	resp, err2 := instancepool.ClientWithResponses.GetInstancePoolWithResponse(ctx, id, reqEditors...)
	if err2 != nil {
		panic(err2)
	}

	return resp.JSON200
}

func (exoplatform8 *ExoPlatform8) InstancePool() *InstancePool {
	return &InstancePool{
		exoplatform8.ClientWithResponses,
	}
}

func (client *Client) ExoPlatform8() *ExoPlatform8 {
	return &ExoPlatform8{
		client.ClientWithResponses,
	}
}
