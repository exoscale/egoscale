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
	testSKSClusterAddons             = []oapi.SksClusterAddons{oapi.SksClusterAddonsExoscaleCloudController}
	testSKSClusterAutoUpgrade        = true
	testSKSClusterCNI                = "calico"
	testSKSClusterCreatedAt, _       = time.Parse(iso8601Format, "2020-11-16T10:41:58Z")
	testSKSClusterDescription        = new(testSuite).randomString(10)
	testSKSClusterEndpoint           = fmt.Sprintf("%s.sks-%s.exo.io", testSKSClusterID, testZone)
	testSKSClusterID                 = new(testSuite).randomID()
	testSKSClusterLabels             = map[string]string{"k1": "v1", "k2": "v2"}
	testSKSClusterName               = new(testSuite).randomString(10)
	testSKSClusterOIDCClientID       = new(testSuite).randomString(10)
	testSKSClusterOIDCGroupsClaim    = new(testSuite).randomString(10)
	testSKSClusterOIDCGroupsPrefix   = new(testSuite).randomString(10)
	testSKSClusterOIDCIssuerURL      = new(testSuite).randomString(10)
	testSKSClusterOIDCRequiredClaim  = map[string]string{"test": new(testSuite).randomString(10)}
	testSKSClusterOIDCUsernameClaim  = new(testSuite).randomString(10)
	testSKSClusterOIDCUsernamePrefix = new(testSuite).randomString(10)
	testSKSClusterServiceLevel       = oapi.SksClusterLevelPro
	testSKSClusterState              = oapi.SksClusterStateRunning
	testSKSClusterVersion            = "1.18.6"
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
					Oidc: &oapi.SksOidc{
						ClientId:     testSKSClusterOIDCClientID,
						GroupsClaim:  &testSKSClusterOIDCGroupsClaim,
						GroupsPrefix: &testSKSClusterOIDCGroupsPrefix,
						IssuerUrl:    testSKSClusterOIDCIssuerURL,
						RequiredClaim: &oapi.SksOidc_RequiredClaim{
							AdditionalProperties: testSKSClusterOIDCRequiredClaim,
						},
						UsernameClaim:  &testSKSClusterOIDCUsernameClaim,
						UsernamePrefix: &testSKSClusterOIDCUsernamePrefix,
					},
					Version: testSKSClusterVersion,
				},
				args.Get(1),
			)
		}).
		Return(
			&oapi.CreateSksClusterResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testSKSClusterID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testSKSClusterID, nil),
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
	},
		CreateSKSClusterWithOIDC(&SKSClusterOIDCConfig{
			ClientID:       &testSKSClusterOIDCClientID,
			GroupsClaim:    &testSKSClusterOIDCGroupsClaim,
			GroupsPrefix:   &testSKSClusterOIDCGroupsPrefix,
			IssuerURL:      &testSKSClusterOIDCIssuerURL,
			RequiredClaim:  &testSKSClusterOIDCRequiredClaim,
			UsernameClaim:  &testSKSClusterOIDCUsernameClaim,
			UsernamePrefix: &testSKSClusterOIDCUsernamePrefix,
		}))
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
					Reference: oapi.NewReference(nil, &testSKSClusterID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testSKSClusterID, nil),
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.DeleteSKSCluster(
		context.Background(),
		testZone,
		&SKSCluster{ID: &testSKSClusterID},
	))
	ts.Require().True(deleted)
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
					KubeletImageGc:     &testSKSNodepoolKubeletImageGc,
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
			KubeletImageGc:       sksNodepoolKubeletImageGcFromAPI(&testSKSNodepoolKubeletImageGc),
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
						KubeletImageGc:     &testSKSNodepoolKubeletImageGc,
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
			KubeletImageGc:       sksNodepoolKubeletImageGcFromAPI(&testSKSNodepoolKubeletImageGc),
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
					Reference: oapi.NewReference(nil, &testSKSClusterID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testSKSClusterID, nil),
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.RotateSKSClusterCCMCredentials(context.Background(), testZone, cluster))
	ts.Require().True(rotated)
}

func (ts *testSuite) TestClient_RotateSKSClusterCSICredentials() {
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
			"RotateSksCsiCredentialsWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testSKSClusterID, args.Get(1))
			rotated = true
		}).
		Return(
			&oapi.RotateSksCsiCredentialsResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testSKSClusterID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testSKSClusterID, nil),
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.RotateSKSClusterCSICredentials(context.Background(), testZone, cluster))
	ts.Require().True(rotated)
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
					Reference: oapi.NewReference(nil, &testSKSClusterID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testSKSClusterID, nil),
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
					Reference: oapi.NewReference(nil, &testSKSClusterID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testSKSClusterID, nil),
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
					Reference: oapi.NewReference(nil, &testSKSClusterID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testSKSClusterID, nil),
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.UpgradeSKSClusterServiceLevel(
		context.Background(),
		testZone,
		&SKSCluster{ID: &testSKSClusterID}))
	ts.Require().True(upgraded)
}
