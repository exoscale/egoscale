package egoscale

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	v2 "github.com/exoscale/egoscale/internal/v2"
)

const iso8601Format = "2006-01-02T15:04:05Z"

var (
	testZone                                         = "ch-gva-2"
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

func TestNetworkLoadBalancer_AddService(t *testing.T) {
	var (
		testOperationID    = "08302193-c7e3-42a6-9b3d-da0b2a536577"
		testOperationState = "success"
		err                error
	)

	mockClient := v2.NewMockClient()
	client := NewClient("x", "x", "x")
	client.v2, err = v2.NewClientWithResponses("", v2.WithHTTPClient(mockClient))
	require.NoError(t, err)

	mockClient.RegisterResponder("POST", "/load-balancer/"+testNLBID+"/service",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(http.StatusOK, v2.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &v2.Reference{Id: &testNLBID},
			})
			if err != nil {
				t.Fatalf("error initializing mock HTTP responder: %s", err)
			}
			return resp, nil
		})

	mockClient.RegisterResponder("GET", "/operation/"+testOperationID,
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(http.StatusOK, v2.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &v2.Reference{Id: &testNLBID},
			})
			if err != nil {
				t.Fatalf("error initializing mock HTTP responder: %s", err)
			}
			return resp, nil
		})

	mockClient.RegisterResponder("GET", "/load-balancer/"+testNLBID, // nolint:dupl
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(http.StatusOK, v2.LoadBalancer{
				Id:          &testNLBID,
				Name:        &testNLBName,
				Description: &testNLBDescription,
				CreatedAt:   &testNLBCreatedAt,
				Ip:          &testNLBIPAddress,
				Services: &[]v2.LoadBalancerService{{
					Id:           &testNLBServiceID,
					Name:         &testNLBServiceName,
					Description:  &testNLBServiceDescription,
					InstancePool: &v2.Resource{Id: &testNLBServiceInstancePoolID},
					Protocol:     &testNLBServiceProtocol,
					Port:         &testNLBServicePort,
					TargetPort:   &testNLBServiceTargetPort,
					Strategy:     &testNLBServiceStrategy,
					Healthcheck: &v2.Healthcheck{
						Mode:     &testNLServiceHealthcheckMode,
						Interval: &testNLBServiceHealthcheckInterval,
						Port:     &testNLBServiceHealthcheckPort,
						Uri:      &testNLBServiceHealthcheckURI,
						TlsSni:   &testNLBServiceHealthcheckTLSSNI,
						Timeout:  &testNLBServiceHealthcheckTimeout,
						Retries:  &testNLBServiceHealthcheckRetries,
					},
					HealthcheckStatus: &[]v2.LoadBalancerServerStatus{
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
				t.Fatalf("error initializing mock HTTP responder: %s", err)
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

		c: client,
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
	require.NoError(t, err)
	require.Equal(t, expected, actual)

}

// UpdateService is not tested as it essentially relies on the already tested GetNetworkLoadBalancer.
func TestNetworkLoadBalancer_UpdateService(t *testing.T) { t.Skip() }

// DeleteService is not tested as it only produces API-side effects.
func TestNetworkLoadBalancer_DeleteService(t *testing.T) { t.Skip() }

// CreateNetworkLoadBalancer is not tested as it essentially relies on the already tested GetNetworkLoadBalancer.
func TestClient_CreateNetworkLoadBalancer(t *testing.T) { t.Skip() }

func TestClient_ListNetworkLoadBalancers(t *testing.T) {
	var err error

	mockClient := v2.NewMockClient()
	client := NewClient("x", "x", "x")
	client.v2, err = v2.NewClientWithResponses("", v2.WithHTTPClient(mockClient))
	require.NoError(t, err)

	mockClient.RegisterResponder("GET", "/load-balancer",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(http.StatusOK, struct {
				LoadBalancers *[]v2.LoadBalancer `json:"load-balancers,omitempty"`
			}{
				LoadBalancers: &[]v2.LoadBalancer{{
					Id:          &testNLBID,
					Name:        &testNLBName,
					Description: &testNLBDescription,
					CreatedAt:   &testNLBCreatedAt,
					Services: &[]v2.LoadBalancerService{{
						Id:           &testNLBServiceID,
						Name:         &testNLBServiceName,
						Description:  &testNLBServiceDescription,
						InstancePool: &v2.Resource{Id: &testNLBServiceInstancePoolID},
						Protocol:     &testNLBServiceProtocol,
						Port:         &testNLBServicePort,
						TargetPort:   &testNLBServiceTargetPort,
						Strategy:     &testNLBServiceStrategy,
						Healthcheck: &v2.Healthcheck{
							Mode:     &testNLServiceHealthcheckMode,
							Interval: &testNLBServiceHealthcheckInterval,
							Port:     &testNLBServiceHealthcheckPort,
							Uri:      &testNLBServiceHealthcheckURI,
							Timeout:  &testNLBServiceHealthcheckTimeout,
							Retries:  &testNLBServiceHealthcheckRetries,
						},
						HealthcheckStatus: &[]v2.LoadBalancerServerStatus{
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
				t.Fatalf("error initializing mock HTTP responder: %s", err)
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

		c:    client,
		zone: testZone,
	}}

	actual, err := client.ListNetworkLoadBalancers(context.Background(), testZone)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestClient_GetNetworkLoadBalancer(t *testing.T) {
	var err error

	mockClient := v2.NewMockClient()
	client := NewClient("x", "x", "x")
	client.v2, err = v2.NewClientWithResponses("", v2.WithHTTPClient(mockClient))
	require.NoError(t, err)

	mockClient.RegisterResponder("GET", "/load-balancer/"+testNLBID, // nolint:dupl
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(http.StatusOK, v2.LoadBalancer{
				Id:          &testNLBID,
				Name:        &testNLBName,
				Description: &testNLBDescription,
				CreatedAt:   &testNLBCreatedAt,
				Ip:          &testNLBIPAddress,
				Services: &[]v2.LoadBalancerService{{
					Id:           &testNLBServiceID,
					Name:         &testNLBServiceName,
					Description:  &testNLBServiceDescription,
					InstancePool: &v2.Resource{Id: &testNLBServiceInstancePoolID},
					Protocol:     &testNLBServiceProtocol,
					Port:         &testNLBServicePort,
					TargetPort:   &testNLBServiceTargetPort,
					Strategy:     &testNLBServiceStrategy,
					Healthcheck: &v2.Healthcheck{
						Mode:     &testNLServiceHealthcheckMode,
						Interval: &testNLBServiceHealthcheckInterval,
						Port:     &testNLBServiceHealthcheckPort,
						Uri:      &testNLBServiceHealthcheckURI,
						Timeout:  &testNLBServiceHealthcheckTimeout,
						Retries:  &testNLBServiceHealthcheckRetries,
					},
					HealthcheckStatus: &[]v2.LoadBalancerServerStatus{
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
				t.Fatalf("error initializing mock HTTP responder: %s", err)
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

		c:    client,
		zone: testZone,
	}

	actual, err := client.GetNetworkLoadBalancer(context.Background(), testZone, expected.ID)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

// UpdateNetworkLoadBalancer is not tested as it essentially relies on the already tested GetNetworkLoadBalancer.
func TestClient_UpdateNetworkLoadBalancer(t *testing.T) { t.Skip() }

// DeleteNetworkLoadBalancer is not tested as it only produces API-side effects.
func TestClient_DeleteNetworkLoadBalancer(t *testing.T) { t.Skip() }
