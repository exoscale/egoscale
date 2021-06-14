package publicapi

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDbaasService_UnmarshalJSON(t *testing.T) {
	var (
		testBackupTime, _               = time.Parse(iso8601Format, "2020-08-12T11:12:36Z")
		testBackupName                  = testBackupTime.Format("2006-01-02_15-04_0.00000000.pghoard")
		testBackupDataSize        int64 = 36259840
		testCreatedAt, _                = time.Parse(iso8601Format, "2020-08-12T11:12:36Z")
		testDescription                 = testRandomString(10)
		testDiskSize              int64 = 10995116277760
		testMaintenanceDOW              = DbaasServiceMaintenanceDowSunday
		testMaintenanceTime             = "01:23:45"
		testName                        = DbaasServiceName("test")
		testNodeCount             int64 = 1
		testNodeCPUCount          int64 = 2
		testNodeMemory            int64 = 2199023255552
		testNodeStateRole               = DbaasNodeStateRoleMaster
		testNodeStateState              = DbaasNodeStateStateRunning
		testPlan                        = "hobbyist-1"
		testState                       = DbaasServiceStateRunning
		testTerminationProtection       = true
		testType                        = DbaasServiceTypeNamePg
		testUserType                    = "primary"
		testUserUsername                = testRandomString(10)
		testUserPassword                = testRandomString(10)
		testURI                         = "postgres://username:password@host:port/dbname?sslmode=required"
		testUpdatedAt, _                = time.Parse(iso8601Format, "2020-08-12T11:12:36Z")

		expected = DbaasService{
			Acl: &[]DbaasServiceAcl{},
			Backups: &[]DbaasServiceBackup{{
				BackupName: testBackupName,
				BackupTime: testBackupTime,
				DataSize:   testBackupDataSize,
			}},
			Components: &[]DbaasServiceComponents{{
				Component: "pg",
				Host:      "host",
				Port:      12345,
				Route:     DbaasServiceComponentsRouteDynamic,
				Usage:     DbaasServiceComponentsUsagePrimary,
			}},
			ConnectionInfo:  &DbaasService_ConnectionInfo{AdditionalProperties: map[string]interface{}{"k": "v"}},
			ConnectionPools: &[]DbaasServiceConnectionPools{},
			CreatedAt:       &testCreatedAt,
			Description:     &testDescription,
			DiskSize:        &testDiskSize,
			Features:        &DbaasService_Features{AdditionalProperties: map[string]interface{}{"k": "v"}},
			Integrations:    &[]DbaasServiceIntegration{},
			Maintenance: &DbaasServiceMaintenance{
				Dow:     testMaintenanceDOW,
				Time:    testMaintenanceTime,
				Updates: []DbaasServiceUpdate{},
			},
			Metadata:     &DbaasService_Metadata{AdditionalProperties: map[string]interface{}{"k": "v"}},
			Name:         testName,
			NodeCount:    &testNodeCount,
			NodeCpuCount: &testNodeCPUCount,
			NodeMemory:   &testNodeMemory,
			NodeStates: &[]DbaasNodeState{{
				Name:            testRandomString(10),
				ProgressUpdates: &[]DbaasNodeStateProgressUpdate{},
				Role:            &testNodeStateRole,
				State:           testNodeStateState,
			}},
			Notifications:         &[]DbaasServiceNotification{},
			Plan:                  testPlan,
			State:                 &testState,
			TerminationProtection: &testTerminationProtection,
			Type:                  testType,
			UpdatedAt:             &testUpdatedAt,
			Uri:                   &testURI,
			UriParams:             &DbaasService_UriParams{AdditionalProperties: map[string]interface{}{"k": "v"}},
			UserConfig:            &DbaasService_UserConfig{AdditionalProperties: map[string]interface{}{"k": "v"}},
			Users: &[]DbaasServiceUser{{
				Password: &testUserPassword,
				Type:     testUserType,
				Username: testUserUsername,
			}},
		}

		actual DbaasService

		jsonDbaasService = `{
  "acl": [],
  "backups": [
    {
      "backup-name": "` + testBackupName + `",
      "backup-time": "` + testBackupTime.Format(iso8601Format) + `",
      "data-size": ` + fmt.Sprint(testBackupDataSize) + `
    }
  ],
  "components": [
    {
      "component": "` + (*expected.Components)[0].Component + `",
      "host": "` + (*expected.Components)[0].Host + `",
      "port": ` + fmt.Sprint((*expected.Components)[0].Port) + `,
      "route": "` + fmt.Sprint((*expected.Components)[0].Route) + `",
      "usage": "` + fmt.Sprint((*expected.Components)[0].Usage) + `"
    }
  ],
  "connection-info": {
    "k": "v"
  },
  "connection-pools": [],
  "created-at": "` + testCreatedAt.Format(iso8601Format) + `",
  "description": "` + testDescription + `",
  "disk-size": ` + fmt.Sprint(testDiskSize) + `,
  "features": {
    "k": "v"
  },
  "integrations": [],
  "maintenance": {
    "dow": "` + fmt.Sprint(testMaintenanceDOW) + `",
    "time": "` + testMaintenanceTime + `",
    "updates": []
  },
  "metadata": {
    "k": "v"
  },
  "name": "` + fmt.Sprint(testName) + `",
  "node-count": ` + fmt.Sprint(testNodeCount) + `,
  "node-cpu-count": ` + fmt.Sprint(testNodeCPUCount) + `,
  "node-memory": ` + fmt.Sprint(testNodeMemory) + `,
  "node-states": [
    {
      "name": "` + (*expected.NodeStates)[0].Name + `",
      "progress-updates": [],
      "role": "` + fmt.Sprint(*(*expected.NodeStates)[0].Role) + `",
      "state": "` + fmt.Sprint((*expected.NodeStates)[0].State) + `"
    }
  ],
  "notifications": [],
  "plan": "` + testPlan + `",
  "state": "` + fmt.Sprint(testState) + `",
  "termination-protection": ` + fmt.Sprint(testTerminationProtection) + `,
  "type": "` + fmt.Sprint(testType) + `",
  "updated-at": "` + testUpdatedAt.Format(iso8601Format) + `",
  "uri": "` + testURI + `",
  "uri-params": {
    "k": "v"
  },
  "user-config": {
    "k": "v"
  },
  "users": [
    {
      "password": "` + testUserPassword + `",
      "type": "` + testUserType + `",
      "username": "` + testUserUsername + `"
    }
  ]
}`
	)

	require.NoError(t, json.Unmarshal([]byte(jsonDbaasService), &actual))
	require.Equal(t, expected, actual)
}

func TestDbaasService_MarshalJSON(t *testing.T) {
	var (
		testBackupTime, _               = time.Parse(iso8601Format, "2020-08-12T11:12:36Z")
		testBackupName                  = testBackupTime.Format("2006-01-02_15-04_0.00000000.pghoard")
		testBackupDataSize        int64 = 36259840
		testCreatedAt, _                = time.Parse(iso8601Format, "2020-08-12T11:12:36Z")
		testDescription                 = "Test Database"
		testDiskSize              int64 = 10995116277760
		testMaintenanceDOW              = DbaasServiceMaintenanceDowSunday
		testMaintenanceTime             = "04:00"
		testName                        = DbaasServiceName("test")
		testNodeCount             int64 = 1
		testNodeCPUCount          int64 = 2
		testNodeMemory            int64 = 2199023255552
		testNodeStateRole               = DbaasNodeStateRoleMaster
		testNodeStateState              = DbaasNodeStateStateRunning
		testPlan                        = "hobbyist-1"
		testState                       = DbaasServiceStateRunning
		testTerminationProtection       = true
		testType                        = DbaasServiceTypeNamePg
		testUserType                    = "primary"
		testUserUsername                = "test-user"
		testUserPassword                = "test-password"
		testURI                         = "postgres://username:password@host:port/dbname?sslmode=required"
		testUpdatedAt, _                = time.Parse(iso8601Format, "2020-08-12T11:12:36Z")

		dbaasService = DbaasService{
			Acl: &[]DbaasServiceAcl{},
			Backups: &[]DbaasServiceBackup{{
				BackupName: testBackupName,
				BackupTime: testBackupTime,
				DataSize:   testBackupDataSize,
			}},
			Components: &[]DbaasServiceComponents{{
				Component: "pg",
				Host:      "host",
				Port:      12345,
				Route:     DbaasServiceComponentsRouteDynamic,
				Usage:     DbaasServiceComponentsUsagePrimary,
			}},
			ConnectionInfo:  &DbaasService_ConnectionInfo{AdditionalProperties: map[string]interface{}{"k": "v"}},
			ConnectionPools: &[]DbaasServiceConnectionPools{},
			CreatedAt:       &testCreatedAt,
			Description:     &testDescription,
			DiskSize:        &testDiskSize,
			Features:        &DbaasService_Features{AdditionalProperties: map[string]interface{}{"k": "v"}},
			Integrations:    &[]DbaasServiceIntegration{},
			Maintenance: &DbaasServiceMaintenance{
				Dow:     testMaintenanceDOW,
				Time:    testMaintenanceTime,
				Updates: []DbaasServiceUpdate{},
			},
			Metadata:     &DbaasService_Metadata{AdditionalProperties: map[string]interface{}{"k": "v"}},
			Name:         testName,
			NodeCount:    &testNodeCount,
			NodeCpuCount: &testNodeCPUCount,
			NodeMemory:   &testNodeMemory,
			NodeStates: &[]DbaasNodeState{{
				Name:            "test-1",
				ProgressUpdates: &[]DbaasNodeStateProgressUpdate{},
				Role:            &testNodeStateRole,
				State:           testNodeStateState,
			}},
			Notifications:         &[]DbaasServiceNotification{},
			Plan:                  testPlan,
			State:                 &testState,
			TerminationProtection: &testTerminationProtection,
			Type:                  testType,
			UpdatedAt:             &testUpdatedAt,
			Uri:                   &testURI,
			UriParams:             &DbaasService_UriParams{AdditionalProperties: map[string]interface{}{"k": "v"}},
			UserConfig:            &DbaasService_UserConfig{AdditionalProperties: map[string]interface{}{"k": "v"}},
			Users: &[]DbaasServiceUser{{
				Password: &testUserPassword,
				Type:     testUserType,
				Username: testUserUsername,
			}},
		}

		expected = []byte(`{` +
			`"acl":[],` +
			`"backups":[{` +
			`"backup-name":"` + testBackupName + `",` +
			`"backup-time":"` + testBackupTime.Format(iso8601Format) + `",` +
			`"data-size":` + fmt.Sprint(testBackupDataSize) +
			`}],` +
			`"components":[{` +
			`"component":"` + (*dbaasService.Components)[0].Component + `",` +
			`"host":"` + (*dbaasService.Components)[0].Host + `",` +
			`"port":` + fmt.Sprint((*dbaasService.Components)[0].Port) + `,` +
			`"route":"` + fmt.Sprint((*dbaasService.Components)[0].Route) + `",` +
			`"usage":"` + fmt.Sprint((*dbaasService.Components)[0].Usage) + `"` +
			`}],` +
			`"connection-info":{"k":"v"},` +
			`"connection-pools":[],` +
			`"created-at":"` + testCreatedAt.Format(iso8601Format) + `",` +
			`"description":"` + testDescription + `",` +
			`"disk-size":` + fmt.Sprint(testDiskSize) + `,` +
			`"features":{"k":"v"},` +
			`"integrations":[],` +
			`"maintenance":{` +
			`"dow":"` + fmt.Sprint(testMaintenanceDOW) + `",` +
			`"time":"` + testMaintenanceTime + `",` +
			`"updates":[]` +
			`},` +
			`"metadata":{"k":"v"},` +
			`"name":"` + fmt.Sprint(testName) + `",` +
			`"node-count":` + fmt.Sprint(testNodeCount) + `,` +
			`"node-cpu-count":` + fmt.Sprint(testNodeCPUCount) + `,` +
			`"node-memory":` + fmt.Sprint(testNodeMemory) + `,` +
			`"node-states":[{` +
			`"name":"` + (*dbaasService.NodeStates)[0].Name + `",` +
			`"progress-updates":[],` +
			`"role":"` + fmt.Sprint(*(*dbaasService.NodeStates)[0].Role) + `",` +
			`"state":"` + fmt.Sprint((*dbaasService.NodeStates)[0].State) + `"` +
			`}],` +
			`"notifications":[],` +
			`"plan":"` + testPlan + `",` +
			`"state":"` + fmt.Sprint(testState) + `",` +
			`"termination-protection":` + fmt.Sprint(testTerminationProtection) + `,` +
			`"type":"` + fmt.Sprint(testType) + `",` +
			`"updated-at":"` + testUpdatedAt.Format(iso8601Format) + `",` +
			`"uri":"` + testURI + `",` +
			`"uri-params":{"k":"v"},` +
			`"user-config":{"k":"v"},` +
			`"users":[{` +
			`"password":"` + testUserPassword + `",` +
			`"type":"` + testUserType + `",` +
			`"username":"` + testUserUsername + `"` +
			`}]` +
			`}`)
	)

	actual, err := json.Marshal(dbaasService)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}
