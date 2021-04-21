package v2

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jarcoal/httpmock"

	papi "github.com/exoscale/egoscale/v2/internal/public-api"
)

var (
	testAntiAffinityGroupDescription = "Test Anti-Affinity Group description"
	testAntiAffinityGroupID          = new(clientTestSuite).randomID()
	testAntiAffinityGroupName        = "test-anti-affinity-group"
)

func (ts *clientTestSuite) TestClient_CreateAntiAffinityGroup() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = "success"
	)

	httpmock.RegisterResponder("POST", "/anti-affinity-group",
		func(req *http.Request) (*http.Response, error) {
			var actual papi.CreateAntiAffinityGroupJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.CreateAntiAffinityGroupJSONRequestBody{
				Description: &testAntiAffinityGroupDescription,
				Name:        testAntiAffinityGroupName,
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testPrivateNetworkID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testAntiAffinityGroupID},
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/anti-affinity-group/%s", testAntiAffinityGroupID), papi.AntiAffinityGroup{
		Description: &testAntiAffinityGroupDescription,
		Id:          &testAntiAffinityGroupID,
		Name:        &testAntiAffinityGroupName,
	})

	expected := &AntiAffinityGroup{
		Description: testAntiAffinityGroupDescription,
		ID:          testAntiAffinityGroupID,
		Name:        testAntiAffinityGroupName,
	}

	actual, err := ts.client.CreateAntiAffinityGroup(context.Background(), testZone, &AntiAffinityGroup{
		Description: testAntiAffinityGroupDescription,
		ID:          testAntiAffinityGroupID,
		Name:        testAntiAffinityGroupName,
	})
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_ListAntiAffinityGroups() {
	ts.mockAPIRequest("GET", "/anti-affinity-group", struct {
		AntiAffinityGroups *[]papi.AntiAffinityGroup `json:"anti-affinity-groups,omitempty"`
	}{
		AntiAffinityGroups: &[]papi.AntiAffinityGroup{{
			Description: &testAntiAffinityGroupDescription,
			Id:          &testAntiAffinityGroupID,
			Name:        &testAntiAffinityGroupName,
		}},
	})

	expected := []*AntiAffinityGroup{{
		Description: testAntiAffinityGroupDescription,
		ID:          testAntiAffinityGroupID,
		Name:        testAntiAffinityGroupName,
	}}

	actual, err := ts.client.ListAntiAffinityGroups(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_GetAntiAffinityGroup() {
	ts.mockAPIRequest("GET", fmt.Sprintf("/anti-affinity-group/%s", testAntiAffinityGroupID), papi.AntiAffinityGroup{
		Description: &testAntiAffinityGroupDescription,
		Id:          &testAntiAffinityGroupID,
		Name:        &testAntiAffinityGroupName,
	})

	expected := &AntiAffinityGroup{
		Description: testAntiAffinityGroupDescription,
		ID:          testAntiAffinityGroupID,
		Name:        testAntiAffinityGroupName,
	}

	actual, err := ts.client.GetAntiAffinityGroup(context.Background(), testZone, expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_DeleteAntiAffinityGroup() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = "success"
		deleted            = false
	)

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("/anti-affinity-group/%s", testAntiAffinityGroupID),
		func(req *http.Request) (*http.Response, error) {
			deleted = true

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testAntiAffinityGroupID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testAntiAffinityGroupID},
	})

	ts.Require().NoError(ts.client.DeleteAntiAffinityGroup(context.Background(), testZone, testAntiAffinityGroupID))
	ts.Require().True(deleted)
}
