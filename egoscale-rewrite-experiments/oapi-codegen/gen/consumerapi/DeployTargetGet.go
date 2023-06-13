package consumerapi

import (
	"context"

	"github.com/sauterp/egoscale-rewrite-experiments/oapi-codegen/gen/oapi"
)

// THIS CODE IS GENERATED

type ExoPlatform3 struct {
	*oapi.ClientWithResponses
}
type DeployTarget struct {
	*oapi.ClientWithResponses
}

func (deploytarget *DeployTarget) Get(ctx context.Context, id string, reqEditors ...oapi.RequestEditorFn) *oapi.DeployTarget {
	resp, err2 := deploytarget.ClientWithResponses.GetDeployTargetWithResponse(ctx, id, reqEditors...)
	if err2 != nil {
		panic(err2)
	}

	return resp.JSON200
}

func (exoplatform3 *ExoPlatform3) DeployTarget() *DeployTarget {
	return &DeployTarget{
		exoplatform3.ClientWithResponses,
	}
}

func (client *Client) ExoPlatform3() *ExoPlatform3 {
	return &ExoPlatform3{
		client.ClientWithResponses,
	}
}
