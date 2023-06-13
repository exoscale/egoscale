package consumerapi

import (
	"context"

	"github.com/sauterp/egoscale-rewrite-experiments/oapi-codegen/gen/oapi"
)

// THIS CODE IS GENERATED

type ExoPlatform22 struct {
	*oapi.ClientWithResponses
}
type SshKey struct {
	*oapi.ClientWithResponses
}

func (sshkey *SshKey) Get(ctx context.Context, name string, reqEditors ...oapi.RequestEditorFn) *oapi.SshKey {
	resp, err2 := sshkey.ClientWithResponses.GetSshKeyWithResponse(ctx, name, reqEditors...)
	if err2 != nil {
		panic(err2)
	}

	return resp.JSON200
}

func (exoplatform22 *ExoPlatform22) SshKey() *SshKey {
	return &SshKey{
		exoplatform22.ClientWithResponses,
	}
}

func (client *Client) ExoPlatform22() *ExoPlatform22 {
	return &ExoPlatform22{
		client.ClientWithResponses,
	}
}
