package v2

import (
	"context"
	"net/http"

	"github.com/stretchr/testify/mock"

	"github.com/exoscale/egoscale/v2/oapi"
)

var (
	testDeployTargetDescription = new(testSuite).randomString(10)
	testDeployTargetID          = new(testSuite).randomID()
	testDeployTargetName        = new(testSuite).randomString(10)
	testDeployTargetType        = "dedicated"
)

func (ts *testSuite) TestClient_FindDeployTarget() {
	ts.mock().
		On("ListDeployTargetsWithResponse",
			mock.Anything,                 // ctx
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.ListDeployTargetsResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &struct {
				DeployTargets *[]oapi.DeployTarget `json:"deploy-targets,omitempty"`
			}{
				DeployTargets: &[]oapi.DeployTarget{{
					Description: &testDeployTargetDescription,
					Id:          testDeployTargetID,
					Name:        &testDeployTargetName,
					Type:        (*oapi.DeployTargetType)(&testDeployTargetType),
				}},
			},
		}, nil)

	ts.mock().
		On("GetDeployTargetWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testDeployTargetID, args.Get(1))
		}).
		Return(&oapi.GetDeployTargetResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.DeployTarget{
				Description: &testDeployTargetDescription,
				Id:          testDeployTargetID,
				Name:        &testDeployTargetName,
				Type:        (*oapi.DeployTargetType)(&testDeployTargetType),
			},
		}, nil)

	expected := &DeployTarget{
		Description: &testDeployTargetDescription,
		ID:          &testDeployTargetID,
		Name:        &testDeployTargetName,
		Type:        &testDeployTargetType,
		Zone:        &testZone,
	}

	actual, err := ts.client.FindDeployTarget(context.Background(), testZone, *expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)

	actual, err = ts.client.FindDeployTarget(context.Background(), testZone, *expected.Name)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_GetDeployTarget() {
	ts.mock().
		On("GetDeployTargetWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testDeployTargetID, args.Get(1))
		}).
		Return(&oapi.GetDeployTargetResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.DeployTarget{
				Description: &testDeployTargetDescription,
				Id:          testDeployTargetID,
				Name:        &testDeployTargetName,
				Type:        (*oapi.DeployTargetType)(&testDeployTargetType),
			},
		}, nil)

	expected := &DeployTarget{
		Description: &testDeployTargetDescription,
		ID:          &testDeployTargetID,
		Name:        &testDeployTargetName,
		Type:        &testDeployTargetType,
		Zone:        &testZone,
	}

	actual, err := ts.client.GetDeployTarget(context.Background(), testZone, *expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_ListDeployTargets() {
	ts.mock().
		On("ListDeployTargetsWithResponse",
			mock.Anything,                 // ctx
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.ListDeployTargetsResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &struct {
				DeployTargets *[]oapi.DeployTarget `json:"deploy-targets,omitempty"`
			}{
				DeployTargets: &[]oapi.DeployTarget{{
					Description: &testDeployTargetDescription,
					Id:          testDeployTargetID,
					Name:        &testDeployTargetName,
					Type:        (*oapi.DeployTargetType)(&testDeployTargetType),
				}},
			},
		}, nil)

	expected := []*DeployTarget{{
		Description: &testDeployTargetDescription,
		ID:          &testDeployTargetID,
		Name:        &testDeployTargetName,
		Type:        &testDeployTargetType,
		Zone:        &testZone,
	}}

	actual, err := ts.client.ListDeployTargets(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}
