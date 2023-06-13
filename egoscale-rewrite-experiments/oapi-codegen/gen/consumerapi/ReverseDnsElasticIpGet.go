package consumerapi

import (
	"context"

	"github.com/sauterp/egoscale-rewrite-experiments/oapi-codegen/gen/oapi"
)

// THIS CODE IS GENERATED

type ExoPlatform16 struct {
	*oapi.ClientWithResponses
}
type ReverseDnsElasticIp struct {
	*oapi.ClientWithResponses
}

func (reversednselasticip *ReverseDnsElasticIp) Get(ctx context.Context, id string, reqEditors ...oapi.RequestEditorFn) *oapi.ReverseDnsRecord {
	resp, err2 := reversednselasticip.ClientWithResponses.GetReverseDnsElasticIpWithResponse(ctx, id, reqEditors...)
	if err2 != nil {
		panic(err2)
	}

	return resp.JSON200
}

func (exoplatform16 *ExoPlatform16) ReverseDnsElasticIp() *ReverseDnsElasticIp {
	return &ReverseDnsElasticIp{
		exoplatform16.ClientWithResponses,
	}
}

func (client *Client) ExoPlatform16() *ExoPlatform16 {
	return &ExoPlatform16{
		client.ClientWithResponses,
	}
}
