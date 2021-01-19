package v2

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSksCluster_UnmarshalJSON(t *testing.T) {
	var (
		testAddons                        = []string{"exoscale-cloud-controller"}
		testCNI                           = "calico"
		testCreatedAt, _                  = time.Parse(iso8601Format, "2020-11-16T10:41:58Z")
		testDescription                   = "Test Cluster description"
		testEndpoint                      = "df421958-3679-4e9c-afb9-02fb6f331301.sks-ch-gva-2.exo.io"
		testID                            = "df421958-3679-4e9c-afb9-02fb6f331301"
		testLevel                         = "pro"
		testName                          = "test-cluster"
		testNodepoolCreatedAt, _          = time.Parse(iso8601Format, "2020-11-18T07:54:36Z")
		testNodepoolDescription           = "Test Nodepool description"
		testNodepoolDiskSize        int64 = 15
		testNodepoolID                    = "6d1eecee-397c-4e16-b103-2d1353bf4ecc"
		testNodepoolInstancePoolID        = "f1f67118-43b6-4632-a709-d55fada62f21"
		testNodepoolInstanceTypeID        = "21624abb-764e-4def-81d7-9fc54b5957fb"
		testNodepoolName                  = "test-nodepool"
		testNodepoolSecurityGroupID       = "efb4f4df-87ce-44e9-b5ee-59a9c1628edf"
		testNodepoolSize            int64 = 3
		testNodepoolState                 = "running"
		testNodepoolTemplateID            = "f270d9a2-db64-4e8e-9cd3-5125887e91aa"
		testNodepoolVersion               = "1.18.6"
		testState                         = "running"
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
  "addons": ["` + testAddons[0] + `"],
  "cni": "` + testCNI + `",
  "created-at": "` + testCreatedAt.Format(iso8601Format) + `",
  "description": "` + testDescription + `",
  "endpoint": "` + testEndpoint + `",
  "id": "` + testID + `",
  "level": "` + testLevel + `",
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
    "state": "` + testNodepoolState + `",
    "template": {"id": "` + testNodepoolTemplateID + `"},
    "version": "` + testNodepoolVersion + `"
  }],
  "state": "` + testState + `",
  "version": "` + testVersion + `"
}`
	)

	require.NoError(t, json.Unmarshal([]byte(jsonSksCluster), &actual))
	require.Equal(t, expected, actual)
}

func TestSksCluster_MarshalJSON(t *testing.T) {
	var (
		testAddons                        = []string{"exoscale-cloud-controller"}
		testCNI                           = "calico"
		testCreatedAt, _                  = time.Parse(iso8601Format, "2020-11-16T10:41:58Z")
		testDescription                   = "Test Cluster description"
		testEndpoint                      = "df421958-3679-4e9c-afb9-02fb6f331301.sks-ch-gva-2.exo.io"
		testID                            = "df421958-3679-4e9c-afb9-02fb6f331301"
		testLevel                         = "pro"
		testName                          = "test-cluster"
		testNodepoolCreatedAt, _          = time.Parse(iso8601Format, "2020-11-18T07:54:36Z")
		testNodepoolDescription           = "Test Nodepool description"
		testNodepoolDiskSize        int64 = 15
		testNodepoolID                    = "6d1eecee-397c-4e16-b103-2d1353bf4ecc"
		testNodepoolInstancePoolID        = "f1f67118-43b6-4632-a709-d55fada62f21"
		testNodepoolInstanceTypeID        = "21624abb-764e-4def-81d7-9fc54b5957fb"
		testNodepoolName                  = "test-nodepool"
		testNodepoolSecurityGroupID       = "efb4f4df-87ce-44e9-b5ee-59a9c1628edf"
		testNodepoolSize            int64 = 3
		testNodepoolState                 = "running"
		testNodepoolTemplateID            = "f270d9a2-db64-4e8e-9cd3-5125887e91aa"
		testNodepoolVersion               = "1.18.6"
		testState                         = "running"
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
			`"addons":["` + testAddons[0] + `"],` +
			`"cni":"` + testCNI + `",` +
			`"created-at":"` + testCreatedAt.Format(iso8601Format) + `",` +
			`"description":"` + testDescription + `",` +
			`"endpoint":"` + testEndpoint + `",` +
			`"id":"` + testID + `",` +
			`"level":"` + testLevel + `",` +
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
			`"state":"` + testNodepoolState + `",` +
			`"template":{"id":"` + testNodepoolTemplateID + `"},` +
			`"version":"` + testNodepoolVersion + `"` +
			`}],` +
			`"state":"` + testState + `",` +
			`"version":"` + testVersion + `"` +
			`}`)
	)

	actual, err := json.Marshal(sksCluster)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}
