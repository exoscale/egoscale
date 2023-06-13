package consumerapi

import (
	"context"

	"github.com/sauterp/egoscale-rewrite-experiments/oapi-codegen/gen/oapi"
)

// THIS CODE IS GENERATED

type ExoPlatform13 struct {
	*oapi.ClientWithResponses
}
type Operation struct {
	*oapi.ClientWithResponses
}

func (operation *Operation) Get(ctx context.Context, id string, reqEditors ...oapi.RequestEditorFn) *oapi.Operation {
	resp, err2 := operation.ClientWithResponses.GetOperationWithResponse(ctx, id, reqEditors...)
	if err2 != nil {
		panic(err2)
	}

	return resp.JSON200
}

func (exoplatform13 *ExoPlatform13) Operation() *Operation {
	return &Operation{
		exoplatform13.ClientWithResponses,
	}
}

func (client *Client) ExoPlatform13() *ExoPlatform13 {
	return &ExoPlatform13{
		client.ClientWithResponses,
	}
}
