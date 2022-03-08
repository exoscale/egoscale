package v2

import (
	"context"
	"net/http"

	"github.com/stretchr/testify/mock"

	"github.com/exoscale/egoscale/v2/oapi"
)

var (
	testAntiAffinityGroupDescription = new(testSuite).randomString(10)
	testAntiAffinityGroupID          = new(testSuite).randomID()
	testAntiAffinityGroupInstanceID  = new(testSuite).randomID()
	testAntiAffinityGroupName        = new(testSuite).randomString(10)
)

func (ts *testSuite) TestClient_CreateAntiAffinityGroup() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
	)

	ts.mock().
		On(
			"CreateAntiAffinityGroupWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				oapi.CreateAntiAffinityGroupJSONRequestBody{
					Description: &testAntiAffinityGroupDescription,
					Name:        testAntiAffinityGroupName,
				},
				args.Get(1),
			)
		}).
		Return(
			&oapi.CreateAntiAffinityGroupResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id: &testOperationID,
					Reference: &struct {
						Command *string `json:"command,omitempty"`
						Id      *string `json:"id,omitempty"` // revive:disable-line
						Link    *string `json:"link,omitempty"`
					}{Id: &testAntiAffinityGroupID},
					State: &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id: &testOperationID,
		Reference: &struct {
			Command *string `json:"command,omitempty"`
			Id      *string `json:"id,omitempty"` // revive:disable-line
			Link    *string `json:"link,omitempty"`
		}{Id: &testAntiAffinityGroupID},
		State: &testOperationState,
	})

	ts.mock().
		On("GetAntiAffinityGroupWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.GetAntiAffinityGroupResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.AntiAffinityGroup{
				Description: &testAntiAffinityGroupDescription,
				Id:          &testAntiAffinityGroupID,
				Name:        &testAntiAffinityGroupName,
			},
		}, nil)

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

func (ts *testSuite) TestClient_DeleteAntiAffinityGroup() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		deleted            = false
	)

	ts.mock().
		On(
			"DeleteAntiAffinityGroupWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testAntiAffinityGroupID, args.Get(1))
			deleted = true
		}).
		Return(
			&oapi.DeleteAntiAffinityGroupResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id: &testOperationID,
					Reference: &struct {
						Command *string `json:"command,omitempty"`
						Id      *string `json:"id,omitempty"` // revive:disable-line
						Link    *string `json:"link,omitempty"`
					}{Id: &testAntiAffinityGroupID},
					State: &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id: &testOperationID,
		Reference: &struct {
			Command *string `json:"command,omitempty"`
			Id      *string `json:"id,omitempty"` // revive:disable-line
			Link    *string `json:"link,omitempty"`
		}{Id: &testAntiAffinityGroupID},
		State: &testOperationState,
	})

	ts.Require().NoError(ts.client.DeleteAntiAffinityGroup(
		context.Background(),
		testZone,
		&AntiAffinityGroup{ID: &testAntiAffinityGroupID},
	))
	ts.Require().True(deleted)
}

func (ts *testSuite) TestClient_FindAntiAffinityGroup() {
	ts.mock().
		On("ListAntiAffinityGroupsWithResponse",
			mock.Anything,                 // ctx
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.ListAntiAffinityGroupsResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &struct {
				AntiAffinityGroups *[]oapi.AntiAffinityGroup `json:"anti-affinity-groups,omitempty"`
			}{
				AntiAffinityGroups: &[]oapi.AntiAffinityGroup{{
					Description: &testAntiAffinityGroupDescription,
					Id:          &testAntiAffinityGroupID,
					Name:        &testAntiAffinityGroupName,
				}},
			},
		}, nil)

	ts.mock().
		On("GetAntiAffinityGroupWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testAntiAffinityGroupID, args.Get(1))
		}).
		Return(&oapi.GetAntiAffinityGroupResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.AntiAffinityGroup{
				Description: &testAntiAffinityGroupDescription,
				Id:          &testAntiAffinityGroupID,
				Instances:   &[]oapi.Instance{{Id: &testAntiAffinityGroupInstanceID}},
				Name:        &testAntiAffinityGroupName,
			},
		}, nil)

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

func (ts *testSuite) TestClient_GetAntiAffinityGroup() {
	ts.mock().
		On("GetAntiAffinityGroupWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testAntiAffinityGroupID, args.Get(1))
		}).
		Return(&oapi.GetAntiAffinityGroupResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.AntiAffinityGroup{
				Description: &testAntiAffinityGroupDescription,
				Id:          &testAntiAffinityGroupID,
				Instances:   &[]oapi.Instance{{Id: &testAntiAffinityGroupInstanceID}},
				Name:        &testAntiAffinityGroupName,
			},
		}, nil)

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

func (ts *testSuite) TestClient_ListAntiAffinityGroups() {
	ts.mock().
		On("ListAntiAffinityGroupsWithResponse",
			mock.Anything,                 // ctx
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.ListAntiAffinityGroupsResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &struct {
				AntiAffinityGroups *[]oapi.AntiAffinityGroup `json:"anti-affinity-groups,omitempty"`
			}{
				AntiAffinityGroups: &[]oapi.AntiAffinityGroup{{
					Description: &testAntiAffinityGroupDescription,
					Id:          &testAntiAffinityGroupID,
					Instances:   &[]oapi.Instance{{Id: &testAntiAffinityGroupInstanceID}},
					Name:        &testAntiAffinityGroupName,
				}},
			},
		}, nil)

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
