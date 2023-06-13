package consumerapi

import (
	"context"

	"github.com/sauterp/egoscale-rewrite-experiments/oapi-codegen/gen/oapi"
)

// THIS CODE IS GENERATED

type ExoPlatform18 struct {
	*oapi.ClientWithResponses
}
type SecurityGroup struct {
	*oapi.ClientWithResponses
}

func (securitygroup *SecurityGroup) Get(ctx context.Context, id string, reqEditors ...oapi.RequestEditorFn) *oapi.SecurityGroup {
	resp, err2 := securitygroup.ClientWithResponses.GetSecurityGroupWithResponse(ctx, id, reqEditors...)
	if err2 != nil {
		panic(err2)
	}

	return resp.JSON200
}

func (exoplatform18 *ExoPlatform18) SecurityGroup() *SecurityGroup {
	return &SecurityGroup{
		exoplatform18.ClientWithResponses,
	}
}

func (client *Client) ExoPlatform18() *ExoPlatform18 {
	return &ExoPlatform18{
		client.ClientWithResponses,
	}
}
