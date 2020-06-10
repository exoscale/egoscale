package v2

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLoadBalancer_UnmarshalJSON(t *testing.T) {
	var (
		testID                               = "c0f306e7-21aa-4b0b-bafd-b86ed31bc2b8"
		testIP                               = "1.2.3.4"
		testName                             = "test-lb"
		testCreatedAt, _                     = time.Parse(iso8601Format, "2020-05-26T12:09:42Z")
		testDescription                      = "Test NLB description"
		testState                            = "running"
		testServiceID                        = "8ec00b67-e7ff-4ce5-b6b9-85fc1e24d878"
		testServiceName                      = "test-service"
		testServiceDescription               = "Test service description"
		testServiceInstancePoolID            = "0b7955e0-7beb-4d3a-dd23-7fe97aa2669f"
		testServiceProtocol                  = "tcp"
		testServicePort                int64 = 1234
		testServiceTargetPort          int64 = 5678
		testServiceStrategy                  = "round-robin"
		testServiceState                     = "running"
		testServiceHealthcheckMode           = "http"
		testServiceHealthcheckPort           = testServiceTargetPort
		testServiceHealthcheckInterval int64 = 10
		testServiceHealthcheckTimeout  int64 = 3
		testServiceHealthcheckRetries  int64 = 1
		testServiceHealthcheckURI            = "/health"

		expectedLoadBalancer = LoadBalancer{
			CreatedAt:   &testCreatedAt,
			Description: &testDescription,
			Id:          &testID,
			Ip:          &testIP,
			Name:        &testName,
			State:       &testState,
			Services: &[]LoadBalancerService{{
				Description:  &testServiceDescription,
				Id:           &testServiceID,
				InstancePool: &Resource{Id: &testServiceInstancePoolID},
				Name:         &testServiceName,
				Port:         &testServicePort,
				Protocol:     &testServiceProtocol,
				State:        &testServiceState,
				Strategy:     &testServiceStrategy,
				TargetPort:   &testServiceTargetPort,
				Healthcheck: &Healthcheck{
					Interval: &testServiceHealthcheckInterval,
					Mode:     &testServiceHealthcheckMode,
					Port:     &testServiceHealthcheckPort,
					Retries:  &testServiceHealthcheckRetries,
					Timeout:  &testServiceHealthcheckTimeout,
					Uri:      &testServiceHealthcheckURI,
				},
				// HealthcheckStatus // FIXME: the API doesn't return this information ATM
			}},
		}

		actualLoadBalancer LoadBalancer

		jsonNLB = `{
  "id": "` + testID + `",
  "ip": "` + testIP + `",
  "description": "` + testDescription + `",
  "state": "` + testState + `",
  "name": "` + testName + `",
  "created-at": "` + testCreatedAt.Format(iso8601Format) + `",
  "services": [
	{
	  "id": "` + testServiceID + `",
	  "description": "` + testServiceDescription + `",
	  "protocol": "` + testServiceProtocol + `",
	  "name": "` + testServiceName + `",
	  "state": "` + testServiceState + `",
	  "target-port": ` + fmt.Sprint(testServiceTargetPort) + `,
	  "port": ` + fmt.Sprint(testServicePort) + `,
	  "instance-pool": {"id": "` + testServiceInstancePoolID + `"},
	  "strategy": "` + testServiceStrategy + `",
	  "healthcheck": {
		"mode": "` + testServiceHealthcheckMode + `",
		"uri": "` + testServiceHealthcheckURI + `",
		"interval": ` + fmt.Sprint(testServiceHealthcheckInterval) + `,
		"timeout": ` + fmt.Sprint(testServiceHealthcheckTimeout) + `,
		"port": ` + fmt.Sprint(testServiceHealthcheckPort) + `,
		"retries": ` + fmt.Sprint(testServiceHealthcheckRetries) + `
	  },
	  "healthcheck-status": null
	}
  ]
}`
	)

	require.NoError(t, json.Unmarshal([]byte(jsonNLB), &actualLoadBalancer))
	require.Equal(t, expectedLoadBalancer, actualLoadBalancer)
}
