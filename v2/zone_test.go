package v2

import (
	"context"
	"net/http"

	"github.com/jarcoal/httpmock"

	papi "github.com/exoscale/egoscale/v2/internal/public-api"
)

func (ts *clientTestSuite) TestClient_ListZones() {
	testZones := []string{
		"at-vie-1",
		"bg-sof-1",
		"ch-dk-2",
		"ch-gva-2",
		"de-fra-1",
		"de-muc-1",
	}

	httpmock.RegisterResponder("GET", "/zone",
		func(_ *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(http.StatusOK, struct {
				Zones *[]papi.Zone `json:"zones,omitempty"`
			}{
				Zones: func() *[]papi.Zone {
					zones := make([]papi.Zone, len(testZones))
					for i := range testZones {
						name := testZones[i]
						zones[i] = papi.Zone{Name: (*papi.ZoneName)(&name)}
					}
					return &zones
				}(),
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}
			return resp, nil
		})

	expected := testZones
	actual, err := ts.client.ListZones(context.Background())
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}
