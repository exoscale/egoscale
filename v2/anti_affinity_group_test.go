package v2

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jarcoal/httpmock"

	"github.com/exoscale/egoscale/v2/oapi"
)

var (
	testAntiAffinityGroupDescription = new(clientTestSuite).randomString(10)
	testAntiAffinityGroupID          = new(clientTestSuite).randomID()
	testAntiAffinityGroupInstanceID  = new(clientTestSuite).randomID()
	testAntiAffinityGroupName        = new(clientTestSuite).randomString(10)
)

func (ts *clientTestSuite) TestClient_CreateAntiAffinityGroup() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
	)

	httpmock.RegisterResponder("POST", "/anti-affinity-group",
		func(req *http.Request) (*http.Response, error) {
			var actual oapi.CreateAntiAffinityGroupJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := oapi.CreateAntiAffinityGroupJSONRequestBody{
				Description: &testAntiAffinityGroupDescription,
				Name:        testAntiAffinityGroupName,
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testPrivateNetworkID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testAntiAffinityGroupID},
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/anti-affinity-group/%s", testAntiAffinityGroupID), oapi.AntiAffinityGroup{
		Description: &testAntiAffinityGroupDescription,
		Id:          &testAntiAffinityGroupID,
		Name:        &testAntiAffinityGroupName,
	})

	expected := &AntiAffinityGroup{
		Description: &testAntiAffinityGroupDescription,
		ID:          &testAntiAffinityGroupID,
		Name:        &testAntiAffinityGroupName,
	}

	actual, err := ts.client.CreateAntiAffinityGroup(context.Background(), testZone, &AntiAffinityGroup{
		Description: &testAntiAffinityGroupDescription,
		ID:          &testAntiAffinityGroupID,
		Name:        &testAntiAffinityGroupName,
	})
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_DeleteAntiAffinityGroup() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		deleted            = false
	)

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("/anti-affinity-group/%s", testAntiAffinityGroupID),
		func(req *http.Request) (*http.Response, error) {
			deleted = true

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testAntiAffinityGroupID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testAntiAffinityGroupID},
	})

	ts.Require().NoError(ts.client.DeleteAntiAffinityGroup(
		context.Background(),
		testZone,
		&AntiAffinityGroup{ID: &testAntiAffinityGroupID},
	))
	ts.Require().True(deleted)
}

func (ts *clientTestSuite) TestClient_FindAntiAffinityGroup() {
	ts.mockAPIRequest("GET", "/anti-affinity-group", struct {
		AntiAffinityGroups *[]oapi.AntiAffinityGroup `json:"anti-affinity-groups,omitempty"`
	}{
		AntiAffinityGroups: &[]oapi.AntiAffinityGroup{{
			Description: &testAntiAffinityGroupDescription,
			Id:          &testAntiAffinityGroupID,
			Name:        &testAntiAffinityGroupName,
		}},
	})

	ts.mockAPIRequest(
		"GET",
		fmt.Sprintf("/anti-affinity-group/%s", testAntiAffinityGroupID),
		oapi.AntiAffinityGroup{
			Description: &testAntiAffinityGroupDescription,
			Id:          &testAntiAffinityGroupID,
			Instances:   &[]oapi.Instance{{Id: &testAntiAffinityGroupInstanceID}},
			Name:        &testAntiAffinityGroupName,
		})

	expected := &AntiAffinityGroup{
		Description: &testAntiAffinityGroupDescription,
		ID:          &testAntiAffinityGroupID,
		InstanceIDs: &[]string{testAntiAffinityGroupInstanceID},
		Name:        &testAntiAffinityGroupName,
	}

	actual, err := ts.client.FindAntiAffinityGroup(context.Background(), testZone, *expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)

	actual, err = ts.client.FindAntiAffinityGroup(context.Background(), testZone, *expected.Name)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_GetAntiAffinityGroup() {
	ts.mockAPIRequest(
		"GET",
		fmt.Sprintf("/anti-affinity-group/%s", testAntiAffinityGroupID),
		oapi.AntiAffinityGroup{
			Description: &testAntiAffinityGroupDescription,
			Id:          &testAntiAffinityGroupID,
			Instances:   &[]oapi.Instance{{Id: &testAntiAffinityGroupInstanceID}},
			Name:        &testAntiAffinityGroupName,
		})

	expected := &AntiAffinityGroup{
		Description: &testAntiAffinityGroupDescription,
		ID:          &testAntiAffinityGroupID,
		InstanceIDs: &[]string{testAntiAffinityGroupInstanceID},
		Name:        &testAntiAffinityGroupName,
	}

	actual, err := ts.client.GetAntiAffinityGroup(context.Background(), testZone, *expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_ListAntiAffinityGroups() {
	ts.mockAPIRequest("GET", "/anti-affinity-group", struct {
		AntiAffinityGroups *[]oapi.AntiAffinityGroup `json:"anti-affinity-groups,omitempty"`
	}{
		AntiAffinityGroups: &[]oapi.AntiAffinityGroup{{
			Description: &testAntiAffinityGroupDescription,
			Id:          &testAntiAffinityGroupID,
			Instances:   &[]oapi.Instance{{Id: &testAntiAffinityGroupInstanceID}},
			Name:        &testAntiAffinityGroupName,
		}},
	})

	expected := []*AntiAffinityGroup{{
		Description: &testAntiAffinityGroupDescription,
		ID:          &testAntiAffinityGroupID,
		InstanceIDs: &[]string{testAntiAffinityGroupInstanceID},
		Name:        &testAntiAffinityGroupName,
	}}

	actual, err := ts.client.ListAntiAffinityGroups(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}
