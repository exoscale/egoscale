package publicapi

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSksCluster_UnmarshalJSON(t *testing.T) {
	var (
		testAddons                        = []SksClusterAddons{SksClusterAddonsExoscaleCloudController}
		testCNI                           = SksClusterCniCalico
		testCreatedAt, _                  = time.Parse(iso8601Format, "2020-11-16T10:41:58Z")
		testDescription                   = "Test Cluster description"
		testEndpoint                      = "df421958-3679-4e9c-afb9-02fb6f331301.sks-ch-gva-2.exo.io"
		testID                            = testRandomID(t)
		testLevel                         = SksClusterLevelPro
		testName                          = "test-cluster"
		testNodepoolCreatedAt, _          = time.Parse(iso8601Format, "2020-11-18T07:54:36Z")
		testNodepoolDescription           = "Test Nodepool description"
		testNodepoolDiskSize        int64 = 15
		testNodepoolID                    = testRandomID(t)
		testNodepoolInstancePoolID        = testRandomID(t)
		testNodepoolInstanceTypeID        = testRandomID(t)
		testNodepoolName                  = "test-nodepool"
		testNodepoolSecurityGroupID       = testRandomID(t)
		testNodepoolSize            int64 = 3
		testNodepoolState                 = SksNodepoolStateRunning
		testNodepoolTemplateID            = testRandomID(t)
		testNodepoolVersion               = "1.18.6"
		testState                         = SksClusterStateRunning
		testVersion                       = "1.18.6"

		expected = SksCluster{
			Addons:      &testAddons,
			Cni:         &testCNI,
			CreatedAt:   &testCreatedAt,
			Description: &testDescription,
			Endpoint:    &testEndpoint,
			Id:          &testID,
			Level:       &testLevel,
			Name:        &testName,
			Nodepools: &[]SksNodepool{{
				CreatedAt:      &testNodepoolCreatedAt,
				Description:    &testNodepoolDescription,
				DiskSize:       &testNodepoolDiskSize,
				Id:             &testNodepoolID,
				InstancePool:   &InstancePool{Id: &testNodepoolInstancePoolID},
				InstanceType:   &InstanceType{Id: &testNodepoolInstanceTypeID},
				Name:           &testNodepoolName,
				SecurityGroups: &[]SecurityGroup{{Id: &testNodepoolSecurityGroupID}},
				Size:           &testNodepoolSize,
				State:          &testNodepoolState,
				Template:       &Template{Id: &testNodepoolTemplateID},
				Version:        &testNodepoolVersion,
			}},
			State:   &testState,
			Version: &testVersion,
		}

		actual SksCluster

		jsonSksCluster = `{
  "addons": ["` + string(testAddons[0]) + `"],
  "cni": "` + string(testCNI) + `",
  "created-at": "` + testCreatedAt.Format(iso8601Format) + `",
  "description": "` + testDescription + `",
  "endpoint": "` + testEndpoint + `",
  "id": "` + testID + `",
  "level": "` + string(testLevel) + `",
  "name": "` + testName + `",
  "nodepools": [{
    "created-at": "` + testNodepoolCreatedAt.Format(iso8601Format) + `",
    "description": "` + testNodepoolDescription + `",
    "disk-size": ` + fmt.Sprint(testNodepoolDiskSize) + `,
    "id": "` + testNodepoolID + `",
    "instance-pool": {"id": "` + testNodepoolInstancePoolID + `"},
    "instance-type": {"id": "` + testNodepoolInstanceTypeID + `"},
    "name": "` + testNodepoolName + `",
    "security-groups": [{"id": "` + testNodepoolSecurityGroupID + `"}],
    "size": ` + fmt.Sprint(testNodepoolSize) + `,
    "state": "` + string(testNodepoolState) + `",
    "template": {"id": "` + testNodepoolTemplateID + `"},
    "version": "` + testNodepoolVersion + `"
  }],
  "state": "` + string(testState) + `",
  "version": "` + testVersion + `"
}`
	)

	require.NoError(t, json.Unmarshal([]byte(jsonSksCluster), &actual))
	require.Equal(t, expected, actual)
}

func TestSksCluster_MarshalJSON(t *testing.T) {
	var (
		testAddons                        = []SksClusterAddons{SksClusterAddonsExoscaleCloudController}
		testCNI                           = SksClusterCniCalico
		testCreatedAt, _                  = time.Parse(iso8601Format, "2020-11-16T10:41:58Z")
		testDescription                   = "Test Cluster description"
		testEndpoint                      = "df421958-3679-4e9c-afb9-02fb6f331301.sks-ch-gva-2.exo.io"
		testID                            = testRandomID(t)
		testLevel                         = SksClusterLevelPro
		testName                          = "test-cluster"
		testNodepoolCreatedAt, _          = time.Parse(iso8601Format, "2020-11-18T07:54:36Z")
		testNodepoolDescription           = "Test Nodepool description"
		testNodepoolDiskSize        int64 = 15
		testNodepoolID                    = testRandomID(t)
		testNodepoolInstancePoolID        = testRandomID(t)
		testNodepoolInstanceTypeID        = testRandomID(t)
		testNodepoolName                  = "test-nodepool"
		testNodepoolSecurityGroupID       = testRandomID(t)
		testNodepoolSize            int64 = 3
		testNodepoolState                 = SksNodepoolStateRunning
		testNodepoolTemplateID            = testRandomID(t)
		testNodepoolVersion               = "1.18.6"
		testState                         = SksClusterStateRunning
		testVersion                       = "1.18.6"

		sksCluster = SksCluster{
			Addons:      &testAddons,
			Cni:         &testCNI,
			CreatedAt:   &testCreatedAt,
			Description: &testDescription,
			Endpoint:    &testEndpoint,
			Id:          &testID,
			Level:       &testLevel,
			Name:        &testName,
			Nodepools: &[]SksNodepool{{
				CreatedAt:      &testNodepoolCreatedAt,
				Description:    &testNodepoolDescription,
				DiskSize:       &testNodepoolDiskSize,
				Id:             &testNodepoolID,
				InstancePool:   &InstancePool{Id: &testNodepoolInstancePoolID},
				InstanceType:   &InstanceType{Id: &testNodepoolInstanceTypeID},
				Name:           &testNodepoolName,
				SecurityGroups: &[]SecurityGroup{{Id: &testNodepoolSecurityGroupID}},
				Size:           &testNodepoolSize,
				State:          &testNodepoolState,
				Template:       &Template{Id: &testNodepoolTemplateID},
				Version:        &testNodepoolVersion,
			}},
			State:   &testState,
			Version: &testVersion,
		}

		expected = []byte(`{` +
			`"addons":["` + string(testAddons[0]) + `"],` +
			`"cni":"` + string(testCNI) + `",` +
			`"created-at":"` + testCreatedAt.Format(iso8601Format) + `",` +
			`"description":"` + testDescription + `",` +
			`"endpoint":"` + testEndpoint + `",` +
			`"id":"` + testID + `",` +
			`"level":"` + string(testLevel) + `",` +
			`"name":"` + testName + `",` +
			`"nodepools":[{` +
			`"created-at":"` + testNodepoolCreatedAt.Format(iso8601Format) + `",` +
			`"description":"` + testNodepoolDescription + `",` +
			`"disk-size":` + fmt.Sprint(testNodepoolDiskSize) + `,` +
			`"id":"` + testNodepoolID + `",` +
			`"instance-pool":{"id":"` + testNodepoolInstancePoolID + `"},` +
			`"instance-type":{"id":"` + testNodepoolInstanceTypeID + `"},` +
			`"name":"` + testNodepoolName + `",` +
			`"security-groups":[{"id":"` + testNodepoolSecurityGroupID + `"}],` +
			`"size":` + fmt.Sprint(testNodepoolSize) + `,` +
			`"state":"` + string(testNodepoolState) + `",` +
			`"template":{"id":"` + testNodepoolTemplateID + `"},` +
			`"version":"` + testNodepoolVersion + `"` +
			`}],` +
			`"state":"` + string(testState) + `",` +
			`"version":"` + testVersion + `"` +
			`}`)
	)

	actual, err := json.Marshal(sksCluster)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}
