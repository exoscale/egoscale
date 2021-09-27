package v2

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/jarcoal/httpmock"

	"github.com/exoscale/egoscale/v2/oapi"
)

var (
	testElasticIPDescription                     = new(clientTestSuite).randomString(10)
	testElasticIPID                              = new(clientTestSuite).randomID()
	testElasticIPAddress                         = "1.2.3.4"
	testElasticIPAddressP                        = net.ParseIP(testElasticIPAddress)
	testElasticIPHealthcheckMode                 = "https"
	testElasticIPHealthcheckPort          uint16 = 8080
	testElasticIPHealthcheckInterval      int64  = 10
	testElasticIPHealthcheckIntervalD            = time.Duration(testElasticIPHealthcheckInterval) * time.Second
	testElasticIPHealthcheckTimeout       int64  = 3
	testElasticIPHealthcheckTimeoutD             = time.Duration(testElasticIPHealthcheckTimeout) * time.Second
	testElasticIPHealthcheckStrikesFail   int64  = 1
	testElasticIPHealthcheckStrikesOK     int64  = 1
	testElasticIPHealthcheckURI                  = new(clientTestSuite).randomString(10)
	testElasticIPHealthcheckTLSSNI               = new(clientTestSuite).randomString(10)
	testElasticIPHealthcheckTLSSkipVerify        = true
)

func (ts *clientTestSuite) TestClient_CreateElasticIP() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
	)

	httpmock.RegisterResponder("POST", "/elastic-ip",
		func(req *http.Request) (*http.Response, error) {
			var actual oapi.CreateElasticIpJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := oapi.CreateElasticIpJSONRequestBody{
				Description: &testElasticIPDescription,
				Healthcheck: &oapi.ElasticIpHealthcheck{
					Interval:      &testElasticIPHealthcheckInterval,
					Mode:          oapi.ElasticIpHealthcheckMode(testElasticIPHealthcheckMode),
					Port:          int64(testElasticIPHealthcheckPort),
					StrikesFail:   &testElasticIPHealthcheckStrikesFail,
					StrikesOk:     &testElasticIPHealthcheckStrikesOK,
					Timeout:       &testElasticIPHealthcheckTimeout,
					TlsSni:        &testElasticIPHealthcheckTLSSNI,
					TlsSkipVerify: &testElasticIPHealthcheckTLSSkipVerify,
					Uri:           &testElasticIPHealthcheckURI,
				},
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testElasticIPID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testElasticIPID},
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/elastic-ip/%s", testElasticIPID), oapi.ElasticIp{
		Description: &testElasticIPDescription,
		Healthcheck: &oapi.ElasticIpHealthcheck{
			Interval:      &testElasticIPHealthcheckInterval,
			Mode:          oapi.ElasticIpHealthcheckMode(testElasticIPHealthcheckMode),
			Port:          int64(testElasticIPHealthcheckPort),
			StrikesFail:   &testElasticIPHealthcheckStrikesFail,
			StrikesOk:     &testElasticIPHealthcheckStrikesOK,
			Timeout:       &testElasticIPHealthcheckTimeout,
			TlsSni:        &testElasticIPHealthcheckTLSSNI,
			TlsSkipVerify: &testElasticIPHealthcheckTLSSkipVerify,
			Uri:           &testElasticIPHealthcheckURI,
		},
		Id: &testElasticIPID,
		Ip: &testElasticIPAddress,
	})

	expected := &ElasticIP{
		Description: &testElasticIPDescription,
		Healthcheck: &ElasticIPHealthcheck{
			Interval:      &testElasticIPHealthcheckIntervalD,
			Mode:          &testElasticIPHealthcheckMode,
			Port:          &testElasticIPHealthcheckPort,
			StrikesFail:   &testElasticIPHealthcheckStrikesFail,
			StrikesOK:     &testElasticIPHealthcheckStrikesOK,
			TLSSNI:        &testElasticIPHealthcheckTLSSNI,
			TLSSkipVerify: &testElasticIPHealthcheckTLSSkipVerify,
			Timeout:       &testElasticIPHealthcheckTimeoutD,
			URI:           &testElasticIPHealthcheckURI,
		},
		ID:        &testElasticIPID,
		IPAddress: &testElasticIPAddressP,
		Zone:      &testZone,
	}

	actual, err := ts.client.CreateElasticIP(context.Background(), testZone, &ElasticIP{
		Description: &testElasticIPDescription,
		Healthcheck: &ElasticIPHealthcheck{
			Interval:      &testElasticIPHealthcheckIntervalD,
			Mode:          &testElasticIPHealthcheckMode,
			Port:          &testElasticIPHealthcheckPort,
			StrikesFail:   &testElasticIPHealthcheckStrikesFail,
			StrikesOK:     &testElasticIPHealthcheckStrikesOK,
			TLSSNI:        &testElasticIPHealthcheckTLSSNI,
			TLSSkipVerify: &testElasticIPHealthcheckTLSSkipVerify,
			Timeout:       &testElasticIPHealthcheckTimeoutD,
			URI:           &testElasticIPHealthcheckURI,
		},
		ID: &testElasticIPID,
	})
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_DeleteElasticIP() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		deleted            = false
	)

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("/elastic-ip/%s", testElasticIPID),
		func(req *http.Request) (*http.Response, error) {
			deleted = true

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testElasticIPID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testElasticIPID},
	})

	ts.Require().NoError(ts.client.DeleteElasticIP(context.Background(), testZone, &ElasticIP{ID: &testElasticIPID}))
	ts.Require().True(deleted)
}

func (ts *clientTestSuite) TestClient_FindElasticIP() {
	ts.mockAPIRequest("GET", "/elastic-ip", struct {
		ElasticIPs *[]oapi.ElasticIp `json:"elastic-ips,omitempty"`
	}{
		ElasticIPs: &[]oapi.ElasticIp{{
			Id: &testElasticIPID,
			Ip: &testElasticIPAddress,
		}},
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/elastic-ip/%s", testElasticIPID), oapi.ElasticIp{
		Id: &testElasticIPID,
		Ip: &testElasticIPAddress,
	})

	expected := &ElasticIP{
		ID:        &testElasticIPID,
		IPAddress: &testElasticIPAddressP,
		Zone:      &testZone,
	}

	actual, err := ts.client.FindElasticIP(context.Background(), testZone, *expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)

	actual, err = ts.client.FindElasticIP(context.Background(), testZone, expected.IPAddress.String())
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_GetElasticIP() {
	ts.mockAPIRequest("GET", fmt.Sprintf("/elastic-ip/%s", testElasticIPID), oapi.ElasticIp{
		Description: &testElasticIPDescription,
		Healthcheck: &oapi.ElasticIpHealthcheck{
			Interval:      &testElasticIPHealthcheckInterval,
			Mode:          oapi.ElasticIpHealthcheckMode(testElasticIPHealthcheckMode),
			Port:          int64(testElasticIPHealthcheckPort),
			StrikesFail:   &testElasticIPHealthcheckStrikesFail,
			StrikesOk:     &testElasticIPHealthcheckStrikesOK,
			Timeout:       &testElasticIPHealthcheckTimeout,
			TlsSni:        &testElasticIPHealthcheckTLSSNI,
			TlsSkipVerify: &testElasticIPHealthcheckTLSSkipVerify,
			Uri:           &testElasticIPHealthcheckURI,
		},
		Id: &testElasticIPID,
		Ip: &testElasticIPAddress,
	})

	expected := &ElasticIP{
		Description: &testElasticIPDescription,
		Healthcheck: &ElasticIPHealthcheck{
			Interval:      &testElasticIPHealthcheckIntervalD,
			Mode:          &testElasticIPHealthcheckMode,
			Port:          &testElasticIPHealthcheckPort,
			StrikesFail:   &testElasticIPHealthcheckStrikesFail,
			StrikesOK:     &testElasticIPHealthcheckStrikesOK,
			TLSSNI:        &testElasticIPHealthcheckTLSSNI,
			TLSSkipVerify: &testElasticIPHealthcheckTLSSkipVerify,
			Timeout:       &testElasticIPHealthcheckTimeoutD,
			URI:           &testElasticIPHealthcheckURI,
		},
		ID:        &testElasticIPID,
		IPAddress: &testElasticIPAddressP,
		Zone:      &testZone,
	}

	actual, err := ts.client.GetElasticIP(context.Background(), testZone, *expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_ListElasticIPs() {
	ts.mockAPIRequest("GET", "/elastic-ip", struct {
		ElasticIPs *[]oapi.ElasticIp `json:"elastic-ips,omitempty"`
	}{
		ElasticIPs: &[]oapi.ElasticIp{{
			Description: &testElasticIPDescription,
			Healthcheck: &oapi.ElasticIpHealthcheck{
				Interval:      &testElasticIPHealthcheckInterval,
				Mode:          oapi.ElasticIpHealthcheckMode(testElasticIPHealthcheckMode),
				Port:          int64(testElasticIPHealthcheckPort),
				StrikesFail:   &testElasticIPHealthcheckStrikesFail,
				StrikesOk:     &testElasticIPHealthcheckStrikesOK,
				Timeout:       &testElasticIPHealthcheckTimeout,
				TlsSni:        &testElasticIPHealthcheckTLSSNI,
				TlsSkipVerify: &testElasticIPHealthcheckTLSSkipVerify,
				Uri:           &testElasticIPHealthcheckURI,
			},
			Id: &testElasticIPID,
			Ip: &testElasticIPAddress,
		}},
	})

	expected := []*ElasticIP{{
		Description: &testElasticIPDescription,
		Healthcheck: &ElasticIPHealthcheck{
			Interval:      &testElasticIPHealthcheckIntervalD,
			Mode:          &testElasticIPHealthcheckMode,
			Port:          &testElasticIPHealthcheckPort,
			StrikesFail:   &testElasticIPHealthcheckStrikesFail,
			StrikesOK:     &testElasticIPHealthcheckStrikesOK,
			TLSSNI:        &testElasticIPHealthcheckTLSSNI,
			TLSSkipVerify: &testElasticIPHealthcheckTLSSkipVerify,
			Timeout:       &testElasticIPHealthcheckTimeoutD,
			URI:           &testElasticIPHealthcheckURI,
		},
		ID:        &testElasticIPID,
		IPAddress: &testElasticIPAddressP,
		Zone:      &testZone,
	}}

	actual, err := ts.client.ListElasticIPs(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_UpdateElasticIP() {
	var (
		testElasticIPDescriptionUpdated              = testElasticIPDescription + "-updated"
		testElasticIPHealthcheckModeUpdated          = oapi.ElasticIpHealthcheckModeTcp
		testElasticIPHealthcheckPortUpdated          = testElasticIPHealthcheckPort + 1
		testElasticIPHealthcheckIntervalUpdated      = testElasticIPHealthcheckInterval + 1
		testElasticIPHealthcheckIntervalDUpdated     = time.Duration(testElasticIPHealthcheckIntervalUpdated) * time.Second
		testElasticIPHealthcheckTimeoutUpdated       = testElasticIPHealthcheckTimeout + 1
		testElasticIPHealthcheckTimeoutDUpdated      = time.Duration(testElasticIPHealthcheckTimeoutUpdated) * time.Second
		testElasticIPHealthcheckStrikesFailUpdated   = testElasticIPHealthcheckStrikesFail + 1
		testElasticIPHealthcheckStrikesOKUpdated     = testElasticIPHealthcheckStrikesOK + 1
		testElasticIPHealthcheckTLSSkipVerifyUpdated = false
		testOperationID                              = ts.randomID()
		testOperationState                           = oapi.OperationStateSuccess
		updated                                      = false
	)

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/elastic-ip/%s", testElasticIPID),
		func(req *http.Request) (*http.Response, error) {
			updated = true

			var actual oapi.UpdateElasticIpJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := oapi.UpdateElasticIpJSONRequestBody{
				Description: &testElasticIPDescriptionUpdated,
				Healthcheck: &oapi.ElasticIpHealthcheck{
					Interval:      &testElasticIPHealthcheckIntervalUpdated,
					Mode:          testElasticIPHealthcheckModeUpdated,
					Port:          int64(testElasticIPHealthcheckPortUpdated),
					StrikesFail:   &testElasticIPHealthcheckStrikesFailUpdated,
					StrikesOk:     &testElasticIPHealthcheckStrikesOKUpdated,
					Timeout:       &testElasticIPHealthcheckTimeoutUpdated,
					TlsSkipVerify: &testElasticIPHealthcheckTLSSkipVerifyUpdated,
				},
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testElasticIPID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testElasticIPID},
	})

	ts.Require().NoError(ts.client.UpdateElasticIP(context.Background(), testZone, &ElasticIP{
		Description: &testElasticIPDescriptionUpdated,
		Healthcheck: &ElasticIPHealthcheck{
			Interval:      &testElasticIPHealthcheckIntervalDUpdated,
			Mode:          (*string)(&testElasticIPHealthcheckModeUpdated),
			Port:          &testElasticIPHealthcheckPortUpdated,
			StrikesFail:   &testElasticIPHealthcheckStrikesFailUpdated,
			StrikesOK:     &testElasticIPHealthcheckStrikesOKUpdated,
			Timeout:       &testElasticIPHealthcheckTimeoutDUpdated,
			TLSSkipVerify: &testElasticIPHealthcheckTLSSkipVerifyUpdated,
		},
		ID:        &testElasticIPID,
		IPAddress: &testElasticIPAddressP,
		Zone:      &testZone,
	}))
	ts.Require().True(updated)
}
