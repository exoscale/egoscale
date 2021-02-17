package v2

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/jarcoal/httpmock"

	papi "github.com/exoscale/egoscale/v2/internal/public-api"
)

const iso8601Format = "2006-01-02T15:04:05Z"

var (
	testNLBID                                        = "9381ab81-59ea-4215-a5a5-781db2fabfe9"
	testNLBName                                      = "test-nlb-name"
	testNLBDescription                               = "test-nlb-description"
	testNLBCreatedAt, _                              = time.Parse(iso8601Format, "2020-05-26T12:09:42Z")
	testNLBIPAddress                                 = "101.102.103.104"
	testNLBState                                     = "running"
	testNLBServiceID                                 = "22d8118c-6585-4dbd-bc67-491256da4e9a"
	testNLBServiceName                               = "test-svc-name"
	testNLBServiceDescription                        = "test-svc-description"
	testNLBServiceInstancePoolID                     = "28c7dd5f-f26f-4ca8-b391-c53795fcee3b"
	testNLBServiceProtocol                           = "tcp"
	testNLBServicePort                         int64 = 80
	testNLBServiceTargetPort                   int64 = 8080
	testNLBServiceStrategy                           = "round-robin"
	testNLBServiceState                              = "running"
	testNLServiceHealthcheckMode                     = "http"
	testNLBServiceHealthcheckPort              int64 = 8080
	testNLBServiceHealthcheckInterval          int64 = 10
	testNLBServiceHealthcheckTimeout           int64 = 3
	testNLBServiceHealthcheckRetries           int64 = 1
	testNLBServiceHealthcheckURI                     = "/health"
	testNLBServiceHealthcheckTLSSNI                  = "example.net"
	testNLBServiceHealthcheckStatus1InstanceIP       = "1.2.3.4"
	testNLBServiceHealthcheckStatus1Status           = "success"
	testNLBServiceHealthcheckStatus2InstanceIP       = "5.6.7.8"
	testNLBServiceHealthcheckStatus2Status           = "success"
)

func (ts *clientTestSuite) TestNetworkLoadBalancer_AddService() {
	var (
		testOperationID    = "08302193-c7e3-42a6-9b3d-da0b2a536577"
		testOperationState = "success"
		err                error
	)

	httpmock.RegisterResponder("POST", "/load-balancer/"+testNLBID+"/service",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testNLBID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "/operation/"+testOperationID,
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testNLBID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "/load-balancer/"+testNLBID, // nolint:dupl
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.LoadBalancer{
				Id:          &testNLBID,
				Name:        &testNLBName,
				Description: &testNLBDescription,
				CreatedAt:   &testNLBCreatedAt,
				Ip:          &testNLBIPAddress,
				Services: &[]papi.LoadBalancerService{{
					Id:           &testNLBServiceID,
					Name:         &testNLBServiceName,
					Description:  &testNLBServiceDescription,
					InstancePool: &papi.InstancePool{Id: &testNLBServiceInstancePoolID},
					Protocol:     &testNLBServiceProtocol,
					Port:         &testNLBServicePort,
					TargetPort:   &testNLBServiceTargetPort,
					Strategy:     &testNLBServiceStrategy,
					Healthcheck: &papi.Healthcheck{
						Mode:     testNLServiceHealthcheckMode,
						Interval: &testNLBServiceHealthcheckInterval,
						Port:     testNLBServiceHealthcheckPort,
						Uri:      &testNLBServiceHealthcheckURI,
						TlsSni:   &testNLBServiceHealthcheckTLSSNI,
						Timeout:  &testNLBServiceHealthcheckTimeout,
						Retries:  &testNLBServiceHealthcheckRetries,
					},
					HealthcheckStatus: &[]papi.LoadBalancerServerStatus{
						{
							PublicIp: &testNLBServiceHealthcheckStatus1InstanceIP,
							Status:   &testNLBServiceHealthcheckStatus1Status,
						},
						{
							PublicIp: &testNLBServiceHealthcheckStatus2InstanceIP,
							Status:   &testNLBServiceHealthcheckStatus2Status,
						},
					},
					State: &testNLBServiceState,
				}},
				State: &testNLBState,
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}
			return resp, nil
		})

	nlb := &NetworkLoadBalancer{
		ID:          testNLBID,
		Name:        testNLBName,
		Description: testNLBDescription,
		CreatedAt:   testNLBCreatedAt,
		IPAddress:   net.ParseIP(testNLBIPAddress),
		State:       testNLBState,

		c: ts.client,
	}

	expected := &NetworkLoadBalancerService{
		ID:             testNLBServiceID,
		Name:           testNLBServiceName,
		Description:    testNLBServiceDescription,
		InstancePoolID: testNLBServiceInstancePoolID,
		Protocol:       testNLBServiceProtocol,
		Port:           uint16(testNLBServicePort),
		TargetPort:     uint16(testNLBServiceTargetPort),
		Strategy:       testNLBServiceStrategy,
		Healthcheck: NetworkLoadBalancerServiceHealthcheck{
			Mode:     testNLServiceHealthcheckMode,
			Port:     uint16(testNLBServiceHealthcheckPort),
			Interval: time.Duration(testNLBServiceHealthcheckInterval) * time.Second,
			Timeout:  time.Duration(testNLBServiceHealthcheckTimeout) * time.Second,
			Retries:  testNLBServiceHealthcheckRetries,
			URI:      testNLBServiceHealthcheckURI,
			TLSSNI:   testNLBServiceHealthcheckTLSSNI,
		},
		HealthcheckStatus: []*NetworkLoadBalancerServerStatus{
			{
				InstanceIP: net.ParseIP(testNLBServiceHealthcheckStatus1InstanceIP),
				Status:     testNLBServiceHealthcheckStatus1Status,
			},
			{
				InstanceIP: net.ParseIP(testNLBServiceHealthcheckStatus2InstanceIP),
				Status:     testNLBServiceHealthcheckStatus2Status,
			},
		},
		State: testNLBServiceState,
	}

	actual, err := nlb.AddService(context.Background(), expected)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestNetworkLoadBalancer_UpdateService() {
	var (
		testNLBServiceNameUpdated        = testNLBServiceName + "-updated"
		testNLBServiceDescriptionUpdated = testNLBServiceDescription + "-updated"
		testOperationID                  = "08302193-c7e3-42a6-9b3d-da0b2a536577"
		testOperationState               = "success"
	)

	nlb := &NetworkLoadBalancer{
		ID:   testNLBID,
		c:    ts.client,
		zone: testZone,

		Services: []*NetworkLoadBalancerService{
			{
				ID:          testNLBServiceID,
				Name:        testNLBServiceName,
				Description: testNLBServiceDescription,
			},
		},
	}

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/load-balancer/%s/service/%s",
		nlb.ID,
		nlb.Services[0].ID),
		func(req *http.Request) (*http.Response, error) {
			var actual papi.UpdateLoadBalancerJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)
			expected := papi.UpdateLoadBalancerJSONRequestBody{
				Name:        &testNLBServiceNameUpdated,
				Description: &testNLBServiceDescriptionUpdated,
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testNLBServiceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	httpmock.RegisterResponder("GET", "/operation/"+testOperationID,
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testNLBServiceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}
			return resp, nil
		})

	nlbServiceUpdated := NetworkLoadBalancerService{
		ID:          nlb.Services[0].ID,
		Name:        testNLBServiceNameUpdated,
		Description: testNLBServiceDescriptionUpdated,
	}
	ts.Require().NoError(nlb.UpdateService(context.Background(), &nlbServiceUpdated))
}

func (ts *clientTestSuite) TestNetworkLoadBalancer_DeleteService() {
	var (
		testOperationID    = "08302193-c7e3-42a6-9b3d-da0b2a536577"
		testOperationState = "success"
	)

	httpmock.RegisterResponder("DELETE",
		fmt.Sprintf("=~^/load-balancer/%s/service/.*", testNLBID),
		func(req *http.Request) (*http.Response, error) {
			ts.Require().Equal(fmt.Sprintf("/load-balancer/%s/service/%s",
				testNLBID, testNLBServiceID), req.URL.String())

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testNLBServiceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	httpmock.RegisterResponder("GET", "/operation/"+testOperationID,
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testNLBServiceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}
			return resp, nil
		})

	nlb := &NetworkLoadBalancer{
		ID:   testNLBID,
		c:    ts.client,
		zone: testZone,

		Services: []*NetworkLoadBalancerService{{ID: testNLBServiceID}},
	}

	ts.Require().NoError(nlb.DeleteService(context.Background(), nlb.Services[0]))
}

func (ts *clientTestSuite) TestClient_CreateNetworkLoadBalancer() {
	var (
		testOperationID    = "08302193-c7e3-42a6-9b3d-da0b2a536577"
		testOperationState = "success"
	)

	httpmock.RegisterResponder("POST", "/load-balancer",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testNLBID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "/operation/"+testOperationID,
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testNLBID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}
			return resp, nil
		})

	httpmock.RegisterResponder("GET", fmt.Sprintf("/load-balancer/%s", testNLBID),
		func(_ *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.LoadBalancer{
				CreatedAt:   &testNLBCreatedAt,
				Description: &testNLBDescription,
				Id:          &testNLBID,
				Name:        &testNLBName,
				State:       &testNLBState,
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}
			return resp, nil
		})

	expected := &NetworkLoadBalancer{
		CreatedAt:   testNLBCreatedAt,
		Description: testNLBDescription,
		ID:          testNLBID,
		Name:        testNLBName,
		Services:    []*NetworkLoadBalancerService{},
		State:       testNLBState,

		c:    ts.client,
		zone: testZone,
	}

	actual, err := ts.client.CreateNetworkLoadBalancer(context.Background(), testZone, &NetworkLoadBalancer{
		Description: testNLBDescription,
		Name:        testNLBName,
	})
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_ListNetworkLoadBalancers() {
	var err error

	httpmock.RegisterResponder("GET", "/load-balancer",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(http.StatusOK, struct {
				LoadBalancers *[]papi.LoadBalancer `json:"load-balancers,omitempty"`
			}{
				LoadBalancers: &[]papi.LoadBalancer{{
					Id:          &testNLBID,
					Name:        &testNLBName,
					Description: &testNLBDescription,
					CreatedAt:   &testNLBCreatedAt,
					Services: &[]papi.LoadBalancerService{{
						Id:           &testNLBServiceID,
						Name:         &testNLBServiceName,
						Description:  &testNLBServiceDescription,
						InstancePool: &papi.InstancePool{Id: &testNLBServiceInstancePoolID},
						Protocol:     &testNLBServiceProtocol,
						Port:         &testNLBServicePort,
						TargetPort:   &testNLBServiceTargetPort,
						Strategy:     &testNLBServiceStrategy,
						Healthcheck: &papi.Healthcheck{
							Mode:     testNLServiceHealthcheckMode,
							Interval: &testNLBServiceHealthcheckInterval,
							Port:     testNLBServiceHealthcheckPort,
							Uri:      &testNLBServiceHealthcheckURI,
							Timeout:  &testNLBServiceHealthcheckTimeout,
							Retries:  &testNLBServiceHealthcheckRetries,
						},
						HealthcheckStatus: &[]papi.LoadBalancerServerStatus{
							{
								PublicIp: &testNLBServiceHealthcheckStatus1InstanceIP,
								Status:   &testNLBServiceHealthcheckStatus1Status,
							},
							{
								PublicIp: &testNLBServiceHealthcheckStatus2InstanceIP,
								Status:   &testNLBServiceHealthcheckStatus2Status,
							},
						},
						State: &testNLBState,
					}},
				}},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}
			return resp, nil
		})

	expected := []*NetworkLoadBalancer{{
		ID:          testNLBID,
		Name:        testNLBName,
		Description: testNLBDescription,
		CreatedAt:   testNLBCreatedAt,
		Services: []*NetworkLoadBalancerService{{
			ID:             testNLBServiceID,
			Name:           testNLBServiceName,
			Description:    testNLBServiceDescription,
			InstancePoolID: testNLBServiceInstancePoolID,
			Protocol:       testNLBServiceProtocol,
			Port:           uint16(testNLBServicePort),
			TargetPort:     uint16(testNLBServiceTargetPort),
			Strategy:       testNLBServiceStrategy,
			Healthcheck: NetworkLoadBalancerServiceHealthcheck{
				Mode:     testNLServiceHealthcheckMode,
				Port:     uint16(testNLBServiceHealthcheckPort),
				Interval: time.Duration(testNLBServiceHealthcheckInterval) * time.Second,
				Timeout:  time.Duration(testNLBServiceHealthcheckTimeout) * time.Second,
				Retries:  testNLBServiceHealthcheckRetries,
				URI:      testNLBServiceHealthcheckURI,
			},
			HealthcheckStatus: []*NetworkLoadBalancerServerStatus{
				{
					InstanceIP: net.ParseIP(testNLBServiceHealthcheckStatus1InstanceIP),
					Status:     testNLBServiceHealthcheckStatus1Status,
				},
				{
					InstanceIP: net.ParseIP(testNLBServiceHealthcheckStatus2InstanceIP),
					Status:     testNLBServiceHealthcheckStatus2Status,
				},
			},
			State: testNLBState,
		}},

		c:    ts.client,
		zone: testZone,
	}}

	actual, err := ts.client.ListNetworkLoadBalancers(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_GetNetworkLoadBalancer() {
	var err error

	httpmock.RegisterResponder("GET", "/load-balancer/"+testNLBID, // nolint:dupl
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.LoadBalancer{
				Id:          &testNLBID,
				Name:        &testNLBName,
				Description: &testNLBDescription,
				CreatedAt:   &testNLBCreatedAt,
				Ip:          &testNLBIPAddress,
				Services: &[]papi.LoadBalancerService{{
					Id:           &testNLBServiceID,
					Name:         &testNLBServiceName,
					Description:  &testNLBServiceDescription,
					InstancePool: &papi.InstancePool{Id: &testNLBServiceInstancePoolID},
					Protocol:     &testNLBServiceProtocol,
					Port:         &testNLBServicePort,
					TargetPort:   &testNLBServiceTargetPort,
					Strategy:     &testNLBServiceStrategy,
					Healthcheck: &papi.Healthcheck{
						Mode:     testNLServiceHealthcheckMode,
						Interval: &testNLBServiceHealthcheckInterval,
						Port:     testNLBServiceHealthcheckPort,
						Uri:      &testNLBServiceHealthcheckURI,
						Timeout:  &testNLBServiceHealthcheckTimeout,
						Retries:  &testNLBServiceHealthcheckRetries,
					},
					HealthcheckStatus: &[]papi.LoadBalancerServerStatus{
						{
							PublicIp: &testNLBServiceHealthcheckStatus1InstanceIP,
							Status:   &testNLBServiceHealthcheckStatus1Status,
						},
						{
							PublicIp: &testNLBServiceHealthcheckStatus2InstanceIP,
							Status:   &testNLBServiceHealthcheckStatus2Status,
						},
					},
					State: &testNLBServiceState,
				}},
				State: &testNLBState,
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}
			return resp, nil
		})

	expected := &NetworkLoadBalancer{
		ID:          testNLBID,
		Name:        testNLBName,
		Description: testNLBDescription,
		CreatedAt:   testNLBCreatedAt,
		IPAddress:   net.ParseIP(testNLBIPAddress),
		State:       testNLBState,
		Services: []*NetworkLoadBalancerService{{
			ID:             testNLBServiceID,
			Name:           testNLBServiceName,
			Description:    testNLBServiceDescription,
			InstancePoolID: testNLBServiceInstancePoolID,
			Protocol:       testNLBServiceProtocol,
			Port:           uint16(testNLBServicePort),
			TargetPort:     uint16(testNLBServiceTargetPort),
			Strategy:       testNLBServiceStrategy,
			Healthcheck: NetworkLoadBalancerServiceHealthcheck{
				Mode:     testNLServiceHealthcheckMode,
				Port:     uint16(testNLBServiceHealthcheckPort),
				Interval: time.Duration(testNLBServiceHealthcheckInterval) * time.Second,
				Timeout:  time.Duration(testNLBServiceHealthcheckTimeout) * time.Second,
				Retries:  testNLBServiceHealthcheckRetries,
				URI:      testNLBServiceHealthcheckURI,
			},
			HealthcheckStatus: []*NetworkLoadBalancerServerStatus{
				{
					InstanceIP: net.ParseIP(testNLBServiceHealthcheckStatus1InstanceIP),
					Status:     testNLBServiceHealthcheckStatus1Status,
				},
				{
					InstanceIP: net.ParseIP(testNLBServiceHealthcheckStatus2InstanceIP),
					Status:     testNLBServiceHealthcheckStatus2Status,
				},
			},
			State: testNLBServiceState,
		}},

		c:    ts.client,
		zone: testZone,
	}

	actual, err := ts.client.GetNetworkLoadBalancer(context.Background(), testZone, expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_UpdateNetworkLoadBalancer() {
	var (
		testNLBNameUpdated        = testNLBName + "-updated"
		testNLBDescriptionUpdated = testNLBDescription + "-updated"
		testOperationID           = "08302193-c7e3-42a6-9b3d-da0b2a536577"
		testOperationState        = "success"
	)

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/load-balancer/%s", testNLBID),
		func(req *http.Request) (*http.Response, error) {
			var actual papi.UpdateLoadBalancerJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)
			expected := papi.UpdateLoadBalancerJSONRequestBody{
				Name:        &testNLBNameUpdated,
				Description: &testNLBDescriptionUpdated,
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testNLBID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	httpmock.RegisterResponder("GET", "/operation/"+testOperationID,
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testNLBID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}
			return resp, nil
		})

	httpmock.RegisterResponder("GET", fmt.Sprintf("/load-balancer/%s", testNLBID),
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.LoadBalancer{
				CreatedAt:   &testNLBCreatedAt,
				Description: &testNLBDescriptionUpdated,
				Id:          &testNLBID,
				Ip:          &testNLBIPAddress,
				Name:        &testNLBNameUpdated,
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}
			return resp, nil
		})

	nlbUpdated := NetworkLoadBalancer{
		ID:          testNLBID,
		Name:        testNLBNameUpdated,
		Description: testNLBDescriptionUpdated,
	}

	_, err := ts.client.UpdateNetworkLoadBalancer(context.Background(), testZone, &nlbUpdated)
	ts.Require().NoError(err)
}

func (ts *clientTestSuite) TestClient_DeleteNetworkLoadBalancer() {
	var (
		testOperationID    = "08302193-c7e3-42a6-9b3d-da0b2a536577"
		testOperationState = "success"
	)

	httpmock.RegisterResponder("DELETE", "=~^/load-balancer/.*",
		func(req *http.Request) (*http.Response, error) {
			ts.Require().Equal(fmt.Sprintf("/load-balancer/%s", testNLBID), req.URL.String())

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testNLBID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	httpmock.RegisterResponder("GET", "/operation/"+testOperationID,
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testNLBID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}
			return resp, nil
		})

	ts.Require().NoError(ts.client.DeleteNetworkLoadBalancer(context.Background(), testZone, testNLBID))
}
