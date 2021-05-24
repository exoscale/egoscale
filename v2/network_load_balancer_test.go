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

var (
	testNLBID                                        = new(clientTestSuite).randomID()
	testNLBName                                      = "test-nlb-name"
	testNLBDescription                               = "test-nlb-description"
	testNLBCreatedAt, _                              = time.Parse(iso8601Format, "2020-05-26T12:09:42Z")
	testNLBIPAddress                                 = "101.102.103.104"
	testNLBState                                     = papi.LoadBalancerStateRunning
	testNLBServiceID                                 = new(clientTestSuite).randomID()
	testNLBServiceName                               = "test-svc-name"
	testNLBServiceDescription                        = new(clientTestSuite).randomID()
	testNLBServiceInstancePoolID                     = new(clientTestSuite).randomID()
	testNLBServiceProtocol                           = papi.LoadBalancerServiceProtocolTcp
	testNLBServicePort                         int64 = 443
	testNLBServiceTargetPort                   int64 = 8443
	testNLBServiceStrategy                           = papi.LoadBalancerServiceStrategyRoundRobin
	testNLBServiceState                              = papi.DbaasServiceStateRunning
	testNLServiceHealthcheckMode                     = papi.LoadBalancerServiceHealthcheckModeHttps
	testNLBServiceHealthcheckPort              int64 = 8080
	testNLBServiceHealthcheckInterval          int64 = 10
	testNLBServiceHealthcheckTimeout           int64 = 3
	testNLBServiceHealthcheckRetries           int64 = 1
	testNLBServiceHealthcheckURI                     = "/health"
	testNLBServiceHealthcheckTLSSNI                  = "example.net"
	testNLBServiceHealthcheckStatus1InstanceIP       = "1.2.3.4"
	testNLBServiceHealthcheckStatus1Status           = papi.LoadBalancerServerStatusStatusSuccess
	testNLBServiceHealthcheckStatus2InstanceIP       = "5.6.7.8"
	testNLBServiceHealthcheckStatus2Status           = papi.LoadBalancerServerStatusStatusSuccess
)

func (ts *clientTestSuite) TestNetworkLoadBalancer_AddService() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = papi.OperationStateSuccess
	)

	httpmock.RegisterResponder("POST", fmt.Sprintf("/load-balancer/%s/service", testNLBID),
		func(req *http.Request) (*http.Response, error) {
			var actual papi.AddServiceToLoadBalancerJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.AddServiceToLoadBalancerJSONRequestBody{
				Description: &testNLBServiceDescription,
				Healthcheck: papi.LoadBalancerServiceHealthcheck{
					Interval: &testNLBServiceHealthcheckInterval,
					Mode:     testNLServiceHealthcheckMode,
					Port:     testNLBServiceHealthcheckPort,
					Retries:  &testNLBServiceHealthcheckRetries,
					Timeout:  &testNLBServiceHealthcheckTimeout,
					TlsSni:   &testNLBServiceHealthcheckTLSSNI,
					Uri:      &testNLBServiceHealthcheckURI,
				},
				InstancePool: papi.InstancePool{Id: &testNLBServiceInstancePoolID},
				Name:         testNLBServiceName,
				Port:         testNLBServicePort,
				Protocol:     papi.AddServiceToLoadBalancerJSONBodyProtocol(testNLBServiceProtocol),
				Strategy:     papi.AddServiceToLoadBalancerJSONBodyStrategy(testNLBServiceStrategy),
				TargetPort:   testNLBServiceTargetPort,
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

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testNLBID},
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/load-balancer/%s", testNLBID), papi.LoadBalancer{
		CreatedAt:   &testNLBCreatedAt,
		Description: &testNLBDescription,
		Id:          &testNLBID,
		Ip:          &testNLBIPAddress,
		Name:        &testNLBName,
		Services: &[]papi.LoadBalancerService{{
			Description: &testNLBServiceDescription,
			Healthcheck: &papi.LoadBalancerServiceHealthcheck{
				Interval: &testNLBServiceHealthcheckInterval,
				Mode:     testNLServiceHealthcheckMode,
				Port:     testNLBServiceHealthcheckPort,
				Retries:  &testNLBServiceHealthcheckRetries,
				Timeout:  &testNLBServiceHealthcheckTimeout,
				TlsSni:   &testNLBServiceHealthcheckTLSSNI,
				Uri:      &testNLBServiceHealthcheckURI,
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
			Id:           &testNLBServiceID,
			InstancePool: &papi.InstancePool{Id: &testNLBServiceInstancePoolID},
			Name:         &testNLBServiceName,
			Port:         &testNLBServicePort,
			Protocol:     &testNLBServiceProtocol,
			Strategy:     &testNLBServiceStrategy,
			TargetPort:   &testNLBServiceTargetPort,
			State:        (*papi.LoadBalancerServiceState)(&testNLBServiceState),
		}},
		State: &testNLBState,
	})

	nlb := &NetworkLoadBalancer{
		CreatedAt:   testNLBCreatedAt,
		Description: testNLBDescription,
		ID:          testNLBID,
		IPAddress:   net.ParseIP(testNLBIPAddress),
		Name:        testNLBName,
		State:       string(testNLBState),

		c: ts.client,
	}

	expected := &NetworkLoadBalancerService{
		Description: testNLBServiceDescription,
		Healthcheck: NetworkLoadBalancerServiceHealthcheck{
			Interval: time.Duration(testNLBServiceHealthcheckInterval) * time.Second,
			Mode:     string(testNLServiceHealthcheckMode),
			Port:     uint16(testNLBServiceHealthcheckPort),
			Retries:  testNLBServiceHealthcheckRetries,
			TLSSNI:   testNLBServiceHealthcheckTLSSNI,
			Timeout:  time.Duration(testNLBServiceHealthcheckTimeout) * time.Second,
			URI:      testNLBServiceHealthcheckURI,
		},
		HealthcheckStatus: []*NetworkLoadBalancerServerStatus{
			{
				InstanceIP: net.ParseIP(testNLBServiceHealthcheckStatus1InstanceIP),
				Status:     string(testNLBServiceHealthcheckStatus1Status),
			},
			{
				InstanceIP: net.ParseIP(testNLBServiceHealthcheckStatus2InstanceIP),
				Status:     string(testNLBServiceHealthcheckStatus2Status),
			},
		},
		ID:             testNLBServiceID,
		InstancePoolID: testNLBServiceInstancePoolID,
		Name:           testNLBServiceName,
		Port:           uint16(testNLBServicePort),
		Protocol:       string(testNLBServiceProtocol),
		Strategy:       string(testNLBServiceStrategy),
		TargetPort:     uint16(testNLBServiceTargetPort),
		State:          string(testNLBServiceState),
	}

	actual, err := nlb.AddService(context.Background(), expected)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestNetworkLoadBalancer_UpdateService() {
	var (
		testNLBServiceNameUpdated                = testNLBServiceName + "-updated"
		testNLBServiceDescriptionUpdated         = testNLBServiceDescription + "-updated"
		testNLBServiceHealthcheckIntervalUpdated = testNLBServiceHealthcheckInterval + 1
		testNLBServiceHealthcheckModeUpdated     = papi.LoadBalancerServiceHealthcheckModeHttp
		testNLBServiceHealthcheckPortUpdated     = testNLBServiceHealthcheckPort + 1
		testNLBServiceHealthcheckRetriesUpdated  = testNLBServiceHealthcheckRetries + 1
		testNLBServiceHealthcheckTLSSNIUpdated   = ""
		testNLBServiceHealthcheckTimeoutUpdated  = testNLBServiceHealthcheckTimeout + 1
		testNLBServiceHealthcheckURIUpdated      = ""
		testOperationID                          = ts.randomID()
		testOperationState                       = papi.OperationStateSuccess
		updated                                  = false
	)

	nlb := &NetworkLoadBalancer{
		ID:   testNLBID,
		c:    ts.client,
		zone: testZone,

		Services: []*NetworkLoadBalancerService{{
			ID:          testNLBServiceID,
			Name:        testNLBServiceName,
			Description: testNLBServiceDescription,
		}},
	}

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/load-balancer/%s/service/%s",
		nlb.ID,
		nlb.Services[0].ID),
		func(req *http.Request) (*http.Response, error) {
			updated = true

			var actual papi.UpdateLoadBalancerServiceJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.UpdateLoadBalancerServiceJSONRequestBody{
				Name:        &testNLBServiceNameUpdated,
				Description: &testNLBServiceDescriptionUpdated,
				Healthcheck: &papi.LoadBalancerServiceHealthcheck{
					Interval: &testNLBServiceHealthcheckIntervalUpdated,
					Mode:     testNLBServiceHealthcheckModeUpdated,
					Port:     testNLBServiceHealthcheckPortUpdated,
					Retries:  &testNLBServiceHealthcheckRetriesUpdated,
					Timeout:  &testNLBServiceHealthcheckTimeoutUpdated,
					Uri:      &testNLBServiceHealthcheckURIUpdated,
				},
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

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testNLBServiceID},
	})

	ts.Require().NoError(nlb.UpdateService(context.Background(), &NetworkLoadBalancerService{
		ID:          nlb.Services[0].ID,
		Name:        testNLBServiceNameUpdated,
		Description: testNLBServiceDescriptionUpdated,
		Healthcheck: NetworkLoadBalancerServiceHealthcheck{
			Interval: time.Duration(testNLBServiceHealthcheckIntervalUpdated) * time.Second,
			Mode:     string(testNLBServiceHealthcheckModeUpdated),
			Port:     uint16(testNLBServiceHealthcheckPortUpdated),
			Retries:  testNLBServiceHealthcheckRetriesUpdated,
			TLSSNI:   testNLBServiceHealthcheckTLSSNIUpdated,
			Timeout:  time.Duration(testNLBServiceHealthcheckTimeoutUpdated) * time.Second,
			URI:      testNLBServiceHealthcheckURIUpdated,
		},
	}))
	ts.Require().True(updated)
}

func (ts *clientTestSuite) TestNetworkLoadBalancer_DeleteService() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = papi.OperationStateSuccess
		deleted            = false
	)

	httpmock.RegisterResponder("DELETE",
		fmt.Sprintf("/load-balancer/%s/service/%s", testNLBID, testNLBServiceID),
		func(req *http.Request) (*http.Response, error) {
			deleted = true

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

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testNLBServiceID},
	})

	nlb := &NetworkLoadBalancer{
		ID:   testNLBID,
		c:    ts.client,
		zone: testZone,

		Services: []*NetworkLoadBalancerService{{ID: testNLBServiceID}},
	}

	ts.Require().NoError(nlb.DeleteService(context.Background(), nlb.Services[0]))
	ts.Require().True(deleted)
}

func (ts *clientTestSuite) TestClient_CreateNetworkLoadBalancer() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = papi.OperationStateSuccess
	)

	httpmock.RegisterResponder("POST", "/load-balancer",
		func(req *http.Request) (*http.Response, error) {
			var actual papi.CreateLoadBalancerJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.CreateLoadBalancerJSONRequestBody{
				Description: &testNLBDescription,
				Name:        testNLBName,
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

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testNLBID},
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/load-balancer/%s", testNLBID), papi.LoadBalancer{
		CreatedAt:   &testNLBCreatedAt,
		Description: &testNLBDescription,
		Id:          &testNLBID,
		Name:        &testNLBName,
		State:       &testNLBState,
	})

	expected := &NetworkLoadBalancer{
		CreatedAt:   testNLBCreatedAt,
		Description: testNLBDescription,
		ID:          testNLBID,
		Name:        testNLBName,
		Services:    []*NetworkLoadBalancerService{},
		State:       string(testNLBState),

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
	ts.mockAPIRequest("GET", "/load-balancer", struct {
		LoadBalancers *[]papi.LoadBalancer `json:"load-balancers,omitempty"`
	}{
		LoadBalancers: &[]papi.LoadBalancer{{
			CreatedAt:   &testNLBCreatedAt,
			Description: &testNLBDescription,
			Id:          &testNLBID,
			Name:        &testNLBName,
			Services: &[]papi.LoadBalancerService{{
				Description: &testNLBServiceDescription,
				Healthcheck: &papi.LoadBalancerServiceHealthcheck{
					Interval: &testNLBServiceHealthcheckInterval,
					Mode:     testNLServiceHealthcheckMode,
					Port:     testNLBServiceHealthcheckPort,
					Retries:  &testNLBServiceHealthcheckRetries,
					Timeout:  &testNLBServiceHealthcheckTimeout,
					Uri:      &testNLBServiceHealthcheckURI,
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
				Id:           &testNLBServiceID,
				InstancePool: &papi.InstancePool{Id: &testNLBServiceInstancePoolID},
				Name:         &testNLBServiceName,
				Port:         &testNLBServicePort,
				Protocol:     &testNLBServiceProtocol,
				State:        (*papi.LoadBalancerServiceState)(&testNLBState),
				Strategy:     &testNLBServiceStrategy,
				TargetPort:   &testNLBServiceTargetPort,
			}},
			State: &testNLBState,
		}},
	})

	expected := []*NetworkLoadBalancer{{
		CreatedAt:   testNLBCreatedAt,
		Description: testNLBDescription,
		ID:          testNLBID,
		Name:        testNLBName,
		Services: []*NetworkLoadBalancerService{{
			Description: testNLBServiceDescription,
			Healthcheck: NetworkLoadBalancerServiceHealthcheck{
				Interval: time.Duration(testNLBServiceHealthcheckInterval) * time.Second,
				Mode:     string(testNLServiceHealthcheckMode),
				Port:     uint16(testNLBServiceHealthcheckPort),
				Retries:  testNLBServiceHealthcheckRetries,
				Timeout:  time.Duration(testNLBServiceHealthcheckTimeout) * time.Second,
				URI:      testNLBServiceHealthcheckURI,
			},
			HealthcheckStatus: []*NetworkLoadBalancerServerStatus{
				{
					InstanceIP: net.ParseIP(testNLBServiceHealthcheckStatus1InstanceIP),
					Status:     string(testNLBServiceHealthcheckStatus1Status),
				},
				{
					InstanceIP: net.ParseIP(testNLBServiceHealthcheckStatus2InstanceIP),
					Status:     string(testNLBServiceHealthcheckStatus2Status),
				},
			},
			ID:             testNLBServiceID,
			InstancePoolID: testNLBServiceInstancePoolID,
			Name:           testNLBServiceName,
			Port:           uint16(testNLBServicePort),
			Protocol:       string(testNLBServiceProtocol),
			Strategy:       string(testNLBServiceStrategy),
			TargetPort:     uint16(testNLBServiceTargetPort),
			State:          string(testNLBState),
		}},
		State: string(testNLBState),

		c:    ts.client,
		zone: testZone,
	}}

	actual, err := ts.client.ListNetworkLoadBalancers(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_GetNetworkLoadBalancer() {
	ts.mockAPIRequest("GET", fmt.Sprintf("/load-balancer/%s", testNLBID), papi.LoadBalancer{
		CreatedAt:   &testNLBCreatedAt,
		Description: &testNLBDescription,
		Id:          &testNLBID,
		Ip:          &testNLBIPAddress,
		Name:        &testNLBName,
		Services: &[]papi.LoadBalancerService{{
			Description: &testNLBServiceDescription,
			Healthcheck: &papi.LoadBalancerServiceHealthcheck{
				Interval: &testNLBServiceHealthcheckInterval,
				Mode:     testNLServiceHealthcheckMode,
				Port:     testNLBServiceHealthcheckPort,
				Retries:  &testNLBServiceHealthcheckRetries,
				Timeout:  &testNLBServiceHealthcheckTimeout,
				Uri:      &testNLBServiceHealthcheckURI,
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
			Id:           &testNLBServiceID,
			InstancePool: &papi.InstancePool{Id: &testNLBServiceInstancePoolID},
			Name:         &testNLBServiceName,
			Port:         &testNLBServicePort,
			Protocol:     &testNLBServiceProtocol,
			Strategy:     &testNLBServiceStrategy,
			TargetPort:   &testNLBServiceTargetPort,
			State:        (*papi.LoadBalancerServiceState)(&testNLBServiceState),
		}},
		State: &testNLBState,
	})

	expected := &NetworkLoadBalancer{
		CreatedAt:   testNLBCreatedAt,
		Description: testNLBDescription,
		ID:          testNLBID,
		IPAddress:   net.ParseIP(testNLBIPAddress),
		Name:        testNLBName,
		State:       string(testNLBState),
		Services: []*NetworkLoadBalancerService{{
			Description: testNLBServiceDescription,
			Healthcheck: NetworkLoadBalancerServiceHealthcheck{
				Interval: time.Duration(testNLBServiceHealthcheckInterval) * time.Second,
				Mode:     string(testNLServiceHealthcheckMode),
				Port:     uint16(testNLBServiceHealthcheckPort),
				Retries:  testNLBServiceHealthcheckRetries,
				Timeout:  time.Duration(testNLBServiceHealthcheckTimeout) * time.Second,
				URI:      testNLBServiceHealthcheckURI,
			},
			HealthcheckStatus: []*NetworkLoadBalancerServerStatus{
				{
					InstanceIP: net.ParseIP(testNLBServiceHealthcheckStatus1InstanceIP),
					Status:     string(testNLBServiceHealthcheckStatus1Status),
				},
				{
					InstanceIP: net.ParseIP(testNLBServiceHealthcheckStatus2InstanceIP),
					Status:     string(testNLBServiceHealthcheckStatus2Status),
				},
			},
			ID:             testNLBServiceID,
			InstancePoolID: testNLBServiceInstancePoolID,
			Name:           testNLBServiceName,
			Port:           uint16(testNLBServicePort),
			Protocol:       string(testNLBServiceProtocol),
			Strategy:       string(testNLBServiceStrategy),
			TargetPort:     uint16(testNLBServiceTargetPort),
			State:          string(testNLBServiceState),
		}},

		c:    ts.client,
		zone: testZone,
	}

	actual, err := ts.client.GetNetworkLoadBalancer(context.Background(), testZone, expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_FindNetworkLoadBalancer() {
	ts.mockAPIRequest("GET", "/load-balancer", struct {
		LoadBalancers *[]papi.LoadBalancer `json:"load-balancers,omitempty"`
	}{
		LoadBalancers: &[]papi.LoadBalancer{{
			CreatedAt: &testNLBCreatedAt,
			Id:        &testNLBID,
			Name:      &testNLBName,
			State:     &testNLBState,
		}},
	})
	ts.mockAPIRequest("GET", fmt.Sprintf("/load-balancer/%s", testNLBID), papi.LoadBalancer{
		CreatedAt: &testNLBCreatedAt,
		Id:        &testNLBID,
		Ip:        &testNLBIPAddress,
		Name:      &testNLBName,
		State:     &testNLBState,
	})

	expected := &NetworkLoadBalancer{
		CreatedAt: testNLBCreatedAt,
		ID:        testNLBID,
		IPAddress: net.ParseIP(testNLBIPAddress),
		Name:      testNLBName,
		Services:  []*NetworkLoadBalancerService{},
		State:     string(testNLBState),

		c:    ts.client,
		zone: testZone,
	}

	actual, err := ts.client.FindNetworkLoadBalancer(context.Background(), testZone, expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)

	actual, err = ts.client.FindNetworkLoadBalancer(context.Background(), testZone, expected.Name)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_UpdateNetworkLoadBalancer() {
	var (
		testNLBNameUpdated        = testNLBName + "-updated"
		testNLBDescriptionUpdated = testNLBDescription + "-updated"
		testOperationID           = ts.randomID()
		testOperationState        = papi.OperationStateSuccess
		updated                   = false
	)

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/load-balancer/%s", testNLBID),
		func(req *http.Request) (*http.Response, error) {
			updated = true

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

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testNLBID},
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/load-balancer/%s", testNLBID), papi.LoadBalancer{
		CreatedAt:   &testNLBCreatedAt,
		Description: &testNLBDescriptionUpdated,
		Id:          &testNLBID,
		Ip:          &testNLBIPAddress,
		Name:        &testNLBNameUpdated,
	})

	nlbUpdated := NetworkLoadBalancer{
		ID:          testNLBID,
		Name:        testNLBNameUpdated,
		Description: testNLBDescriptionUpdated,
	}

	ts.Require().NoError(ts.client.UpdateNetworkLoadBalancer(context.Background(), testZone, &nlbUpdated))
	ts.Require().True(updated)
}

func (ts *clientTestSuite) TestClient_DeleteNetworkLoadBalancer() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = papi.OperationStateSuccess
		deleted            = false
	)

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("/load-balancer/%s", testNLBID),
		func(req *http.Request) (*http.Response, error) {
			deleted = true

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

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testNLBID},
	})

	ts.Require().NoError(ts.client.DeleteNetworkLoadBalancer(context.Background(), testZone, testNLBID))
	ts.Require().True(deleted)
}
