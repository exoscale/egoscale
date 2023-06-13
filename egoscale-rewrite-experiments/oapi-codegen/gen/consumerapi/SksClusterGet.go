package consumerapi

import (
	"context"

	"github.com/sauterp/egoscale-rewrite-experiments/oapi-codegen/gen/oapi"
)

// THIS CODE IS GENERATED

type ExoPlatform19 struct {
	*oapi.ClientWithResponses
}
type SksCluster struct {
	*oapi.ClientWithResponses
}

func (skscluster *SksCluster) Get(ctx context.Context, id string, reqEditors ...oapi.RequestEditorFn) *oapi.SksCluster {
	resp, err2 := skscluster.ClientWithResponses.GetSksClusterWithResponse(ctx, id, reqEditors...)
	if err2 != nil {
		panic(err2)
	}

	return resp.JSON200
}

func (exoplatform19 *ExoPlatform19) SksCluster() *SksCluster {
	return &SksCluster{
		exoplatform19.ClientWithResponses,
	}
}

func (client *Client) ExoPlatform19() *ExoPlatform19 {
	return &ExoPlatform19{
		client.ClientWithResponses,
	}
}
