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
	testNLBID                                          = new(testSuite).randomID()
	testNLBName                                        = new(testSuite).randomString(10)
	testNLBDescription                                 = new(testSuite).randomString(10)
	testNLBCreatedAt, _                                = time.Parse(iso8601Format, "2020-05-26T12:09:42Z")
	testNLBIPAddress                                   = "101.102.103.104"
	testNLBIPAddressP                                  = net.ParseIP("101.102.103.104")
	testNLBLabels                                      = map[string]string{"k1": "v1", "k2": "v2"}
	testNLBState                                       = oapi.LoadBalancerStateRunning
	testNLBServiceID                                   = new(testSuite).randomID()
	testNLBServiceName                                 = new(testSuite).randomString(10)
	testNLBServiceDescription                          = new(testSuite).randomID()
	testNLBServiceInstancePoolID                       = new(testSuite).randomID()
	testNLBServiceProtocol                             = oapi.LoadBalancerServiceProtocolTcp
	testNLBServicePort                          uint16 = 443
	testNLBServiceTargetPort                    uint16 = 8443
	testNLBServiceStrategy                             = oapi.LoadBalancerServiceStrategyRoundRobin
	testNLBServiceState                                = oapi.LoadBalancerServiceStateRunning
	testNLServiceHealthcheckMode                       = oapi.LoadBalancerServiceHealthcheckModeHttps
	testNLBServiceHealthcheckPort               uint16 = 8080
	testNLBServiceHealthcheckInterval           int64  = 10
	testNLBServiceHealthcheckIntervalD                 = time.Duration(testNLBServiceHealthcheckInterval) * time.Second
	testNLBServiceHealthcheckTimeout            int64  = 3
	testNLBServiceHealthcheckTimeoutD                  = time.Duration(testNLBServiceHealthcheckTimeout) * time.Second
	testNLBServiceHealthcheckRetries            int64  = 1
	testNLBServiceHealthcheckURI                       = new(testSuite).randomString(10)
	testNLBServiceHealthcheckTLSSNI                    = new(testSuite).randomString(10)
	testNLBServiceHealthcheckStatus1InstanceIP         = "1.2.3.4"
	testNLBServiceHealthcheckStatus1InstanceIPP        = net.ParseIP("1.2.3.4")
	testNLBServiceHealthcheckStatus1Status             = oapi.LoadBalancerServerStatusStatusSuccess
	testNLBServiceHealthcheckStatus2InstanceIP         = "5.6.7.8"
	testNLBServiceHealthcheckStatus2InstanceIPP        = net.ParseIP("5.6.7.8")
	testNLBServiceHealthcheckStatus2Status             = oapi.LoadBalancerServerStatusStatusSuccess
)

func (ts *testSuite) TestClient_CreateNetworkLoadBalancer() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
	)

	ts.mock().
		On(
			"CreateLoadBalancerWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				oapi.CreateLoadBalancerJSONRequestBody{
					Description: &testNLBDescription,
					Labels:      &oapi.Labels{AdditionalProperties: testNLBLabels},
					Name:        testNLBName,
				},
				args.Get(1),
			)
		}).
		Return(
			&oapi.CreateLoadBalancerResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: &oapi.Reference{Id: &testNLBID},
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: &oapi.Reference{Id: &testNLBID},
		State:     &testOperationState,
	})

	ts.mock().
		On("GetLoadBalancerWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.GetLoadBalancerResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.LoadBalancer{
				CreatedAt:   &testNLBCreatedAt,
				Description: &testNLBDescription,
				Id:          &testNLBID,
				Labels:      &oapi.Labels{AdditionalProperties: testNLBLabels},
				Name:        &testNLBName,
				State:       &testNLBState,
			},
		}, nil)

	expected := &NetworkLoadBalancer{
		CreatedAt:   &testNLBCreatedAt,
		Description: &testNLBDescription,
		ID:          &testNLBID,
		Labels:      &testNLBLabels,
		Name:        &testNLBName,
		Services:    []*NetworkLoadBalancerService{},
		State:       (*string)(&testNLBState),
		Zone:        &testZone,
	}

	actual, err := ts.client.CreateNetworkLoadBalancer(context.Background(), testZone, &NetworkLoadBalancer{
		Description: &testNLBDescription,
		Labels:      &testNLBLabels,
		Name:        &testNLBName,
	})
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_CreateNetworkLoadBalancerService() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
	)

	ts.mock().
		On(
			"AddServiceToLoadBalancerWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testNLBID, args.Get(1))
			ts.Require().Equal(
				oapi.AddServiceToLoadBalancerJSONRequestBody{
					Description: &testNLBServiceDescription,
					Healthcheck: oapi.LoadBalancerServiceHealthcheck{
						Interval: &testNLBServiceHealthcheckInterval,
						Mode:     &testNLServiceHealthcheckMode,
						Port:     func() *int64 { v := int64(testNLBServiceHealthcheckPort); return &v }(),
						Retries:  &testNLBServiceHealthcheckRetries,
						Timeout:  &testNLBServiceHealthcheckTimeout,
						TlsSni:   &testNLBServiceHealthcheckTLSSNI,
						Uri:      &testNLBServiceHealthcheckURI,
					},
					InstancePool: oapi.InstancePool{Id: &testNLBServiceInstancePoolID},
					Name:         testNLBServiceName,
					Port:         int64(testNLBServicePort),
					Protocol:     oapi.AddServiceToLoadBalancerJSONBodyProtocol(testNLBServiceProtocol),
					Strategy:     oapi.AddServiceToLoadBalancerJSONBodyStrategy(testNLBServiceStrategy),
					TargetPort:   int64(testNLBServiceTargetPort),
				},
				args.Get(2),
			)
		}).
		Return(
			&oapi.AddServiceToLoadBalancerResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: &oapi.Reference{Id: &testNLBID},
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: &oapi.Reference{Id: &testNLBID},
		State:     &testOperationState,
	})

	ts.mock().
		On("GetLoadBalancerWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.GetLoadBalancerResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.LoadBalancer{
				CreatedAt:   &testNLBCreatedAt,
				Description: &testNLBDescription,
				Id:          &testNLBID,
				Ip:          &testNLBIPAddress,
				Name:        &testNLBName,
				Services: &[]oapi.LoadBalancerService{{
					Description: &testNLBServiceDescription,
					Healthcheck: &oapi.LoadBalancerServiceHealthcheck{
						Interval: &testNLBServiceHealthcheckInterval,
						Mode:     &testNLServiceHealthcheckMode,
						Port:     func() *int64 { v := int64(testNLBServiceHealthcheckPort); return &v }(),
						Retries:  &testNLBServiceHealthcheckRetries,
						Timeout:  &testNLBServiceHealthcheckTimeout,
						TlsSni:   &testNLBServiceHealthcheckTLSSNI,
						Uri:      &testNLBServiceHealthcheckURI,
					},
					HealthcheckStatus: &[]oapi.LoadBalancerServerStatus{
						{
							PublicIp: &testNLBServiceHealthcheckStatus1InstanceIP,
							Status:   &testNLBServiceHealthcheckStatus1Status,
						},
						{
							PublicIp: &testNLBServiceHealthcheckStatus2InstanceIP,
							Status:   &testNLBServiceHealthcheckStatus2Status,
						},
					},
					Id:           &testNLBServiceID,
					InstancePool: &oapi.InstancePool{Id: &testNLBServiceInstancePoolID},
					Name:         &testNLBServiceName,
					Port:         func() *int64 { v := int64(testNLBServicePort); return &v }(),
					Protocol:     &testNLBServiceProtocol,
					Strategy:     &testNLBServiceStrategy,
					TargetPort:   func() *int64 { v := int64(testNLBServiceTargetPort); return &v }(),
					State:        (*oapi.LoadBalancerServiceState)(&testNLBServiceState),
				}},
				State: &testNLBState,
			},
		}, nil)

	nlb := &NetworkLoadBalancer{
		CreatedAt:   &testNLBCreatedAt,
		Description: &testNLBDescription,
		ID:          &testNLBID,
		IPAddress:   &testNLBIPAddressP,
		Name:        &testNLBName,
		State:       (*string)(&testNLBState),
		Zone:        &testZone,
	}

	expected := &NetworkLoadBalancerService{
		Description: &testNLBServiceDescription,
		Healthcheck: &NetworkLoadBalancerServiceHealthcheck{
			Interval: &testNLBServiceHealthcheckIntervalD,
			Mode:     (*string)(&testNLServiceHealthcheckMode),
			Port:     &testNLBServiceHealthcheckPort,
			Retries:  &testNLBServiceHealthcheckRetries,
			TLSSNI:   &testNLBServiceHealthcheckTLSSNI,
			Timeout:  &testNLBServiceHealthcheckTimeoutD,
			URI:      &testNLBServiceHealthcheckURI,
		},
		HealthcheckStatus: []*NetworkLoadBalancerServerStatus{
			{
				InstanceIP: &testNLBServiceHealthcheckStatus1InstanceIPP,
				Status:     (*string)(&testNLBServiceHealthcheckStatus1Status),
			},
			{
				InstanceIP: &testNLBServiceHealthcheckStatus2InstanceIPP,
				Status:     (*string)(&testNLBServiceHealthcheckStatus2Status),
			},
		},
		ID:             &testNLBServiceID,
		InstancePoolID: &testNLBServiceInstancePoolID,
		Name:           &testNLBServiceName,
		Port:           &testNLBServicePort,
		Protocol:       (*string)(&testNLBServiceProtocol),
		Strategy:       (*string)(&testNLBServiceStrategy),
		TargetPort:     &testNLBServiceTargetPort,
		State:          (*string)(&testNLBServiceState),
	}

	actual, err := ts.client.CreateNetworkLoadBalancerService(context.Background(), testZone, nlb, expected)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_DeleteNetworkLoadBalancer() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		deleted            = false
	)

	ts.mock().
		On(
			"DeleteLoadBalancerWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testNLBID, args.Get(1))
			deleted = true
		}).
		Return(
			&oapi.DeleteLoadBalancerResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: &oapi.Reference{Id: &testNLBID},
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: &oapi.Reference{Id: &testNLBID},
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.DeleteNetworkLoadBalancer(
		context.Background(),
		testZone,
		&NetworkLoadBalancer{ID: &testNLBID},
	))
	ts.Require().True(deleted)
}

func (ts *testSuite) TestClient_DeleteNetworkLoadBalancerService() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		deleted            = false
	)

	ts.mock().
		On(
			"DeleteLoadBalancerServiceWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // serviceId
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testNLBID, args.Get(1))
			ts.Require().Equal(testNLBServiceID, args.Get(2))
			deleted = true
		}).
		Return(
			&oapi.DeleteLoadBalancerServiceResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: &oapi.Reference{Id: &testNLBID},
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: &oapi.Reference{Id: &testNLBID},
		State:     &testOperationState,
	})

	nlb := &NetworkLoadBalancer{
		ID:       &testNLBID,
		Services: []*NetworkLoadBalancerService{{ID: &testNLBServiceID}},
	}

	ts.Require().NoError(ts.client.DeleteNetworkLoadBalancerService(
		context.Background(),
		testZone,
		nlb,
		nlb.Services[0],
	))
	ts.Require().True(deleted)
}

func (ts *testSuite) TestClient_FindNetworkLoadBalancer() {
	ts.mock().
		On("ListLoadBalancersWithResponse",
			mock.Anything,                 // ctx
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.ListLoadBalancersResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &struct {
				LoadBalancers *[]oapi.LoadBalancer `json:"load-balancers,omitempty"`
			}{
				LoadBalancers: &[]oapi.LoadBalancer{{
					CreatedAt: &testNLBCreatedAt,
					Id:        &testNLBID,
					Name:      &testNLBName,
					State:     &testNLBState,
				}},
			},
		}, nil)

	ts.mock().
		On("GetLoadBalancerWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testNLBID, args.Get(1))
		}).
		Return(&oapi.GetLoadBalancerResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.LoadBalancer{
				CreatedAt: &testNLBCreatedAt,
				Id:        &testNLBID,
				Ip:        &testNLBIPAddress,
				Name:      &testNLBName,
				State:     &testNLBState,
			},
		}, nil)

	expected := &NetworkLoadBalancer{
		CreatedAt: &testNLBCreatedAt,
		ID:        &testNLBID,
		IPAddress: &testNLBIPAddressP,
		Name:      &testNLBName,
		Services:  []*NetworkLoadBalancerService{},
		State:     (*string)(&testNLBState),
		Zone:      &testZone,
	}

	actual, err := ts.client.FindNetworkLoadBalancer(context.Background(), testZone, *expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)

	actual, err = ts.client.FindNetworkLoadBalancer(context.Background(), testZone, *expected.Name)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_GetNetworkLoadBalancer() {
	ts.mock().
		On("GetLoadBalancerWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testNLBID, args.Get(1))
		}).
		Return(&oapi.GetLoadBalancerResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.LoadBalancer{
				CreatedAt:   &testNLBCreatedAt,
				Description: &testNLBDescription,
				Id:          &testNLBID,
				Ip:          &testNLBIPAddress,
				Name:        &testNLBName,
				Services: &[]oapi.LoadBalancerService{{
					Description: &testNLBServiceDescription,
					Healthcheck: &oapi.LoadBalancerServiceHealthcheck{
						Interval: &testNLBServiceHealthcheckInterval,
						Mode:     &testNLServiceHealthcheckMode,
						Port:     func() *int64 { v := int64(testNLBServiceHealthcheckPort); return &v }(),
						Retries:  &testNLBServiceHealthcheckRetries,
						Timeout:  &testNLBServiceHealthcheckTimeout,
						Uri:      &testNLBServiceHealthcheckURI,
					},
					HealthcheckStatus: &[]oapi.LoadBalancerServerStatus{
						{
							PublicIp: &testNLBServiceHealthcheckStatus1InstanceIP,
							Status:   &testNLBServiceHealthcheckStatus1Status,
						},
						{
							PublicIp: &testNLBServiceHealthcheckStatus2InstanceIP,
							Status:   &testNLBServiceHealthcheckStatus2Status,
						},
					},
					Id:           &testNLBServiceID,
					InstancePool: &oapi.InstancePool{Id: &testNLBServiceInstancePoolID},
					Name:         &testNLBServiceName,
					Port:         func() *int64 { v := int64(testNLBServicePort); return &v }(),
					Protocol:     &testNLBServiceProtocol,
					State:        (*oapi.LoadBalancerServiceState)(&testNLBServiceState),
					Strategy:     &testNLBServiceStrategy,
					TargetPort:   func() *int64 { v := int64(testNLBServiceTargetPort); return &v }(),
				}},
				State: &testNLBState,
			},
		}, nil)

	expected := &NetworkLoadBalancer{
		CreatedAt:   &testNLBCreatedAt,
		Description: &testNLBDescription,
		ID:          &testNLBID,
		IPAddress:   &testNLBIPAddressP,
		Name:        &testNLBName,
		Services: []*NetworkLoadBalancerService{{
			Description: &testNLBServiceDescription,
			Healthcheck: &NetworkLoadBalancerServiceHealthcheck{
				Interval: &testNLBServiceHealthcheckIntervalD,
				Mode:     (*string)(&testNLServiceHealthcheckMode),
				Port:     &testNLBServiceHealthcheckPort,
				Retries:  &testNLBServiceHealthcheckRetries,
				Timeout:  &testNLBServiceHealthcheckTimeoutD,
				URI:      &testNLBServiceHealthcheckURI,
			},
			HealthcheckStatus: []*NetworkLoadBalancerServerStatus{
				{
					InstanceIP: &testNLBServiceHealthcheckStatus1InstanceIPP,
					Status:     (*string)(&testNLBServiceHealthcheckStatus1Status),
				},
				{
					InstanceIP: &testNLBServiceHealthcheckStatus2InstanceIPP,
					Status:     (*string)(&testNLBServiceHealthcheckStatus2Status),
				},
			},
			ID:             &testNLBServiceID,
			InstancePoolID: &testNLBServiceInstancePoolID,
			Name:           &testNLBServiceName,
			Port:           &testNLBServicePort,
			Protocol:       (*string)(&testNLBServiceProtocol),
			State:          (*string)(&testNLBServiceState),
			Strategy:       (*string)(&testNLBServiceStrategy),
			TargetPort:     &testNLBServiceTargetPort,
		}},
		State: (*string)(&testNLBState),
		Zone:  &testZone,
	}

	actual, err := ts.client.GetNetworkLoadBalancer(context.Background(), testZone, *expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_ListNetworkLoadBalancers() {
	ts.mock().
		On("ListLoadBalancersWithResponse",
			mock.Anything,                 // ctx
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.ListLoadBalancersResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &struct {
				LoadBalancers *[]oapi.LoadBalancer `json:"load-balancers,omitempty"`
			}{
				LoadBalancers: &[]oapi.LoadBalancer{{
					CreatedAt:   &testNLBCreatedAt,
					Description: &testNLBDescription,
					Id:          &testNLBID,
					Name:        &testNLBName,
					Services: &[]oapi.LoadBalancerService{{
						Description: &testNLBServiceDescription,
						Healthcheck: &oapi.LoadBalancerServiceHealthcheck{
							Interval: &testNLBServiceHealthcheckInterval,
							Mode:     &testNLServiceHealthcheckMode,
							Port:     func() *int64 { v := int64(testNLBServiceHealthcheckPort); return &v }(),
							Retries:  &testNLBServiceHealthcheckRetries,
							Timeout:  &testNLBServiceHealthcheckTimeout,
							Uri:      &testNLBServiceHealthcheckURI,
						},
						HealthcheckStatus: &[]oapi.LoadBalancerServerStatus{
							{
								PublicIp: &testNLBServiceHealthcheckStatus1InstanceIP,
								Status:   &testNLBServiceHealthcheckStatus1Status,
							},
							{
								PublicIp: &testNLBServiceHealthcheckStatus2InstanceIP,
								Status:   &testNLBServiceHealthcheckStatus2Status,
							},
						},
						Id:           &testNLBServiceID,
						InstancePool: &oapi.InstancePool{Id: &testNLBServiceInstancePoolID},
						Name:         &testNLBServiceName,
						Port:         func() *int64 { v := int64(testNLBServicePort); return &v }(),
						Protocol:     &testNLBServiceProtocol,
						State:        (*oapi.LoadBalancerServiceState)(&testNLBState),
						Strategy:     &testNLBServiceStrategy,
						TargetPort:   func() *int64 { v := int64(testNLBServiceTargetPort); return &v }(),
					}},
					State: &testNLBState,
				}},
			},
		}, nil)

	expected := []*NetworkLoadBalancer{{
		CreatedAt:   &testNLBCreatedAt,
		Description: &testNLBDescription,
		ID:          &testNLBID,
		Name:        &testNLBName,
		Services: []*NetworkLoadBalancerService{{
			Description: &testNLBServiceDescription,
			Healthcheck: &NetworkLoadBalancerServiceHealthcheck{
				Interval: &testNLBServiceHealthcheckIntervalD,
				Mode:     (*string)(&testNLServiceHealthcheckMode),
				Port:     &testNLBServiceHealthcheckPort,
				Retries:  &testNLBServiceHealthcheckRetries,
				Timeout:  &testNLBServiceHealthcheckTimeoutD,
				URI:      &testNLBServiceHealthcheckURI,
			},
			HealthcheckStatus: []*NetworkLoadBalancerServerStatus{
				{
					InstanceIP: &testNLBServiceHealthcheckStatus1InstanceIPP,
					Status:     (*string)(&testNLBServiceHealthcheckStatus1Status),
				},
				{
					InstanceIP: &testNLBServiceHealthcheckStatus2InstanceIPP,
					Status:     (*string)(&testNLBServiceHealthcheckStatus2Status),
				},
			},
			ID:             &testNLBServiceID,
			InstancePoolID: &testNLBServiceInstancePoolID,
			Name:           &testNLBServiceName,
			Port:           &testNLBServicePort,
			Protocol:       (*string)(&testNLBServiceProtocol),
			State:          (*string)(&testNLBState),
			Strategy:       (*string)(&testNLBServiceStrategy),
			TargetPort:     &testNLBServiceTargetPort,
		}},
		State: (*string)(&testNLBState),
		Zone:  &testZone,
	}}

	actual, err := ts.client.ListNetworkLoadBalancers(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_UpdateNetworkLoadBalancer() {
	var (
		testNLBDescriptionUpdated = testNLBDescription + "-updated"
		testNLBLabelsUpdated      = map[string]string{"k3": "v3"}
		testNLBNameUpdated        = testNLBName + "-updated"
		testOperationID           = ts.randomID()
		testOperationState        = oapi.OperationStateSuccess
		updated                   = false
	)

	ts.mock().
		On(
			"UpdateLoadBalancerWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testNLBID, args.Get(1))
			ts.Require().Equal(
				oapi.UpdateLoadBalancerJSONRequestBody{
					Description: &testNLBDescriptionUpdated,
					Labels:      &oapi.Labels{AdditionalProperties: testNLBLabelsUpdated},
					Name:        &testNLBNameUpdated,
				},
				args.Get(2),
			)
			updated = true
		}).
		Return(
			&oapi.UpdateLoadBalancerResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: &oapi.Reference{Id: &testNLBID},
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: &oapi.Reference{Id: &testNLBID},
		State:     &testOperationState,
	})

	ts.mock().
		On("GetLoadBalancerWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.GetLoadBalancerResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.LoadBalancer{
				CreatedAt:   &testNLBCreatedAt,
				Description: &testNLBDescriptionUpdated,
				Id:          &testNLBID,
				Ip:          &testNLBIPAddress,
				Labels:      &oapi.Labels{AdditionalProperties: testNLBLabelsUpdated},
				Name:        &testNLBNameUpdated,
			},
		}, nil)

	nlbUpdated := NetworkLoadBalancer{
		Description: &testNLBDescriptionUpdated,
		ID:          &testNLBID,
		Labels:      &testNLBLabelsUpdated,
		Name:        &testNLBNameUpdated,
	}

	ts.Require().NoError(ts.client.UpdateNetworkLoadBalancer(context.Background(), testZone, &nlbUpdated))
	ts.Require().True(updated)
}

func (ts *testSuite) TestClient_UpdateNetworkLoadBalancerService() {
	var (
		testNLBServiceNameUpdated                 = testNLBServiceName + "-updated"
		testNLBServiceDescriptionUpdated          = testNLBServiceDescription + "-updated"
		testNLBServiceHealthcheckModeUpdated      = oapi.LoadBalancerServiceHealthcheckModeHttp
		testNLBServiceHealthcheckPortUpdated      = testNLBServiceHealthcheckPort + 1
		testNLBServiceHealthcheckRetriesUpdated   = testNLBServiceHealthcheckRetries + 1
		testNLBServiceHealthcheckTLSSNIUpdated    = ""
		testNLBServiceHealthcheckIntervalUpdated  = testNLBServiceHealthcheckInterval + 1
		testNLBServiceHealthcheckIntervalDUpdated = time.Duration(testNLBServiceHealthcheckIntervalUpdated) * time.Second
		testNLBServiceHealthcheckTimeoutUpdated   = testElasticIPHealthcheckTimeout + 1
		testNLBServiceHealthcheckTimeoutDUpdated  = time.Duration(testNLBServiceHealthcheckTimeoutUpdated) * time.Second
		testNLBServiceHealthcheckURIUpdated       = ""
		testOperationID                           = ts.randomID()
		testOperationState                        = oapi.OperationStateSuccess
		updated                                   = false
	)

	nlb := &NetworkLoadBalancer{
		ID: &testNLBID,
		Services: []*NetworkLoadBalancerService{{
			ID:          &testNLBServiceID,
			Name:        &testNLBServiceName,
			Description: &testNLBServiceDescription,
		}},
	}

	ts.mock().
		On(
			"UpdateLoadBalancerServiceWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // serviceId
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testNLBID, args.Get(1))
			ts.Require().Equal(testNLBServiceID, args.Get(2))
			ts.Require().Equal(
				oapi.UpdateLoadBalancerServiceJSONRequestBody{
					Name:        &testNLBServiceNameUpdated,
					Description: &testNLBServiceDescriptionUpdated,
					Healthcheck: &oapi.LoadBalancerServiceHealthcheck{
						Interval: &testNLBServiceHealthcheckIntervalUpdated,
						Mode:     &testNLBServiceHealthcheckModeUpdated,
						Port:     func() *int64 { v := int64(testNLBServiceHealthcheckPortUpdated); return &v }(),
						Retries:  &testNLBServiceHealthcheckRetriesUpdated,
						Timeout:  &testNLBServiceHealthcheckTimeoutUpdated,
						TlsSni:   &testNLBServiceHealthcheckTLSSNIUpdated,
						Uri:      &testNLBServiceHealthcheckURIUpdated,
					},
				},
				args.Get(3),
			)
			updated = true
		}).
		Return(
			&oapi.UpdateLoadBalancerServiceResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: &oapi.Reference{Id: &testNLBID},
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: &oapi.Reference{Id: &testNLBID},
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.UpdateNetworkLoadBalancerService(context.Background(), testZone, nlb,
		&NetworkLoadBalancerService{
			ID:          nlb.Services[0].ID,
			Name:        &testNLBServiceNameUpdated,
			Description: &testNLBServiceDescriptionUpdated,
			Healthcheck: &NetworkLoadBalancerServiceHealthcheck{
				Interval: &testNLBServiceHealthcheckIntervalDUpdated,
				Mode:     (*string)(&testNLBServiceHealthcheckModeUpdated),
				Port:     &testNLBServiceHealthcheckPortUpdated,
				Retries:  &testNLBServiceHealthcheckRetriesUpdated,
				TLSSNI:   &testNLBServiceHealthcheckTLSSNIUpdated,
				Timeout:  &testNLBServiceHealthcheckTimeoutDUpdated,
				URI:      &testNLBServiceHealthcheckURIUpdated,
			},
		}))
	ts.Require().True(updated)
}
