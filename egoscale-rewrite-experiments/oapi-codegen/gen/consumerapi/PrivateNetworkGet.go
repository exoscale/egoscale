package consumerapi

import (
	"context"

	"github.com/sauterp/egoscale-rewrite-experiments/oapi-codegen/gen/oapi"
)

// THIS CODE IS GENERATED

type ExoPlatform14 struct {
	*oapi.ClientWithResponses
}
type PrivateNetwork struct {
	*oapi.ClientWithResponses
}

func (privatenetwork *PrivateNetwork) Get(ctx context.Context, id string, reqEditors ...oapi.RequestEditorFn) *oapi.PrivateNetwork {
	resp, err2 := privatenetwork.ClientWithResponses.GetPrivateNetworkWithResponse(ctx, id, reqEditors...)
	if err2 != nil {
		panic(err2)
	}

	return resp.JSON200
}

func (exoplatform14 *ExoPlatform14) PrivateNetwork() *PrivateNetwork {
	return &PrivateNetwork{
		exoplatform14.ClientWithResponses,
	}
}

func (client *Client) ExoPlatform14() *ExoPlatform14 {
	return &ExoPlatform14{
		client.ClientWithResponses,
	}
}
