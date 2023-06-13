package consumerapi

import (
	"context"

	"github.com/sauterp/egoscale-rewrite-experiments/oapi-codegen/gen/oapi"
)

// THIS CODE IS GENERATED

type ExoPlatform6 struct {
	*oapi.ClientWithResponses
}
type ElasticIp struct {
	*oapi.ClientWithResponses
}

func (elasticip *ElasticIp) Get(ctx context.Context, id string, reqEditors ...oapi.RequestEditorFn) *oapi.ElasticIp {
	resp, err2 := elasticip.ClientWithResponses.GetElasticIpWithResponse(ctx, id, reqEditors...)
	if err2 != nil {
		panic(err2)
	}

	return resp.JSON200
}

func (exoplatform6 *ExoPlatform6) ElasticIp() *ElasticIp {
	return &ElasticIp{
		exoplatform6.ClientWithResponses,
	}
}

func (client *Client) ExoPlatform6() *ExoPlatform6 {
	return &ExoPlatform6{
		client.ClientWithResponses,
	}
}
