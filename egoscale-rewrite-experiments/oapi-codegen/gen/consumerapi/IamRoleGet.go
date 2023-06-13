package consumerapi

import (
	"context"

	"github.com/sauterp/egoscale-rewrite-experiments/oapi-codegen/gen/oapi"
)

// THIS CODE IS GENERATED

type ExoPlatform7 struct {
	*oapi.ClientWithResponses
}
type IamRole struct {
	*oapi.ClientWithResponses
}

func (iamrole *IamRole) Get(ctx context.Context, id string, reqEditors ...oapi.RequestEditorFn) *oapi.IamRole {
	resp, err2 := iamrole.ClientWithResponses.GetIamRoleWithResponse(ctx, id, reqEditors...)
	if err2 != nil {
		panic(err2)
	}

	return resp.JSON200
}

func (exoplatform7 *ExoPlatform7) IamRole() *IamRole {
	return &IamRole{
		exoplatform7.ClientWithResponses,
	}
}

func (client *Client) ExoPlatform7() *ExoPlatform7 {
	return &ExoPlatform7{
		client.ClientWithResponses,
	}
}
