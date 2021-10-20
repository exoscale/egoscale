package v2

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/exoscale/egoscale/v2/oapi"
	"github.com/jarcoal/httpmock"
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
	testDatabaseServiceName                        = new(clientTestSuite).randomString(10)
	testDatabaseServiceNodeCPUCount          int64 = 2
	testDatabaseServiceNodeCount             int64 = 1
	testDatabaseServiceNodeMemory            int64 = 2199023255552
	testDatabaseServiceState                       = oapi.EnumServiceStateRunning
	testDatabaseServiceTerminationProtection       = true
	testDatabaseServiceType                        = oapi.DbaasServiceTypeName("pg")
	testDatabaseServiceTypeDefaultVersion          = "13"
	testDatabaseServiceTypeDescription             = new(clientTestSuite).randomString(10)
	testDatabaseServiceTypeLatestVersion           = "13.3"
	testDatabaseServiceTypeName                    = oapi.DbaasServiceTypeName("pg")
	testDatabaseServiceUpdatedAt, _                = time.Parse(iso8601Format, "2020-08-12T11:12:36Z")
)

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

func (ts *clientTestSuite) TestClient_GetDatabaseServiceType() {
	ts.mockAPIRequest("GET",
		fmt.Sprintf("/dbaas-service-type/%s", testDatabaseServiceTypeName),
		oapi.DbaasServiceType{
			DefaultVersion: &testDatabaseServiceTypeDefaultVersion,
			Description:    &testDatabaseServiceTypeDescription,
			LatestVersion:  &testDatabaseServiceTypeLatestVersion,
			Name:           &testDatabaseServiceTypeName,
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
			UserConfigSchema: &map[string]interface{}{"k": "v"},
		})

	expected := &DatabaseServiceType{
		DefaultVersion: &testDatabaseServiceTypeDefaultVersion,
		Description:    &testDatabaseServiceTypeDescription,
		LatestVersion:  &testDatabaseServiceTypeLatestVersion,
		Name:           (*string)(&testDatabaseServiceTypeName),
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
			UserConfigSchema: &map[string]interface{}{"k": "v"},
		}},
	})

	expected := []*DatabaseServiceType{{
		DefaultVersion: &testDatabaseServiceTypeDefaultVersion,
		Description:    &testDatabaseServiceTypeDescription,
		LatestVersion:  &testDatabaseServiceTypeLatestVersion,
		Name:           (*string)(&testDatabaseServiceTypeName),
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

func (ts *clientTestSuite) TestClient_ListDatabaseServices() {
	ts.mockAPIRequest("GET", "/dbaas-service", struct {
		DatabaseServices *[]oapi.DbaasServiceCommon `json:"dbaas-services,omitempty"`
	}{
		DatabaseServices: &[]oapi.DbaasServiceCommon{{
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
	})

	expected := []*DatabaseService{{
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
	}}

	actual, err := ts.client.ListDatabaseServices(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}
