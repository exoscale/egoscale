package v2

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/exoscale/egoscale/v2/oapi"
	"github.com/jarcoal/httpmock"
)

var (
	testDatabasePlanBackupConfigInterval     int64 = 24
	testDatabasePlanBackupConfigMaxCount     int64 = 2
	testDatabasePlanBackupConfigRecoveryMode       = "pitr"
	testDatabasePlanDiskSpace                int64 = 10737418240
	testDatabasePlanName                           = "hobbyist-1"
	testDatabasePlanNodeCPUCount             int64 = 2
	testDatabasePlanNodeCount                int64 = 1
	testDatabasePlanNodeMemory               int64 = 2147483648
	testDatabaseServiceBackupDataSize        int64 = 36259840
	testDatabaseServiceBackupName                  = testDatabaseServiceBackupTime.Format("2006-01-02_15-04_0.00000000.pghoard")
	testDatabaseServiceBackupTime, _               = time.Parse(iso8601Format, "2020-08-12T11:12:36Z")
	testDatabaseServiceComponent                   = "pgbouncer"
	testDatabaseServiceComponentHost               = new(clientTestSuite).randomString(30)
	testDatabaseServiceComponentPort         int64 = 12345
	testDatabaseServiceComponentRoute              = oapi.DbaasServiceComponentsRouteDynamic
	testDatabaseServiceComponentUsage              = oapi.DbaasServiceComponentsUsagePrimary
	testDatabaseServiceCreatedAt, _                = time.Parse(iso8601Format, "2020-08-12T11:12:36Z")
	testDatabaseServiceDescription                 = new(clientTestSuite).randomString(10)
	testDatabaseServiceDiskSize              int64 = 10995116277760
	testDatabaseServiceMaintenanceDOW              = oapi.DbaasServiceMaintenanceDowSunday
	testDatabaseServiceMaintenanceTime             = "01:23:45"
	testDatabaseServiceName                        = new(clientTestSuite).randomString(10)
	testDatabaseServiceNodeCPUCount          int64 = 2
	testDatabaseServiceNodeCount             int64 = 1
	testDatabaseServiceNodeMemory            int64 = 2199023255552
	testDatabaseServiceNodeStateRole               = oapi.DbaasNodeStateRoleMaster
	testDatabaseServiceNodeStateState              = oapi.DbaasNodeStateStateRunning
	testDatabaseServiceState                       = oapi.DbaasServiceStateRunning
	testDatabaseServiceTerminationProtection       = true
	testDatabaseServiceType                        = oapi.DbaasServiceTypeName("pg")
	testDatabaseServiceTypeDefaultVersion          = "13"
	testDatabaseServiceTypeDescription             = new(clientTestSuite).randomString(10)
	testDatabaseServiceTypeLatestVersion           = "13.3"
	testDatabaseServiceTypeName                    = oapi.DbaasServiceTypeName("pg")
	testDatabaseServiceURI                         = "postgres://username:password@host:port/dbname?sslmode=required"
	testDatabaseServiceUpdatedAt, _                = time.Parse(iso8601Format, "2020-08-12T11:12:36Z")
	testDatabaseServiceUserPassword                = new(clientTestSuite).randomString(10)
	testDatabaseServiceUserType                    = "primary"
	testDatabaseServiceUserUsername                = new(clientTestSuite).randomString(10)
)

func (ts *clientTestSuite) TestClient_CreateDatabaseService() {
	var (
		testID             = ts.randomID()
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
	)

	httpmock.RegisterResponder("POST", "/dbaas-service",
		func(req *http.Request) (*http.Response, error) {
			var actual oapi.CreateDbaasServiceJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := oapi.CreateDbaasServiceJSONRequestBody{
				Maintenance: &struct {
					Dow  oapi.CreateDbaasServiceJSONBodyMaintenanceDow `json:"dow"`
					Time string                                        `json:"time"`
				}{
					Dow:  oapi.CreateDbaasServiceJSONBodyMaintenanceDow(testDatabaseServiceMaintenanceDOW),
					Time: testDatabaseServiceMaintenanceTime,
				},
				Name:                  oapi.DbaasServiceName(testDatabaseServiceName),
				Plan:                  testDatabasePlanName,
				TerminationProtection: &testDatabaseServiceTerminationProtection,
				Type:                  testDatabaseServiceType,
				UserConfig: &oapi.CreateDbaasServiceJSONBody_UserConfig{
					AdditionalProperties: map[string]interface{}{"k": "v"},
				},
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/dbaas-service/%s", testDatabaseServiceName), oapi.DbaasService{
		Acl: &[]oapi.DbaasServiceAcl{},
		Backups: &[]oapi.DbaasServiceBackup{{
			BackupName: testDatabaseServiceBackupName,
			BackupTime: testDatabaseServiceBackupTime,
			DataSize:   testDatabaseServiceBackupDataSize,
		}},
		Components: &[]oapi.DbaasServiceComponents{{
			Component: testDatabaseServiceComponent,
			Host:      testDatabaseServiceComponentHost,
			Port:      testDatabaseServiceComponentPort,
			Route:     testDatabaseServiceComponentRoute,
			Usage:     testDatabaseServiceComponentUsage,
		}},
		ConnectionInfo:  &oapi.DbaasService_ConnectionInfo{AdditionalProperties: map[string]interface{}{"k": "v"}},
		ConnectionPools: &[]oapi.DbaasServiceConnectionPools{},
		CreatedAt:       &testDatabaseServiceCreatedAt,
		Description:     &testDatabaseServiceDescription,
		DiskSize:        &testDatabaseServiceDiskSize,
		Features:        &oapi.DbaasService_Features{AdditionalProperties: map[string]interface{}{"k": "v"}},
		Integrations:    &[]oapi.DbaasServiceIntegration{},
		Maintenance: &oapi.DbaasServiceMaintenance{
			Dow:     testDatabaseServiceMaintenanceDOW,
			Time:    testDatabaseServiceMaintenanceTime,
			Updates: []oapi.DbaasServiceUpdate{},
		},
		Metadata:     &oapi.DbaasService_Metadata{AdditionalProperties: map[string]interface{}{"k": "v"}},
		Name:         oapi.DbaasServiceName(testDatabaseServiceName),
		NodeCount:    &testDatabaseServiceNodeCount,
		NodeCpuCount: &testDatabaseServiceNodeCPUCount,
		NodeMemory:   &testDatabaseServiceNodeMemory,
		NodeStates: &[]oapi.DbaasNodeState{{
			Name:            ts.randomString(10),
			ProgressUpdates: &[]oapi.DbaasNodeStateProgressUpdate{},
			Role:            &testDatabaseServiceNodeStateRole,
			State:           testDatabaseServiceNodeStateState,
		}},
		Notifications:         &[]oapi.DbaasServiceNotification{},
		Plan:                  testDatabasePlanName,
		State:                 &testDatabaseServiceState,
		TerminationProtection: &testDatabaseServiceTerminationProtection,
		Type:                  testDatabaseServiceType,
		UpdatedAt:             &testDatabaseServiceUpdatedAt,
		Uri:                   &testDatabaseServiceURI,
		UriParams:             &oapi.DbaasService_UriParams{AdditionalProperties: map[string]interface{}{"k": "v"}},
		UserConfig:            &oapi.DbaasService_UserConfig{AdditionalProperties: map[string]interface{}{"k": "v"}},
		Users: &[]oapi.DbaasServiceUser{{
			Password: &testDatabaseServiceUserPassword,
			Type:     testDatabaseServiceUserType,
			Username: testDatabaseServiceUserUsername,
		}},
	})

	expected := &DatabaseService{
		Backups: []*DatabaseServiceBackup{{
			Name: &testDatabaseServiceBackupName,
			Size: &testDatabaseServiceBackupDataSize,
			Date: &testDatabaseServiceBackupTime,
		}},
		Components: []*DatabaseServiceComponent{{
			Name: &testDatabaseServiceComponent,
			Info: map[string]interface{}{
				"host":  testDatabaseServiceComponentHost,
				"port":  testDatabaseServiceComponentPort,
				"route": testDatabaseServiceComponentRoute,
				"usage": testDatabaseServiceComponentUsage,
			},
		}},
		ConnectionInfo: map[string]interface{}{"k": "v"},
		CreatedAt:      &testDatabaseServiceCreatedAt,
		DiskSize:       &testDatabaseServiceDiskSize,
		Features:       map[string]interface{}{"k": "v"},
		Maintenance: &DatabaseServiceMaintenance{
			DOW:  string(testDatabaseServiceMaintenanceDOW),
			Time: testDatabaseServiceMaintenanceTime,
		},
		Metadata:              map[string]interface{}{"k": "v"},
		Name:                  &testDatabaseServiceName,
		Nodes:                 &testDatabaseServiceNodeCount,
		NodeCPUs:              &testDatabaseServiceNodeCPUCount,
		NodeMemory:            &testDatabaseServiceNodeMemory,
		Plan:                  &testDatabasePlanName,
		State:                 (*string)(&testDatabaseServiceState),
		TerminationProtection: &testDatabaseServiceTerminationProtection,
		Type:                  (*string)(&testDatabaseServiceType),
		UpdatedAt:             &testDatabaseServiceUpdatedAt,
		URI:                   func() *url.URL { u, _ := url.Parse(testDatabaseServiceURI); return u }(),
		UserConfig:            &map[string]interface{}{"k": "v"},
		Users: []*DatabaseServiceUser{{
			Password: &testDatabaseServiceUserPassword,
			Type:     &testDatabaseServiceUserType,
			UserName: &testDatabaseServiceUserUsername,
		}},
	}

	actual, err := ts.client.CreateDatabaseService(context.Background(), testZone, &DatabaseService{
		Maintenance: &DatabaseServiceMaintenance{
			DOW:  string(testDatabaseServiceMaintenanceDOW),
			Time: testDatabaseServiceMaintenanceTime,
		},
		Name:                  &testDatabaseServiceName,
		Plan:                  &testDatabasePlanName,
		TerminationProtection: &testDatabaseServiceTerminationProtection,
		Type:                  (*string)(&testDatabaseServiceType),
		UserConfig:            &map[string]interface{}{"k": "v"},
	})
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_DeleteDatabaseService() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		deleted            = false
	)

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("/dbaas-service/%s", testDatabaseServiceName),
		func(req *http.Request) (*http.Response, error) {
			deleted = true

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testDatabaseServiceName},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testDatabaseServiceName},
	})

	ts.Require().NoError(ts.client.DeleteDatabaseService(
		context.Background(),
		testZone,
		&DatabaseService{Name: &testDatabaseServiceName}),
	)
	ts.Require().True(deleted)
}

func (ts *clientTestSuite) TestClient_GetDatabaseCACertificate() {
	testCACertificate := `-----BEGIN CERTIFICATE-----
` + ts.randomString(1000) +
		`-----END CERTIFICATE-----
`

	ts.mockAPIRequest("GET", "/dbaas-ca-certificate", struct {
		Certificate *string "json:\"certificate,omitempty\""
	}{
		Certificate: &testCACertificate,
	})

	actual, err := ts.client.GetDatabaseCACertificate(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(testCACertificate, actual)
}

func (ts *clientTestSuite) TestClient_GetDatabaseService() {
	ts.mockAPIRequest("GET",
		fmt.Sprintf("/dbaas-service/%s", testDatabaseServiceName),
		oapi.DbaasService{
			Acl: &[]oapi.DbaasServiceAcl{},
			Backups: &[]oapi.DbaasServiceBackup{{
				BackupName: testDatabaseServiceBackupName,
				BackupTime: testDatabaseServiceBackupTime,
				DataSize:   testDatabaseServiceBackupDataSize,
			}},
			Components: &[]oapi.DbaasServiceComponents{{
				Component: testDatabaseServiceComponent,
				Host:      testDatabaseServiceComponentHost,
				Port:      testDatabaseServiceComponentPort,
				Route:     testDatabaseServiceComponentRoute,
				Usage:     testDatabaseServiceComponentUsage,
			}},
			ConnectionInfo:  &oapi.DbaasService_ConnectionInfo{AdditionalProperties: map[string]interface{}{"k": "v"}},
			ConnectionPools: &[]oapi.DbaasServiceConnectionPools{},
			CreatedAt:       &testDatabaseServiceCreatedAt,
			Description:     &testDatabaseServiceDescription,
			DiskSize:        &testDatabaseServiceDiskSize,
			Features:        &oapi.DbaasService_Features{AdditionalProperties: map[string]interface{}{"k": "v"}},
			Integrations:    &[]oapi.DbaasServiceIntegration{},
			Maintenance: &oapi.DbaasServiceMaintenance{
				Dow:     testDatabaseServiceMaintenanceDOW,
				Time:    testDatabaseServiceMaintenanceTime,
				Updates: []oapi.DbaasServiceUpdate{},
			},
			Metadata:     &oapi.DbaasService_Metadata{AdditionalProperties: map[string]interface{}{"k": "v"}},
			Name:         oapi.DbaasServiceName(testDatabaseServiceName),
			NodeCount:    &testDatabaseServiceNodeCount,
			NodeCpuCount: &testDatabaseServiceNodeCPUCount,
			NodeMemory:   &testDatabaseServiceNodeMemory,
			NodeStates: &[]oapi.DbaasNodeState{{
				Name:            ts.randomString(10),
				ProgressUpdates: &[]oapi.DbaasNodeStateProgressUpdate{},
				Role:            &testDatabaseServiceNodeStateRole,
				State:           testDatabaseServiceNodeStateState,
			}},
			Notifications:         &[]oapi.DbaasServiceNotification{},
			Plan:                  testDatabasePlanName,
			State:                 &testDatabaseServiceState,
			TerminationProtection: &testDatabaseServiceTerminationProtection,
			Type:                  testDatabaseServiceType,
			UpdatedAt:             &testDatabaseServiceUpdatedAt,
			Uri:                   &testDatabaseServiceURI,
			UriParams:             &oapi.DbaasService_UriParams{AdditionalProperties: map[string]interface{}{"k": "v"}},
			UserConfig:            &oapi.DbaasService_UserConfig{AdditionalProperties: map[string]interface{}{"k": "v"}},
			Users: &[]oapi.DbaasServiceUser{{
				Password: &testDatabaseServiceUserPassword,
				Type:     testDatabaseServiceUserType,
				Username: testDatabaseServiceUserUsername,
			}},
		})

	expected := &DatabaseService{
		Backups: []*DatabaseServiceBackup{{
			Name: &testDatabaseServiceBackupName,
			Size: &testDatabaseServiceBackupDataSize,
			Date: &testDatabaseServiceBackupTime,
		}},
		Components: []*DatabaseServiceComponent{{
			Name: &testDatabaseServiceComponent,
			Info: map[string]interface{}{
				"host":  testDatabaseServiceComponentHost,
				"port":  testDatabaseServiceComponentPort,
				"route": testDatabaseServiceComponentRoute,
				"usage": testDatabaseServiceComponentUsage,
			},
		}},
		ConnectionInfo: map[string]interface{}{"k": "v"},
		CreatedAt:      &testDatabaseServiceCreatedAt,
		DiskSize:       &testDatabaseServiceDiskSize,
		Features:       map[string]interface{}{"k": "v"},
		Maintenance: &DatabaseServiceMaintenance{
			DOW:  string(testDatabaseServiceMaintenanceDOW),
			Time: testDatabaseServiceMaintenanceTime,
		},
		Metadata:              map[string]interface{}{"k": "v"},
		Name:                  &testDatabaseServiceName,
		Nodes:                 &testDatabaseServiceNodeCount,
		NodeCPUs:              &testDatabaseServiceNodeCPUCount,
		NodeMemory:            &testDatabaseServiceNodeMemory,
		Plan:                  &testDatabasePlanName,
		State:                 (*string)(&testDatabaseServiceState),
		TerminationProtection: &testDatabaseServiceTerminationProtection,
		Type:                  (*string)(&testDatabaseServiceType),
		UpdatedAt:             &testDatabaseServiceUpdatedAt,
		URI:                   func() *url.URL { u, _ := url.Parse(testDatabaseServiceURI); return u }(),
		UserConfig:            &map[string]interface{}{"k": "v"},
		Users: []*DatabaseServiceUser{{
			Password: &testDatabaseServiceUserPassword,
			Type:     &testDatabaseServiceUserType,
			UserName: &testDatabaseServiceUserUsername,
		}},
	}

	actual, err := ts.client.GetDatabaseService(context.Background(), testZone, *expected.Name)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_GetDatabaseServiceType() {
	ts.mockAPIRequest("GET",
		fmt.Sprintf("/dbaas-service-type/%s", testDatabaseServiceTypeName),
		oapi.DbaasServiceType{
			DefaultVersion: &testDatabaseServiceTypeDefaultVersion,
			Description:    &testDatabaseServiceTypeDescription,
			LatestVersion:  &testDatabaseServiceTypeLatestVersion,
			Name:           &testDatabaseServiceTypeName,
			Plans: &[]oapi.DbaasPlan{{
				BackupConfig: &oapi.DbaasBackupConfig{
					Interval:     &testDatabasePlanBackupConfigInterval,
					MaxCount:     &testDatabasePlanBackupConfigMaxCount,
					RecoveryMode: &testDatabasePlanBackupConfigRecoveryMode,
				},
				DiskSpace:    &testDatabasePlanDiskSpace,
				Name:         &testDatabasePlanName,
				NodeCount:    &testDatabasePlanNodeCount,
				NodeCpuCount: &testDatabasePlanNodeCPUCount,
				NodeMemory:   &testDatabasePlanNodeMemory,
			}},
			UserConfigSchema: &oapi.DbaasServiceType_UserConfigSchema{
				AdditionalProperties: map[string]interface{}{"k": "v"},
			},
		})

	expected := &DatabaseServiceType{
		DefaultVersion: &testDatabaseServiceTypeDefaultVersion,
		Description:    &testDatabaseServiceTypeDescription,
		LatestVersion:  &testDatabaseServiceTypeLatestVersion,
		Name:           (*string)(&testDatabaseServiceTypeName),
		Plans: []*DatabasePlan{{
			BackupConfig: &DatabaseBackupConfig{
				Interval:     &testDatabasePlanBackupConfigInterval,
				MaxCount:     &testDatabasePlanBackupConfigMaxCount,
				RecoveryMode: &testDatabasePlanBackupConfigRecoveryMode,
			},
			DiskSpace:  &testDatabasePlanDiskSpace,
			Name:       &testDatabasePlanName,
			Nodes:      &testDatabasePlanNodeCount,
			NodeCPUs:   &testDatabasePlanNodeCPUCount,
			NodeMemory: &testDatabasePlanNodeMemory,
		}},
		UserConfigSchema: map[string]interface{}{"k": "v"},
	}

	actual, err := ts.client.GetDatabaseServiceType(context.Background(), testZone, *expected.Name)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_ListDatabaseServiceTypes() {
	ts.mockAPIRequest("GET", "/dbaas-service-type", struct {
		DatabaseServiceTypes *[]oapi.DbaasServiceType `json:"dbaas-service-types,omitempty"`
	}{
		DatabaseServiceTypes: &[]oapi.DbaasServiceType{{
			DefaultVersion: &testDatabaseServiceTypeDefaultVersion,
			Description:    &testDatabaseServiceTypeDescription,
			LatestVersion:  &testDatabaseServiceTypeLatestVersion,
			Name:           &testDatabaseServiceTypeName,
			Plans: &[]oapi.DbaasPlan{{
				BackupConfig: &oapi.DbaasBackupConfig{
					Interval:     &testDatabasePlanBackupConfigInterval,
					MaxCount:     &testDatabasePlanBackupConfigMaxCount,
					RecoveryMode: &testDatabasePlanBackupConfigRecoveryMode,
				},
				DiskSpace:    &testDatabasePlanDiskSpace,
				Name:         &testDatabasePlanName,
				NodeCount:    &testDatabasePlanNodeCount,
				NodeCpuCount: &testDatabasePlanNodeCPUCount,
				NodeMemory:   &testDatabasePlanNodeMemory,
			}},
			UserConfigSchema: &oapi.DbaasServiceType_UserConfigSchema{
				AdditionalProperties: map[string]interface{}{"k": "v"},
			},
		}},
	})

	expected := []*DatabaseServiceType{{
		DefaultVersion: &testDatabaseServiceTypeDefaultVersion,
		Description:    &testDatabaseServiceTypeDescription,
		LatestVersion:  &testDatabaseServiceTypeLatestVersion,
		Name:           (*string)(&testDatabaseServiceTypeName),
		Plans: []*DatabasePlan{{
			BackupConfig: &DatabaseBackupConfig{
				Interval:     &testDatabasePlanBackupConfigInterval,
				MaxCount:     &testDatabasePlanBackupConfigMaxCount,
				RecoveryMode: &testDatabasePlanBackupConfigRecoveryMode,
			},
			DiskSpace:  &testDatabasePlanDiskSpace,
			Name:       &testDatabasePlanName,
			Nodes:      &testDatabasePlanNodeCount,
			NodeCPUs:   &testDatabasePlanNodeCPUCount,
			NodeMemory: &testDatabasePlanNodeMemory,
		}},
		UserConfigSchema: map[string]interface{}{"k": "v"},
	}}

	actual, err := ts.client.ListDatabaseServiceTypes(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_ListDatabaseServices() {
	ts.mockAPIRequest("GET", "/dbaas-service", struct {
		DatabaseServices *[]oapi.DbaasService `json:"dbaas-services,omitempty"`
	}{
		DatabaseServices: &[]oapi.DbaasService{{
			Acl: &[]oapi.DbaasServiceAcl{},
			Backups: &[]oapi.DbaasServiceBackup{{
				BackupName: testDatabaseServiceBackupName,
				BackupTime: testDatabaseServiceBackupTime,
				DataSize:   testDatabaseServiceBackupDataSize,
			}},
			Components: &[]oapi.DbaasServiceComponents{{
				Component: testDatabaseServiceComponent,
				Host:      testDatabaseServiceComponentHost,
				Port:      testDatabaseServiceComponentPort,
				Route:     testDatabaseServiceComponentRoute,
				Usage:     testDatabaseServiceComponentUsage,
			}},
			ConnectionInfo:  &oapi.DbaasService_ConnectionInfo{AdditionalProperties: map[string]interface{}{"k": "v"}},
			ConnectionPools: &[]oapi.DbaasServiceConnectionPools{},
			CreatedAt:       &testDatabaseServiceCreatedAt,
			Description:     &testDatabaseServiceDescription,
			DiskSize:        &testDatabaseServiceDiskSize,
			Features:        &oapi.DbaasService_Features{AdditionalProperties: map[string]interface{}{"k": "v"}},
			Integrations:    &[]oapi.DbaasServiceIntegration{},
			Maintenance: &oapi.DbaasServiceMaintenance{
				Dow:     testDatabaseServiceMaintenanceDOW,
				Time:    testDatabaseServiceMaintenanceTime,
				Updates: []oapi.DbaasServiceUpdate{},
			},
			Metadata:     &oapi.DbaasService_Metadata{AdditionalProperties: map[string]interface{}{"k": "v"}},
			Name:         oapi.DbaasServiceName(testDatabaseServiceName),
			NodeCount:    &testDatabaseServiceNodeCount,
			NodeCpuCount: &testDatabaseServiceNodeCPUCount,
			NodeMemory:   &testDatabaseServiceNodeMemory,
			NodeStates: &[]oapi.DbaasNodeState{{
				Name:            ts.randomString(10),
				ProgressUpdates: &[]oapi.DbaasNodeStateProgressUpdate{},
				Role:            &testDatabaseServiceNodeStateRole,
				State:           testDatabaseServiceNodeStateState,
			}},
			Notifications:         &[]oapi.DbaasServiceNotification{},
			Plan:                  testDatabasePlanName,
			State:                 &testDatabaseServiceState,
			TerminationProtection: &testDatabaseServiceTerminationProtection,
			Type:                  testDatabaseServiceType,
			UpdatedAt:             &testDatabaseServiceUpdatedAt,
			Uri:                   &testDatabaseServiceURI,
			UriParams:             &oapi.DbaasService_UriParams{AdditionalProperties: map[string]interface{}{"k": "v"}},
			UserConfig:            &oapi.DbaasService_UserConfig{AdditionalProperties: map[string]interface{}{"k": "v"}},
			Users: &[]oapi.DbaasServiceUser{{
				Password: &testDatabaseServiceUserPassword,
				Type:     testDatabaseServiceUserType,
				Username: testDatabaseServiceUserUsername,
			}},
		}},
	})

	expected := []*DatabaseService{{
		Backups: []*DatabaseServiceBackup{{
			Name: &testDatabaseServiceBackupName,
			Size: &testDatabaseServiceBackupDataSize,
			Date: &testDatabaseServiceBackupTime,
		}},
		Components: []*DatabaseServiceComponent{{
			Name: &testDatabaseServiceComponent,
			Info: map[string]interface{}{
				"host":  testDatabaseServiceComponentHost,
				"port":  testDatabaseServiceComponentPort,
				"route": testDatabaseServiceComponentRoute,
				"usage": testDatabaseServiceComponentUsage,
			},
		}},
		ConnectionInfo: map[string]interface{}{"k": "v"},
		CreatedAt:      &testDatabaseServiceCreatedAt,
		DiskSize:       &testDatabaseServiceDiskSize,
		Features:       map[string]interface{}{"k": "v"},
		Maintenance: &DatabaseServiceMaintenance{
			DOW:  string(testDatabaseServiceMaintenanceDOW),
			Time: testDatabaseServiceMaintenanceTime,
		},
		Metadata:              map[string]interface{}{"k": "v"},
		Name:                  &testDatabaseServiceName,
		Nodes:                 &testDatabaseServiceNodeCount,
		NodeCPUs:              &testDatabaseServiceNodeCPUCount,
		NodeMemory:            &testDatabaseServiceNodeMemory,
		Plan:                  &testDatabasePlanName,
		State:                 (*string)(&testDatabaseServiceState),
		TerminationProtection: &testDatabaseServiceTerminationProtection,
		Type:                  (*string)(&testDatabaseServiceType),
		UpdatedAt:             &testDatabaseServiceUpdatedAt,
		URI:                   func() *url.URL { u, _ := url.Parse(testDatabaseServiceURI); return u }(),
		UserConfig:            &map[string]interface{}{"k": "v"},
		Users: []*DatabaseServiceUser{{
			Password: &testDatabaseServiceUserPassword,
			Type:     &testDatabaseServiceUserType,
			UserName: &testDatabaseServiceUserUsername,
		}},
	}}

	actual, err := ts.client.ListDatabaseServices(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_UpdateDatabaseService() {
	var (
		testID             = ts.randomID()
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		updated            = false
	)

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/dbaas-service/%s", testDatabaseServiceName),
		func(req *http.Request) (*http.Response, error) {
			updated = true

			var actual oapi.UpdateDbaasServiceJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := oapi.UpdateDbaasServiceJSONRequestBody{
				Maintenance: &struct {
					Dow  oapi.UpdateDbaasServiceJSONBodyMaintenanceDow `json:"dow"`
					Time string                                        `json:"time"`
				}{
					Dow:  oapi.UpdateDbaasServiceJSONBodyMaintenanceDow(testDatabaseServiceMaintenanceDOW),
					Time: testDatabaseServiceMaintenanceTime,
				},
				Plan:                  &testDatabasePlanName,
				TerminationProtection: &testDatabaseServiceTerminationProtection,
				UserConfig: &oapi.UpdateDbaasServiceJSONBody_UserConfig{
					AdditionalProperties: map[string]interface{}{"k": "v"},
				},
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testDatabaseServiceName},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testID},
	})

	ts.Require().NoError(ts.client.UpdateDatabaseService(context.Background(), testZone, &DatabaseService{
		Maintenance: &DatabaseServiceMaintenance{
			DOW:  string(testDatabaseServiceMaintenanceDOW),
			Time: testDatabaseServiceMaintenanceTime,
		},
		Name:                  &testDatabaseServiceName,
		Plan:                  &testDatabasePlanName,
		TerminationProtection: &testDatabaseServiceTerminationProtection,
		Type:                  (*string)(&testDatabaseServiceType),
		UserConfig:            &map[string]interface{}{"k": "v"},
	}))
	ts.Require().True(updated)
}
