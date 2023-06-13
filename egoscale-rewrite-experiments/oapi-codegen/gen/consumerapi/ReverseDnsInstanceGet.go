package consumerapi

import (
	"context"

	"github.com/sauterp/egoscale-rewrite-experiments/oapi-codegen/gen/oapi"
)

// THIS CODE IS GENERATED

type ExoPlatform17 struct {
	*oapi.ClientWithResponses
}
type ReverseDnsInstance struct {
	*oapi.ClientWithResponses
}

func (reversednsinstance *ReverseDnsInstance) Get(ctx context.Context, id string, reqEditors ...oapi.RequestEditorFn) *oapi.ReverseDnsRecord {
	resp, err2 := reversednsinstance.ClientWithResponses.GetReverseDnsInstanceWithResponse(ctx, id, reqEditors...)
	if err2 != nil {
		panic(err2)
	}

	return resp.JSON200
}

func (exoplatform17 *ExoPlatform17) ReverseDnsInstance() *ReverseDnsInstance {
	return &ReverseDnsInstance{
		exoplatform17.ClientWithResponses,
	}
}

func (client *Client) ExoPlatform17() *ExoPlatform17 {
	return &ExoPlatform17{
		client.ClientWithResponses,
	}
}
