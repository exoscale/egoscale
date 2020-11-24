package egoscale

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	v2 "github.com/exoscale/egoscale/internal/v2"
)

var (
	testSKSClusterCreatedAt, _          = time.Parse(iso8601Format, "2020-11-16T10:41:58Z")
	testSKSClusterDescription           = "Test Cluster description"
	testSKSClusterEndpoint              = "df421958-3679-4e9c-afb9-02fb6f331301.sks-ch-gva-2.exo.io"
	testSKSClusterID                    = "df421958-3679-4e9c-afb9-02fb6f331301"
	testSKSClusterName                  = "test-cluster"
	testSKSClusterState                 = "running"
	testSKSClusterVersion               = "1.18.6"
	testSKSNodepoolCreatedAt, _         = time.Parse(iso8601Format, "2020-11-18T07:54:36Z")
	testSKSNodepoolDescription          = "Test Nodepool description"
	testSKSNodepoolDiskSize       int64 = 15
	testSKSNodepoolID                   = "6d1eecee-397c-4e16-b103-2d1353bf4ecc"
	testSKSNodepoolInstancePoolID       = "f1f67118-43b6-4632-a709-d55fada62f21"
	testSKSNodepoolInstanceTypeID       = "21624abb-764e-4def-81d7-9fc54b5957fb"
	testSKSNodepoolName                 = "test-nodepool"
	testSKSNodepoolSize           int64 = 3
	// testSKSNodepoolSecurityGroupID       = "efb4f4df-87ce-44e9-b5ee-59a9c1628edf" // TODO: the API doesn't return this field ATM
	testSKSNodepoolState      = "running"
	testSKSNodepoolTemplateID = "f270d9a2-db64-4e8e-9cd3-5125887e91aa"
	testSKSNodepoolVersion    = "1.18.6"
)

func TestSKSCluster_RequestKubeconfig(t *testing.T) {
	var (
		testRequestUser  = "test-user"
		testRequestRoles = []string{"system:masters"}
		testKubeconfig   = base64.StdEncoding.EncodeToString([]byte("test"))
		err              error
	)

	mockClient := v2.NewMockClient()
	client := NewClient("x", "x", "x")
	client.v2, err = v2.NewClientWithResponses("", v2.WithHTTPClient(mockClient))
	require.NoError(t, err)

	cluster := &SKSCluster{
		ID:   testSKSClusterID,
		c:    client,
		zone: testZone,
	}

	mockClient.RegisterResponder("POST", "/sks-cluster-kubeconfig/"+cluster.ID,
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(http.StatusOK, struct {
				Kubeconfig string `json:"kubeconfig,omitempty"`
			}{
				Kubeconfig: testKubeconfig,
			})
			if err != nil {
				t.Fatalf("error initializing mock HTTP responder: %s", err)
			}
			return resp, nil
		})

	actual, err := cluster.RequestKubeconfig(testRequestUser, testRequestRoles, time.Hour)
	require.NoError(t, err)
	require.Equal(t, testKubeconfig, actual)
}

func TestSKSCluster_AddNodepool(t *testing.T) {
	var (
		testOperationID    = "08302193-c7e3-42a6-9b3d-da0b2a536577"
		testOperationState = "success"
		err                error
	)

	mockClient := v2.NewMockClient()
	client := NewClient("x", "x", "x")
	client.v2, err = v2.NewClientWithResponses("", v2.WithHTTPClient(mockClient))
	require.NoError(t, err)

	mockClient.RegisterResponder("POST", "/sks-cluster/"+testSKSClusterID+"/nodepool",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(http.StatusOK, v2.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &v2.Reference{Id: &testSKSNodepoolID},
			})
			if err != nil {
				t.Fatalf("error initializing mock HTTP responder: %s", err)
			}
			return resp, nil
		})

	mockClient.RegisterResponder("GET", "/operation/"+testOperationID,
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(http.StatusOK, v2.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &v2.Reference{Id: &testSKSNodepoolID},
			})
			if err != nil {
				t.Fatalf("error initializing mock HTTP responder: %s", err)
			}
			return resp, nil
		})

	mockClient.RegisterResponder("GET", fmt.Sprintf("/sks-cluster/%s/nodepool/%s",
		testSKSClusterID, testSKSNodepoolID),
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(http.StatusOK, v2.SksNodepool{
				CreatedAt:    &testSKSNodepoolCreatedAt,
				Description:  &testSKSNodepoolDescription,
				DiskSize:     &testSKSNodepoolDiskSize,
				Id:           &testSKSNodepoolID,
				InstancePool: &v2.Resource{Id: &testSKSNodepoolInstancePoolID},
				InstanceType: &v2.InstanceType{Id: &testSKSNodepoolInstanceTypeID},
				Name:         &testSKSNodepoolName,
				// SecurityGroups: nil, // TODO: the API doesn't return this field ATM
				Size:     &testSKSNodepoolSize,
				State:    &testSKSNodepoolState,
				Template: &v2.Template{Id: &testSKSNodepoolTemplateID},
				Version:  &testSKSNodepoolVersion,
			})
			if err != nil {
				t.Fatalf("error initializing mock HTTP responder: %s", err)
			}
			return resp, nil
		})

	cluster := &SKSCluster{
		ID: testSKSClusterID,

		c:    client,
		zone: testZone,
	}

	expected := &SKSNodepool{
		ID:             testSKSNodepoolID,
		Name:           testSKSNodepoolName,
		Description:    testSKSNodepoolDescription,
		CreatedAt:      testSKSNodepoolCreatedAt,
		InstancePoolID: testSKSNodepoolInstancePoolID,
		InstanceTypeID: testSKSNodepoolInstanceTypeID,
		TemplateID:     testSKSNodepoolTemplateID,
		DiskSize:       testSKSNodepoolDiskSize,
		// SecurityGroupIDs: nil, // TODO: the API doesn't return this field ATM
		Version: testSKSNodepoolVersion,
		Size:    testSKSNodepoolSize,
		State:   testSKSNodepoolState,
	}

	actual, err := cluster.AddNodepool(context.Background(), expected)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

// UpdateNodepool is not tested as it only produces API-side effects.
func TestSKSCluster_UpdateNodepool(t *testing.T) { t.Skip() }

// ScaleNodepool is not tested as it only produces API-side effects.
func TestSKSCluster_ScaleNodepool(t *testing.T) { t.Skip() }

// DeleteNodepool is not tested as it only produces API-side effects.
func TestSKSCluster_DeleteNodepool(t *testing.T) { t.Skip() }

// CreateSKSCluster is not tested as it essentially relies on the already tested GetSKSCluster.
func TestClient_CreateSKSCluster(t *testing.T) { t.Skip() }

func TestClient_ListSKSClusters(t *testing.T) {
	var err error

	mockClient := v2.NewMockClient()
	client := NewClient("x", "x", "x")
	client.v2, err = v2.NewClientWithResponses("", v2.WithHTTPClient(mockClient))
	require.NoError(t, err)

	mockClient.RegisterResponder("GET", "/sks-cluster",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(http.StatusOK, struct {
				SksClusters *[]v2.SksCluster `json:"sks-clusters,omitempty"`
			}{
				SksClusters: &[]v2.SksCluster{{
					CreatedAt:   &testSKSClusterCreatedAt,
					Description: &testSKSClusterDescription,
					Endpoint:    &testSKSClusterEndpoint,
					Id:          &testSKSClusterID,
					Name:        &testSKSClusterName,
					Nodepools: &[]v2.SksNodepool{{
						CreatedAt:    &testSKSNodepoolCreatedAt,
						Description:  &testSKSNodepoolDescription,
						DiskSize:     &testSKSNodepoolDiskSize,
						Id:           &testSKSNodepoolID,
						InstancePool: &v2.Resource{Id: &testSKSNodepoolInstancePoolID},
						InstanceType: &v2.InstanceType{Id: &testSKSNodepoolInstanceTypeID},
						Name:         &testSKSNodepoolName,
						// SecurityGroups: nil, // TODO: the API doesn't return this field ATM
						Size:     &testSKSNodepoolSize,
						State:    &testSKSNodepoolState,
						Template: &v2.Template{Id: &testSKSNodepoolTemplateID},
						Version:  &testSKSNodepoolVersion,
					}},
					State:   &testSKSClusterState,
					Version: &testSKSClusterVersion,
				}},
			})
			if err != nil {
				t.Fatalf("error initializing mock HTTP responder: %s", err)
			}
			return resp, nil
		})

	expected := []*SKSCluster{{
		CreatedAt:   testSKSClusterCreatedAt,
		Description: testSKSClusterDescription,
		Endpoint:    testSKSClusterEndpoint,
		ID:          testSKSClusterID,
		Name:        testSKSClusterName,
		Nodepools: []*SKSNodepool{{
			CreatedAt:      testSKSNodepoolCreatedAt,
			Description:    testSKSNodepoolDescription,
			DiskSize:       testSKSNodepoolDiskSize,
			ID:             testSKSNodepoolID,
			InstancePoolID: testSKSNodepoolInstancePoolID,
			InstanceTypeID: testSKSNodepoolInstanceTypeID,
			Name:           testSKSNodepoolName,
			// SecurityGroups: nil, // TODO: the API doesn't return this field ATM
			Size:       testSKSNodepoolSize,
			State:      testSKSClusterState,
			TemplateID: testSKSNodepoolTemplateID,
			Version:    testSKSNodepoolVersion,
		}},
		State:   testSKSClusterState,
		Version: testSKSClusterVersion,

		c:    client,
		zone: testZone,
	}}

	actual, err := client.ListSKSClusters(context.Background(), testZone)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

// UpdateSKSCluster is not tested as it essentially relies on the already tested GetSKSCluster.
func TestClient_UpdateSKSCluster(t *testing.T) { t.Skip() }

// DeleteSKSCluster is not tested as it only produces API-side effects.
func TestClient_DeleteSKSCluster(t *testing.T) { t.Skip() }
