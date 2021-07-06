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
	testInstancePoolName                      = new(clientTestSuite).randomString(10)
	testInstancePoolPrivateNetworkID          = new(clientTestSuite).randomID()
	testInstancePoolSSHKey                    = new(clientTestSuite).randomString(10)
	testInstancePoolSecurityGroupID           = new(clientTestSuite).randomID()
	testInstancePoolSize                int64 = 3
	testInstancePoolState                     = papi.InstancePoolStateRunning
	testInstancePoolTemplateID                = new(clientTestSuite).randomID()
	testInstancePoolUserData                  = "I2Nsb3VkLWNvbmZpZwphcHRfdXBncmFkZTogdHJ1ZQ=="
)

func (ts *clientTestSuite) TestInstancePool_AntiAffinityGroups() {
	ts.mockAPIRequest(
		"GET",
		fmt.Sprintf("/anti-affinity-group/%s", testAntiAffinityGroupID),
		papi.AntiAffinityGroup{
			Id:   &testAntiAffinityGroupID,
			Name: &testAntiAffinityGroupName,
		},
	)

	expected := []*AntiAffinityGroup{{
		ID:   &testAntiAffinityGroupID,
		Name: &testAntiAffinityGroupName,
	}}

	instancePool := &InstancePool{
		AntiAffinityGroupIDs: &[]string{testAntiAffinityGroupID},

		c:    ts.client,
		zone: testZone,
	}

	actual, err := instancePool.AntiAffinityGroups(context.Background())
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestInstancePool_ElasticIPs() {
	ts.mockAPIRequest(
		"GET",
		fmt.Sprintf("/elastic-ip/%s", testElasticIPID),
		papi.ElasticIp{
			Id: &testElasticIPID,
			Ip: &testElasticIPAddress,
		},
	)

	expected := []*ElasticIP{{
		ID:        &testElasticIPID,
		IPAddress: &testElasticIPAddressP,

		c:    ts.client,
		zone: testZone,
	}}

	instancePool := &InstancePool{
		ElasticIPIDs: &[]string{testElasticIPID},

		c:    ts.client,
		zone: testZone,
	}

	actual, err := instancePool.ElasticIPs(context.Background())
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestInstancePool_Instances() {
	ts.mockAPIRequest("GET", fmt.Sprintf("/instance/%s", testInstanceID), papi.Instance{
		CreatedAt:    &testInstanceCreatedAt,
		DiskSize:     &testInstanceDiskSize,
		Id:           &testInstanceID,
		InstanceType: &papi.InstanceType{Id: &testInstanceInstanceTypeID},
		Name:         &testInstanceName,
		State:        &testInstanceState,
		Template:     &papi.Template{Id: &testInstanceTemplateID},
	})

	expected := []*Instance{{
		CreatedAt:      &testInstanceCreatedAt,
		DiskSize:       &testInstanceDiskSize,
		ID:             &testInstanceID,
		InstanceTypeID: &testInstanceInstanceTypeID,
		Name:           &testInstanceName,
		State:          (*string)(&testInstanceState),
		TemplateID:     &testInstanceTemplateID,

		c:    ts.client,
		zone: testZone,
	}}

	instancePool := &InstancePool{
		InstanceIDs: &[]string{testInstanceID},

		c:    ts.client,
		zone: testZone,
	}

	actual, err := instancePool.Instances(context.Background())
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestInstancePool_PrivateNetworks() {
	ts.mockAPIRequest("GET", fmt.Sprintf("/private-network/%s", testPrivateNetworkID), papi.PrivateNetwork{
		Id:   &testPrivateNetworkID,
		Name: &testPrivateNetworkName,
	})

	expected := []*PrivateNetwork{{
		ID:   &testPrivateNetworkID,
		Name: &testPrivateNetworkName,

		c:    ts.client,
		zone: testZone,
	}}

	instancePool := &InstancePool{
		PrivateNetworkIDs: &[]string{testPrivateNetworkID},

		c:    ts.client,
		zone: testZone,
	}

	actual, err := instancePool.PrivateNetworks(context.Background())
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestInstancePool_SecurityGroups() {
	ts.mockAPIRequest("GET", fmt.Sprintf("/security-group/%s", testSecurityGroupID), papi.SecurityGroup{
		Id:   &testSecurityGroupID,
		Name: &testSecurityGroupName,
	})

	expected := []*SecurityGroup{
		{
			ID:   &testSecurityGroupID,
			Name: &testSecurityGroupName,

			c:    ts.client,
			zone: testZone,
		},
	}

	instancePool := &InstancePool{
		SecurityGroupIDs: &[]string{testSecurityGroupID},

		c:    ts.client,
		zone: testZone,
	}

	actual, err := instancePool.SecurityGroups(context.Background())
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestInstancePool_Scale() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = papi.OperationStateSuccess
		testScaleSize      = testInstancePoolSize * 2
		scaled             = false
	)

	instancePool := &InstancePool{
		ID:   &testInstancePoolID,
		c:    ts.client,
		zone: testZone,

		InstanceIDs: &[]string{testInstancePoolID},
	}

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/instance-pool/%s:scale", *instancePool.ID),
		func(req *http.Request) (*http.Response, error) {
			scaled = true

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
	ts.Require().True(scaled)
}

func (ts *clientTestSuite) TestInstancePool_EvictMembers() {
	var (
		testOperationID     = ts.randomID()
		testOperationState  = papi.OperationStateSuccess
		testEvictedMemberID = ts.randomID()
		evicted             = false
	)

	instancePool := &InstancePool{
		ID:   &testInstancePoolID,
		c:    ts.client,
		zone: testZone,
	}

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/instance-pool/%s:evict", *instancePool.ID),
		func(req *http.Request) (*http.Response, error) {
			evicted = true

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
	ts.Require().True(evicted)
}

func (ts *clientTestSuite) TestClient_CreateInstancePool() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = papi.OperationStateSuccess
	)

	httpmock.RegisterResponder("POST", "/instance-pool",
		func(req *http.Request) (*http.Response, error) {
			var actual papi.CreateInstancePoolJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.CreateInstancePoolJSONRequestBody{
				AntiAffinityGroups: &[]papi.AntiAffinityGroup{{Id: &testInstancePoolAntiAffinityGroupID}},
				DeployTarget:       &papi.DeployTarget{Id: &testInstancePoolDeployTargetID},
				Description:        &testInstancePoolDescription,
				DiskSize:           testInstancePoolDiskSize,
				ElasticIps:         &[]papi.ElasticIp{{Id: &testInstancePoolElasticIPID}},
				InstancePrefix:     &testInstancePoolInstancePrefix,
				InstanceType:       papi.InstanceType{Id: &testInstancePoolInstanceTypeID},
				Ipv6Enabled:        &testInstancePoolIPv6Enabled,
				Labels:             &papi.Labels{AdditionalProperties: testInstancePoolLabels},
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
		DeployTarget:       &papi.DeployTarget{Id: &testInstancePoolDeployTargetID},
		Description:        &testInstancePoolDescription,
		DiskSize:           &testInstancePoolDiskSize,
		ElasticIps:         &[]papi.ElasticIp{{Id: &testInstancePoolElasticIPID}},
		Id:                 &testInstancePoolID,
		InstancePrefix:     &testInstancePoolInstancePrefix,
		InstanceType:       &papi.InstanceType{Id: &testInstancePoolInstanceTypeID},
		Instances:          &[]papi.Instance{{Id: &testInstancePoolInstanceID}},
		Ipv6Enabled:        &testInstancePoolIPv6Enabled,
		Labels:             &papi.Labels{AdditionalProperties: testInstancePoolLabels},
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
		ManagerID:            &testInstancePoolManagerID,
		Name:                 &testInstancePoolName,
		PrivateNetworkIDs:    &[]string{testInstancePoolPrivateNetworkID},
		SSHKey:               &testInstancePoolSSHKey,
		SecurityGroupIDs:     &[]string{testInstancePoolSecurityGroupID},
		Size:                 &testInstancePoolSize,
		State:                (*string)(&testInstancePoolState),
		TemplateID:           &testInstancePoolTemplateID,
		UserData:             &testInstancePoolUserData,

		c:    ts.client,
		zone: testZone,
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
			Labels:             &papi.Labels{AdditionalProperties: testInstancePoolLabels},
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
		AntiAffinityGroupIDs: &[]string{testInstancePoolAntiAffinityGroupID},
		Description:          &testInstancePoolDescription,
		DiskSize:             &testInstancePoolDiskSize,
		ElasticIPIDs:         &[]string{testInstancePoolElasticIPID},
		ID:                   &testInstancePoolID,
		IPv6Enabled:          &testInstancePoolIPv6Enabled,
		InstanceIDs:          &[]string{testInstancePoolInstanceID},
		InstanceTypeID:       &testInstancePoolInstanceTypeID,
		Labels:               &testInstancePoolLabels,
		ManagerID:            &testInstancePoolManagerID,
		Name:                 &testInstancePoolName,
		PrivateNetworkIDs:    &[]string{testInstancePoolPrivateNetworkID},
		SSHKey:               &testInstancePoolSSHKey,
		SecurityGroupIDs:     &[]string{testInstancePoolSecurityGroupID},
		Size:                 &testInstancePoolSize,
		State:                (*string)(&testInstancePoolState),
		TemplateID:           &testInstancePoolTemplateID,
		UserData:             &testInstancePoolUserData,

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
		Labels:             &papi.Labels{AdditionalProperties: testInstancePoolLabels},
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
		AntiAffinityGroupIDs: &[]string{testInstancePoolAntiAffinityGroupID},
		Description:          &testInstancePoolDescription,
		DiskSize:             &testInstancePoolDiskSize,
		ElasticIPIDs:         &[]string{testInstancePoolElasticIPID},
		ID:                   &testInstancePoolID,
		IPv6Enabled:          &testInstancePoolIPv6Enabled,
		InstanceIDs:          &[]string{testInstancePoolInstanceID},
		InstanceTypeID:       &testInstancePoolInstanceTypeID,
		Labels:               &testInstancePoolLabels,
		ManagerID:            &testInstancePoolManagerID,
		Name:                 &testInstancePoolName,
		PrivateNetworkIDs:    &[]string{testInstancePoolPrivateNetworkID},
		SSHKey:               &testInstancePoolSSHKey,
		SecurityGroupIDs:     &[]string{testInstancePoolSecurityGroupID},
		Size:                 &testInstancePoolSize,
		State:                (*string)(&testInstancePoolState),
		TemplateID:           &testInstancePoolTemplateID,
		UserData:             &testInstancePoolUserData,

		c:    ts.client,
		zone: testZone,
	}

	actual, err := ts.client.GetInstancePool(context.Background(), testZone, *expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_FindInstancePool() {
	ts.mockAPIRequest("GET", "/instance-pool", struct {
		InstancePools *[]papi.InstancePool `json:"instance-pools,omitempty"`
	}{
		InstancePools: &[]papi.InstancePool{{
			DiskSize:     &testInstancePoolDiskSize,
			Id:           &testInstancePoolID,
			InstanceType: &papi.InstanceType{Id: &testInstancePoolInstanceTypeID},
			Manager:      &papi.Manager{Id: &testInstancePoolManagerID},
			Name:         &testInstancePoolName,
			Size:         &testInstancePoolSize,
			State:        &testInstancePoolState,
			Template:     &papi.Template{Id: &testInstancePoolTemplateID},
		}},
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/instance-pool/%s", testInstancePoolID), papi.InstancePool{
		DiskSize:     &testInstancePoolDiskSize,
		Id:           &testInstancePoolID,
		InstanceType: &papi.InstanceType{Id: &testInstancePoolInstanceTypeID},
		Manager:      &papi.Manager{Id: &testInstancePoolManagerID},
		Name:         &testInstancePoolName,
		Size:         &testInstancePoolSize,
		State:        &testInstancePoolState,
		Template:     &papi.Template{Id: &testInstancePoolTemplateID},
	})

	expected := &InstancePool{
		DiskSize:       &testInstancePoolDiskSize,
		ID:             &testInstancePoolID,
		InstanceTypeID: &testInstancePoolInstanceTypeID,
		ManagerID:      &testInstancePoolManagerID,
		Name:           &testInstancePoolName,
		Size:           &testInstancePoolSize,
		State:          (*string)(&testInstancePoolState),
		TemplateID:     &testInstancePoolTemplateID,

		c:    ts.client,
		zone: testZone,
	}

	actual, err := ts.client.FindInstancePool(context.Background(), testZone, *expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)

	actual, err = ts.client.FindInstancePool(context.Background(), testZone, *expected.Name)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
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
		testOperationState                         = papi.OperationStateSuccess
		updated                                    = false
	)

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/instance-pool/%s", testInstancePoolID),
		func(req *http.Request) (*http.Response, error) {
			updated = true

			var actual papi.UpdateInstancePoolJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.UpdateInstancePoolJSONRequestBody{
				AntiAffinityGroups: &[]papi.AntiAffinityGroup{{Id: &testInstancePoolAntiAffinityGroupIDUpdated}},
				DeployTarget:       &papi.DeployTarget{Id: &testInstancePoolDeployTargetIDUpdated},
				Description:        &testInstancePoolDescriptionUpdated,
				DiskSize:           &testInstancePoolDiskSizeUpdated,
				ElasticIps:         &[]papi.ElasticIp{{Id: &testInstancePoolElasticIPIDUpdated}},
				InstancePrefix:     &testInstancePoolInstancePrefixUpdated,
				InstanceType:       &papi.InstanceType{Id: &testInstancePoolInstanceTypeIDUpdated},
				Ipv6Enabled:        &testInstancePoolIPv6EnabledUpdated,
				Labels:             &papi.Labels{AdditionalProperties: testInstancePoolLabelsUpdated},
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

func (ts *clientTestSuite) TestClient_DeleteInstancePool() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = papi.OperationStateSuccess
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
