package consumerapi

import (
	"context"

	"github.com/sauterp/egoscale-rewrite-experiments/oapi-codegen/gen/oapi"
)

// THIS CODE IS GENERATED

type ExoPlatform15 struct {
	*oapi.ClientWithResponses
}
type Quota struct {
	*oapi.ClientWithResponses
}

func (quota *Quota) Get(ctx context.Context, entity string, reqEditors ...oapi.RequestEditorFn) *oapi.Quota {
	resp, err2 := quota.ClientWithResponses.GetQuotaWithResponse(ctx, entity, reqEditors...)
	if err2 != nil {
		panic(err2)
	}

	return resp.JSON200
}

func (exoplatform15 *ExoPlatform15) Quota() *Quota {
	return &Quota{
		exoplatform15.ClientWithResponses,
	}
}

func (client *Client) ExoPlatform15() *ExoPlatform15 {
	return &ExoPlatform15{
		client.ClientWithResponses,
	}
}
