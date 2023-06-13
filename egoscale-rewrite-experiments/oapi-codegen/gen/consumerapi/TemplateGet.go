package consumerapi

import (
	"context"

	"github.com/sauterp/egoscale-rewrite-experiments/oapi-codegen/gen/oapi"
)

// THIS CODE IS GENERATED

type ExoPlatform23 struct {
	*oapi.ClientWithResponses
}
type Template struct {
	*oapi.ClientWithResponses
}

func (template *Template) Get(ctx context.Context, id string, reqEditors ...oapi.RequestEditorFn) *oapi.Template {
	resp, err2 := template.ClientWithResponses.GetTemplateWithResponse(ctx, id, reqEditors...)
	if err2 != nil {
		panic(err2)
	}

	return resp.JSON200
}

func (exoplatform23 *ExoPlatform23) Template() *Template {
	return &Template{
		exoplatform23.ClientWithResponses,
	}
}

func (client *Client) ExoPlatform23() *ExoPlatform23 {
	return &ExoPlatform23{
		client.ClientWithResponses,
	}
}
