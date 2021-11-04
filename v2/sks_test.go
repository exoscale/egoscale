package v2

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/exoscale/egoscale/v2/oapi"
)

var (
	testSKSClusterAddons                     = []oapi.SksClusterAddons{oapi.SksClusterAddonsExoscaleCloudController}
	testSKSClusterAutoUpgrade                = true
	testSKSClusterCNI                        = "calico"
	testSKSClusterCreatedAt, _               = time.Parse(iso8601Format, "2020-11-16T10:41:58Z")
	testSKSClusterDescription                = new(testSuite).randomString(10)
	testSKSClusterEndpoint                   = fmt.Sprintf("%s.sks-%s.exo.io", testSKSClusterID, testZone)
	testSKSClusterID                         = new(testSuite).randomID()
	testSKSClusterLabels                     = map[string]string{"k1": "v1", "k2": "v2"}
	testSKSClusterName                       = new(testSuite).randomString(10)
	testSKSClusterServiceLevel               = oapi.SksClusterLevelPro
	testSKSClusterState                      = oapi.SksClusterStateRunning
	testSKSClusterVersion                    = "1.18.6"
	testSKSNodepoolAddons                    = []oapi.SksNodepoolAddons{oapi.SksNodepoolAddonsLinbit}
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

func (ts *testSuite) TestClient_CreateSKSCluster() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
	)

	ts.mock().
		On(
			"CreateSksClusterWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				oapi.CreateSksClusterJSONRequestBody{
					Addons: &[]oapi.CreateSksClusterJSONBodyAddons{
						oapi.CreateSksClusterJSONBodyAddons(testSKSClusterAddons[0]),
					},
					AutoUpgrade: &testSKSClusterAutoUpgrade,
					Cni:         (*oapi.CreateSksClusterJSONBodyCni)(&testSKSClusterCNI),
					Description: &testSKSClusterDescription,
					Labels:      &oapi.Labels{AdditionalProperties: testSKSClusterLabels},
					Level:       oapi.CreateSksClusterJSONBodyLevel(testSKSClusterServiceLevel),
					Name:        testSKSClusterName,
					Version:     testSKSClusterVersion,
				},
				args.Get(1),
			)
		}).
		Return(
			&oapi.CreateSksClusterResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: &oapi.Reference{Id: &testSKSClusterID},
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: &oapi.Reference{Id: &testSKSClusterID},
		State:     &testOperationState,
	})

	ts.mock().
		On("GetSksClusterWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.GetSksClusterResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.SksCluster{
				Addons:      &testSKSClusterAddons,
				AutoUpgrade: &testSKSClusterAutoUpgrade,
				Cni:         (*oapi.SksClusterCni)(&testSKSClusterCNI),
				CreatedAt:   &testSKSClusterCreatedAt,
				Description: &testSKSClusterDescription,
				Endpoint:    &testSKSClusterEndpoint,
				Id:          &testSKSClusterID,
				Labels:      &oapi.Labels{AdditionalProperties: testSKSClusterLabels},
				Level:       &testSKSClusterServiceLevel,
				Name:        &testSKSClusterName,
				State:       &testSKSClusterState,
				Version:     &testSKSClusterVersion,
			},
		}, nil)

	expected := &SKSCluster{
		AddOns:       &[]string{string(testSKSClusterAddons[0])},
		AutoUpgrade:  &testSKSClusterAutoUpgrade,
		CNI:          &testSKSClusterCNI,
		CreatedAt:    &testSKSClusterCreatedAt,
		Description:  &testSKSClusterDescription,
		Endpoint:     &testSKSClusterEndpoint,
		ID:           &testSKSClusterID,
		Labels:       &testSKSClusterLabels,
		Name:         &testSKSClusterName,
		Nodepools:    []*SKSNodepool{},
		ServiceLevel: (*string)(&testSKSClusterServiceLevel),
		State:        (*string)(&testSKSClusterState),
		Version:      &testSKSClusterVersion,
		Zone:         &testZone,
	}

	actual, err := ts.client.CreateSKSCluster(context.Background(), testZone, &SKSCluster{
		AddOns:       &[]string{string(testSKSClusterAddons[0])},
		AutoUpgrade:  &testSKSClusterAutoUpgrade,
		CNI:          &testSKSClusterCNI,
		Description:  &testSKSClusterDescription,
		Labels:       &testSKSClusterLabels,
		Name:         &testSKSClusterName,
		ServiceLevel: (*string)(&testSKSClusterServiceLevel),
		Version:      &testSKSClusterVersion,
	})
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

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
					Reference: &oapi.Reference{Id: &testSKSNodepoolID},
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: &oapi.Reference{Id: &testSKSNodepoolID},
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

func (ts *testSuite) TestClient_DeleteSKSCluster() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		deleted            = false
	)

	ts.mock().
		On(
			"DeleteSksClusterWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testSKSClusterID, args.Get(1))
			deleted = true
		}).
		Return(
			&oapi.DeleteSksClusterResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: &oapi.Reference{Id: &testSKSClusterID},
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: &oapi.Reference{Id: &testSKSClusterID},
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.DeleteSKSCluster(
		context.Background(),
		testZone,
		&SKSCluster{ID: &testSKSClusterID},
	))
	ts.Require().True(deleted)
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
					Reference: &oapi.Reference{Id: &testSKSNodepoolID},
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: &oapi.Reference{Id: &testSKSNodepoolID},
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
					Reference: &oapi.Reference{Id: &testSKSNodepoolID},
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: &oapi.Reference{Id: &testSKSNodepoolID},
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

func (ts *testSuite) TestClient_FindSKSCluster() {
	ts.mock().
		On("ListSksClustersWithResponse",
			mock.Anything,                 // ctx
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.ListSksClustersResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &struct {
				SksClusters *[]oapi.SksCluster `json:"sks-clusters,omitempty"`
			}{
				SksClusters: &[]oapi.SksCluster{{
					CreatedAt: &testSKSClusterCreatedAt,
					Endpoint:  &testSKSClusterEndpoint,
					Id:        &testSKSClusterID,
					Level:     &testSKSClusterServiceLevel,
					Name:      &testSKSClusterName,
					State:     &testSKSClusterState,
					Version:   &testSKSClusterVersion,
				}},
			},
		}, nil)

	ts.mock().
		On("GetSksClusterWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testSKSClusterID, args.Get(1))
		}).
		Return(&oapi.GetSksClusterResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.SksCluster{
				CreatedAt: &testSKSClusterCreatedAt,
				Endpoint:  &testSKSClusterEndpoint,
				Id:        &testSKSClusterID,
				Level:     &testSKSClusterServiceLevel,
				Name:      &testSKSClusterName,
				State:     &testSKSClusterState,
				Version:   &testSKSClusterVersion,
			},
		}, nil)

	expected := &SKSCluster{
		CreatedAt:    &testSKSClusterCreatedAt,
		Endpoint:     &testSKSClusterEndpoint,
		ID:           &testSKSClusterID,
		Name:         &testSKSClusterName,
		Nodepools:    []*SKSNodepool{},
		ServiceLevel: (*string)(&testSKSClusterServiceLevel),
		State:        (*string)(&testSKSClusterState),
		Version:      &testSKSClusterVersion,
		Zone:         &testZone,
	}

	actual, err := ts.client.FindSKSCluster(context.Background(), testZone, *expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)

	actual, err = ts.client.FindSKSCluster(context.Background(), testZone, *expected.Name)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_GetSKSCluster() {
	ts.mock().
		On("GetSksClusterWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testSKSClusterID, args.Get(1))
		}).
		Return(&oapi.GetSksClusterResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.SksCluster{
				Addons:      &testSKSClusterAddons,
				AutoUpgrade: &testSKSClusterAutoUpgrade,
				Cni:         (*oapi.SksClusterCni)(&testSKSClusterCNI),
				CreatedAt:   &testSKSClusterCreatedAt,
				Description: &testSKSClusterDescription,
				Endpoint:    &testSKSClusterEndpoint,
				Id:          &testSKSClusterID,
				Labels:      &oapi.Labels{AdditionalProperties: testSKSClusterLabels},
				Level:       &testSKSClusterServiceLevel,
				Name:        &testSKSClusterName,
				Nodepools: &[]oapi.SksNodepool{{
					AntiAffinityGroups: &[]oapi.AntiAffinityGroup{{Id: &testSKSNodepoolAntiAffinityGroupID}},
					CreatedAt:          &testSKSNodepoolCreatedAt,
					Description:        &testSKSNodepoolDescription,
					DiskSize:           &testSKSNodepoolDiskSize,
					Id:                 &testSKSNodepoolID,
					InstancePool:       &oapi.InstancePool{Id: &testSKSNodepoolInstancePoolID},
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
				}},
				State:   &testSKSClusterState,
				Version: &testSKSClusterVersion,
			},
		}, nil)

	expected := &SKSCluster{
		AddOns:      &[]string{string(testSKSClusterAddons[0])},
		AutoUpgrade: &testSKSClusterAutoUpgrade,
		CNI:         &testSKSClusterCNI,
		CreatedAt:   &testSKSClusterCreatedAt,
		Description: &testSKSClusterDescription,
		Endpoint:    &testSKSClusterEndpoint,
		ID:          &testSKSClusterID,
		Labels:      &testSKSClusterLabels,
		Name:        &testSKSClusterName,
		Nodepools: []*SKSNodepool{{
			AntiAffinityGroupIDs: &[]string{testSKSNodepoolAntiAffinityGroupID},
			CreatedAt:            &testSKSNodepoolCreatedAt,
			Description:          &testSKSNodepoolDescription,
			DiskSize:             &testSKSNodepoolDiskSize,
			ID:                   &testSKSNodepoolID,
			InstancePoolID:       &testSKSNodepoolInstancePoolID,
			InstanceTypeID:       &testSKSNodepoolInstanceTypeID,
			Labels:               &testSKSNodepoolLabels,
			Name:                 &testSKSNodepoolName,
			PrivateNetworkIDs:    &[]string{testSKSNodepoolPrivateNetworkID},
			SecurityGroupIDs:     &[]string{testSKSNodepoolSecurityGroupID},
			Size:                 &testSKSNodepoolSize,
			State:                (*string)(&testSKSClusterState),
			Taints: &map[string]*SKSNodepoolTaint{
				testSKSNodepoolTaintKey: {
					Effect: string(testSKSNodepoolTaintEffect),
					Value:  testSKSNodepoolTaintValue,
				},
			},
			TemplateID: &testSKSNodepoolTemplateID,
			Version:    &testSKSNodepoolVersion,
		}},
		ServiceLevel: (*string)(&testSKSClusterServiceLevel),
		State:        (*string)(&testSKSClusterState),
		Version:      &testSKSClusterVersion,
		Zone:         &testZone,
	}

	actual, err := ts.client.GetSKSCluster(context.Background(), testZone, *expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_GetSKSClusterAuthorityCert() {
	var (
		testAuthority   = "aggregation"
		testCertificate = base64.StdEncoding.EncodeToString([]byte("test"))
	)

	cluster := &SKSCluster{
		ID: &testSKSClusterID,
	}

	ts.mock().
		On(
			"GetSksClusterAuthorityCertWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // authority
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testSKSClusterID, args.Get(1))
			ts.Require().Equal(oapi.GetSksClusterAuthorityCertParamsAuthority(testAuthority), args.Get(2))
		}).
		Return(
			&oapi.GetSksClusterAuthorityCertResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &struct {
					Cacert *string `json:"cacert,omitempty"`
				}{
					Cacert: &testCertificate,
				},
			},
			nil,
		)

	actual, err := ts.client.GetSKSClusterAuthorityCert(context.Background(), testZone, cluster, testAuthority)
	ts.Require().NoError(err)
	ts.Require().Equal(testCertificate, actual)
}

func (ts *testSuite) TestClient_GetSKSClusterKubeconfig() {
	var (
		testRequestUser   = "test-user"
		testRequestGroups = []string{"system:masters"}
		testKubeconfig    = base64.StdEncoding.EncodeToString([]byte("test"))
		testTTL           = time.Hour
		testTTLSec        = int64(testTTL.Seconds())
	)

	cluster := &SKSCluster{
		ID: &testSKSClusterID,
	}

	ts.mock().
		On(
			"GenerateSksClusterKubeconfigWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testSKSClusterID, args.Get(1))
			ts.Require().Equal(oapi.GenerateSksClusterKubeconfigJSONRequestBody{
				Groups: &testRequestGroups,
				Ttl:    &testTTLSec,
				User:   &testRequestUser,
			}, args.Get(2))
		}).
		Return(
			&oapi.GenerateSksClusterKubeconfigResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &struct {
					Kubeconfig *string `json:"kubeconfig,omitempty"`
				}{
					Kubeconfig: &testKubeconfig,
				},
			},
			nil,
		)

	actual, err := ts.client.GetSKSClusterKubeconfig(
		context.Background(),
		testZone,
		cluster,
		testRequestUser,
		testRequestGroups,
		testTTL,
	)
	ts.Require().NoError(err)
	ts.Require().Equal(testKubeconfig, actual)
}

func (ts *testSuite) TestClient_ListSKSClusters() {
	ts.mock().
		On("ListSksClustersWithResponse",
			mock.Anything,                 // ctx
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.ListSksClustersResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &struct {
				SksClusters *[]oapi.SksCluster `json:"sks-clusters,omitempty"`
			}{
				SksClusters: &[]oapi.SksCluster{{
					Addons:      &testSKSClusterAddons,
					AutoUpgrade: &testSKSClusterAutoUpgrade,
					Cni:         (*oapi.SksClusterCni)(&testSKSClusterCNI),
					CreatedAt:   &testSKSClusterCreatedAt,
					Description: &testSKSClusterDescription,
					Endpoint:    &testSKSClusterEndpoint,
					Id:          &testSKSClusterID,
					Labels:      &oapi.Labels{AdditionalProperties: testSKSClusterLabels},
					Level:       &testSKSClusterServiceLevel,
					Name:        &testSKSClusterName,
					Nodepools: &[]oapi.SksNodepool{{
						AntiAffinityGroups: &[]oapi.AntiAffinityGroup{{Id: &testSKSNodepoolAntiAffinityGroupID}},
						CreatedAt:          &testSKSNodepoolCreatedAt,
						Description:        &testSKSNodepoolDescription,
						DiskSize:           &testSKSNodepoolDiskSize,
						Id:                 &testSKSNodepoolID,
						InstancePool:       &oapi.InstancePool{Id: &testSKSNodepoolInstancePoolID},
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
					}},
					State:   &testSKSClusterState,
					Version: &testSKSClusterVersion,
				}},
			},
		}, nil)

	expected := []*SKSCluster{{
		AddOns:      &[]string{string(testSKSClusterAddons[0])},
		AutoUpgrade: &testSKSClusterAutoUpgrade,
		CNI:         &testSKSClusterCNI,
		CreatedAt:   &testSKSClusterCreatedAt,
		Description: &testSKSClusterDescription,
		Endpoint:    &testSKSClusterEndpoint,
		ID:          &testSKSClusterID,
		Labels:      &testSKSClusterLabels,
		Name:        &testSKSClusterName,
		Nodepools: []*SKSNodepool{{
			AntiAffinityGroupIDs: &[]string{testSKSNodepoolAntiAffinityGroupID},
			CreatedAt:            &testSKSNodepoolCreatedAt,
			Description:          &testSKSNodepoolDescription,
			DiskSize:             &testSKSNodepoolDiskSize,
			ID:                   &testSKSNodepoolID,
			InstancePoolID:       &testSKSNodepoolInstancePoolID,
			InstanceTypeID:       &testSKSNodepoolInstanceTypeID,
			Labels:               &testSKSNodepoolLabels,
			Name:                 &testSKSNodepoolName,
			PrivateNetworkIDs:    &[]string{testSKSNodepoolPrivateNetworkID},
			SecurityGroupIDs:     &[]string{testSKSNodepoolSecurityGroupID},
			Size:                 &testSKSNodepoolSize,
			State:                (*string)(&testSKSClusterState),
			Taints: &map[string]*SKSNodepoolTaint{
				testSKSNodepoolTaintKey: {
					Effect: string(testSKSNodepoolTaintEffect),
					Value:  testSKSNodepoolTaintValue,
				},
			},
			TemplateID: &testSKSNodepoolTemplateID,
			Version:    &testSKSNodepoolVersion,
		}},
		ServiceLevel: (*string)(&testSKSClusterServiceLevel),
		State:        (*string)(&testSKSClusterState),
		Version:      &testSKSClusterVersion,
		Zone:         &testZone,
	}}

	actual, err := ts.client.ListSKSClusters(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_ListSKSClusterVersions() {
	var (
		testSKSClusterVersions = []string{
			"1.20.0",
			"1.18.6",
		}
		err error
	)

	ts.mock().
		On(
			"ListSksClusterVersionsWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // params
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(&oapi.ListSksClusterVersionsParams{
				IncludeDeprecated: func() *string { v := "true"; return &v }(),
			}, args.Get(1))
		}).
		Return(
			&oapi.ListSksClusterVersionsResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &struct {
					SksClusterVersions *[]string `json:"sks-cluster-versions,omitempty"`
				}{
					SksClusterVersions: &testSKSClusterVersions,
				},
			},
			nil,
		)

	expected := testSKSClusterVersions
	actual, err := ts.client.ListSKSClusterVersions(
		context.Background(),
		ListSKSClusterVersionsWithDeprecated(true),
	)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_RotateSKSClusterCCMCredentials() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		rotated            = false
	)

	cluster := &SKSCluster{
		ID: &testSKSClusterID,
	}

	ts.mock().
		On(
			"RotateSksCcmCredentialsWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testSKSClusterID, args.Get(1))
			rotated = true
		}).
		Return(
			&oapi.RotateSksCcmCredentialsResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: &oapi.Reference{Id: &testSKSClusterID},
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: &oapi.Reference{Id: &testSKSClusterID},
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.RotateSKSClusterCCMCredentials(context.Background(), testZone, cluster))
	ts.Require().True(rotated)
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
					Reference: &oapi.Reference{Id: &testSKSNodepoolID},
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: &oapi.Reference{Id: &testSKSNodepoolID},
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

func (ts *testSuite) TestClient_UpdateSKSCluster() {
	var (
		testSKSClusterAutoUpgradeUpdated = false
		testSKSClusterDescriptionUpdated = testSKSClusterDescription + "-updated"
		testSKSClusterLabelsUpdated      = map[string]string{"k3": "v3"}
		testSKSClusterNameUpdated        = testSKSClusterName + "-updated"
		testOperationID                  = ts.randomID()
		testOperationState               = oapi.OperationStateSuccess
		updated                          = false
	)

	ts.mock().
		On(
			"UpdateSksClusterWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				oapi.UpdateSksClusterJSONRequestBody{
					AutoUpgrade: &testSKSClusterAutoUpgradeUpdated,
					Labels:      &oapi.Labels{AdditionalProperties: testSKSClusterLabelsUpdated},
					Name:        &testSKSClusterNameUpdated,
					Description: &testSKSClusterDescriptionUpdated,
				},
				args.Get(2),
			)
			updated = true
		}).
		Return(
			&oapi.UpdateSksClusterResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: &oapi.Reference{Id: &testSKSClusterID},
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: &oapi.Reference{Id: &testSKSClusterID},
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.UpdateSKSCluster(context.Background(), testZone, &SKSCluster{
		AutoUpgrade: &testSKSClusterAutoUpgradeUpdated,
		ID:          &testSKSClusterID,
		Labels:      &testSKSClusterLabelsUpdated,
		Name:        &testSKSClusterNameUpdated,
		Description: &testSKSClusterDescriptionUpdated,
	}))
	ts.Require().True(updated)
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
					Reference: &oapi.Reference{Id: &testSKSNodepoolID},
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: &oapi.Reference{Id: &testSKSNodepoolID},
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

func (ts *testSuite) TestClient_UpgradeSKSCluster() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		upgraded           = false
	)

	ts.mock().
		On(
			"UpgradeSksClusterWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testSKSClusterID, args.Get(1))
			ts.Require().Equal(
				oapi.UpgradeSksClusterJSONRequestBody{Version: testSKSClusterVersion},
				args.Get(2),
			)
			upgraded = true
		}).
		Return(
			&oapi.UpgradeSksClusterResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: &oapi.Reference{Id: &testSKSClusterID},
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: &oapi.Reference{Id: &testSKSClusterID},
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.UpgradeSKSCluster(
		context.Background(),
		testZone,
		&SKSCluster{ID: &testSKSClusterID},
		testSKSClusterVersion))
	ts.Require().True(upgraded)
}

func (ts *testSuite) TestClient_UpgradeSKSClusterServiceLevel() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		upgraded           = false
	)

	ts.mock().
		On(
			"UpgradeSksClusterServiceLevelWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testSKSClusterID, args.Get(1))
			upgraded = true
		}).
		Return(
			&oapi.UpgradeSksClusterServiceLevelResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: &oapi.Reference{Id: &testSKSClusterID},
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: &oapi.Reference{Id: &testSKSClusterID},
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.UpgradeSKSClusterServiceLevel(
		context.Background(),
		testZone,
		&SKSCluster{ID: &testSKSClusterID}))
	ts.Require().True(upgraded)
}
