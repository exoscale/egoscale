package v2

import (
	"context"
	"net/http"
	"time"

	"github.com/exoscale/egoscale/v2/oapi"
	"github.com/stretchr/testify/mock"
)

var (
	testDatabasePlanAuthorized                     = true
	testDatabasePlanBackupConfigInterval     int64 = 24
	testDatabasePlanBackupConfigMaxCount     int64 = 2
	testDatabasePlanBackupConfigRecoveryMode       = "pitr"
	testDatabasePlanDiskSpace                int64 = 10737418240
	testDatabasePlanName                           = "hobbyist-1"
	testDatabasePlanNodeCPUCount             int64 = 2
	testDatabasePlanNodeCount                int64 = 1
	testDatabasePlanNodeMemory               int64 = 2147483648
	testDatabaseServiceCreatedAt, _                = time.Parse(iso8601Format, "2020-08-12T11:12:36Z")
	testDatabaseServiceDiskSize              int64 = 10995116277760
	testDatabaseServiceName                        = new(testSuite).randomString(10)
	testDatabaseServiceNodeCPUCount          int64 = 2
	testDatabaseServiceNodeCount             int64 = 1
	testDatabaseServiceNodeMemory            int64 = 2199023255552
	testDatabaseServiceState                       = oapi.EnumServiceStateRunning
	testDatabaseServiceTerminationProtection       = true
	testDatabaseServiceType                        = oapi.DbaasServiceTypeName("pg")
	testDatabaseServiceTypeDefaultVersion          = "13"
	testDatabaseServiceTypeDescription             = new(testSuite).randomString(10)
	testDatabaseServiceTypeAvailableVersions       = []string{"12", "13"}
	testDatabaseServiceTypeName                    = oapi.DbaasServiceTypeName("pg")
	testDatabaseServiceUpdatedAt, _                = time.Parse(iso8601Format, "2020-08-12T11:12:36Z")
	testDatabaseServiceMigrationMethod             = "replication"
	testDatabaseMigrationError                     = "Failed to run query on source db: 'connection refused'"
	testDatabaseMigrationDetailsStatusOapi         = oapi.EnumMigrationStatusFailed
	testDatabaseMigrationDetailsStatus             = DatabaseMigrationStatusFailed
	testDatabaseMigrationStatus                    = "failed"
	testDatabaseMasterLastIoSecondsAgo             = int64(12)
	testDatabaseMasterLinkStatusOapi               = oapi.EnumMasterLinkStatusUp
	testDatabaseMasterLinkStatus                   = MasterLinkStatusUp
)

func (ts *testSuite) TestClient_DeleteDatabaseService() {
	deleted := false

	ts.mock().
		On(
			"DeleteDbaasServiceWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // name
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testDatabaseServiceName, args.Get(1))
			deleted = true
		}).
		Return(
			&oapi.DeleteDbaasServiceResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200:      new(oapi.Operation),
			},
			nil,
		)

	ts.Require().NoError(ts.client.DeleteDatabaseService(
		context.Background(),
		testZone,
		&DatabaseService{Name: &testDatabaseServiceName}),
	)
	ts.Require().True(deleted)
}

func (ts *testSuite) TestClient_GetDatabaseCACertificate() {
	testCACertificate := `-----BEGIN CERTIFICATE-----
` + ts.randomString(1000) +
		`-----END CERTIFICATE-----
`

	ts.mock().
		On(
			"GetDbaasCaCertificateWithResponse",
			mock.Anything,                 // ctx
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(
			&oapi.GetDbaasCaCertificateResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &struct {
					Certificate *string "json:\"certificate,omitempty\""
				}{
					Certificate: &testCACertificate,
				},
			},
			nil,
		)

	actual, err := ts.client.GetDatabaseCACertificate(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(testCACertificate, actual)
}

func (ts *testSuite) TestClient_GetDatabaseServiceType() {
	ts.mock().
		On("GetDbaasServiceTypeWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // name
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(string(testDatabaseServiceTypeName), args.Get(1))
		}).
		Return(&oapi.GetDbaasServiceTypeResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.DbaasServiceType{
				AvailableVersions: &testDatabaseServiceTypeAvailableVersions,
				DefaultVersion:    &testDatabaseServiceTypeDefaultVersion,
				Description:       &testDatabaseServiceTypeDescription,
				Name:              &testDatabaseServiceTypeName,
				Plans: &[]oapi.DbaasPlan{{
					Authorized: &testDatabasePlanAuthorized,
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
			},
		}, nil)

	expected := &DatabaseServiceType{
		AvailableVersions: &testDatabaseServiceTypeAvailableVersions,
		DefaultVersion:    &testDatabaseServiceTypeDefaultVersion,
		Description:       &testDatabaseServiceTypeDescription,
		Name:              (*string)(&testDatabaseServiceTypeName),
		Plans: []*DatabasePlan{{
			Authorized: &testDatabasePlanAuthorized,
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
	}

	actual, err := ts.client.GetDatabaseServiceType(context.Background(), testZone, *expected.Name)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_ListDatabaseServiceTypes() {
	ts.mock().
		On("ListDbaasServiceTypesWithResponse",
			mock.Anything,                 // ctx
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.ListDbaasServiceTypesResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &struct {
				DbaasServiceTypes *[]oapi.DbaasServiceType `json:"dbaas-service-types,omitempty"`
			}{
				DbaasServiceTypes: &[]oapi.DbaasServiceType{{
					AvailableVersions: &testDatabaseServiceTypeAvailableVersions,
					DefaultVersion:    &testDatabaseServiceTypeDefaultVersion,
					Description:       &testDatabaseServiceTypeDescription,
					Name:              &testDatabaseServiceTypeName,
					Plans: &[]oapi.DbaasPlan{{
						Authorized: &testDatabasePlanAuthorized,
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
				}},
			},
		}, nil)

	expected := []*DatabaseServiceType{{
		AvailableVersions: &testDatabaseServiceTypeAvailableVersions,
		DefaultVersion:    &testDatabaseServiceTypeDefaultVersion,
		Description:       &testDatabaseServiceTypeDescription,
		Name:              (*string)(&testDatabaseServiceTypeName),
		Plans: []*DatabasePlan{{
			Authorized: &testDatabasePlanAuthorized,
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
	}}

	actual, err := ts.client.ListDatabaseServiceTypes(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_ListDatabaseServices() {
	ts.mock().
		On("ListDbaasServicesWithResponse",
			mock.Anything,                 // ctx
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.ListDbaasServicesResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &struct {
				DbaasServices *[]oapi.DbaasServiceCommon `json:"dbaas-services,omitempty"`
			}{
				DbaasServices: &[]oapi.DbaasServiceCommon{{
					CreatedAt:             &testDatabaseServiceCreatedAt,
					DiskSize:              &testDatabaseServiceDiskSize,
					Name:                  oapi.DbaasServiceName(testDatabaseServiceName),
					NodeCount:             &testDatabaseServiceNodeCount,
					NodeCpuCount:          &testDatabaseServiceNodeCPUCount,
					NodeMemory:            &testDatabaseServiceNodeMemory,
					Notifications:         &[]oapi.DbaasServiceNotification{},
					Plan:                  testDatabasePlanName,
					State:                 &testDatabaseServiceState,
					TerminationProtection: &testDatabaseServiceTerminationProtection,
					Type:                  testDatabaseServiceType,
					UpdatedAt:             &testDatabaseServiceUpdatedAt,
				}},
			},
		}, nil)

	expected := DatabaseService{
		CreatedAt:             &testDatabaseServiceCreatedAt,
		DiskSize:              &testDatabaseServiceDiskSize,
		Name:                  &testDatabaseServiceName,
		Nodes:                 &testDatabaseServiceNodeCount,
		NodeCPUs:              &testDatabaseServiceNodeCPUCount,
		NodeMemory:            &testDatabaseServiceNodeMemory,
		Notifications:         []*DatabaseServiceNotification{},
		Plan:                  &testDatabasePlanName,
		State:                 (*string)(&testDatabaseServiceState),
		TerminationProtection: &testDatabaseServiceTerminationProtection,
		Type:                  (*string)(&testDatabaseServiceType),
		UpdatedAt:             &testDatabaseServiceUpdatedAt,
		Zone:                  &testZone,
	}

	list, err := ts.client.ListDatabaseServices(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal([]*DatabaseService{&expected}, list)

	found, err := ts.client.FindDatabaseService(context.Background(), testZone, testDatabaseServiceName)
	ts.Require().NoError(err)
	ts.Require().Equal(&expected, found)
}

func (ts *testSuite) TestClient_GetDatabaseMigrationStatus() {
	ts.mock().
		On("GetDbaasMigrationStatusWithResponse",
			mock.Anything, // ctx
			oapi.DbaasServiceName(testDatabaseServiceName),
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.GetDbaasMigrationStatusResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.DbaasMigrationStatus{
				Details: &[]struct {
					Dbname *string                   `json:"dbname,omitempty"`
					Error  *string                   `json:"error,omitempty"`
					Method *string                   `json:"method,omitempty"`
					Status *oapi.EnumMigrationStatus `json:"status,omitempty"`
				}{{
					Dbname: &testDatabaseServiceName,
					Error:  &testDatabaseMigrationError,
					Method: &testDatabaseServiceMigrationMethod,
					Status: &testDatabaseMigrationDetailsStatusOapi,
				}},

				Error:                  &testDatabaseMigrationError,
				MasterLastIoSecondsAgo: &testDatabaseMasterLastIoSecondsAgo,
				MasterLinkStatus:       &testDatabaseMasterLinkStatusOapi,
				Method:                 &testDatabaseServiceMigrationMethod,
				Status:                 &testDatabaseMigrationStatus,
			},
		}, nil)

	expected := &DatabaseMigrationStatus{
		Details: []DatabaseMigrationStatusDetails{
			{
				DBName: &testDatabaseServiceName,
				Error:  &testDatabaseMigrationError,
				Method: &testDatabaseServiceMigrationMethod,
				Status: &testDatabaseMigrationDetailsStatus,
			},
		},
		Error:                  &testDatabaseMigrationError,
		MasterLastIOSecondsAgo: &testDatabaseMasterLastIoSecondsAgo,
		MasterLinkStatus:       &testDatabaseMasterLinkStatus,
		Method:                 &testDatabaseServiceMigrationMethod,
		Status:                 &testDatabaseMigrationStatus,
	}

	actual, err := ts.client.GetDatabaseMigrationStatus(context.Background(), testZone, testDatabaseServiceName)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}
