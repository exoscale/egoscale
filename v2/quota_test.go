package v2

import (
	"context"
	"net/http"

	"github.com/exoscale/egoscale/v2/oapi"
	"github.com/stretchr/testify/mock"
)

var (
	testQuotaResource       = "instance"
	testQuotaLimit    int64 = 2
	testQuotaUsage    int64 = 1
)

func (ts *testSuite) TestClient_ListQuotas() {
	ts.mock().
		On("ListQuotasWithResponse",
			mock.Anything,                 // ctx
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.ListQuotasResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &struct {
				Quotas *[]oapi.Quota `json:"quotas,omitempty"`
			}{
				Quotas: &[]oapi.Quota{{
					Limit:    &testQuotaLimit,
					Resource: &testQuotaResource,
					Usage:    &testQuotaUsage,
				}},
			},
		}, nil)

	expected := []*Quota{{
		Resource: &testQuotaResource,
		Usage:    &testQuotaUsage,
		Limit:    &testQuotaLimit,
	}}

	actual, err := ts.client.ListQuotas(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_GetQuota() {
	ts.mock().
		On("GetQuotaWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.GetQuotaResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.Quota{
				Limit:    &testQuotaLimit,
				Resource: &testQuotaResource,
				Usage:    &testQuotaUsage,
			},
		}, nil)

	expected := &Quota{
		Resource: &testQuotaResource,
		Usage:    &testQuotaUsage,
		Limit:    &testQuotaLimit,
	}

	actual, err := ts.client.GetQuota(context.Background(), testZone, *expected.Resource)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}
