package v2

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jarcoal/httpmock"

	papi "github.com/exoscale/egoscale/v2/internal/public-api"
)

var (
	testInstancePoolAntiAffinityGroupID       = new(clientTestSuite).randomID()
	testInstancePoolDescription               = "Test Instance Pool description"
	testInstancePoolDiskSize            int64 = 10
	testInstancePoolElasticIPID               = new(clientTestSuite).randomID()
	testInstancePoolID                        = new(clientTestSuite).randomID()
	testInstancePoolIPv6Enabled               = true
	testInstancePoolInstanceID                = new(clientTestSuite).randomID()
	testInstancePoolInstanceTypeID            = new(clientTestSuite).randomID()
	testInstancePoolManagerID                 = new(clientTestSuite).randomID()
	testInstancePoolName                      = "test-instancepool"
	testInstancePoolPrivateNetworkID          = new(clientTestSuite).randomID()
	testInstancePoolSecurityGroupID           = new(clientTestSuite).randomID()
	testInstancePoolSize                int64 = 3
	testInstancePoolState                     = "running"
	testInstancePoolSSHKey                    = "test-ssh-key"
	testInstancePoolTemplateID                = new(clientTestSuite).randomID()
	testInstancePoolUserData                  = "I2Nsb3VkLWNvbmZpZwphcHRfdXBncmFkZTogdHJ1ZQ=="
)

func (ts *clientTestSuite) TestInstancePool_Scale() {
	var (
		testOperationID          = ts.randomID()
		testOperationState       = "success"
		testScaleSize      int64 = testInstancePoolSize * 2
	)

	instancePool := &InstancePool{
		ID:   testInstancePoolID,
		c:    ts.client,
		zone: testZone,

		InstanceIDs: []string{testInstancePoolID},
	}

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/instance-pool/%s:scale", instancePool.ID),
		func(req *http.Request) (*http.Response, error) {
			var actual papi.ScaleInstancePoolJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.ScaleInstancePoolJSONRequestBody{Size: testScaleSize}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testInstancePoolID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testInstancePoolID},
	})

	ts.Require().NoError(instancePool.Scale(context.Background(), testScaleSize))
}

func (ts *clientTestSuite) TestInstancePool_EvictMembers() {
	var (
		testOperationID     = ts.randomID()
		testOperationState  = "success"
		testEvictedMemberID = ts.randomID()
	)

	instancePool := &InstancePool{
		ID:   testInstancePoolID,
		c:    ts.client,
		zone: testZone,
	}

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/instance-pool/%s:evict", instancePool.ID),
		func(req *http.Request) (*http.Response, error) {
			var actual papi.EvictInstancePoolMembersJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.EvictInstancePoolMembersJSONRequestBody{Instances: &[]string{testEvictedMemberID}}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testInstancePoolID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testInstancePoolID},
	})

	ts.Require().NoError(instancePool.EvictMembers(context.Background(), []string{testEvictedMemberID}))
}

func (ts *clientTestSuite) TestInstancePool_ResetField() {
	var (
		testResetField     = "description"
		testOperationID    = ts.randomID()
		testOperationState = "success"
	)

	httpmock.RegisterResponder("DELETE", "=~^/instance-pool/.*",
		func(req *http.Request) (*http.Response, error) {
			ts.Require().Equal(
				fmt.Sprintf("/instance-pool/%s/%s", testInstancePoolID, testResetField),
				req.URL.String())

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testInstancePoolID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testInstancePoolID},
	})

	instancePool := &InstancePool{
		ID:   testInstancePoolID,
		c:    ts.client,
		zone: testZone,
	}

	ts.Require().NoError(instancePool.ResetField(context.Background(), &instancePool.Description))
}

func (ts *clientTestSuite) TestClient_CreateInstancePool() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = "success"
	)

	httpmock.RegisterResponder("POST", "/instance-pool",
		func(req *http.Request) (*http.Response, error) {
			var actual papi.CreateInstancePoolJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.CreateInstancePoolJSONRequestBody{
				AntiAffinityGroups: &[]papi.AntiAffinityGroup{{Id: &testInstancePoolAntiAffinityGroupID}},
				Description:        &testInstancePoolDescription,
				DiskSize:           testInstancePoolDiskSize,
				ElasticIps:         &[]papi.ElasticIp{{Id: &testInstancePoolElasticIPID}},
				InstanceType:       papi.InstanceType{Id: &testInstancePoolInstanceTypeID},
				Ipv6Enabled:        &testInstancePoolIPv6Enabled,
				Name:               testInstancePoolName,
				PrivateNetworks:    &[]papi.PrivateNetwork{{Id: &testInstancePoolPrivateNetworkID}},
				SecurityGroups:     &[]papi.SecurityGroup{{Id: &testInstancePoolSecurityGroupID}},
				Size:               testInstancePoolSize,
				SshKey:             &papi.SshKey{Name: &testInstancePoolSSHKey},
				Template:           papi.Template{Id: &testInstancePoolTemplateID},
				UserData:           &testInstancePoolUserData,
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testInstancePoolID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testInstancePoolID},
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/instance-pool/%s", testInstancePoolID), papi.InstancePool{
		AntiAffinityGroups: &[]papi.AntiAffinityGroup{{Id: &testInstancePoolAntiAffinityGroupID}},
		Description:        &testInstancePoolDescription,
		DiskSize:           &testInstancePoolDiskSize,
		ElasticIps:         &[]papi.ElasticIp{{Id: &testInstancePoolElasticIPID}},
		Id:                 &testInstancePoolID,
		InstanceType:       &papi.InstanceType{Id: &testInstancePoolInstanceTypeID},
		Instances:          &[]papi.Instance{{Id: &testInstancePoolInstanceID}},
		Ipv6Enabled:        &testInstancePoolIPv6Enabled,
		Manager:            &papi.Manager{Id: &testInstancePoolManagerID},
		Name:               &testInstancePoolName,
		PrivateNetworks:    &[]papi.PrivateNetwork{{Id: &testInstancePoolPrivateNetworkID}},
		SecurityGroups:     &[]papi.SecurityGroup{{Id: &testInstancePoolSecurityGroupID}},
		Size:               &testInstancePoolSize,
		SshKey:             &papi.SshKey{Name: &testInstancePoolSSHKey},
		State:              &testInstancePoolState,
		Template:           &papi.Template{Id: &testInstancePoolTemplateID},
		UserData:           &testInstancePoolUserData,
	})

	expected := &InstancePool{
		AntiAffinityGroupIDs: []string{testInstancePoolAntiAffinityGroupID},
		Description:          testInstancePoolDescription,
		DiskSize:             testInstancePoolDiskSize,
		ElasticIPIDs:         []string{testInstancePoolElasticIPID},
		ID:                   testInstancePoolID,
		IPv6Enabled:          testInstancePoolIPv6Enabled,
		InstanceIDs:          []string{testInstancePoolInstanceID},
		InstanceTypeID:       testInstancePoolInstanceTypeID,
		ManagerID:            testInstancePoolManagerID,
		Name:                 testInstancePoolName,
		PrivateNetworkIDs:    []string{testInstancePoolPrivateNetworkID},
		SSHKey:               testInstancePoolSSHKey,
		SecurityGroupIDs:     []string{testInstancePoolSecurityGroupID},
		Size:                 testInstancePoolSize,
		State:                testInstancePoolState,
		TemplateID:           testInstancePoolTemplateID,
		UserData:             testInstancePoolUserData,

		c:    ts.client,
		zone: testZone,
	}

	actual, err := ts.client.CreateInstancePool(context.Background(), testZone, &InstancePool{
		AntiAffinityGroupIDs: []string{testInstancePoolAntiAffinityGroupID},
		Description:          testInstancePoolDescription,
		DiskSize:             testInstancePoolDiskSize,
		ElasticIPIDs:         []string{testInstancePoolElasticIPID},
		IPv6Enabled:          testInstancePoolIPv6Enabled,
		InstanceTypeID:       testInstancePoolInstanceTypeID,
		Name:                 testInstancePoolName,
		PrivateNetworkIDs:    []string{testInstancePoolPrivateNetworkID},
		SSHKey:               testInstancePoolSSHKey,
		SecurityGroupIDs:     []string{testInstancePoolSecurityGroupID},
		Size:                 testInstancePoolSize,
		TemplateID:           testInstancePoolTemplateID,
		UserData:             testInstancePoolUserData,
	})
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_ListInstancePools() {
	ts.mockAPIRequest("GET", "/instance-pool", struct {
		InstancePools *[]papi.InstancePool `json:"instance-pools,omitempty"`
	}{
		InstancePools: &[]papi.InstancePool{{
			AntiAffinityGroups: &[]papi.AntiAffinityGroup{{Id: &testInstancePoolAntiAffinityGroupID}},
			Description:        &testInstancePoolDescription,
			DiskSize:           &testInstancePoolDiskSize,
			ElasticIps:         &[]papi.ElasticIp{{Id: &testInstancePoolElasticIPID}},
			Id:                 &testInstancePoolID,
			InstanceType:       &papi.InstanceType{Id: &testInstancePoolInstanceTypeID},
			Instances:          &[]papi.Instance{{Id: &testInstancePoolInstanceID}},
			Ipv6Enabled:        &testInstancePoolIPv6Enabled,
			Manager:            &papi.Manager{Id: &testInstancePoolManagerID},
			Name:               &testInstancePoolName,
			PrivateNetworks:    &[]papi.PrivateNetwork{{Id: &testInstancePoolPrivateNetworkID}},
			SecurityGroups:     &[]papi.SecurityGroup{{Id: &testInstancePoolSecurityGroupID}},
			Size:               &testInstancePoolSize,
			SshKey:             &papi.SshKey{Name: &testInstancePoolSSHKey},
			State:              &testInstancePoolState,
			Template:           &papi.Template{Id: &testInstancePoolTemplateID},
			UserData:           &testInstancePoolUserData,
		}},
	})

	expected := []*InstancePool{{
		AntiAffinityGroupIDs: []string{testInstancePoolAntiAffinityGroupID},
		Description:          testInstancePoolDescription,
		DiskSize:             testInstancePoolDiskSize,
		ElasticIPIDs:         []string{testInstancePoolElasticIPID},
		ID:                   testInstancePoolID,
		IPv6Enabled:          testInstancePoolIPv6Enabled,
		InstanceIDs:          []string{testInstancePoolInstanceID},
		InstanceTypeID:       testInstancePoolInstanceTypeID,
		ManagerID:            testInstancePoolManagerID,
		Name:                 testInstancePoolName,
		PrivateNetworkIDs:    []string{testInstancePoolPrivateNetworkID},
		SSHKey:               testInstancePoolSSHKey,
		SecurityGroupIDs:     []string{testInstancePoolSecurityGroupID},
		Size:                 testInstancePoolSize,
		State:                testInstancePoolState,
		TemplateID:           testInstancePoolTemplateID,
		UserData:             testInstancePoolUserData,

		c:    ts.client,
		zone: testZone,
	}}

	actual, err := ts.client.ListInstancePools(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_GetInstancePool() {
	ts.mockAPIRequest("GET", fmt.Sprintf("/instance-pool/%s", testInstancePoolID), papi.InstancePool{
		AntiAffinityGroups: &[]papi.AntiAffinityGroup{{Id: &testInstancePoolAntiAffinityGroupID}},
		Description:        &testInstancePoolDescription,
		DiskSize:           &testInstancePoolDiskSize,
		ElasticIps:         &[]papi.ElasticIp{{Id: &testInstancePoolElasticIPID}},
		Id:                 &testInstancePoolID,
		InstanceType:       &papi.InstanceType{Id: &testInstancePoolInstanceTypeID},
		Instances:          &[]papi.Instance{{Id: &testInstancePoolInstanceID}},
		Ipv6Enabled:        &testInstancePoolIPv6Enabled,
		Manager:            &papi.Manager{Id: &testInstancePoolManagerID},
		Name:               &testInstancePoolName,
		PrivateNetworks:    &[]papi.PrivateNetwork{{Id: &testInstancePoolPrivateNetworkID}},
		SecurityGroups:     &[]papi.SecurityGroup{{Id: &testInstancePoolSecurityGroupID}},
		Size:               &testInstancePoolSize,
		SshKey:             &papi.SshKey{Name: &testInstancePoolSSHKey},
		State:              &testInstancePoolState,
		Template:           &papi.Template{Id: &testInstancePoolTemplateID},
		UserData:           &testInstancePoolUserData,
	})

	expected := &InstancePool{
		AntiAffinityGroupIDs: []string{testInstancePoolAntiAffinityGroupID},
		Description:          testInstancePoolDescription,
		DiskSize:             testInstancePoolDiskSize,
		ElasticIPIDs:         []string{testInstancePoolElasticIPID},
		ID:                   testInstancePoolID,
		IPv6Enabled:          testInstancePoolIPv6Enabled,
		InstanceIDs:          []string{testInstancePoolInstanceID},
		InstanceTypeID:       testInstancePoolInstanceTypeID,
		ManagerID:            testInstancePoolManagerID,
		Name:                 testInstancePoolName,
		PrivateNetworkIDs:    []string{testInstancePoolPrivateNetworkID},
		SSHKey:               testInstancePoolSSHKey,
		SecurityGroupIDs:     []string{testInstancePoolSecurityGroupID},
		Size:                 testInstancePoolSize,
		State:                testInstancePoolState,
		TemplateID:           testInstancePoolTemplateID,
		UserData:             testInstancePoolUserData,

		c:    ts.client,
		zone: testZone,
	}

	actual, err := ts.client.GetInstancePool(context.Background(), testZone, expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_UpdateInstancePool() {
	var (
		testInstancePoolAntiAffinityGroupIDUpdated       = new(clientTestSuite).randomID()
		testInstancePoolDescriptionUpdated               = testInstancePoolDescription + "-updated"
		testInstancePoolDiskSizeUpdated            int64 = testInstancePoolDiskSize * 2
		testInstancePoolElasticIPIDUpdated               = new(clientTestSuite).randomID()
		testInstancePoolIPv6EnabledUpdated               = true
		testInstancePoolInstanceTypeIDUpdated            = new(clientTestSuite).randomID()
		testInstancePoolNameUpdated                      = testInstancePoolName + "-updated"
		testInstancePoolPrivateNetworkIDUpdated          = new(clientTestSuite).randomID()
		testInstancePoolSecurityGroupIDUpdated           = new(clientTestSuite).randomID()
		testInstancePoolSSHKeyUpdated                    = testInstancePoolSSHKey + "-updated"
		testInstancePoolTemplateIDUpdated                = new(clientTestSuite).randomID()
		testInstancePoolUserDataUpdated                  = testInstancePoolUserData + "-updated"
		testOperationID                                  = ts.randomID()
		testOperationState                               = "success"
		updated                                          = false
	)

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/instance-pool/%s", testInstancePoolID),
		func(req *http.Request) (*http.Response, error) {
			updated = true

			var actual papi.UpdateInstancePoolJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.UpdateInstancePoolJSONRequestBody{
				AntiAffinityGroups: &[]papi.AntiAffinityGroup{{Id: &testInstancePoolAntiAffinityGroupIDUpdated}},
				Description:        &testInstancePoolDescriptionUpdated,
				DiskSize:           &testInstancePoolDiskSizeUpdated,
				ElasticIps:         &[]papi.ElasticIp{{Id: &testInstancePoolElasticIPIDUpdated}},
				InstanceType:       &papi.InstanceType{Id: &testInstancePoolInstanceTypeIDUpdated},
				Ipv6Enabled:        &testInstancePoolIPv6EnabledUpdated,
				Name:               &testInstancePoolNameUpdated,
				PrivateNetworks:    &[]papi.PrivateNetwork{{Id: &testInstancePoolPrivateNetworkIDUpdated}},
				SecurityGroups:     &[]papi.SecurityGroup{{Id: &testInstancePoolSecurityGroupIDUpdated}},
				SshKey:             &papi.SshKey{Name: &testInstancePoolSSHKeyUpdated},
				Template:           &papi.Template{Id: &testInstancePoolTemplateIDUpdated},
				UserData:           &testInstancePoolUserDataUpdated,
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testInstancePoolID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testInstancePoolID},
	})

	ts.Require().NoError(ts.client.UpdateInstancePool(context.Background(), testZone, &InstancePool{
		AntiAffinityGroupIDs: []string{testInstancePoolAntiAffinityGroupIDUpdated},
		Description:          testInstancePoolDescriptionUpdated,
		DiskSize:             testInstancePoolDiskSizeUpdated,
		ElasticIPIDs:         []string{testInstancePoolElasticIPIDUpdated},
		ID:                   testInstancePoolID,
		IPv6Enabled:          testInstancePoolIPv6EnabledUpdated,
		InstanceIDs:          []string{testInstancePoolInstanceTypeIDUpdated},
		InstanceTypeID:       testInstancePoolInstanceTypeIDUpdated,
		Name:                 testInstancePoolNameUpdated,
		PrivateNetworkIDs:    []string{testInstancePoolPrivateNetworkIDUpdated},
		SSHKey:               testInstancePoolSSHKeyUpdated,
		SecurityGroupIDs:     []string{testInstancePoolSecurityGroupIDUpdated},
		TemplateID:           testInstancePoolTemplateIDUpdated,
		UserData:             testInstancePoolUserDataUpdated,
	}))
	ts.Require().True(updated)
}

func (ts *clientTestSuite) TestClient_DeleteInstancePool() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = "success"
		deleted            = false
	)

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("/instance-pool/%s", testInstancePoolID),
		func(req *http.Request) (*http.Response, error) {
			deleted = true

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testInstancePoolID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testInstancePoolID},
	})

	ts.Require().NoError(ts.client.DeleteInstancePool(context.Background(), testZone, testInstancePoolID))
	ts.Require().True(deleted)
}
