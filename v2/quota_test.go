package v2

import (
	"context"
	"fmt"

	"github.com/exoscale/egoscale/v2/oapi"
)

var (
	testQuotaResource       = "instance"
	testQuotaLimit    int64 = 2
	testQuotaUsage    int64 = 1
)

func (ts *clientTestSuite) TestClient_ListQuotas() {
	ts.mockAPIRequest("GET", "/quota", struct {
		Quotas *[]oapi.Quota `json:"quotas,omitempty"`
	}{
		Quotas: &[]oapi.Quota{{
			Limit:    &testQuotaLimit,
			Resource: &testQuotaResource,
			Usage:    &testQuotaUsage,
		}},
	})

	expected := []*Quota{{
		Resource: &testQuotaResource,
		Usage:    &testQuotaUsage,
		Limit:    &testQuotaLimit,
	}}

	actual, err := ts.client.ListQuotas(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_GetQuota() {
	ts.mockAPIRequest("GET", fmt.Sprintf("/quota/%s", testQuotaResource), oapi.Quota{
		Limit:    &testQuotaLimit,
		Resource: &testQuotaResource,
		Usage:    &testQuotaUsage,
	})

	expected := &Quota{
		Resource: &testQuotaResource,
		Usage:    &testQuotaUsage,
		Limit:    &testQuotaLimit,
	}

	actual, err := ts.client.GetQuota(context.Background(), testZone, *expected.Resource)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}
