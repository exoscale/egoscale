package consumerapi

import (
	"context"

	"github.com/sauterp/egoscale-rewrite-experiments/oapi-codegen/gen/oapi"
)

// THIS CODE IS GENERATED

type ExoPlatform5 struct {
	*oapi.ClientWithResponses
}
type DnsDomain struct {
	*oapi.ClientWithResponses
}

func (dnsdomain *DnsDomain) Get(ctx context.Context, id string, reqEditors ...oapi.RequestEditorFn) *oapi.DnsDomain {
	resp, err2 := dnsdomain.ClientWithResponses.GetDnsDomainWithResponse(ctx, id, reqEditors...)
	if err2 != nil {
		panic(err2)
	}

	return resp.JSON200
}

func (exoplatform5 *ExoPlatform5) DnsDomain() *DnsDomain {
	return &DnsDomain{
		exoplatform5.ClientWithResponses,
	}
}

func (client *Client) ExoPlatform5() *ExoPlatform5 {
	return &ExoPlatform5{
		client.ClientWithResponses,
	}
}
