package v2

import (
	"context"
	"net/http"

	"github.com/stretchr/testify/mock"

	"github.com/exoscale/egoscale/v2/oapi"
)

func (ts *testSuite) TestClient_ListZones() {
	testZones := []string{
		"at-vie-1",
		"bg-sof-1",
		"ch-dk-2",
		"ch-gva-2",
		"de-fra-1",
		"de-muc-1",
	}

	ts.mock().
		On("ListZonesWithResponse",
			mock.Anything,                 // ctx
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.ListZonesResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &struct {
				Zones *[]oapi.Zone `json:"zones,omitempty"`
			}{
				Zones: func() *[]oapi.Zone {
					zones := make([]oapi.Zone, len(testZones))
					for i := range testZones {
						name := testZones[i]
						zones[i] = oapi.Zone{Name: (*oapi.ZoneName)(&name)}
					}
					return &zones
				}(),
			},
		}, nil)

	expected := testZones
	actual, err := ts.client.ListZones(context.Background())
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}
