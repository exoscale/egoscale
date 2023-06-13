package consumerapi

import (
	"context"

	"github.com/sauterp/egoscale-rewrite-experiments/oapi-codegen/gen/oapi"
)

// THIS CODE IS GENERATED

type ExoPlatform12 struct {
	*oapi.ClientWithResponses
}
type LoadBalancerService struct {
	*oapi.ClientWithResponses
}

func (loadbalancerservice *LoadBalancerService) Get(ctx context.Context, id string, serviceId string, reqEditors ...oapi.RequestEditorFn) *oapi.LoadBalancerService {
	resp, err2 := loadbalancerservice.ClientWithResponses.GetLoadBalancerServiceWithResponse(ctx, id, serviceId, reqEditors...)
	if err2 != nil {
		panic(err2)
	}

	return resp.JSON200
}

func (exoplatform12 *ExoPlatform12) LoadBalancerService() *LoadBalancerService {
	return &LoadBalancerService{
		exoplatform12.ClientWithResponses,
	}
}

func (client *Client) ExoPlatform12() *ExoPlatform12 {
	return &ExoPlatform12{
		client.ClientWithResponses,
	}
}
