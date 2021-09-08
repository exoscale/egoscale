package v2

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/jarcoal/httpmock"

	"github.com/exoscale/egoscale/v2/oapi"
)

var (
	testSKSClusterAddons                     = []oapi.SksClusterAddons{oapi.SksClusterAddonsExoscaleCloudController}
	testSKSClusterAutoUpgrade                = true
	testSKSClusterCNI                        = "calico"
	testSKSClusterCreatedAt, _               = time.Parse(iso8601Format, "2020-11-16T10:41:58Z")
	testSKSClusterDescription                = new(clientTestSuite).randomString(10)
	testSKSClusterEndpoint                   = fmt.Sprintf("%s.sks-%s.exo.io", testSKSClusterID, testZone)
	testSKSClusterID                         = new(clientTestSuite).randomID()
	testSKSClusterLabels                     = map[string]string{"k1": "v1", "k2": "v2"}
	testSKSClusterName                       = new(clientTestSuite).randomString(10)
	testSKSClusterServiceLevel               = oapi.SksClusterLevelPro
	testSKSClusterState                      = oapi.SksClusterStateRunning
	testSKSClusterVersion                    = "1.18.6"
	testSKSNodepoolAddons                    = []oapi.SksNodepoolAddons{oapi.SksNodepoolAddonsLinbit}
	testSKSNodepoolAntiAffinityGroupID       = new(clientTestSuite).randomID()
	testSKSNodepoolCreatedAt, _              = time.Parse(iso8601Format, "2020-11-18T07:54:36Z")
	testSKSNodepoolDeployTargetID            = new(clientTestSuite).randomID()
	testSKSNodepoolDescription               = new(clientTestSuite).randomString(10)
	testSKSNodepoolDiskSize            int64 = 15
	testSKSNodepoolID                        = new(clientTestSuite).randomID()
	testSKSNodepoolInstancePoolID            = new(clientTestSuite).randomID()
	testSKSNodepoolInstancePrefix            = new(clientTestSuite).randomString(10)
	testSKSNodepoolInstanceTypeID            = new(clientTestSuite).randomID()
	testSKSNodepoolLabels                    = map[string]string{"k1": "v1", "k2": "v2"}
	testSKSNodepoolName                      = new(clientTestSuite).randomString(10)
	testSKSNodepoolPrivateNetworkID          = new(clientTestSuite).randomID()
	testSKSNodepoolSecurityGroupID           = new(clientTestSuite).randomID()
	testSKSNodepoolSize                int64 = 3
	testSKSNodepoolState                     = oapi.SksNodepoolStateRunning
	testSKSNodepoolTemplateID                = new(clientTestSuite).randomID()
	testSKSNodepoolVersion                   = "1.18.6"
)

func (ts *clientTestSuite) TestClient_CreateSKSCluster() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
	)

	httpmock.RegisterResponder("POST", "/sks-cluster",
		func(req *http.Request) (*http.Response, error) {
			var actual oapi.CreateSksClusterJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := oapi.CreateSksClusterJSONRequestBody{
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
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testSKSClusterID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testSKSClusterID},
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/sks-cluster/%s", testSKSClusterID), oapi.SksCluster{
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
	})

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

func (ts *clientTestSuite) TestCLient_CreateSKSNodepool() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
	)

	httpmock.RegisterResponder("POST", fmt.Sprintf("/sks-cluster/%s/nodepool", testSKSClusterID),
		func(req *http.Request) (*http.Response, error) {
			var actual oapi.CreateSksNodepoolJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := oapi.CreateSksNodepoolJSONRequestBody{
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
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testSKSNodepoolID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testSKSNodepoolID},
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/sks-cluster/%s/nodepool/%s",
		testSKSClusterID, testSKSNodepoolID),
		oapi.SksNodepool{
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
			Template:           &oapi.Template{Id: &testSKSNodepoolTemplateID},
			Version:            &testSKSNodepoolVersion,
		})

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
		TemplateID:           &testSKSNodepoolTemplateID,
		Version:              &testSKSNodepoolVersion,
	}

	actual, err := ts.client.CreateSKSNodepool(context.Background(), testZone, cluster, expected)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_DeleteSKSCluster() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		deleted            = false
	)

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("/sks-cluster/%s", testSKSClusterID),
		func(req *http.Request) (*http.Response, error) {
			deleted = true

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testSKSClusterID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testSKSClusterID},
	})

	ts.Require().NoError(ts.client.DeleteSKSCluster(
		context.Background(),
		testZone,
		&SKSCluster{ID: &testSKSClusterID},
	))
	ts.Require().True(deleted)
}

func (ts *clientTestSuite) TestClient_DeleteSKSNodepool() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
	)

	httpmock.RegisterResponder("DELETE",
		fmt.Sprintf("=~^/sks-cluster/%s/nodepool/.*", testSKSClusterID),
		func(req *http.Request) (*http.Response, error) {
			ts.Require().Equal(fmt.Sprintf("/sks-cluster/%s/nodepool/%s",
				testSKSClusterID, testSKSNodepoolID), req.URL.String())

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testSKSClusterID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testSKSNodepoolID},
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
}

func (ts *clientTestSuite) TestClient_EvictSKSNodepoolMembers() {
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

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/sks-cluster/%s/nodepool/%s:evict",
		*cluster.ID,
		*cluster.Nodepools[0].ID),
		func(req *http.Request) (*http.Response, error) {
			evicted = true

			var actual oapi.EvictSksNodepoolMembersJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := oapi.EvictSksNodepoolMembersJSONRequestBody{Instances: &[]string{testEvictedMemberID}}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testSKSNodepoolID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testSKSNodepoolID},
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

func (ts *clientTestSuite) TestClient_FindSKSCluster() {
	ts.mockAPIRequest("GET", "/sks-cluster", struct {
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
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/sks-cluster/%s", testSKSClusterID), oapi.SksCluster{
		CreatedAt: &testSKSClusterCreatedAt,
		Endpoint:  &testSKSClusterEndpoint,
		Id:        &testSKSClusterID,
		Level:     &testSKSClusterServiceLevel,
		Name:      &testSKSClusterName,
		State:     &testSKSClusterState,
		Version:   &testSKSClusterVersion,
	})

	expected := &SKSCluster{
		CreatedAt:    &testSKSClusterCreatedAt,
		Endpoint:     &testSKSClusterEndpoint,
		ID:           &testSKSClusterID,
		Name:         &testSKSClusterName,
		Nodepools:    []*SKSNodepool{},
		ServiceLevel: (*string)(&testSKSClusterServiceLevel),
		State:        (*string)(&testSKSClusterState),
		Version:      &testSKSClusterVersion,
	}

	actual, err := ts.client.FindSKSCluster(context.Background(), testZone, *expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)

	actual, err = ts.client.FindSKSCluster(context.Background(), testZone, *expected.Name)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_GetSKSCluster() {
	ts.mockAPIRequest("GET", fmt.Sprintf("/sks-cluster/%s", testSKSClusterID), oapi.SksCluster{
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
			Template:           &oapi.Template{Id: &testSKSNodepoolTemplateID},
			Version:            &testSKSNodepoolVersion,
		}},
		State:   &testSKSClusterState,
		Version: &testSKSClusterVersion,
	})

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
			TemplateID:           &testSKSNodepoolTemplateID,
			Version:              &testSKSNodepoolVersion,
		}},
		ServiceLevel: (*string)(&testSKSClusterServiceLevel),
		State:        (*string)(&testSKSClusterState),
		Version:      &testSKSClusterVersion,
	}

	actual, err := ts.client.GetSKSCluster(context.Background(), testZone, *expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_GetSKSClusterAuthorityCert() {
	var (
		testAuthority   = "aggregation"
		testCertificate = base64.StdEncoding.EncodeToString([]byte("test"))
	)

	cluster := &SKSCluster{
		ID: &testSKSClusterID,
	}

	ts.mockAPIRequest("GET",
		fmt.Sprintf("/sks-cluster/%s/authority/%s/cert", *cluster.ID, testAuthority),
		struct {
			Cacert string `json:"cacert,omitempty"`
		}{
			Cacert: testCertificate,
		})

	actual, err := ts.client.GetSKSClusterAuthorityCert(context.Background(), testZone, cluster, testAuthority)
	ts.Require().NoError(err)
	ts.Require().Equal(testCertificate, actual)
}

func (ts *clientTestSuite) TestClient_GetSKSClusterKubeconfig() {
	var (
		testRequestUser   = "test-user"
		testRequestGroups = []string{"system:masters"}
		testKubeconfig    = base64.StdEncoding.EncodeToString([]byte("test"))
	)

	cluster := &SKSCluster{
		ID: &testSKSClusterID,
	}

	ts.mockAPIRequest("POST", fmt.Sprintf("/sks-cluster-kubeconfig/%s", *cluster.ID), struct {
		Kubeconfig string `json:"kubeconfig,omitempty"`
	}{
		Kubeconfig: testKubeconfig,
	})

	actual, err := ts.client.GetSKSClusterKubeconfig(
		context.Background(),
		testZone,
		cluster,
		testRequestUser,
		testRequestGroups,
		time.Hour,
	)
	ts.Require().NoError(err)
	ts.Require().Equal(testKubeconfig, actual)
}

func (ts *clientTestSuite) TestClient_ListSKSClusters() {
	ts.mockAPIRequest("GET", "/sks-cluster", struct {
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
				Template:           &oapi.Template{Id: &testSKSNodepoolTemplateID},
				Version:            &testSKSNodepoolVersion,
			}},
			State:   &testSKSClusterState,
			Version: &testSKSClusterVersion,
		}},
	})

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
			TemplateID:           &testSKSNodepoolTemplateID,
			Version:              &testSKSNodepoolVersion,
		}},
		ServiceLevel: (*string)(&testSKSClusterServiceLevel),
		State:        (*string)(&testSKSClusterState),
		Version:      &testSKSClusterVersion,
	}}

	actual, err := ts.client.ListSKSClusters(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_ListSKSClusterVersions() {
	var (
		testSKSClusterVersions = []string{
			"1.20.0",
			"1.18.6",
		}
		err error
	)

	ts.mockAPIRequest("GET", "/sks-cluster-version", struct {
		SksClusterVersions *[]string `json:"sks-cluster-versions,omitempty"`
	}{
		SksClusterVersions: &testSKSClusterVersions,
	})

	expected := testSKSClusterVersions
	actual, err := ts.client.ListSKSClusterVersions(context.Background())
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_RotateSKSClusterCCMCredentials() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		rotated            = false
	)

	cluster := &SKSCluster{
		ID: &testSKSClusterID,
	}

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/sks-cluster/%s/rotate-ccm-credentials", *cluster.ID),
		func(req *http.Request) (*http.Response, error) {
			rotated = true

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testSKSClusterID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testSKSNodepoolID},
	})

	ts.Require().NoError(ts.client.RotateSKSClusterCCMCredentials(context.Background(), testZone, cluster))
	ts.Require().True(rotated)
}

func (ts *clientTestSuite) TestClient_ScaleSKSNodepool() {
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

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/sks-cluster/%s/nodepool/%s:scale",
		*cluster.ID,
		*cluster.Nodepools[0].ID),
		func(req *http.Request) (*http.Response, error) {
			scaled = true

			var actual oapi.ScaleSksNodepoolJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := oapi.ScaleSksNodepoolJSONRequestBody{Size: testScaleSize}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testSKSNodepoolID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testSKSNodepoolID},
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

func (ts *clientTestSuite) TestClient_UpdateSKSCluster() {
	var (
		testSKSClusterAutoUpgradeUpdated = false
		testSKSClusterDescriptionUpdated = testSKSClusterDescription + "-updated"
		testSKSClusterLabelsUpdated      = map[string]string{"k3": "v3"}
		testSKSClusterNameUpdated        = testSKSClusterName + "-updated"
		testOperationID                  = ts.randomID()
		testOperationState               = oapi.OperationStateSuccess
		updated                          = false
	)

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/sks-cluster/%s", testSKSClusterID),
		func(req *http.Request) (*http.Response, error) {
			updated = true

			var actual oapi.UpdateSksClusterJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := oapi.UpdateSksClusterJSONRequestBody{
				AutoUpgrade: &testSKSClusterAutoUpgradeUpdated,
				Labels:      &oapi.Labels{AdditionalProperties: testSKSClusterLabelsUpdated},
				Name:        &testSKSClusterNameUpdated,
				Description: &testSKSClusterDescriptionUpdated,
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testSKSClusterID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testSKSClusterID},
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

func (ts *clientTestSuite) TestClient_UpdateSKSNodepool() {
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

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/sks-cluster/%s/nodepool/%s",
		*cluster.ID,
		*cluster.Nodepools[0].ID),
		func(req *http.Request) (*http.Response, error) {
			updated = true

			var actual oapi.UpdateSksNodepoolJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := oapi.UpdateSksNodepoolJSONRequestBody{
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
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testSKSNodepoolID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testSKSNodepoolID},
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
	}))
	ts.Require().True(updated)
}

func (ts *clientTestSuite) TestClient_UgradeSKSCluster() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		upgraded           = false
	)

	httpmock.RegisterResponder("PUT",
		fmt.Sprintf("/sks-cluster/%s/upgrade", testSKSClusterID),
		func(req *http.Request) (*http.Response, error) {
			upgraded = true

			var actual oapi.UpgradeSksClusterJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := oapi.UpgradeSksClusterJSONRequestBody{Version: testSKSClusterVersion}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testSKSNodepoolID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testSKSNodepoolID},
	})

	ts.Require().NoError(ts.client.UpgradeSKSCluster(
		context.Background(),
		testZone,
		&SKSCluster{ID: &testSKSClusterID},
		testSKSClusterVersion))
	ts.Require().True(upgraded)
}

func (ts *clientTestSuite) TestClient_UgradeSKSClusterServiceLevel() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		upgraded           = false
	)

	httpmock.RegisterResponder("PUT",
		fmt.Sprintf("/sks-cluster/%s/upgrade-service-level", testSKSClusterID),
		func(req *http.Request) (*http.Response, error) {
			upgraded = true

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testSKSNodepoolID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testSKSNodepoolID},
	})

	ts.Require().NoError(ts.client.UpgradeSKSClusterServiceLevel(
		context.Background(),
		testZone,
		&SKSCluster{ID: &testSKSClusterID}))
	ts.Require().True(upgraded)
}
