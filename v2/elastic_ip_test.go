package v2

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/exoscale/egoscale/v2/oapi"
)

var (
	testElasticIPDescription                     = new(testSuite).randomString(10)
	testElasticIPID                              = new(testSuite).randomID()
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
	testElasticIPHealthcheckURI                  = new(testSuite).randomString(10)
	testElasticIPHealthcheckTLSSNI               = new(testSuite).randomString(10)
	testElasticIPHealthcheckTLSSkipVerify        = true
)

func (ts *testSuite) TestClient_CreateElasticIP() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
	)

	ts.mock().
		On(
			"CreateElasticIpWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				oapi.CreateElasticIpJSONRequestBody{
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
				},
				args.Get(1),
			)
		}).
		Return(
			&oapi.CreateElasticIpResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: &oapi.Reference{Id: &testElasticIPID},
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: &oapi.Reference{Id: &testElasticIPID},
		State:     &testOperationState,
	})

	ts.mock().
		On("GetElasticIpWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.GetElasticIpResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.ElasticIp{
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
			},
		}, nil)

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

func (ts *testSuite) TestClient_DeleteElasticIP() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		deleted            = false
	)

	ts.mock().
		On(
			"DeleteElasticIpWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testElasticIPID, args.Get(1))
			deleted = true
		}).
		Return(
			&oapi.DeleteElasticIpResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: &oapi.Reference{Id: &testElasticIPID},
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: &oapi.Reference{Id: &testElasticIPID},
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.DeleteElasticIP(context.Background(), testZone, &ElasticIP{ID: &testElasticIPID}))
	ts.Require().True(deleted)
}

func (ts *testSuite) TestClient_FindElasticIP() {
	ts.mock().
		On("ListElasticIpsWithResponse",
			mock.Anything,                 // ctx
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.ListElasticIpsResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &struct {
				ElasticIps *[]oapi.ElasticIp `json:"elastic-ips,omitempty"`
			}{
				ElasticIps: &[]oapi.ElasticIp{{
					Id: &testElasticIPID,
					Ip: &testElasticIPAddress,
				}},
			},
		}, nil)

	ts.mock().
		On("GetElasticIpWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testElasticIPID, args.Get(1))
		}).
		Return(&oapi.GetElasticIpResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.ElasticIp{
				Id: &testElasticIPID,
				Ip: &testElasticIPAddress,
			},
		}, nil)

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

func (ts *testSuite) TestClient_GetElasticIP() {
	ts.mock().
		On("GetElasticIpWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testElasticIPID, args.Get(1))
		}).
		Return(&oapi.GetElasticIpResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.ElasticIp{
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
			},
		}, nil)

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

func (ts *testSuite) TestClient_ListElasticIPs() {
	ts.mock().
		On("ListElasticIpsWithResponse",
			mock.Anything,                 // ctx
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.ListElasticIpsResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &struct {
				ElasticIps *[]oapi.ElasticIp `json:"elastic-ips,omitempty"`
			}{
				ElasticIps: &[]oapi.ElasticIp{{
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
			},
		}, nil)

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

func (ts *testSuite) TestClient_UpdateElasticIP() {
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

	ts.mock().
		On(
			"UpdateElasticIpWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				oapi.UpdateElasticIpJSONRequestBody{
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
				},
				args.Get(2),
			)
			updated = true
		}).
		Return(
			&oapi.UpdateElasticIpResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: &oapi.Reference{Id: &testElasticIPID},
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: &oapi.Reference{Id: &testElasticIPID},
		State:     &testOperationState,
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
