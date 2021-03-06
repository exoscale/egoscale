package v2

import (
	"context"
	"fmt"

	papi "github.com/exoscale/egoscale/v2/internal/public-api"
)

var (
	testDeployTargetDescription = new(clientTestSuite).randomString(10)
	testDeployTargetID          = new(clientTestSuite).randomID()
	testDeployTargetName        = new(clientTestSuite).randomString(10)
	testDeployTargetType        = "dedicated"
)

func (ts *clientTestSuite) TestClient_ListDeployTargets() {
	ts.mockAPIRequest("GET", "/deploy-target", struct {
		DeployTargets *[]papi.DeployTarget `json:"deploy-targets,omitempty"`
	}{
		DeployTargets: &[]papi.DeployTarget{{
			Description: &testDeployTargetDescription,
			Id:          &testDeployTargetID,
			Name:        &testDeployTargetName,
			Type: func() *papi.DeployTargetType {
				v := papi.DeployTargetType(testDeployTargetType)
				return &v
			}(),
		}},
	})

	expected := []*DeployTarget{{
		Description: &testDeployTargetDescription,
		ID:          &testDeployTargetID,
		Name:        &testDeployTargetName,
		Type:        &testDeployTargetType,
	}}

	actual, err := ts.client.ListDeployTargets(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_GetDeployTarget() {
	ts.mockAPIRequest(
		"GET",
		fmt.Sprintf("/deploy-target/%s", testDeployTargetID),
		papi.DeployTarget{
			Description: &testDeployTargetDescription,
			Id:          &testDeployTargetID,
			Name:        &testDeployTargetName,
			Type: func() *papi.DeployTargetType {
				v := papi.DeployTargetType(testDeployTargetType)
				return &v
			}(),
		})

	expected := &DeployTarget{
		Description: &testDeployTargetDescription,
		ID:          &testDeployTargetID,
		Name:        &testDeployTargetName,
		Type:        &testDeployTargetType,
	}

	actual, err := ts.client.GetDeployTarget(context.Background(), testZone, *expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_FindDeployTarget() {
	ts.mockAPIRequest("GET", "/deploy-target", struct {
		DeployTargets *[]papi.DeployTarget `json:"deploy-targets,omitempty"`
	}{
		DeployTargets: &[]papi.DeployTarget{{
			Description: &testDeployTargetDescription,
			Id:          &testDeployTargetID,
			Name:        &testDeployTargetName,
			Type: func() *papi.DeployTargetType {
				v := papi.DeployTargetType(testDeployTargetType)
				return &v
			}(),
		}},
	})

	ts.mockAPIRequest(
		"GET",
		fmt.Sprintf("/deploy-target/%s", testDeployTargetID),
		papi.DeployTarget{
			Description: &testDeployTargetDescription,
			Id:          &testDeployTargetID,
			Name:        &testDeployTargetName,
			Type: func() *papi.DeployTargetType {
				v := papi.DeployTargetType(testDeployTargetType)
				return &v
			}(),
		})

	expected := &DeployTarget{
		Description: &testDeployTargetDescription,
		ID:          &testDeployTargetID,
		Name:        &testDeployTargetName,
		Type:        &testDeployTargetType,
	}

	actual, err := ts.client.FindDeployTarget(context.Background(), testZone, *expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)

	actual, err = ts.client.FindDeployTarget(context.Background(), testZone, *expected.Name)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}
