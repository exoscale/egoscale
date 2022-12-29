package v2

import (
	"context"
	"net/http"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/exoscale/egoscale/v2/oapi"
)

var (
	testSKSNodepoolAddons                    = []oapi.SksNodepoolAddons{oapi.SksNodepoolAddonsStorageLvm}
	testSKSNodepoolAntiAffinityGroupID       = new(testSuite).randomID()
	testSKSNodepoolCreatedAt, _              = time.Parse(iso8601Format, "2020-11-18T07:54:36Z")
	testSKSNodepoolDeployTargetID            = new(testSuite).randomID()
	testSKSNodepoolDescription               = new(testSuite).randomString(10)
	testSKSNodepoolDiskSize            int64 = 15
	testSKSNodepoolID                        = new(testSuite).randomID()
	testSKSNodepoolInstancePoolID            = new(testSuite).randomID()
	testSKSNodepoolInstancePrefix            = new(testSuite).randomString(10)
	testSKSNodepoolInstanceTypeID            = new(testSuite).randomID()
	testSKSNodepoolLabels                    = map[string]string{"k1": "v1", "k2": "v2"}
	testSKSNodepoolName                      = new(testSuite).randomString(10)
	testSKSNodepoolPrivateNetworkID          = new(testSuite).randomID()
	testSKSNodepoolSecurityGroupID           = new(testSuite).randomID()
	testSKSNodepoolSize                int64 = 3
	testSKSNodepoolState                     = oapi.SksNodepoolStateRunning
	testSKSNodepoolTaintEffect               = oapi.SksNodepoolTaintEffectNoExecute
	testSKSNodepoolTaintKey                  = new(testSuite).randomString(10)
	testSKSNodepoolTaintValue                = new(testSuite).randomString(10)
	testSKSNodepoolTemplateID                = new(testSuite).randomID()
	testSKSNodepoolVersion                   = "1.18.6"
)

func (ts *testSuite) TestCLient_CreateSKSNodepool() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
	)

	ts.mock().
		On(
			"CreateSksNodepoolWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testSKSClusterID, args.Get(1))
			ts.Require().Equal(
				oapi.CreateSksNodepoolJSONRequestBody{
					Addons: &[]oapi.CreateSksNodepoolJSONBodyAddons{
						oapi.CreateSksNodepoolJSONBodyAddons(testSKSNodepoolAddons[0]),
					},
					AntiAffinityGroups: &[]oapi.AntiAffinityGroup{{Id: &testSKSNodepoolAntiAffinityGroupID}},
					DeployTarget:       &oapi.DeployTarget{Id: &testSKSNodepoolDeployTargetID},
					Description:        &testSKSNodepoolDescription,
					DiskSize:           testSKSNodepoolDiskSize,
					InstancePrefix:     &testSKSNodepoolInstancePrefix,
					InstanceType:       oapi.InstanceType{Id: &testSKSNodepoolInstanceTypeID},
					Labels:             &oapi.Labels{AdditionalProperties: testSKSNodepoolLabels},
					Name:               testSKSNodepoolName,
					PrivateNetworks:    &[]oapi.PrivateNetwork{{Id: &testSKSNodepoolPrivateNetworkID}},
					SecurityGroups:     &[]oapi.SecurityGroup{{Id: &testSKSNodepoolSecurityGroupID}},
					Size:               testSKSNodepoolSize,
					Taints: &oapi.SksNodepoolTaints{
						AdditionalProperties: map[string]oapi.SksNodepoolTaint{
							testSKSNodepoolTaintKey: {
								Effect: testSKSNodepoolTaintEffect,
								Value:  testSKSNodepoolTaintValue,
							},
						},
					},
				},
				args.Get(2),
			)
		}).
		Return(
			&oapi.CreateSksNodepoolResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testSKSNodepoolID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testSKSNodepoolID, nil),
		State:     &testOperationState,
	})

	ts.mock().
		On("GetSksNodepoolWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // sksNodepoolId
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.GetSksNodepoolResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.SksNodepool{
				Addons:             &testSKSNodepoolAddons,
				AntiAffinityGroups: &[]oapi.AntiAffinityGroup{{Id: &testSKSNodepoolAntiAffinityGroupID}},
				CreatedAt:          &testSKSNodepoolCreatedAt,
				DeployTarget:       &oapi.DeployTarget{Id: &testSKSNodepoolDeployTargetID},
				Description:        &testSKSNodepoolDescription,
				DiskSize:           &testSKSNodepoolDiskSize,
				Id:                 &testSKSNodepoolID,
				InstancePool:       &oapi.InstancePool{Id: &testSKSNodepoolInstancePoolID},
				InstancePrefix:     &testSKSNodepoolInstancePrefix,
				InstanceType:       &oapi.InstanceType{Id: &testSKSNodepoolInstanceTypeID},
				Labels:             &oapi.Labels{AdditionalProperties: testSKSNodepoolLabels},
				Name:               &testSKSNodepoolName,
				PrivateNetworks:    &[]oapi.PrivateNetwork{{Id: &testSKSNodepoolPrivateNetworkID}},
				SecurityGroups:     &[]oapi.SecurityGroup{{Id: &testSKSNodepoolSecurityGroupID}},
				Size:               &testSKSNodepoolSize,
				State:              &testSKSNodepoolState,
				Taints: &oapi.SksNodepoolTaints{
					AdditionalProperties: map[string]oapi.SksNodepoolTaint{
						testSKSNodepoolTaintKey: {
							Effect: testSKSNodepoolTaintEffect,
							Value:  testSKSNodepoolTaintValue,
						},
					},
				},
				Template: &oapi.Template{Id: &testSKSNodepoolTemplateID},
				Version:  &testSKSNodepoolVersion,
			},
		}, nil)

	cluster := &SKSCluster{
		ID: &testSKSClusterID,
	}

	expected := &SKSNodepool{
		AddOns:               &[]string{string(testSKSNodepoolAddons[0])},
		AntiAffinityGroupIDs: &[]string{testSKSNodepoolAntiAffinityGroupID},
		CreatedAt:            &testSKSNodepoolCreatedAt,
		DeployTargetID:       &testSKSNodepoolDeployTargetID,
		Description:          &testSKSNodepoolDescription,
		DiskSize:             &testSKSNodepoolDiskSize,
		ID:                   &testSKSNodepoolID,
		InstancePoolID:       &testSKSNodepoolInstancePoolID,
		InstancePrefix:       &testSKSNodepoolInstancePrefix,
		InstanceTypeID:       &testSKSNodepoolInstanceTypeID,
		Labels:               &testSKSNodepoolLabels,
		Name:                 &testSKSNodepoolName,
		PrivateNetworkIDs:    &[]string{testSKSNodepoolPrivateNetworkID},
		SecurityGroupIDs:     &[]string{testSKSNodepoolSecurityGroupID},
		Size:                 &testSKSNodepoolSize,
		State:                (*string)(&testSKSNodepoolState),
		Taints: &map[string]*SKSNodepoolTaint{
			testSKSNodepoolTaintKey: {
				Effect: string(testSKSNodepoolTaintEffect),
				Value:  testSKSNodepoolTaintValue,
			},
		},
		TemplateID: &testSKSNodepoolTemplateID,
		Version:    &testSKSNodepoolVersion,
	}

	actual, err := ts.client.CreateSKSNodepool(context.Background(), testZone, cluster, expected)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_DeleteSKSNodepool() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		deleted            = false
	)

	ts.mock().
		On(
			"DeleteSksNodepoolWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // sksNodepoolId
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testSKSClusterID, args.Get(1))
			ts.Require().Equal(testSKSNodepoolID, args.Get(2))
			deleted = true
		}).
		Return(
			&oapi.DeleteSksNodepoolResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testSKSNodepoolID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testSKSNodepoolID, nil),
		State:     &testOperationState,
	})

	cluster := &SKSCluster{
		ID:        &testSKSClusterID,
		Nodepools: []*SKSNodepool{{ID: &testSKSNodepoolID}},
	}

	ts.Require().NoError(ts.client.DeleteSKSNodepool(
		context.Background(),
		testZone,
		cluster,
		cluster.Nodepools[0]),
	)
	ts.Require().True(deleted)
}

func (ts *testSuite) TestClient_EvictSKSNodepoolMembers() {
	var (
		testOperationID     = ts.randomID()
		testOperationState  = oapi.OperationStateSuccess
		testEvictedMemberID = ts.randomID()
		evicted             = false
	)

	cluster := &SKSCluster{
		ID:        &testSKSClusterID,
		Nodepools: []*SKSNodepool{{ID: &testSKSNodepoolID}},
	}

	ts.mock().
		On(
			"EvictSksNodepoolMembersWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // sksNodepoolId
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testSKSClusterID, args.Get(1))
			ts.Require().Equal(testSKSNodepoolID, args.Get(2))
			ts.Require().Equal(
				oapi.EvictSksNodepoolMembersJSONRequestBody{Instances: &[]string{testEvictedMemberID}},
				args.Get(3),
			)
			evicted = true
		}).
		Return(
			&oapi.EvictSksNodepoolMembersResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testSKSNodepoolID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testSKSNodepoolID, nil),
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.EvictSKSNodepoolMembers(
		context.Background(),
		testZone,
		cluster,
		cluster.Nodepools[0],
		[]string{testEvictedMemberID}),
	)
	ts.Require().True(evicted)
}

func (ts *testSuite) TestClient_ScaleSKSNodepool() {
	var (
		testOperationID          = ts.randomID()
		testOperationState       = oapi.OperationStateSuccess
		testScaleSize      int64 = 3
		scaled                   = false
	)

	cluster := &SKSCluster{
		ID:        &testSKSClusterID,
		Nodepools: []*SKSNodepool{{ID: &testSKSNodepoolID}},
	}

	ts.mock().
		On(
			"ScaleSksNodepoolWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // sksNodepoolId
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testSKSClusterID, args.Get(1))
			ts.Require().Equal(testSKSNodepoolID, args.Get(2))
			ts.Require().Equal(oapi.ScaleSksNodepoolJSONRequestBody{Size: testScaleSize}, args.Get(3))
			scaled = true
		}).
		Return(
			&oapi.ScaleSksNodepoolResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testSKSNodepoolID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testSKSNodepoolID, nil),
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.ScaleSKSNodepool(
		context.Background(),
		testZone,
		cluster,
		cluster.Nodepools[0],
		testScaleSize,
	))
	ts.Require().True(scaled)
}

func (ts *testSuite) TestClient_UpdateSKSNodepool() {
	var (
		testOperationID                           = ts.randomID()
		testSKSNodepoolAntiAffinityGroupIDUpdated = ts.randomID()
		testSKSNodepoolDeployTargetIDUpdated      = ts.randomID()
		testSKSNodepoolDescriptionUpdated         = testSKSNodepoolDescription + "-updated"
		testSKSNodepoolDiskSizeUpdated            = testSKSNodepoolDiskSize + 1
		testSKSNodepoolInstancePrefixUpdated      = testSKSNodepoolInstancePrefix + "-updated"
		testSKSNodepoolInstanceTypeIDUpdated      = testSKSNodepoolInstanceTypeID + "-updated"
		testSKSNodepoolLabelsUpdated              = map[string]string{"k3": "v3"}
		testSKSNodepoolNameUpdated                = testSKSNodepoolName + "-updated"
		testSKSNodepoolPrivateNetworkIDUpdated    = ts.randomID()
		testSKSNodepoolSecurityGroupIDUpdated     = ts.randomID()
		testOperationState                        = oapi.OperationStateSuccess
		updated                                   = false
	)

	cluster := &SKSCluster{
		ID: &testSKSClusterID,
		Nodepools: []*SKSNodepool{{
			DeployTargetID: &testSKSNodepoolDeployTargetID,
			Description:    &testSKSNodepoolDescription,
			DiskSize:       &testSKSNodepoolDiskSize,
			ID:             &testSKSNodepoolID,
			InstancePrefix: &testSKSNodepoolInstancePrefix,
			InstanceTypeID: &testSKSNodepoolInstanceTypeID,
			Labels:         &testSKSNodepoolLabels,
			Name:           &testSKSNodepoolName,
		}},
	}

	ts.mock().
		On(
			"UpdateSksNodepoolWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // sksNodepoolId
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testSKSClusterID, args.Get(1))
			ts.Require().Equal(testSKSNodepoolID, args.Get(2))
			ts.Require().Equal(
				oapi.UpdateSksNodepoolJSONRequestBody{
					AntiAffinityGroups: &[]oapi.AntiAffinityGroup{{Id: &testSKSNodepoolAntiAffinityGroupIDUpdated}},
					DeployTarget:       &oapi.DeployTarget{Id: &testSKSNodepoolDeployTargetIDUpdated},
					Description:        &testSKSNodepoolDescriptionUpdated,
					DiskSize:           &testSKSNodepoolDiskSizeUpdated,
					InstancePrefix:     &testSKSNodepoolInstancePrefixUpdated,
					InstanceType:       &oapi.InstanceType{Id: &testSKSNodepoolInstanceTypeIDUpdated},
					Labels:             &oapi.Labels{AdditionalProperties: testSKSNodepoolLabelsUpdated},
					Name:               &testSKSNodepoolNameUpdated,
					PrivateNetworks:    &[]oapi.PrivateNetwork{{Id: &testSKSNodepoolPrivateNetworkIDUpdated}},
					SecurityGroups:     &[]oapi.SecurityGroup{{Id: &testSKSNodepoolSecurityGroupIDUpdated}},
					Taints: &oapi.SksNodepoolTaints{
						AdditionalProperties: map[string]oapi.SksNodepoolTaint{
							testSKSNodepoolTaintKey: {
								Effect: testSKSNodepoolTaintEffect,
								Value:  testSKSNodepoolTaintValue,
							},
						},
					},
				},
				args.Get(3),
			)
			updated = true
		}).
		Return(
			&oapi.UpdateSksNodepoolResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testSKSNodepoolID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testSKSNodepoolID, nil),
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.UpdateSKSNodepool(context.Background(), testZone, cluster, &SKSNodepool{
		AntiAffinityGroupIDs: &[]string{testSKSNodepoolAntiAffinityGroupIDUpdated},
		DeployTargetID:       &testSKSNodepoolDeployTargetIDUpdated,
		Description:          &testSKSNodepoolDescriptionUpdated,
		DiskSize:             &testSKSNodepoolDiskSizeUpdated,
		ID:                   cluster.Nodepools[0].ID,
		InstancePrefix:       &testSKSNodepoolInstancePrefixUpdated,
		InstanceTypeID:       &testSKSNodepoolInstanceTypeIDUpdated,
		Labels:               &testSKSNodepoolLabelsUpdated,
		Name:                 &testSKSNodepoolNameUpdated,
		PrivateNetworkIDs:    &[]string{testSKSNodepoolPrivateNetworkIDUpdated},
		SecurityGroupIDs:     &[]string{testSKSNodepoolSecurityGroupIDUpdated},
		Taints: &map[string]*SKSNodepoolTaint{
			testSKSNodepoolTaintKey: {
				Effect: string(testSKSNodepoolTaintEffect),
				Value:  testSKSNodepoolTaintValue,
			},
		},
	}))
	ts.Require().True(updated)
}
