package v2

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSksNodepool_UnmarshalJSON(t *testing.T) {
	var (
		testID                   = "c19542b7-d269-4bd4-bf7c-2cae36d066d3"
		testName                 = "test-nodepool"
		testCreatedAt, _         = time.Parse(iso8601Format, "2020-08-12T11:12:36Z")
		testDescription          = "Test Nodepool description"
		testSize           int64 = 3
		testDiskSize       int64 = 15
		testInstancePoolID       = "f1f67118-43b6-4632-a709-d55fada62f21"
		testInstanceTypeID       = "21624abb-764e-4def-81d7-9fc54b5957fb"
		// testSecurityGroupID       = "efb4f4df-87ce-44e9-b5ee-59a9c1628edf"
		testState      = "running"
		testTemplateID = "f270d9a2-db64-4e8e-9cd3-5125887e91aa"
		testVersion    = "1.18.6"

		expected = SksNodepool{
			CreatedAt:    &testCreatedAt,
			Description:  &testDescription,
			DiskSize:     &testDiskSize,
			Id:           &testID,
			InstancePool: &Resource{Id: &testInstancePoolID},
			InstanceType: &InstanceType{Id: &testInstanceTypeID},
			Name:         &testName,
			// SecurityGroups: nil, // TODO: the API doesn't return this field ATM
			Size:     &testSize,
			State:    &testState,
			Template: &Template{Id: &testTemplateID},
			Version:  &testVersion,
		}

		actual SksNodepool

		jsonSksNodepool = `{
  "created-at": "` + testCreatedAt.Format(iso8601Format) + `",
  "description": "` + testDescription + `",
  "disk-size": ` + fmt.Sprint(testDiskSize) + `,
  "id": "` + testID + `",
  "instance-pool": {"id": "` + testInstancePoolID + `"},
  "instance-type": {"id": "` + testInstanceTypeID + `"},
  "name": "` + testName + `",
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
		testID                   = "c19542b7-d269-4bd4-bf7c-2cae36d066d3"
		testName                 = "test-nodepool"
		testCreatedAt, _         = time.Parse(iso8601Format, "2020-08-12T11:12:36Z")
		testDescription          = "Test Nodepool description"
		testSize           int64 = 3
		testDiskSize       int64 = 15
		testInstancePoolID       = "f1f67118-43b6-4632-a709-d55fada62f21"
		testInstanceTypeID       = "21624abb-764e-4def-81d7-9fc54b5957fb"
		// testSecurityGroupID       = "efb4f4df-87ce-44e9-b5ee-59a9c1628edf"
		testState      = "running"
		testTemplateID = "f270d9a2-db64-4e8e-9cd3-5125887e91aa"
		testVersion    = "1.18.6"

		sksNodepool = SksNodepool{
			CreatedAt:    &testCreatedAt,
			Description:  &testDescription,
			DiskSize:     &testDiskSize,
			Id:           &testID,
			InstancePool: &Resource{Id: &testInstancePoolID},
			InstanceType: &InstanceType{Id: &testInstanceTypeID},
			Name:         &testName,
			// SecurityGroups: nil, // TODO: the API doesn't return this field ATM
			Size:     &testSize,
			State:    &testState,
			Template: &Template{Id: &testTemplateID},
			Version:  &testVersion,
		}

		expected = []byte(`{` +
			`"created-at":"` + testCreatedAt.Format(iso8601Format) + `",` +
			`"description":"` + testDescription + `",` +
			`"disk-size":` + fmt.Sprint(testDiskSize) + `,` +
			`"id":"` + testID + `",` +
			`"instance-pool":{"id":"` + testInstancePoolID + `"},` +
			`"instance-type":{"id":"` + testInstanceTypeID + `"},` +
			`"name":"` + testName + `",` +
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
