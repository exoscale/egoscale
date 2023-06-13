package consumerapi

import (
	"context"

	"github.com/sauterp/egoscale-rewrite-experiments/oapi-codegen/gen/oapi"
)

// THIS CODE IS GENERATED

type ExoPlatform20 struct {
	*oapi.ClientWithResponses
}
type SksNodepool struct {
	*oapi.ClientWithResponses
}

func (sksnodepool *SksNodepool) Get(ctx context.Context, id string, sksNodepoolId string, reqEditors ...oapi.RequestEditorFn) *oapi.SksNodepool {
	resp, err2 := sksnodepool.ClientWithResponses.GetSksNodepoolWithResponse(ctx, id, sksNodepoolId, reqEditors...)
	if err2 != nil {
		panic(err2)
	}

	return resp.JSON200
}

func (exoplatform20 *ExoPlatform20) SksNodepool() *SksNodepool {
	return &SksNodepool{
		exoplatform20.ClientWithResponses,
	}
}

func (client *Client) ExoPlatform20() *ExoPlatform20 {
	return &ExoPlatform20{
		client.ClientWithResponses,
	}
}
