package v2

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/jarcoal/httpmock"

	papi "github.com/exoscale/egoscale/v2/internal/public-api"
)

var (
	testSKSClusterAddons                     = []papi.SksClusterAddons{papi.SksClusterAddonsExoscaleCloudController}
	testSKSClusterAutoUpgrade                = true
	testSKSClusterCNI                        = "calico"
	testSKSClusterCreatedAt, _               = time.Parse(iso8601Format, "2020-11-16T10:41:58Z")
	testSKSClusterDescription                = new(clientTestSuite).randomString(10)
	testSKSClusterEndpoint                   = fmt.Sprintf("%s.sks-%s.exo.io", testSKSClusterID, testZone)
	testSKSClusterID                         = new(clientTestSuite).randomID()
	testSKSClusterLabels                     = map[string]string{"k1": "v1", "k2": "v2"}
	testSKSClusterName                       = new(clientTestSuite).randomString(10)
	testSKSClusterServiceLevel               = papi.SksClusterLevelPro
	testSKSClusterState                      = papi.SksClusterStateRunning
	testSKSClusterVersion                    = "1.18.6"
	testSKSNodepoolAddons                    = []papi.SksNodepoolAddons{papi.SksNodepoolAddonsLinbit}
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
	testSKSNodepoolState                     = papi.SksNodepoolStateRunning
	testSKSNodepoolTemplateID                = new(clientTestSuite).randomID()
	testSKSNodepoolVersion                   = "1.18.6"
)

func (ts *clientTestSuite) TestSKSNodepool_AntiAffinityGroups() {
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

	sksNodepool := &SKSNodepool{
		AntiAffinityGroupIDs: &[]string{testAntiAffinityGroupID},

		c:    ts.client,
		zone: testZone,
	}

	actual, err := sksNodepool.AntiAffinityGroups(context.Background())
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestSKSNodepool_PrivateNetworks() {
	ts.mockAPIRequest("GET", fmt.Sprintf("/private-network/%s", testPrivateNetworkID),
		papi.PrivateNetwork{
			Id:   &testPrivateNetworkID,
			Name: &testPrivateNetworkName,
		})

	expected := []*PrivateNetwork{
		{
			ID:   &testPrivateNetworkID,
			Name: &testPrivateNetworkName,

			c:    ts.client,
			zone: testZone,
		},
	}

	sksNodepool := &SKSNodepool{
		PrivateNetworkIDs: &[]string{testPrivateNetworkID},

		c:    ts.client,
		zone: testZone,
	}

	actual, err := sksNodepool.PrivateNetworks(context.Background())
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestSKSNodepool_SecurityGroups() {
	ts.mockAPIRequest("GET", fmt.Sprintf("/security-group/%s", testSecurityGroupID),
		papi.SecurityGroup{
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

	sksNodepool := &SKSNodepool{
		SecurityGroupIDs: &[]string{testSecurityGroupID},

		c:    ts.client,
		zone: testZone,
	}

	actual, err := sksNodepool.SecurityGroups(context.Background())
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestSKSCluster_RotateCCMCredentials() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = papi.OperationStateSuccess
		rotated            = false
	)

	cluster := &SKSCluster{
		ID:   &testSKSClusterID,
		c:    ts.client,
		zone: testZone,
	}

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/sks-cluster/%s/rotate-ccm-credentials", *cluster.ID),
		func(req *http.Request) (*http.Response, error) {
			rotated = true

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testSKSClusterID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testSKSNodepoolID},
	})

	ts.Require().NoError(cluster.RotateCCMCredentials(context.Background()))
	ts.Require().True(rotated)
}

func (ts *clientTestSuite) TestSKSCluster_AuthorityCert() {
	var (
		testAuthority   = "aggregation"
		testCertificate = base64.StdEncoding.EncodeToString([]byte("test"))
	)

	cluster := &SKSCluster{
		ID:   &testSKSClusterID,
		c:    ts.client,
		zone: testZone,
	}

	ts.mockAPIRequest("GET",
		fmt.Sprintf("/sks-cluster/%s/authority/%s/cert", *cluster.ID, testAuthority),
		struct {
			Cacert string `json:"cacert,omitempty"`
		}{
			Cacert: testCertificate,
		})

	actual, err := cluster.AuthorityCert(context.Background(), testAuthority)
	ts.Require().NoError(err)
	ts.Require().Equal(testCertificate, actual)
}

func (ts *clientTestSuite) TestSKSCluster_RequestKubeconfig() {
	var (
		testRequestUser   = "test-user"
		testRequestGroups = []string{"system:masters"}
		testKubeconfig    = base64.StdEncoding.EncodeToString([]byte("test"))
	)

	cluster := &SKSCluster{
		ID:   &testSKSClusterID,
		c:    ts.client,
		zone: testZone,
	}

	ts.mockAPIRequest("POST", fmt.Sprintf("/sks-cluster-kubeconfig/%s", *cluster.ID), struct {
		Kubeconfig string `json:"kubeconfig,omitempty"`
	}{
		Kubeconfig: testKubeconfig,
	})

	actual, err := cluster.RequestKubeconfig(context.Background(), testRequestUser, testRequestGroups, time.Hour)
	ts.Require().NoError(err)
	ts.Require().Equal(testKubeconfig, actual)
}

func (ts *clientTestSuite) TestSKSCluster_AddNodepool() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = papi.OperationStateSuccess
	)

	httpmock.RegisterResponder("POST", fmt.Sprintf("/sks-cluster/%s/nodepool", testSKSClusterID),
		func(req *http.Request) (*http.Response, error) {
			var actual papi.CreateSksNodepoolJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.CreateSksNodepoolJSONRequestBody{
				Addons: &[]papi.CreateSksNodepoolJSONBodyAddons{
					papi.CreateSksNodepoolJSONBodyAddons(testSKSNodepoolAddons[0]),
				},
				AntiAffinityGroups: &[]papi.AntiAffinityGroup{{Id: &testSKSNodepoolAntiAffinityGroupID}},
				DeployTarget:       &papi.DeployTarget{Id: &testSKSNodepoolDeployTargetID},
				Description:        &testSKSNodepoolDescription,
				DiskSize:           testSKSNodepoolDiskSize,
				InstancePrefix:     &testSKSNodepoolInstancePrefix,
				InstanceType:       papi.InstanceType{Id: &testSKSNodepoolInstanceTypeID},
				Labels:             &papi.Labels{AdditionalProperties: testSKSNodepoolLabels},
				Name:               testSKSNodepoolName,
				PrivateNetworks:    &[]papi.PrivateNetwork{{Id: &testSKSNodepoolPrivateNetworkID}},
				SecurityGroups:     &[]papi.SecurityGroup{{Id: &testSKSNodepoolSecurityGroupID}},
				Size:               testSKSNodepoolSize,
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testSKSNodepoolID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testSKSNodepoolID},
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/sks-cluster/%s/nodepool/%s",
		testSKSClusterID, testSKSNodepoolID),
		papi.SksNodepool{
			Addons:             &testSKSNodepoolAddons,
			AntiAffinityGroups: &[]papi.AntiAffinityGroup{{Id: &testSKSNodepoolAntiAffinityGroupID}},
			CreatedAt:          &testSKSNodepoolCreatedAt,
			DeployTarget:       &papi.DeployTarget{Id: &testSKSNodepoolDeployTargetID},
			Description:        &testSKSNodepoolDescription,
			DiskSize:           &testSKSNodepoolDiskSize,
			Id:                 &testSKSNodepoolID,
			InstancePool:       &papi.InstancePool{Id: &testSKSNodepoolInstancePoolID},
			InstancePrefix:     &testSKSNodepoolInstancePrefix,
			InstanceType:       &papi.InstanceType{Id: &testSKSNodepoolInstanceTypeID},
			Labels:             &papi.Labels{AdditionalProperties: testSKSNodepoolLabels},
			Name:               &testSKSNodepoolName,
			PrivateNetworks:    &[]papi.PrivateNetwork{{Id: &testSKSNodepoolPrivateNetworkID}},
			SecurityGroups:     &[]papi.SecurityGroup{{Id: &testSKSNodepoolSecurityGroupID}},
			Size:               &testSKSNodepoolSize,
			State:              &testSKSNodepoolState,
			Template:           &papi.Template{Id: &testSKSNodepoolTemplateID},
			Version:            &testSKSNodepoolVersion,
		})

	cluster := &SKSCluster{
		ID: &testSKSClusterID,

		c:    ts.client,
		zone: testZone,
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

		c:    ts.client,
		zone: testZone,
	}

	actual, err := cluster.AddNodepool(context.Background(), expected)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestSKSCluster_UpdateNodepool() {
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
		testOperationState                        = papi.OperationStateSuccess
		updated                                   = false
	)

	cluster := &SKSCluster{
		ID:   &testSKSClusterID,
		c:    ts.client,
		zone: testZone,

		Nodepools: []*SKSNodepool{
			{
				DeployTargetID: &testSKSNodepoolDeployTargetID,
				Description:    &testSKSNodepoolDescription,
				DiskSize:       &testSKSNodepoolDiskSize,
				ID:             &testSKSNodepoolID,
				InstancePrefix: &testSKSNodepoolInstancePrefix,
				InstanceTypeID: &testSKSNodepoolInstanceTypeID,
				Labels:         &testSKSNodepoolLabels,
				Name:           &testSKSNodepoolName,
			},
		},
	}

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/sks-cluster/%s/nodepool/%s",
		*cluster.ID,
		*cluster.Nodepools[0].ID),
		func(req *http.Request) (*http.Response, error) {
			updated = true

			var actual papi.UpdateSksNodepoolJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.UpdateSksNodepoolJSONRequestBody{
				AntiAffinityGroups: &[]papi.AntiAffinityGroup{{Id: &testSKSNodepoolAntiAffinityGroupIDUpdated}},
				DeployTarget:       &papi.DeployTarget{Id: &testSKSNodepoolDeployTargetIDUpdated},
				Description:        &testSKSNodepoolDescriptionUpdated,
				DiskSize:           &testSKSNodepoolDiskSizeUpdated,
				InstancePrefix:     &testSKSNodepoolInstancePrefixUpdated,
				InstanceType:       &papi.InstanceType{Id: &testSKSNodepoolInstanceTypeIDUpdated},
				Labels:             &papi.Labels{AdditionalProperties: testSKSNodepoolLabelsUpdated},
				Name:               &testSKSNodepoolNameUpdated,
				PrivateNetworks:    &[]papi.PrivateNetwork{{Id: &testSKSNodepoolPrivateNetworkIDUpdated}},
				SecurityGroups:     &[]papi.SecurityGroup{{Id: &testSKSNodepoolSecurityGroupIDUpdated}},
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testSKSNodepoolID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testSKSNodepoolID},
	})

	nodepoolUpdated := SKSNodepool{
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
	}
	ts.Require().NoError(cluster.UpdateNodepool(context.Background(), &nodepoolUpdated))
	ts.Require().True(updated)
}

func (ts *clientTestSuite) TestSKSCluster_ScaleNodepool() {
	var (
		testOperationID          = ts.randomID()
		testOperationState       = papi.OperationStateSuccess
		testScaleSize      int64 = 3
		scaled                   = false
	)

	cluster := &SKSCluster{
		ID:   &testSKSClusterID,
		c:    ts.client,
		zone: testZone,

		Nodepools: []*SKSNodepool{{ID: &testSKSNodepoolID}},
	}

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/sks-cluster/%s/nodepool/%s:scale",
		*cluster.ID,
		*cluster.Nodepools[0].ID),
		func(req *http.Request) (*http.Response, error) {
			scaled = true

			var actual papi.ScaleSksNodepoolJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.ScaleSksNodepoolJSONRequestBody{Size: testScaleSize}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testSKSNodepoolID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testSKSNodepoolID},
	})

	ts.Require().NoError(cluster.ScaleNodepool(context.Background(), cluster.Nodepools[0], testScaleSize))
	ts.Require().True(scaled)
}

func (ts *clientTestSuite) TestSKSCluster_EvictNodepoolMembers() {
	var (
		testOperationID     = ts.randomID()
		testOperationState  = papi.OperationStateSuccess
		testEvictedMemberID = ts.randomID()
		evicted             = false
	)

	cluster := &SKSCluster{
		ID:   &testSKSClusterID,
		c:    ts.client,
		zone: testZone,

		Nodepools: []*SKSNodepool{{ID: &testSKSNodepoolID}},
	}

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/sks-cluster/%s/nodepool/%s:evict",
		*cluster.ID,
		*cluster.Nodepools[0].ID),
		func(req *http.Request) (*http.Response, error) {
			evicted = true

			var actual papi.EvictSksNodepoolMembersJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.EvictSksNodepoolMembersJSONRequestBody{Instances: &[]string{testEvictedMemberID}}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testSKSNodepoolID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testSKSNodepoolID},
	})

	ts.Require().NoError(cluster.EvictNodepoolMembers(
		context.Background(),
		cluster.Nodepools[0],
		[]string{testEvictedMemberID}))
	ts.Require().True(evicted)
}

func (ts *clientTestSuite) TestSKSCluster_DeleteNodepool() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = papi.OperationStateSuccess
	)

	httpmock.RegisterResponder("DELETE",
		fmt.Sprintf("=~^/sks-cluster/%s/nodepool/.*", testSKSClusterID),
		func(req *http.Request) (*http.Response, error) {
			ts.Require().Equal(fmt.Sprintf("/sks-cluster/%s/nodepool/%s",
				testSKSClusterID, testSKSNodepoolID), req.URL.String())

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testSKSClusterID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testSKSNodepoolID},
	})

	cluster := &SKSCluster{
		ID:   &testSKSClusterID,
		c:    ts.client,
		zone: testZone,

		Nodepools: []*SKSNodepool{{ID: &testSKSNodepoolID}},
	}

	ts.Require().NoError(cluster.DeleteNodepool(context.Background(), cluster.Nodepools[0]))
}

func (ts *clientTestSuite) TestClient_CreateSKSCluster() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = papi.OperationStateSuccess
	)

	httpmock.RegisterResponder("POST", "/sks-cluster",
		func(req *http.Request) (*http.Response, error) {
			var actual papi.CreateSksClusterJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.CreateSksClusterJSONRequestBody{
				Addons: &[]papi.CreateSksClusterJSONBodyAddons{
					papi.CreateSksClusterJSONBodyAddons(testSKSClusterAddons[0]),
				},
				AutoUpgrade: &testSKSClusterAutoUpgrade,
				Cni:         (*papi.CreateSksClusterJSONBodyCni)(&testSKSClusterCNI),
				Description: &testSKSClusterDescription,
				Labels:      &papi.Labels{AdditionalProperties: testSKSClusterLabels},
				Level:       papi.CreateSksClusterJSONBodyLevel(testSKSClusterServiceLevel),
				Name:        testSKSClusterName,
				Version:     testSKSClusterVersion,
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testSKSClusterID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testSKSClusterID},
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/sks-cluster/%s", testSKSClusterID), papi.SksCluster{
		Addons:      &testSKSClusterAddons,
		AutoUpgrade: &testSKSClusterAutoUpgrade,
		Cni:         (*papi.SksClusterCni)(&testSKSClusterCNI),
		CreatedAt:   &testSKSClusterCreatedAt,
		Description: &testSKSClusterDescription,
		Endpoint:    &testSKSClusterEndpoint,
		Id:          &testSKSClusterID,
		Labels:      &papi.Labels{AdditionalProperties: testSKSClusterLabels},
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

		c:    ts.client,
		zone: testZone,
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

func (ts *clientTestSuite) TestClient_ListSKSClusters() {
	ts.mockAPIRequest("GET", "/sks-cluster", struct {
		SksClusters *[]papi.SksCluster `json:"sks-clusters,omitempty"`
	}{
		SksClusters: &[]papi.SksCluster{{
			Addons:      &testSKSClusterAddons,
			AutoUpgrade: &testSKSClusterAutoUpgrade,
			Cni:         (*papi.SksClusterCni)(&testSKSClusterCNI),
			CreatedAt:   &testSKSClusterCreatedAt,
			Description: &testSKSClusterDescription,
			Endpoint:    &testSKSClusterEndpoint,
			Id:          &testSKSClusterID,
			Labels:      &papi.Labels{AdditionalProperties: testSKSClusterLabels},
			Level:       &testSKSClusterServiceLevel,
			Name:        &testSKSClusterName,
			Nodepools: &[]papi.SksNodepool{{
				AntiAffinityGroups: &[]papi.AntiAffinityGroup{{Id: &testSKSNodepoolAntiAffinityGroupID}},
				CreatedAt:          &testSKSNodepoolCreatedAt,
				Description:        &testSKSNodepoolDescription,
				DiskSize:           &testSKSNodepoolDiskSize,
				Id:                 &testSKSNodepoolID,
				InstancePool:       &papi.InstancePool{Id: &testSKSNodepoolInstancePoolID},
				InstanceType:       &papi.InstanceType{Id: &testSKSNodepoolInstanceTypeID},
				Labels:             &papi.Labels{AdditionalProperties: testSKSNodepoolLabels},
				Name:               &testSKSNodepoolName,
				PrivateNetworks:    &[]papi.PrivateNetwork{{Id: &testSKSNodepoolPrivateNetworkID}},
				SecurityGroups:     &[]papi.SecurityGroup{{Id: &testSKSNodepoolSecurityGroupID}},
				Size:               &testSKSNodepoolSize,
				State:              &testSKSNodepoolState,
				Template:           &papi.Template{Id: &testSKSNodepoolTemplateID},
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

			c:    ts.client,
			zone: testZone,
		}},
		ServiceLevel: (*string)(&testSKSClusterServiceLevel),
		State:        (*string)(&testSKSClusterState),
		Version:      &testSKSClusterVersion,

		c:    ts.client,
		zone: testZone,
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

func (ts *clientTestSuite) TestClient_GetSKSCluster() {
	ts.mockAPIRequest("GET", fmt.Sprintf("/sks-cluster/%s", testSKSClusterID), papi.SksCluster{
		Addons:      &testSKSClusterAddons,
		AutoUpgrade: &testSKSClusterAutoUpgrade,
		Cni:         (*papi.SksClusterCni)(&testSKSClusterCNI),
		CreatedAt:   &testSKSClusterCreatedAt,
		Description: &testSKSClusterDescription,
		Endpoint:    &testSKSClusterEndpoint,
		Id:          &testSKSClusterID,
		Labels:      &papi.Labels{AdditionalProperties: testSKSClusterLabels},
		Level:       &testSKSClusterServiceLevel,
		Name:        &testSKSClusterName,
		Nodepools: &[]papi.SksNodepool{{
			AntiAffinityGroups: &[]papi.AntiAffinityGroup{{Id: &testSKSNodepoolAntiAffinityGroupID}},
			CreatedAt:          &testSKSNodepoolCreatedAt,
			Description:        &testSKSNodepoolDescription,
			DiskSize:           &testSKSNodepoolDiskSize,
			Id:                 &testSKSNodepoolID,
			InstancePool:       &papi.InstancePool{Id: &testSKSNodepoolInstancePoolID},
			InstanceType:       &papi.InstanceType{Id: &testSKSNodepoolInstanceTypeID},
			Labels:             &papi.Labels{AdditionalProperties: testSKSNodepoolLabels},
			Name:               &testSKSNodepoolName,
			PrivateNetworks:    &[]papi.PrivateNetwork{{Id: &testSKSNodepoolPrivateNetworkID}},
			SecurityGroups:     &[]papi.SecurityGroup{{Id: &testSKSNodepoolSecurityGroupID}},
			Size:               &testSKSNodepoolSize,
			State:              &testSKSNodepoolState,
			Template:           &papi.Template{Id: &testSKSNodepoolTemplateID},
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

			c:    ts.client,
			zone: testZone,
		}},
		ServiceLevel: (*string)(&testSKSClusterServiceLevel),
		State:        (*string)(&testSKSClusterState),
		Version:      &testSKSClusterVersion,

		c:    ts.client,
		zone: testZone,
	}

	actual, err := ts.client.GetSKSCluster(context.Background(), testZone, *expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_FindSKSCluster() {
	ts.mockAPIRequest("GET", "/sks-cluster", struct {
		SksClusters *[]papi.SksCluster `json:"sks-clusters,omitempty"`
	}{
		SksClusters: &[]papi.SksCluster{{
			CreatedAt: &testSKSClusterCreatedAt,
			Endpoint:  &testSKSClusterEndpoint,
			Id:        &testSKSClusterID,
			Level:     &testSKSClusterServiceLevel,
			Name:      &testSKSClusterName,
			State:     &testSKSClusterState,
			Version:   &testSKSClusterVersion,
		}},
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/sks-cluster/%s", testSKSClusterID), papi.SksCluster{
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

		c:    ts.client,
		zone: testZone,
	}

	actual, err := ts.client.FindSKSCluster(context.Background(), testZone, *expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)

	actual, err = ts.client.FindSKSCluster(context.Background(), testZone, *expected.Name)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_UpdateSKSCluster() {
	var (
		testSKSClusterAutoUpgradeUpdated = false
		testSKSClusterDescriptionUpdated = testSKSClusterDescription + "-updated"
		testSKSClusterLabelsUpdated      = map[string]string{"k3": "v3"}
		testSKSClusterNameUpdated        = testSKSClusterName + "-updated"
		testOperationID                  = ts.randomID()
		testOperationState               = papi.OperationStateSuccess
		updated                          = false
	)

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/sks-cluster/%s", testSKSClusterID),
		func(req *http.Request) (*http.Response, error) {
			updated = true

			var actual papi.UpdateSksClusterJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.UpdateSksClusterJSONRequestBody{
				AutoUpgrade: &testSKSClusterAutoUpgradeUpdated,
				Labels:      &papi.Labels{AdditionalProperties: testSKSClusterLabelsUpdated},
				Name:        &testSKSClusterNameUpdated,
				Description: &testSKSClusterDescriptionUpdated,
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testSKSClusterID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testSKSClusterID},
	})

	clusterUpdated := SKSCluster{
		AutoUpgrade: &testSKSClusterAutoUpgradeUpdated,
		ID:          &testSKSClusterID,
		Labels:      &testSKSClusterLabelsUpdated,
		Name:        &testSKSClusterNameUpdated,
		Description: &testSKSClusterDescriptionUpdated,
	}
	ts.Require().NoError(ts.client.UpdateSKSCluster(context.Background(), testZone, &clusterUpdated))
	ts.Require().True(updated)
}

func (ts *clientTestSuite) TestClient_UgradeSKSCluster() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = papi.OperationStateSuccess
		upgraded           = false
	)

	httpmock.RegisterResponder("PUT",
		fmt.Sprintf("/sks-cluster/%s/upgrade", testSKSClusterID),
		func(req *http.Request) (*http.Response, error) {
			upgraded = true

			var actual papi.UpgradeSksClusterJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.UpgradeSksClusterJSONRequestBody{Version: testSKSClusterVersion}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testSKSNodepoolID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testSKSNodepoolID},
	})

	ts.Require().NoError(ts.client.UpgradeSKSCluster(
		context.Background(),
		testZone,
		testSKSClusterID,
		testSKSClusterVersion))
	ts.Require().True(upgraded)
}

func (ts *clientTestSuite) TestClient_DeleteSKSCluster() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = papi.OperationStateSuccess
		deleted            = false
	)

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("/sks-cluster/%s", testSKSClusterID),
		func(req *http.Request) (*http.Response, error) {
			deleted = true

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testSKSClusterID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testSKSClusterID},
	})

	ts.Require().NoError(ts.client.DeleteSKSCluster(context.Background(), testZone, testSKSClusterID))
	ts.Require().True(deleted)
}
