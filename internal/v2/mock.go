package v2

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
	ClientWithResponsesInterface
}

func (c *MockClient) ListLoadBalancersWithResponse(ctx context.Context) (*ListLoadBalancersResponse, error) {
	args := c.Called(ctx)
	return args.Get(0).(*ListLoadBalancersResponse), args.Error(1)
}

func (c *MockClient) GetLoadBalancerWithResponse(ctx context.Context, id string) (*GetLoadBalancerResponse, error) {
	args := c.Called(ctx, id)
	return args.Get(0).(*GetLoadBalancerResponse), args.Error(1)
}

func (c *MockClient) OperationPoller(_, _ string) PollFunc {
	panic("not implemented")
}
