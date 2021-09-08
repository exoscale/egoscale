package oapi

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSksNodepool_UnmarshalJSON(t *testing.T) {
	var (
		testAddons                    = []SksNodepoolAddons{SksNodepoolAddonsLinbit}
		testAntiAffinityGroupID       = testRandomID(t)
		testCreatedAt, _              = time.Parse(iso8601Format, "2020-08-12T11:12:36Z")
		testDeployTargetID            = testRandomID(t)
		testDescription               = "Test Nodepool description"
		testDiskSize            int64 = 15
		testID                        = testRandomID(t)
		testInstancePoolID            = testRandomID(t)
		testInstancePrefix            = "test-nodepool"
		testInstanceTypeID            = testRandomID(t)
		testLabels                    = map[string]string{"k1": "v1", "k2": "v2"}
		testName                      = "test-nodepool"
		testPrivateNetworkID          = testRandomID(t)
		testSecurityGroupID           = testRandomID(t)
		testSize                int64 = 3
		testState                     = SksNodepoolStateRunning
		testTemplateID                = testRandomID(t)
		testVersion                   = "1.18.6"

		expected = SksNodepool{
			Addons:             &testAddons,
			AntiAffinityGroups: &[]AntiAffinityGroup{{Id: &testAntiAffinityGroupID}},
			CreatedAt:          &testCreatedAt,
			DeployTarget:       &DeployTarget{Id: &testDeployTargetID},
			Description:        &testDescription,
			DiskSize:           &testDiskSize,
			Id:                 &testID,
			InstancePrefix:     &testInstancePrefix,
			InstancePool:       &InstancePool{Id: &testInstancePoolID},
			InstanceType:       &InstanceType{Id: &testInstanceTypeID},
			Labels:             &Labels{AdditionalProperties: testLabels},
			Name:               &testName,
			PrivateNetworks:    &[]PrivateNetwork{{Id: &testPrivateNetworkID}},
			SecurityGroups:     &[]SecurityGroup{{Id: &testSecurityGroupID}},
			Size:               &testSize,
			State:              &testState,
			Template:           &Template{Id: &testTemplateID},
			Version:            &testVersion,
		}

		actual SksNodepool

		jsonSksNodepool = `{
  "addons": ["` + string(testAddons[0]) + `"],
  "anti-affinity-groups": [{"id":"` + testAntiAffinityGroupID + `"}],
  "created-at": "` + testCreatedAt.Format(iso8601Format) + `",
  "deploy-target": {"id": "` + testDeployTargetID + `"},
  "description": "` + testDescription + `",
  "disk-size": ` + fmt.Sprint(testDiskSize) + `,
  "id": "` + testID + `",
  "instance-pool": {"id": "` + testInstancePoolID + `"},
  "instance-prefix": "` + testInstancePrefix + `",
  "instance-type": {"id": "` + testInstanceTypeID + `"},
  "labels": {"k1": "` + testLabels["k1"] + `", "k2": "` + testLabels["k2"] + `"},
  "name": "` + testName + `",
  "private-networks": [{"id":"` + testPrivateNetworkID + `"}],
  "security-groups": [{"id":"` + testSecurityGroupID + `"}],
  "size": ` + fmt.Sprint(testSize) + `,
  "state": "` + string(testState) + `",
  "template": {"id": "` + testTemplateID + `"},
  "version": "` + testVersion + `"
}`
	)

	require.NoError(t, json.Unmarshal([]byte(jsonSksNodepool), &actual))
	require.Equal(t, expected, actual)
}

func TestSksNodepool_MarshalJSON(t *testing.T) {
	var (
		testAddons                    = []SksNodepoolAddons{SksNodepoolAddonsLinbit}
		testAntiAffinityGroupID       = testRandomID(t)
		testCreatedAt, _              = time.Parse(iso8601Format, "2020-08-12T11:12:36Z")
		testDeployTargetID            = testRandomID(t)
		testDescription               = "Test Nodepool description"
		testDiskSize            int64 = 15
		testID                        = testRandomID(t)
		testInstancePoolID            = testRandomID(t)
		testInstancePrefix            = "test-nodepool"
		testInstanceTypeID            = testRandomID(t)
		testLabels                    = map[string]string{"k1": "v1", "k2": "v2"}
		testName                      = "test-nodepool"
		testPrivateNetworkID          = testRandomID(t)
		testSecurityGroupID           = testRandomID(t)
		testSize                int64 = 3
		testState                     = SksNodepoolStateRunning
		testTemplateID                = testRandomID(t)
		testVersion                   = "1.18.6"

		sksNodepool = SksNodepool{
			Addons:             &testAddons,
			AntiAffinityGroups: &[]AntiAffinityGroup{{Id: &testAntiAffinityGroupID}},
			CreatedAt:          &testCreatedAt,
			DeployTarget:       &DeployTarget{Id: &testDeployTargetID},
			Description:        &testDescription,
			DiskSize:           &testDiskSize,
			Id:                 &testID,
			InstancePrefix:     &testInstancePrefix,
			InstancePool:       &InstancePool{Id: &testInstancePoolID},
			InstanceType:       &InstanceType{Id: &testInstanceTypeID},
			Labels:             &Labels{AdditionalProperties: testLabels},
			Name:               &testName,
			PrivateNetworks:    &[]PrivateNetwork{{Id: &testPrivateNetworkID}},
			SecurityGroups:     &[]SecurityGroup{{Id: &testSecurityGroupID}},
			Size:               &testSize,
			State:              &testState,
			Template:           &Template{Id: &testTemplateID},
			Version:            &testVersion,
		}

		expected = []byte(`{` +
			`"addons":["` + string(testAddons[0]) + `"],` +
			`"anti-affinity-groups":[{"id":"` + testAntiAffinityGroupID + `"}],` +
			`"created-at":"` + testCreatedAt.Format(iso8601Format) + `",` +
			`"deploy-target":{"id":"` + testDeployTargetID + `"},` +
			`"description":"` + testDescription + `",` +
			`"disk-size":` + fmt.Sprint(testDiskSize) + `,` +
			`"id":"` + testID + `",` +
			`"instance-pool":{"id":"` + testInstancePoolID + `"},` +
			`"instance-prefix":"` + testInstancePrefix + `",` +
			`"instance-type":{"id":"` + testInstanceTypeID + `"},` +
			`"labels":{"k1":"` + testLabels["k1"] + `","k2":"` + testLabels["k2"] + `"},` +
			`"name":"` + testName + `",` +
			`"private-networks":[{"id":"` + testPrivateNetworkID + `"}],` +
			`"security-groups":[{"id":"` + testSecurityGroupID + `"}],` +
			`"size":` + fmt.Sprint(testSize) + `,` +
			`"state":"` + string(testState) + `",` +
			`"template":{"id":"` + testTemplateID + `"},` +
			`"version":"` + testVersion + `"` +
			`}`)
	)

	actual, err := json.Marshal(sksNodepool)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}
