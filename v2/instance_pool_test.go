package v2

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jarcoal/httpmock"

	"github.com/exoscale/egoscale/v2/oapi"
)

var (
	testInstancePoolAntiAffinityGroupID       = new(clientTestSuite).randomID()
	testInstancePoolDeployTargetID            = new(clientTestSuite).randomID()
	testInstancePoolDescription               = new(clientTestSuite).randomString(10)
	testInstancePoolDiskSize            int64 = 10
	testInstancePoolElasticIPID               = new(clientTestSuite).randomID()
	testInstancePoolID                        = new(clientTestSuite).randomID()
	testInstancePoolIPv6Enabled               = true
	testInstancePoolInstanceID                = new(clientTestSuite).randomID()
	testInstancePoolInstancePrefix            = new(clientTestSuite).randomString(10)
	testInstancePoolInstanceTypeID            = new(clientTestSuite).randomID()
	testInstancePoolLabels                    = map[string]string{"k1": "v1", "k2": "v2"}
	testInstancePoolManagerID                 = new(clientTestSuite).randomID()
	testInstancePoolManagerType               = oapi.ManagerTypeSksNodepool
	testInstancePoolName                      = new(clientTestSuite).randomString(10)
	testInstancePoolPrivateNetworkID          = new(clientTestSuite).randomID()
	testInstancePoolSSHKey                    = new(clientTestSuite).randomString(10)
	testInstancePoolSecurityGroupID           = new(clientTestSuite).randomID()
	testInstancePoolSize                int64 = 3
	testInstancePoolState                     = oapi.InstancePoolStateRunning
	testInstancePoolTemplateID                = new(clientTestSuite).randomID()
	testInstancePoolUserData                  = "I2Nsb3VkLWNvbmZpZwphcHRfdXBncmFkZTogdHJ1ZQ=="
)

func (ts *clientTestSuite) TestClient_CreateInstancePool() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
	)

	httpmock.RegisterResponder("POST", "/instance-pool",
		func(req *http.Request) (*http.Response, error) {
			var actual oapi.CreateInstancePoolJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := oapi.CreateInstancePoolJSONRequestBody{
				AntiAffinityGroups: &[]oapi.AntiAffinityGroup{{Id: &testInstancePoolAntiAffinityGroupID}},
				DeployTarget:       &oapi.DeployTarget{Id: &testInstancePoolDeployTargetID},
				Description:        &testInstancePoolDescription,
				DiskSize:           testInstancePoolDiskSize,
				ElasticIps:         &[]oapi.ElasticIp{{Id: &testInstancePoolElasticIPID}},
				InstancePrefix:     &testInstancePoolInstancePrefix,
				InstanceType:       oapi.InstanceType{Id: &testInstancePoolInstanceTypeID},
				Ipv6Enabled:        &testInstancePoolIPv6Enabled,
				Labels:             &oapi.Labels{AdditionalProperties: testInstancePoolLabels},
				Name:               testInstancePoolName,
				PrivateNetworks:    &[]oapi.PrivateNetwork{{Id: &testInstancePoolPrivateNetworkID}},
				SecurityGroups:     &[]oapi.SecurityGroup{{Id: &testInstancePoolSecurityGroupID}},
				Size:               testInstancePoolSize,
				SshKey:             &oapi.SshKey{Name: &testInstancePoolSSHKey},
				Template:           oapi.Template{Id: &testInstancePoolTemplateID},
				UserData:           &testInstancePoolUserData,
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testInstancePoolID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testInstancePoolID},
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/instance-pool/%s", testInstancePoolID), oapi.InstancePool{
		AntiAffinityGroups: &[]oapi.AntiAffinityGroup{{Id: &testInstancePoolAntiAffinityGroupID}},
		DeployTarget:       &oapi.DeployTarget{Id: &testInstancePoolDeployTargetID},
		Description:        &testInstancePoolDescription,
		DiskSize:           &testInstancePoolDiskSize,
		ElasticIps:         &[]oapi.ElasticIp{{Id: &testInstancePoolElasticIPID}},
		Id:                 &testInstancePoolID,
		InstancePrefix:     &testInstancePoolInstancePrefix,
		InstanceType:       &oapi.InstanceType{Id: &testInstancePoolInstanceTypeID},
		Instances:          &[]oapi.Instance{{Id: &testInstancePoolInstanceID}},
		Ipv6Enabled:        &testInstancePoolIPv6Enabled,
		Labels:             &oapi.Labels{AdditionalProperties: testInstancePoolLabels},
		Manager:            &oapi.Manager{Id: &testInstancePoolManagerID, Type: &testInstancePoolManagerType},
		Name:               &testInstancePoolName,
		PrivateNetworks:    &[]oapi.PrivateNetwork{{Id: &testInstancePoolPrivateNetworkID}},
		SecurityGroups:     &[]oapi.SecurityGroup{{Id: &testInstancePoolSecurityGroupID}},
		Size:               &testInstancePoolSize,
		SshKey:             &oapi.SshKey{Name: &testInstancePoolSSHKey},
		State:              &testInstancePoolState,
		Template:           &oapi.Template{Id: &testInstancePoolTemplateID},
		UserData:           &testInstancePoolUserData,
	})

	expected := &InstancePool{
		AntiAffinityGroupIDs: &[]string{testInstancePoolAntiAffinityGroupID},
		DeployTargetID:       &testInstancePoolDeployTargetID,
		Description:          &testInstancePoolDescription,
		DiskSize:             &testInstancePoolDiskSize,
		ElasticIPIDs:         &[]string{testInstancePoolElasticIPID},
		ID:                   &testInstancePoolID,
		IPv6Enabled:          &testInstancePoolIPv6Enabled,
		InstanceIDs:          &[]string{testInstancePoolInstanceID},
		InstancePrefix:       &testInstancePoolInstancePrefix,
		InstanceTypeID:       &testInstancePoolInstanceTypeID,
		Labels:               &testInstanceLabels,
		Manager:              &InstancePoolManager{ID: testInstancePoolManagerID, Type: string(testInstancePoolManagerType)},
		Name:                 &testInstancePoolName,
		PrivateNetworkIDs:    &[]string{testInstancePoolPrivateNetworkID},
		SSHKey:               &testInstancePoolSSHKey,
		SecurityGroupIDs:     &[]string{testInstancePoolSecurityGroupID},
		Size:                 &testInstancePoolSize,
		State:                (*string)(&testInstancePoolState),
		TemplateID:           &testInstancePoolTemplateID,
		UserData:             &testInstancePoolUserData,
		Zone:                 &testZone,
	}

	actual, err := ts.client.CreateInstancePool(context.Background(), testZone, &InstancePool{
		AntiAffinityGroupIDs: &[]string{testInstancePoolAntiAffinityGroupID},
		DeployTargetID:       &testInstancePoolDeployTargetID,
		Description:          &testInstancePoolDescription,
		DiskSize:             &testInstancePoolDiskSize,
		ElasticIPIDs:         &[]string{testInstancePoolElasticIPID},
		IPv6Enabled:          &testInstancePoolIPv6Enabled,
		InstancePrefix:       &testInstancePoolInstancePrefix,
		InstanceTypeID:       &testInstancePoolInstanceTypeID,
		Labels:               &testInstanceLabels,
		Name:                 &testInstancePoolName,
		PrivateNetworkIDs:    &[]string{testInstancePoolPrivateNetworkID},
		SSHKey:               &testInstancePoolSSHKey,
		SecurityGroupIDs:     &[]string{testInstancePoolSecurityGroupID},
		Size:                 &testInstancePoolSize,
		TemplateID:           &testInstancePoolTemplateID,
		UserData:             &testInstancePoolUserData,
	})
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_DeleteInstancePool() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		deleted            = false
	)

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("/instance-pool/%s", testInstancePoolID),
		func(req *http.Request) (*http.Response, error) {
			deleted = true

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testInstancePoolID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testInstancePoolID},
	})

	ts.Require().NoError(ts.client.DeleteInstancePool(
		context.Background(),
		testZone,
		&InstancePool{ID: &testInstancePoolID},
	))
	ts.Require().True(deleted)
}

func (ts *clientTestSuite) TestClient_EvictInstancePooltMembers() {
	var (
		testOperationID     = ts.randomID()
		testOperationState  = oapi.OperationStateSuccess
		testEvictedMemberID = ts.randomID()
		evicted             = false
	)

	instancePool := &InstancePool{
		ID: &testInstancePoolID,
	}

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/instance-pool/%s:evict", *instancePool.ID),
		func(req *http.Request) (*http.Response, error) {
			evicted = true

			var actual oapi.EvictInstancePoolMembersJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := oapi.EvictInstancePoolMembersJSONRequestBody{Instances: &[]string{testEvictedMemberID}}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testInstancePoolID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testInstancePoolID},
	})

	ts.Require().NoError(ts.client.EvictInstancePoolMembers(
		context.Background(),
		testZone,
		instancePool,
		[]string{testEvictedMemberID},
	))
	ts.Require().True(evicted)
}

func (ts *clientTestSuite) TestClient_FindInstancePool() {
	ts.mockAPIRequest("GET", "/instance-pool", struct {
		InstancePools *[]oapi.InstancePool `json:"instance-pools,omitempty"`
	}{
		InstancePools: &[]oapi.InstancePool{{
			DiskSize:     &testInstancePoolDiskSize,
			Id:           &testInstancePoolID,
			InstanceType: &oapi.InstanceType{Id: &testInstancePoolInstanceTypeID},
			Manager:      &oapi.Manager{Id: &testInstancePoolManagerID, Type: &testInstancePoolManagerType},
			Name:         &testInstancePoolName,
			Size:         &testInstancePoolSize,
			State:        &testInstancePoolState,
			Template:     &oapi.Template{Id: &testInstancePoolTemplateID},
		}},
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/instance-pool/%s", testInstancePoolID), oapi.InstancePool{
		DiskSize:     &testInstancePoolDiskSize,
		Id:           &testInstancePoolID,
		InstanceType: &oapi.InstanceType{Id: &testInstancePoolInstanceTypeID},
		Manager:      &oapi.Manager{Id: &testInstancePoolManagerID, Type: &testInstancePoolManagerType},
		Name:         &testInstancePoolName,
		Size:         &testInstancePoolSize,
		State:        &testInstancePoolState,
		Template:     &oapi.Template{Id: &testInstancePoolTemplateID},
	})

	expected := &InstancePool{
		DiskSize:       &testInstancePoolDiskSize,
		ID:             &testInstancePoolID,
		InstanceTypeID: &testInstancePoolInstanceTypeID,
		Manager: &InstancePoolManager{
			ID:   testInstancePoolManagerID,
			Type: string(testInstancePoolManagerType),
		},
		Name:       &testInstancePoolName,
		Size:       &testInstancePoolSize,
		State:      (*string)(&testInstancePoolState),
		TemplateID: &testInstancePoolTemplateID,
		Zone:       &testZone,
	}

	actual, err := ts.client.FindInstancePool(context.Background(), testZone, *expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)

	actual, err = ts.client.FindInstancePool(context.Background(), testZone, *expected.Name)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_GetInstancePool() {
	ts.mockAPIRequest("GET", fmt.Sprintf("/instance-pool/%s", testInstancePoolID), oapi.InstancePool{
		AntiAffinityGroups: &[]oapi.AntiAffinityGroup{{Id: &testInstancePoolAntiAffinityGroupID}},
		Description:        &testInstancePoolDescription,
		DiskSize:           &testInstancePoolDiskSize,
		ElasticIps:         &[]oapi.ElasticIp{{Id: &testInstancePoolElasticIPID}},
		Id:                 &testInstancePoolID,
		InstanceType:       &oapi.InstanceType{Id: &testInstancePoolInstanceTypeID},
		Instances:          &[]oapi.Instance{{Id: &testInstancePoolInstanceID}},
		Ipv6Enabled:        &testInstancePoolIPv6Enabled,
		Labels:             &oapi.Labels{AdditionalProperties: testInstancePoolLabels},
		Manager:            &oapi.Manager{Id: &testInstancePoolManagerID, Type: &testInstancePoolManagerType},
		Name:               &testInstancePoolName,
		PrivateNetworks:    &[]oapi.PrivateNetwork{{Id: &testInstancePoolPrivateNetworkID}},
		SecurityGroups:     &[]oapi.SecurityGroup{{Id: &testInstancePoolSecurityGroupID}},
		Size:               &testInstancePoolSize,
		SshKey:             &oapi.SshKey{Name: &testInstancePoolSSHKey},
		State:              &testInstancePoolState,
		Template:           &oapi.Template{Id: &testInstancePoolTemplateID},
		UserData:           &testInstancePoolUserData,
	})

	expected := &InstancePool{
		AntiAffinityGroupIDs: &[]string{testInstancePoolAntiAffinityGroupID},
		Description:          &testInstancePoolDescription,
		DiskSize:             &testInstancePoolDiskSize,
		ElasticIPIDs:         &[]string{testInstancePoolElasticIPID},
		ID:                   &testInstancePoolID,
		IPv6Enabled:          &testInstancePoolIPv6Enabled,
		InstanceIDs:          &[]string{testInstancePoolInstanceID},
		InstanceTypeID:       &testInstancePoolInstanceTypeID,
		Labels:               &testInstancePoolLabels,
		Manager: &InstancePoolManager{
			ID:   testInstancePoolManagerID,
			Type: string(testInstancePoolManagerType),
		},
		Name:              &testInstancePoolName,
		PrivateNetworkIDs: &[]string{testInstancePoolPrivateNetworkID},
		SSHKey:            &testInstancePoolSSHKey,
		SecurityGroupIDs:  &[]string{testInstancePoolSecurityGroupID},
		Size:              &testInstancePoolSize,
		State:             (*string)(&testInstancePoolState),
		TemplateID:        &testInstancePoolTemplateID,
		UserData:          &testInstancePoolUserData,
		Zone:              &testZone,
	}

	actual, err := ts.client.GetInstancePool(context.Background(), testZone, *expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_ListInstancePools() {
	ts.mockAPIRequest("GET", "/instance-pool", struct {
		InstancePools *[]oapi.InstancePool `json:"instance-pools,omitempty"`
	}{
		InstancePools: &[]oapi.InstancePool{{
			AntiAffinityGroups: &[]oapi.AntiAffinityGroup{{Id: &testInstancePoolAntiAffinityGroupID}},
			Description:        &testInstancePoolDescription,
			DiskSize:           &testInstancePoolDiskSize,
			ElasticIps:         &[]oapi.ElasticIp{{Id: &testInstancePoolElasticIPID}},
			Id:                 &testInstancePoolID,
			InstanceType:       &oapi.InstanceType{Id: &testInstancePoolInstanceTypeID},
			Instances:          &[]oapi.Instance{{Id: &testInstancePoolInstanceID}},
			Ipv6Enabled:        &testInstancePoolIPv6Enabled,
			Labels:             &oapi.Labels{AdditionalProperties: testInstancePoolLabels},
			Manager:            &oapi.Manager{Id: &testInstancePoolManagerID, Type: &testInstancePoolManagerType},
			Name:               &testInstancePoolName,
			PrivateNetworks:    &[]oapi.PrivateNetwork{{Id: &testInstancePoolPrivateNetworkID}},
			SecurityGroups:     &[]oapi.SecurityGroup{{Id: &testInstancePoolSecurityGroupID}},
			Size:               &testInstancePoolSize,
			SshKey:             &oapi.SshKey{Name: &testInstancePoolSSHKey},
			State:              &testInstancePoolState,
			Template:           &oapi.Template{Id: &testInstancePoolTemplateID},
			UserData:           &testInstancePoolUserData,
		}},
	})

	expected := []*InstancePool{{
		AntiAffinityGroupIDs: &[]string{testInstancePoolAntiAffinityGroupID},
		Description:          &testInstancePoolDescription,
		DiskSize:             &testInstancePoolDiskSize,
		ElasticIPIDs:         &[]string{testInstancePoolElasticIPID},
		ID:                   &testInstancePoolID,
		IPv6Enabled:          &testInstancePoolIPv6Enabled,
		InstanceIDs:          &[]string{testInstancePoolInstanceID},
		InstanceTypeID:       &testInstancePoolInstanceTypeID,
		Labels:               &testInstancePoolLabels,
		Manager: &InstancePoolManager{
			ID:   testInstancePoolManagerID,
			Type: string(testInstancePoolManagerType),
		},
		Name:              &testInstancePoolName,
		PrivateNetworkIDs: &[]string{testInstancePoolPrivateNetworkID},
		SSHKey:            &testInstancePoolSSHKey,
		SecurityGroupIDs:  &[]string{testInstancePoolSecurityGroupID},
		Size:              &testInstancePoolSize,
		State:             (*string)(&testInstancePoolState),
		TemplateID:        &testInstancePoolTemplateID,
		UserData:          &testInstancePoolUserData,
		Zone:              &testZone,
	}}

	actual, err := ts.client.ListInstancePools(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_ScaleInstancePool() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		testScaleSize      = testInstancePoolSize * 2
		scaled             = false
	)

	instancePool := &InstancePool{
		ID:          &testInstancePoolID,
		InstanceIDs: &[]string{testInstancePoolID},
	}

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/instance-pool/%s:scale", *instancePool.ID),
		func(req *http.Request) (*http.Response, error) {
			scaled = true

			var actual oapi.ScaleInstancePoolJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := oapi.ScaleInstancePoolJSONRequestBody{Size: testScaleSize}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testInstancePoolID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testInstancePoolID},
	})

	ts.Require().NoError(ts.client.ScaleInstancePool(context.Background(), testZone, instancePool, testScaleSize))
	ts.Require().True(scaled)
}

func (ts *clientTestSuite) TestClient_UpdateInstancePool() {
	var (
		testInstancePoolAntiAffinityGroupIDUpdated = new(clientTestSuite).randomID()
		testInstancePoolDeployTargetIDUpdated      = new(clientTestSuite).randomID()
		testInstancePoolDescriptionUpdated         = testInstancePoolDescription + "-updated"
		testInstancePoolDiskSizeUpdated            = testInstancePoolDiskSize * 2
		testInstancePoolElasticIPIDUpdated         = new(clientTestSuite).randomID()
		testInstancePoolIPv6EnabledUpdated         = true
		testInstancePoolInstancePrefixUpdated      = testInstancePoolInstancePrefix + "-updated"
		testInstancePoolInstanceTypeIDUpdated      = new(clientTestSuite).randomID()
		testInstancePoolLabelsUpdated              = map[string]string{"k3": "v3"}
		testInstancePoolNameUpdated                = testInstancePoolName + "-updated"
		testInstancePoolPrivateNetworkIDUpdated    = new(clientTestSuite).randomID()
		testInstancePoolSecurityGroupIDUpdated     = new(clientTestSuite).randomID()
		testInstancePoolSSHKeyUpdated              = testInstancePoolSSHKey + "-updated"
		testInstancePoolTemplateIDUpdated          = new(clientTestSuite).randomID()
		testInstancePoolUserDataUpdated            = testInstancePoolUserData + "-updated"
		testOperationID                            = ts.randomID()
		testOperationState                         = oapi.OperationStateSuccess
		updated                                    = false
	)

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/instance-pool/%s", testInstancePoolID),
		func(req *http.Request) (*http.Response, error) {
			updated = true

			var actual oapi.UpdateInstancePoolJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := oapi.UpdateInstancePoolJSONRequestBody{
				AntiAffinityGroups: &[]oapi.AntiAffinityGroup{{Id: &testInstancePoolAntiAffinityGroupIDUpdated}},
				DeployTarget:       &oapi.DeployTarget{Id: &testInstancePoolDeployTargetIDUpdated},
				Description:        &testInstancePoolDescriptionUpdated,
				DiskSize:           &testInstancePoolDiskSizeUpdated,
				ElasticIps:         &[]oapi.ElasticIp{{Id: &testInstancePoolElasticIPIDUpdated}},
				InstancePrefix:     &testInstancePoolInstancePrefixUpdated,
				InstanceType:       &oapi.InstanceType{Id: &testInstancePoolInstanceTypeIDUpdated},
				Ipv6Enabled:        &testInstancePoolIPv6EnabledUpdated,
				Labels:             &oapi.Labels{AdditionalProperties: testInstancePoolLabelsUpdated},
				Name:               &testInstancePoolNameUpdated,
				PrivateNetworks:    &[]oapi.PrivateNetwork{{Id: &testInstancePoolPrivateNetworkIDUpdated}},
				SecurityGroups:     &[]oapi.SecurityGroup{{Id: &testInstancePoolSecurityGroupIDUpdated}},
				SshKey:             &oapi.SshKey{Name: &testInstancePoolSSHKeyUpdated},
				Template:           &oapi.Template{Id: &testInstancePoolTemplateIDUpdated},
				UserData:           &testInstancePoolUserDataUpdated,
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testInstancePoolID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testInstancePoolID},
	})

	ts.Require().NoError(ts.client.UpdateInstancePool(context.Background(), testZone, &InstancePool{
		AntiAffinityGroupIDs: &[]string{testInstancePoolAntiAffinityGroupIDUpdated},
		DeployTargetID:       &testInstancePoolDeployTargetIDUpdated,
		Description:          &testInstancePoolDescriptionUpdated,
		DiskSize:             &testInstancePoolDiskSizeUpdated,
		ElasticIPIDs:         &[]string{testInstancePoolElasticIPIDUpdated},
		ID:                   &testInstancePoolID,
		IPv6Enabled:          &testInstancePoolIPv6EnabledUpdated,
		InstanceIDs:          &[]string{testInstancePoolInstanceTypeIDUpdated},
		InstancePrefix:       &testInstancePoolInstancePrefixUpdated,
		InstanceTypeID:       &testInstancePoolInstanceTypeIDUpdated,
		Labels:               &testInstancePoolLabelsUpdated,
		Name:                 &testInstancePoolNameUpdated,
		PrivateNetworkIDs:    &[]string{testInstancePoolPrivateNetworkIDUpdated},
		SSHKey:               &testInstancePoolSSHKeyUpdated,
		SecurityGroupIDs:     &[]string{testInstancePoolSecurityGroupIDUpdated},
		TemplateID:           &testInstancePoolTemplateIDUpdated,
		UserData:             &testInstancePoolUserDataUpdated,
	}))
	ts.Require().True(updated)
}
