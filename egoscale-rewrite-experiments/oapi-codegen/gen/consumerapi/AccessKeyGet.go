package consumerapi

import (
	"context"

	"github.com/sauterp/egoscale-rewrite-experiments/oapi-codegen/gen/oapi"
)

// THIS CODE IS GENERATED

type ExoPlatform0 struct {
	*oapi.ClientWithResponses
}
type AccessKey struct {
	*oapi.ClientWithResponses
}

func (accesskey *AccessKey) Get(ctx context.Context, key string, reqEditors ...oapi.RequestEditorFn) *oapi.AccessKey {
	resp, err2 := accesskey.ClientWithResponses.GetAccessKeyWithResponse(ctx, key, reqEditors...)
	if err2 != nil {
		panic(err2)
	}

	return resp.JSON200
}

func (exoplatform0 *ExoPlatform0) AccessKey() *AccessKey {
	return &AccessKey{
		exoplatform0.ClientWithResponses,
	}
}

func (client *Client) ExoPlatform0() *ExoPlatform0 {
	return &ExoPlatform0{
		client.ClientWithResponses,
	}
}
