package egoscale

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	v2 "github.com/exoscale/egoscale/internal/v2"
)

var (
	testZone                                         = "ch-gva-2"
	testNLBID                                        = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
	testNLBName                                      = "test-nlb-name"
	testNLBDescription                               = "test-nlb-description"
	testNLBCreatedAt                                 = time.Now()
	testNLBState                                     = "running"
	testNLBServiceID                                 = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
	testNLBServiceName                               = "test-svc-name"
	testNLBServiceDescription                        = "test-svc-description"
	testNLBServiceInstancePoolID                     = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
	testNLBServiceProtocol                           = "tcp"
	testNLBServicePort                         int64 = 80
	testNLBServiceTargetPort                   int64 = 8080
	testNLBServiceStrategy                           = "round-robin"
	testNLServiceHealthcheckMode                     = "http"
	testNLBServiceHealthcheckPort              int64 = 8080
	testNLBServiceHealthcheckInterval          int64 = 10
	testNLBServiceHealthcheckTimeout           int64 = 3
	testNLBServiceHealthcheckRetries           int64 = 1
	testNLBServiceHealthcheckURI                     = "/health"
	testNLBServiceHealthcheckStatus1InstanceIP       = "1.2.3.4"
	testNLBServiceHealthcheckStatus1Status           = "success"
	testNLBServiceHealthcheckStatus2InstanceIP       = "5.6.7.8"
	testNLBServiceHealthcheckStatus2Status           = "success"
)

func TestNetworkLoadBalancer_AddService(t *testing.T) { t.Skip() }

func TestNetworkLoadBalancer_UpdateService(t *testing.T) { t.Skip() }

func TestNetworkLoadBalancer_DeleteService(t *testing.T) { t.Skip() }

func TestClient_CreateNetworkLoadBalancer(t *testing.T) { t.Skip() }

func TestClient_ListNetworkLoadBalancers(t *testing.T) {
	v2MockClient := new(v2.MockClient)
	client := NewClient("x", "x", "x")
	client.v2 = v2MockClient

	v2MockClient.
		On("ListLoadBalancersWithResponse", mock.Anything, mock.Anything).
		Return(&v2.ListLoadBalancersResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &struct {
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
			},
		}, nil)

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

	nlb, err := client.ListNetworkLoadBalancers(context.Background(), testZone)
	require.NoError(t, err)
	require.Equal(t, expected, nlb)
}

func TestClient_GetNetworkLoadBalancer(t *testing.T) {
	v2MockClient := new(v2.MockClient)
	client := NewClient("x", "x", "x")
	client.v2 = v2MockClient

	v2MockClient.
		On("GetLoadBalancerWithResponse", mock.Anything, mock.Anything).
		Return(&v2.GetLoadBalancerResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &v2.LoadBalancer{
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
			},
		}, nil)

	expected := &NetworkLoadBalancer{
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
	}

	nlb, err := client.GetNetworkLoadBalancer(context.Background(), testZone, expected.ID)
	require.NoError(t, err)
	require.Equal(t, expected, nlb)
}

func TestClient_UpdateNetworkLoadBalancer(t *testing.T) { t.Skip() }

func TestClient_DeleteNetworkLoadBalancer(t *testing.T) { t.Skip() }
