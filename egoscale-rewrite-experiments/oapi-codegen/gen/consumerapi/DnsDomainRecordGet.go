package consumerapi

import (
	"context"

	"github.com/sauterp/egoscale-rewrite-experiments/oapi-codegen/gen/oapi"
)

// THIS CODE IS GENERATED

type ExoPlatform4 struct {
	*oapi.ClientWithResponses
}
type DnsDomainRecord struct {
	*oapi.ClientWithResponses
}

func (dnsdomainrecord *DnsDomainRecord) Get(ctx context.Context, domainId string, recordId string, reqEditors ...oapi.RequestEditorFn) *oapi.DnsDomainRecord {
	resp, err2 := dnsdomainrecord.ClientWithResponses.GetDnsDomainRecordWithResponse(ctx, domainId, recordId, reqEditors...)
	if err2 != nil {
		panic(err2)
	}

	return resp.JSON200
}

func (exoplatform4 *ExoPlatform4) DnsDomainRecord() *DnsDomainRecord {
	return &DnsDomainRecord{
		exoplatform4.ClientWithResponses,
	}
}

func (client *Client) ExoPlatform4() *ExoPlatform4 {
	return &ExoPlatform4{
		client.ClientWithResponses,
	}
}
