package publicapi

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSksNodepool_UnmarshalJSON(t *testing.T) {
	var (
		testAntiAffinityGroupID       = testRandomID(t)
		testCreatedAt, _              = time.Parse(iso8601Format, "2020-08-12T11:12:36Z")
		testDescription               = "Test Nodepool description"
		testDiskSize            int64 = 15
		testID                        = testRandomID(t)
		testInstancePoolID            = testRandomID(t)
		testInstanceTypeID            = testRandomID(t)
		testName                      = "test-nodepool"
		testSecurityGroupID           = testRandomID(t)
		testSize                int64 = 3
		testState                     = "running"
		testTemplateID                = testRandomID(t)
		testVersion                   = "1.18.6"

		expected = SksNodepool{
			AntiAffinityGroups: &[]AntiAffinityGroup{{Id: &testAntiAffinityGroupID}},
			CreatedAt:          &testCreatedAt,
			Description:        &testDescription,
			DiskSize:           &testDiskSize,
			Id:                 &testID,
			InstancePool:       &InstancePool{Id: &testInstancePoolID},
			InstanceType:       &InstanceType{Id: &testInstanceTypeID},
			Name:               &testName,
			SecurityGroups:     &[]SecurityGroup{{Id: &testSecurityGroupID}},
			Size:               &testSize,
			State:              &testState,
			Template:           &Template{Id: &testTemplateID},
			Version:            &testVersion,
		}

		actual SksNodepool

		jsonSksNodepool = `{
  "anti-affinity-groups": [{"id":"` + testAntiAffinityGroupID + `"}],
  "created-at": "` + testCreatedAt.Format(iso8601Format) + `",
  "description": "` + testDescription + `",
  "disk-size": ` + fmt.Sprint(testDiskSize) + `,
  "id": "` + testID + `",
  "instance-pool": {"id": "` + testInstancePoolID + `"},
  "instance-type": {"id": "` + testInstanceTypeID + `"},
  "name": "` + testName + `",
  "security-groups": [{"id":"` + testSecurityGroupID + `"}],
  "size": ` + fmt.Sprint(testSize) + `,
  "state": "` + testState + `",
  "template": {"id": "` + testTemplateID + `"},
  "version": "` + testVersion + `"
}`
	)

	require.NoError(t, json.Unmarshal([]byte(jsonSksNodepool), &actual))
	require.Equal(t, expected, actual)
}

func TestSksNodepool_MarshalJSON(t *testing.T) {
	var (
		testAntiAffinityGroupID       = testRandomID(t)
		testCreatedAt, _              = time.Parse(iso8601Format, "2020-08-12T11:12:36Z")
		testDescription               = "Test Nodepool description"
		testDiskSize            int64 = 15
		testID                        = testRandomID(t)
		testInstancePoolID            = testRandomID(t)
		testInstanceTypeID            = testRandomID(t)
		testName                      = "test-nodepool"
		testSecurityGroupID           = testRandomID(t)
		testSize                int64 = 3
		testState                     = "running"
		testTemplateID                = testRandomID(t)
		testVersion                   = "1.18.6"

		sksNodepool = SksNodepool{
			AntiAffinityGroups: &[]AntiAffinityGroup{{Id: &testAntiAffinityGroupID}},
			CreatedAt:          &testCreatedAt,
			Description:        &testDescription,
			DiskSize:           &testDiskSize,
			Id:                 &testID,
			InstancePool:       &InstancePool{Id: &testInstancePoolID},
			InstanceType:       &InstanceType{Id: &testInstanceTypeID},
			Name:               &testName,
			SecurityGroups:     &[]SecurityGroup{{Id: &testSecurityGroupID}},
			Size:               &testSize,
			State:              &testState,
			Template:           &Template{Id: &testTemplateID},
			Version:            &testVersion,
		}

		expected = []byte(`{` +
			`"anti-affinity-groups":[{"id":"` + testAntiAffinityGroupID + `"}],` +
			`"created-at":"` + testCreatedAt.Format(iso8601Format) + `",` +
			`"description":"` + testDescription + `",` +
			`"disk-size":` + fmt.Sprint(testDiskSize) + `,` +
			`"id":"` + testID + `",` +
			`"instance-pool":{"id":"` + testInstancePoolID + `"},` +
			`"instance-type":{"id":"` + testInstanceTypeID + `"},` +
			`"name":"` + testName + `",` +
			`"security-groups":[{"id":"` + testSecurityGroupID + `"}],` +
			`"size":` + fmt.Sprint(testSize) + `,` +
			`"state":"` + testState + `",` +
			`"template":{"id":"` + testTemplateID + `"},` +
			`"version":"` + testVersion + `"` +
			`}`)
	)

	actual, err := json.Marshal(sksNodepool)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}
