package v2

import (
	"context"
	"net/http"

	"github.com/stretchr/testify/mock"

	"github.com/exoscale/egoscale/v2/oapi"
)

var (
	testInstancePoolAntiAffinityGroupID       = new(testSuite).randomID()
	testInstancePoolDeployTargetID            = new(testSuite).randomID()
	testInstancePoolDescription               = new(testSuite).randomString(10)
	testInstancePoolDiskSize            int64 = 10
	testInstancePoolElasticIPID               = new(testSuite).randomID()
	testInstancePoolID                        = new(testSuite).randomID()
	testInstancePoolIPv6Enabled               = true
	testInstancePoolInstanceID                = new(testSuite).randomID()
	testInstancePoolInstancePrefix            = new(testSuite).randomString(10)
	testInstancePoolInstanceTypeID            = new(testSuite).randomID()
	testInstancePoolLabels                    = map[string]string{"k1": "v1", "k2": "v2"}
	testInstancePoolManagerID                 = new(testSuite).randomID()
	testInstancePoolManagerType               = oapi.ManagerTypeSksNodepool
	testInstancePoolName                      = new(testSuite).randomString(10)
	testInstancePoolPrivateNetworkID          = new(testSuite).randomID()
	testInstancePoolSSHKey                    = new(testSuite).randomString(10)
	testInstancePoolSecurityGroupID           = new(testSuite).randomID()
	testInstancePoolSize                int64 = 3
	testInstancePoolState                     = oapi.InstancePoolStateRunning
	testInstancePoolTemplateID                = new(testSuite).randomID()
	testInstancePoolUserData                  = "I2Nsb3VkLWNvbmZpZwphcHRfdXBncmFkZTogdHJ1ZQ=="
)

func (ts *testSuite) TestClient_CreateInstancePool() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
	)

	ts.mock().
		On(
			"CreateInstancePoolWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				oapi.CreateInstancePoolJSONRequestBody{
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
				},
				args.Get(1),
			)
		}).
		Return(
			&oapi.CreateInstancePoolResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testInstancePoolID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testInstancePoolID, nil),
		State:     &testOperationState,
	})

	ts.mock().
		On("GetInstancePoolWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.GetInstancePoolResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.InstancePool{
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
			},
		}, nil)

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

func (ts *testSuite) TestClient_DeleteInstancePool() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		deleted            = false
	)

	ts.mock().
		On(
			"DeleteInstancePoolWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testInstancePoolID, args.Get(1))
			deleted = true
		}).
		Return(
			&oapi.DeleteInstancePoolResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testInstancePoolID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testInstancePoolID, nil),
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.DeleteInstancePool(
		context.Background(),
		testZone,
		&InstancePool{ID: &testInstancePoolID},
	))
	ts.Require().True(deleted)
}

func (ts *testSuite) TestClient_EvictInstancePooltMembers() {
	var (
		testOperationID     = ts.randomID()
		testOperationState  = oapi.OperationStateSuccess
		testEvictedMemberID = ts.randomID()
		evicted             = false
	)

	instancePool := &InstancePool{
		ID: &testInstancePoolID,
	}

	ts.mock().
		On(
			"EvictInstancePoolMembersWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testInstancePoolID, args.Get(1))
			ts.Require().Equal(
				oapi.EvictInstancePoolMembersJSONRequestBody{Instances: &[]string{testEvictedMemberID}},
				args.Get(2),
			)
			evicted = true
		}).
		Return(
			&oapi.EvictInstancePoolMembersResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testInstancePoolID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testInstancePoolID, nil),
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.EvictInstancePoolMembers(
		context.Background(),
		testZone,
		instancePool,
		[]string{testEvictedMemberID},
	))
	ts.Require().True(evicted)
}

func (ts *testSuite) TestClient_FindInstancePool() {
	ts.mock().
		On("ListInstancePoolsWithResponse",
			mock.Anything,                 // ctx
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.ListInstancePoolsResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &struct {
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
			},
		}, nil)

	ts.mock().
		On("GetInstancePoolWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testInstancePoolID, args.Get(1))
		}).
		Return(&oapi.GetInstancePoolResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.InstancePool{
				DiskSize:     &testInstancePoolDiskSize,
				Id:           &testInstancePoolID,
				InstanceType: &oapi.InstanceType{Id: &testInstancePoolInstanceTypeID},
				Manager:      &oapi.Manager{Id: &testInstancePoolManagerID, Type: &testInstancePoolManagerType},
				Name:         &testInstancePoolName,
				Size:         &testInstancePoolSize,
				State:        &testInstancePoolState,
				Template:     &oapi.Template{Id: &testInstancePoolTemplateID},
			},
		}, nil)

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

func (ts *testSuite) TestClient_GetInstancePool() {
	ts.mock().
		On("GetInstancePoolWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testInstancePoolID, args.Get(1))
		}).
		Return(&oapi.GetInstancePoolResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.InstancePool{
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
			},
		}, nil)

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

func (ts *testSuite) TestClient_ListInstancePools() {
	ts.mock().
		On("ListInstancePoolsWithResponse",
			mock.Anything,                 // ctx
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.ListInstancePoolsResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &struct {
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
			},
		}, nil)

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

func (ts *testSuite) TestClient_ScaleInstancePool() {
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

	ts.mock().
		On(
			"ScaleInstancePoolWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testInstancePoolID, args.Get(1))
			ts.Require().Equal(oapi.ScaleInstancePoolJSONRequestBody{Size: testScaleSize}, args.Get(2))
			scaled = true
		}).
		Return(
			&oapi.ScaleInstancePoolResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testInstancePoolID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testInstancePoolID, nil),
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.ScaleInstancePool(context.Background(), testZone, instancePool, testScaleSize))
	ts.Require().True(scaled)
}

func (ts *testSuite) TestClient_UpdateInstancePool() {
	var (
		testInstancePoolAntiAffinityGroupIDUpdated = new(testSuite).randomID()
		testInstancePoolDeployTargetIDUpdated      = new(testSuite).randomID()
		testInstancePoolDescriptionUpdated         = testInstancePoolDescription + "-updated"
		testInstancePoolDiskSizeUpdated            = testInstancePoolDiskSize * 2
		testInstancePoolElasticIPIDUpdated         = new(testSuite).randomID()
		testInstancePoolIPv6EnabledUpdated         = true
		testInstancePoolInstancePrefixUpdated      = testInstancePoolInstancePrefix + "-updated"
		testInstancePoolInstanceTypeIDUpdated      = new(testSuite).randomID()
		testInstancePoolLabelsUpdated              = map[string]string{"k3": "v3"}
		testInstancePoolNameUpdated                = testInstancePoolName + "-updated"
		testInstancePoolPrivateNetworkIDUpdated    = new(testSuite).randomID()
		testInstancePoolSecurityGroupIDUpdated     = new(testSuite).randomID()
		testInstancePoolSSHKeyUpdated              = testInstancePoolSSHKey + "-updated"
		testInstancePoolTemplateIDUpdated          = new(testSuite).randomID()
		testInstancePoolUserDataUpdated            = testInstancePoolUserData + "-updated"
		testOperationID                            = ts.randomID()
		testOperationState                         = oapi.OperationStateSuccess
		updated                                    = false
	)

	ts.mock().
		On(
			"UpdateInstancePoolWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testInstancePoolID, args.Get(1))
			ts.Require().Equal(
				oapi.UpdateInstancePoolJSONRequestBody{
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
				},
				args.Get(2),
			)
			updated = true
		}).
		Return(
			&oapi.UpdateInstancePoolResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testInstancePoolID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testInstancePoolID, nil),
		State:     &testOperationState,
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
