package consumerapi

import (
	"context"

	"github.com/sauterp/egoscale-rewrite-experiments/oapi-codegen/gen/oapi"
)

// THIS CODE IS GENERATED

type ExoPlatform11 struct {
	*oapi.ClientWithResponses
}
type LoadBalancer struct {
	*oapi.ClientWithResponses
}

func (loadbalancer *LoadBalancer) Get(ctx context.Context, id string, reqEditors ...oapi.RequestEditorFn) *oapi.LoadBalancer {
	resp, err2 := loadbalancer.ClientWithResponses.GetLoadBalancerWithResponse(ctx, id, reqEditors...)
	if err2 != nil {
		panic(err2)
	}

	return resp.JSON200
}

func (exoplatform11 *ExoPlatform11) LoadBalancer() *LoadBalancer {
	return &LoadBalancer{
		exoplatform11.ClientWithResponses,
	}
}

func (client *Client) ExoPlatform11() *ExoPlatform11 {
	return &ExoPlatform11{
		client.ClientWithResponses,
	}
}
