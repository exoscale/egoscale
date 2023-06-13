package consumerapi

import (
	"context"

	"github.com/sauterp/egoscale-rewrite-experiments/oapi-codegen/gen/oapi"
)

// THIS CODE IS GENERATED

type ExoPlatform1 struct {
	*oapi.ClientWithResponses
}
type AntiAffinityGroup struct {
	*oapi.ClientWithResponses
}

func (antiaffinitygroup *AntiAffinityGroup) Get(ctx context.Context, id string, reqEditors ...oapi.RequestEditorFn) *oapi.AntiAffinityGroup {
	resp, err2 := antiaffinitygroup.ClientWithResponses.GetAntiAffinityGroupWithResponse(ctx, id, reqEditors...)
	if err2 != nil {
		panic(err2)
	}

	return resp.JSON200
}

func (exoplatform1 *ExoPlatform1) AntiAffinityGroup() *AntiAffinityGroup {
	return &AntiAffinityGroup{
		exoplatform1.ClientWithResponses,
	}
}

func (client *Client) ExoPlatform1() *ExoPlatform1 {
	return &ExoPlatform1{
		client.ClientWithResponses,
	}
}
