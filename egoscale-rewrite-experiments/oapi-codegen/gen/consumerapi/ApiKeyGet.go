package consumerapi

import (
	"context"

	"github.com/sauterp/egoscale-rewrite-experiments/oapi-codegen/gen/oapi"
)

// THIS CODE IS GENERATED

type ExoPlatform2 struct {
	*oapi.ClientWithResponses
}
type ApiKey struct {
	*oapi.ClientWithResponses
}

func (apikey *ApiKey) Get(ctx context.Context, id string, reqEditors ...oapi.RequestEditorFn) *oapi.IamApiKey {
	resp, err2 := apikey.ClientWithResponses.GetApiKeyWithResponse(ctx, id, reqEditors...)
	if err2 != nil {
		panic(err2)
	}

	return resp.JSON200
}

func (exoplatform2 *ExoPlatform2) ApiKey() *ApiKey {
	return &ApiKey{
		exoplatform2.ClientWithResponses,
	}
}

func (client *Client) ExoPlatform2() *ExoPlatform2 {
	return &ExoPlatform2{
		client.ClientWithResponses,
	}
}
