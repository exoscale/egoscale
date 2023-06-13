package consumerapi

import (
	"context"

	"github.com/sauterp/egoscale-rewrite-experiments/oapi-codegen/gen/oapi"
)

// THIS CODE IS GENERATED

type ComputeAPI struct {
	*oapi.ClientWithResponses
}
type Instance struct {
	*oapi.ClientWithResponses
}

func (instance *Instance) Get(ctx context.Context, id string, reqEditors ...oapi.RequestEditorFn) *oapi.Instance {
	resp, err2 := instance.ClientWithResponses.GetInstanceWithResponse(ctx, id, reqEditors...)
	if err2 != nil {
		panic(err2)
	}

	return resp.JSON200
}

func (computeapi *ComputeAPI) Instance() *Instance {
	return &Instance{
		computeapi.ClientWithResponses,
	}
}

func (client *Client) ComputeAPI() *ComputeAPI {
	return &ComputeAPI{
		client.ClientWithResponses,
	}
}
