package publicapi

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLoadBalancer_UnmarshalJSON(t *testing.T) {
	var (
		testID                                       = testRandomID(t)
		testIP                                       = "1.2.3.4"
		testName                                     = "test-lb"
		testCreatedAt, _                             = time.Parse(iso8601Format, "2020-05-26T12:09:42Z")
		testDescription                              = "Test NLB description"
		testState                                    = LoadBalancerStateRunning
		testServiceID                                = testRandomID(t)
		testServiceName                              = "test-service"
		testServiceDescription                       = "Test service description"
		testServiceInstancePoolID                    = testRandomID(t)
		testServiceProtocol                          = LoadBalancerServiceProtocolTcp
		testServicePort                        int64 = 1234
		testServiceTargetPort                  int64 = 5678
		testServiceStrategy                          = LoadBalancerServiceStrategyRoundRobin
		testServiceState                             = LoadBalancerServiceStateRunning
		testServiceHealthcheckMode                   = LoadBalancerServiceHealthcheckModeHttps
		testServiceHealthcheckPort                   = testServiceTargetPort
		testServiceHealthcheckInterval         int64 = 10
		testServiceHealthcheckTimeout          int64 = 3
		testServiceHealthcheckRetries          int64 = 1
		testServiceHealthcheckURI                    = "/health"
		testServiceHealthcheckTLSSNI                 = "example.net"
		testServiceHealthcheckStatusInstanceIP       = "5.6.7.8"
		testServiceHealthcheckStatusStatus           = LoadBalancerServerStatusStatusSuccess

		expected = LoadBalancer{
			CreatedAt:   &testCreatedAt,
			Description: &testDescription,
			Id:          &testID,
			Ip:          &testIP,
			Name:        &testName,
			Services: &[]LoadBalancerService{{
				Description: &testServiceDescription,
				Healthcheck: &LoadBalancerServiceHealthcheck{
					Interval: &testServiceHealthcheckInterval,
					Mode:     testServiceHealthcheckMode,
					Port:     testServiceHealthcheckPort,
					Retries:  &testServiceHealthcheckRetries,
					Timeout:  &testServiceHealthcheckTimeout,
					TlsSni:   &testServiceHealthcheckTLSSNI,
					Uri:      &testServiceHealthcheckURI,
				},
				HealthcheckStatus: &[]LoadBalancerServerStatus{{
					PublicIp: &testServiceHealthcheckStatusInstanceIP,
					Status:   &testServiceHealthcheckStatusStatus,
				}},
				Id:           &testServiceID,
				InstancePool: &InstancePool{Id: &testServiceInstancePoolID},
				Name:         &testServiceName,
				Port:         &testServicePort,
				Protocol:     &testServiceProtocol,
				State:        &testServiceState,
				Strategy:     &testServiceStrategy,
				TargetPort:   &testServiceTargetPort,
			}},
			State: &testState,
		}

		actual LoadBalancer

		jsonNLB = `{
  "id": "` + testID + `",
  "ip": "` + testIP + `",
  "description": "` + testDescription + `",
  "state": "` + string(testState) + `",
  "name": "` + testName + `",
  "created-at": "` + testCreatedAt.Format(iso8601Format) + `",
  "services": [
	{
	  "id": "` + testServiceID + `",
	  "description": "` + testServiceDescription + `",
	  "protocol": "` + string(testServiceProtocol) + `",
	  "name": "` + testServiceName + `",
	  "state": "` + string(testServiceState) + `",
	  "target-port": ` + fmt.Sprint(testServiceTargetPort) + `,
	  "port": ` + fmt.Sprint(testServicePort) + `,
	  "instance-pool": {"id": "` + testServiceInstancePoolID + `"},
	  "strategy": "` + string(testServiceStrategy) + `",
	  "healthcheck": {
		"mode": "` + string(testServiceHealthcheckMode) + `",
		"uri": "` + testServiceHealthcheckURI + `",
		"interval": ` + fmt.Sprint(testServiceHealthcheckInterval) + `,
		"timeout": ` + fmt.Sprint(testServiceHealthcheckTimeout) + `,
		"port": ` + fmt.Sprint(testServiceHealthcheckPort) + `,
		"retries": ` + fmt.Sprint(testServiceHealthcheckRetries) + `,
		"tls-sni": "` + fmt.Sprint(testServiceHealthcheckTLSSNI) + `"
	  },
	  "healthcheck-status": [
        {
          "public-ip": "` + testServiceHealthcheckStatusInstanceIP + `",
          "status": "` + string(testServiceHealthcheckStatusStatus) + `"
        }
      ]
	}
  ]
}`
	)

	require.NoError(t, json.Unmarshal([]byte(jsonNLB), &actual))
	require.Equal(t, expected, actual)
}

func TestLoadBalancer_MarshalJSON(t *testing.T) {
	var (
		testID                                       = testRandomID(t)
		testIP                                       = "1.2.3.4"
		testName                                     = "test-lb"
		testCreatedAt, _                             = time.Parse(iso8601Format, "2020-05-26T12:09:42Z")
		testDescription                              = "Test NLB description"
		testState                                    = LoadBalancerStateRunning
		testServiceID                                = testRandomID(t)
		testServiceName                              = "test-service"
		testServiceDescription                       = "Test service description"
		testServiceInstancePoolID                    = testRandomID(t)
		testServiceProtocol                          = LoadBalancerServiceProtocolTcp
		testServicePort                        int64 = 1234
		testServiceTargetPort                  int64 = 5678
		testServiceStrategy                          = LoadBalancerServiceStrategyRoundRobin
		testServiceState                             = LoadBalancerServiceStateRunning
		testServiceHealthcheckMode                   = LoadBalancerServiceHealthcheckModeHttps
		testServiceHealthcheckPort                   = testServiceTargetPort
		testServiceHealthcheckInterval         int64 = 10
		testServiceHealthcheckTimeout          int64 = 3
		testServiceHealthcheckRetries          int64 = 1
		testServiceHealthcheckURI                    = "/health"
		testServiceHealthcheckTLSSNI                 = "example.net"
		testServiceHealthcheckStatusInstanceIP       = "5.6.7.8"
		testServiceHealthcheckStatusStatus           = LoadBalancerServerStatusStatusSuccess

		lb = LoadBalancer{
			CreatedAt:   &testCreatedAt,
			Description: &testDescription,
			Id:          &testID,
			Ip:          &testIP,
			Name:        &testName,
			State:       &testState,
			Services: &[]LoadBalancerService{{
				Description: &testServiceDescription,
				Healthcheck: &LoadBalancerServiceHealthcheck{
					Interval: &testServiceHealthcheckInterval,
					Mode:     testServiceHealthcheckMode,
					Port:     testServiceHealthcheckPort,
					Retries:  &testServiceHealthcheckRetries,
					Timeout:  &testServiceHealthcheckTimeout,
					TlsSni:   &testServiceHealthcheckTLSSNI,
					Uri:      &testServiceHealthcheckURI,
				},
				HealthcheckStatus: &[]LoadBalancerServerStatus{{
					PublicIp: &testServiceHealthcheckStatusInstanceIP,
					Status:   &testServiceHealthcheckStatusStatus,
				}},
				Id:           &testServiceID,
				InstancePool: &InstancePool{Id: &testServiceInstancePoolID},
				Name:         &testServiceName,
				Port:         &testServicePort,
				Protocol:     &testServiceProtocol,
				State:        &testServiceState,
				Strategy:     &testServiceStrategy,
				TargetPort:   &testServiceTargetPort,
			}},
		}

		expected = []byte(`{` +
			`"created-at":"` + testCreatedAt.Format(iso8601Format) + `",` +
			`"description":"` + testDescription + `",` +
			`"id":"` + testID + `",` +
			`"ip":"` + testIP + `",` +
			`"name":"` + testName + `",` +
			`"services":[{` +
			`"description":"` + testServiceDescription + `",` +
			`"healthcheck":` +
			`{` +
			`"interval":` + fmt.Sprint(testServiceHealthcheckInterval) + `,` +
			`"mode":"` + string(testServiceHealthcheckMode) + `",` +
			`"port":` + fmt.Sprint(testServiceHealthcheckPort) + `,` +
			`"retries":` + fmt.Sprint(testServiceHealthcheckRetries) + `,` +
			`"timeout":` + fmt.Sprint(testServiceHealthcheckTimeout) + `,` +
			`"tls-sni":"` + testServiceHealthcheckTLSSNI + `",` +
			`"uri":"` + testServiceHealthcheckURI + `"` +
			`},` +
			`"healthcheck-status":[` +
			`{` +
			`"public-ip":"` + testServiceHealthcheckStatusInstanceIP + `",` +
			`"status":"` + string(testServiceHealthcheckStatusStatus) + `"` +
			`}` +
			`],` +
			`"id":"` + testServiceID + `",` +
			`"instance-pool":{"id":"` + testServiceInstancePoolID + `"},` +
			`"name":"` + testServiceName + `",` +
			`"port":` + fmt.Sprint(testServicePort) + `,` +
			`"protocol":"` + string(testServiceProtocol) + `",` +
			`"state":"` + string(testServiceState) + `",` +
			`"strategy":"` + string(testServiceStrategy) + `",` +
			`"target-port":` + fmt.Sprint(testServiceTargetPort) +
			`}` +
			`],` +
			`"state":"` + string(testState) + `"` +
			`}`)
	)

	actual, err := json.Marshal(lb)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}
