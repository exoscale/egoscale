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
	testNLBID           = new(testSuite).randomID()
	testNLBName         = new(testSuite).randomString(10)
	testNLBDescription  = new(testSuite).randomString(10)
	testNLBCreatedAt, _ = time.Parse(iso8601Format, "2020-05-26T12:09:42Z")
	testNLBIPAddress    = "101.102.103.104"
	testNLBIPAddressP   = net.ParseIP("101.102.103.104")
	testNLBLabels       = map[string]string{"k1": "v1", "k2": "v2"}
	testNLBState        = oapi.LoadBalancerStateRunning
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
					Id: &testOperationID,
					Reference: &struct {
						Command *string `json:"command,omitempty"`
						Id      *string `json:"id,omitempty"`
						Link    *string `json:"link,omitempty"`
					}{Id: &testNLBID},
					State: &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id: &testOperationID,
		Reference: &struct {
			Command *string `json:"command,omitempty"`
			Id      *string `json:"id,omitempty"`
			Link    *string `json:"link,omitempty"`
		}{Id: &testNLBID},
		State: &testOperationState,
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
					Id: &testOperationID,
					Reference: &struct {
						Command *string `json:"command,omitempty"`
						Id      *string `json:"id,omitempty"`
						Link    *string `json:"link,omitempty"`
					}{Id: &testNLBID},
					State: &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id: &testOperationID,
		Reference: &struct {
			Command *string `json:"command,omitempty"`
			Id      *string `json:"id,omitempty"`
			Link    *string `json:"link,omitempty"`
		}{Id: &testNLBID},
		State: &testOperationState,
	})

	ts.Require().NoError(ts.client.DeleteNetworkLoadBalancer(
		context.Background(),
		testZone,
		&NetworkLoadBalancer{ID: &testNLBID},
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
					Id: &testOperationID,
					Reference: &struct {
						Command *string `json:"command,omitempty"`
						Id      *string `json:"id,omitempty"`
						Link    *string `json:"link,omitempty"`
					}{Id: &testNLBID},
					State: &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id: &testOperationID,
		Reference: &struct {
			Command *string `json:"command,omitempty"`
			Id      *string `json:"id,omitempty"`
			Link    *string `json:"link,omitempty"`
		}{Id: &testNLBID},
		State: &testOperationState,
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
