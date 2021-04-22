package publicapi

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestInstance_UnmarshalJSON(t *testing.T) {
	var (
		testAntiAffinityGroupID       = testRandomID(t)
		testCreatedAt, _              = time.Parse(iso8601Format, "2020-08-12T11:12:36Z")
		testDiskSize            int64 = 15
		testElasticIPID               = testRandomID(t)
		testID                        = testRandomID(t)
		testInstanceTypeID            = testRandomID(t)
		testIpv6Address               = "2001:db8:abcd::1"
		testManagerID                 = testRandomID(t)
		testName                      = "test-instance"
		testPrivateNetworkID          = testRandomID(t)
		testSecurityGroupID           = testRandomID(t)
		testSnapshotID                = testRandomID(t)
		testSSHKeyName                = testRandomID(t)
		testState                     = "running"
		testTemplateID                = testRandomID(t)
		testUserData                  = "I2Nsb3VkLWNvbmZpZwphcHRfdXBncmFkZTogdHJ1ZQ=="

		expected = Instance{
			AntiAffinityGroups: &[]AntiAffinityGroup{{Id: &testAntiAffinityGroupID}},
			CreatedAt:          &testCreatedAt,
			DiskSize:           &testDiskSize,
			ElasticIps:         &[]ElasticIp{{Id: &testElasticIPID}},
			Id:                 &testID,
			InstanceType:       &InstanceType{Id: &testInstanceTypeID},
			Ipv6Address:        &testIpv6Address,
			Manager:            &Manager{Id: &testManagerID},
			Name:               &testName,
			PrivateNetworks:    &[]PrivateNetwork{{Id: &testPrivateNetworkID}},
			SecurityGroups:     &[]SecurityGroup{{Id: &testSecurityGroupID}},
			Snapshots:          &[]Snapshot{{Id: &testSnapshotID}},
			SshKey:             &SshKey{Name: &testSSHKeyName},
			State:              &testState,
			Template:           &Template{Id: &testTemplateID},
			UserData:           &testUserData,
		}

		actual Instance

		jsonInstance = `{
  "anti-affinity-groups": [{"id":"` + testAntiAffinityGroupID + `"}],
  "created-at": "` + testCreatedAt.Format(iso8601Format) + `",
  "disk-size": ` + fmt.Sprint(testDiskSize) + `,
  "elastic-ips": [{"id":"` + testElasticIPID + `"}],
  "id": "` + testID + `",
  "instance-type": {"id": "` + testInstanceTypeID + `"},
  "ipv6-address": "` + fmt.Sprint(testIpv6Address) + `",
  "manager": {"id": "` + testManagerID + `"},
  "name": "` + testName + `",
  "private-networks": [{"id":"` + testPrivateNetworkID + `"}],
  "security-groups": [{"id":"` + testSecurityGroupID + `"}],
  "snapshots": [{"id":"` + testSnapshotID + `"}],
  "ssh-key": {"name": "` + testSSHKeyName + `"},
  "state": "` + testState + `",
  "template": {"id": "` + testTemplateID + `"},
  "user-data": "` + testUserData + `"
}`
	)

	require.NoError(t, json.Unmarshal([]byte(jsonInstance), &actual))
	require.Equal(t, expected, actual)
}

func TestInstance_MarshalJSON(t *testing.T) {
	var (
		testAntiAffinityGroupID       = testRandomID(t)
		testCreatedAt, _              = time.Parse(iso8601Format, "2020-08-12T11:12:36Z")
		testDiskSize            int64 = 15
		testElasticIPID               = testRandomID(t)
		testID                        = testRandomID(t)
		testInstanceTypeID            = testRandomID(t)
		testIpv6Address               = "2001:db8:abcd::1"
		testManagerID                 = testRandomID(t)
		testName                      = "test-instance"
		testPrivateNetworkID          = testRandomID(t)
		testSecurityGroupID           = testRandomID(t)
		testSnapshotID                = testRandomID(t)
		testSSHKeyName                = testRandomID(t)
		testState                     = "running"
		testTemplateID                = testRandomID(t)
		testUserData                  = "I2Nsb3VkLWNvbmZpZwphcHRfdXBncmFkZTogdHJ1ZQ=="

		instance = Instance{
			AntiAffinityGroups: &[]AntiAffinityGroup{{Id: &testAntiAffinityGroupID}},
			CreatedAt:          &testCreatedAt,
			DiskSize:           &testDiskSize,
			ElasticIps:         &[]ElasticIp{{Id: &testElasticIPID}},
			Id:                 &testID,
			InstanceType:       &InstanceType{Id: &testInstanceTypeID},
			Ipv6Address:        &testIpv6Address,
			Manager:            &Manager{Id: &testManagerID},
			Name:               &testName,
			PrivateNetworks:    &[]PrivateNetwork{{Id: &testPrivateNetworkID}},
			SecurityGroups:     &[]SecurityGroup{{Id: &testSecurityGroupID}},
			Snapshots:          &[]Snapshot{{Id: &testSnapshotID}},
			SshKey:             &SshKey{Name: &testSSHKeyName},
			State:              &testState,
			Template:           &Template{Id: &testTemplateID},
			UserData:           &testUserData,
		}

		expected = []byte(`{` +
			`"anti-affinity-groups":[{"id":"` + testAntiAffinityGroupID + `"}],` +
			`"created-at":"` + testCreatedAt.Format(iso8601Format) + `",` +
			`"disk-size":` + fmt.Sprint(testDiskSize) + `,` +
			`"elastic-ips":[{"id":"` + testElasticIPID + `"}],` +
			`"id":"` + testID + `",` +
			`"instance-type":{"id":"` + testInstanceTypeID + `"},` +
			`"ipv6-address":"` + fmt.Sprint(testIpv6Address) + `",` +
			`"manager":{"id":"` + testManagerID + `"},` +
			`"name":"` + testName + `",` +
			`"private-networks":[{"id":"` + testPrivateNetworkID + `"}],` +
			`"security-groups":[{"id":"` + testSecurityGroupID + `"}],` +
			`"snapshots":[{"id":"` + testSnapshotID + `"}],` +
			`"ssh-key":{"name":"` + testSSHKeyName + `"},` +
			`"state":"` + testState + `",` +
			`"template":{"id":"` + testTemplateID + `"},` +
			`"user-data":"` + testUserData + `"` +
			`}`)
	)

	actual, err := json.Marshal(instance)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}
