package consumerapi

import (
	"context"

	"github.com/sauterp/egoscale-rewrite-experiments/oapi-codegen/gen/oapi"
)

// THIS CODE IS GENERATED

type ExoPlatform21 struct {
	*oapi.ClientWithResponses
}
type Snapshot struct {
	*oapi.ClientWithResponses
}

func (snapshot *Snapshot) Get(ctx context.Context, id string, reqEditors ...oapi.RequestEditorFn) *oapi.Snapshot {
	resp, err2 := snapshot.ClientWithResponses.GetSnapshotWithResponse(ctx, id, reqEditors...)
	if err2 != nil {
		panic(err2)
	}

	return resp.JSON200
}

func (exoplatform21 *ExoPlatform21) Snapshot() *Snapshot {
	return &Snapshot{
		exoplatform21.ClientWithResponses,
	}
}

func (client *Client) ExoPlatform21() *ExoPlatform21 {
	return &ExoPlatform21{
		client.ClientWithResponses,
	}
}
