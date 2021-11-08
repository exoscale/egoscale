package v2

import (
	"context"
	"net/http"

	"github.com/stretchr/testify/mock"

	"github.com/exoscale/egoscale/v2/oapi"
)

var (
	testSecurityGroupDescription = new(testSuite).randomString(10)
	testSecurityGroupID          = new(testSuite).randomID()
	testSecurityGroupName        = new(testSuite).randomString(10)
)

func (ts *testSuite) TestClient_CreateSecurityGroup() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
	)

	ts.mock().
		On(
			"CreateSecurityGroupWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				oapi.CreateSecurityGroupJSONRequestBody{
					Description: &testSecurityGroupDescription,
					Name:        testSecurityGroupName,
				},
				args.Get(1),
			)
		}).
		Return(
			&oapi.CreateSecurityGroupResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: &oapi.Reference{Id: &testSecurityGroupID},
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: &oapi.Reference{Id: &testSecurityGroupID},
		State:     &testOperationState,
	})

	ts.mock().
		On("GetSecurityGroupWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.GetSecurityGroupResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.SecurityGroup{
				Description: &testSecurityGroupDescription,
				Id:          &testSecurityGroupID,
				Name:        &testSecurityGroupName,
			},
		}, nil)

	expected := &SecurityGroup{
		Description: &testSecurityGroupDescription,
		ID:          &testSecurityGroupID,
		Name:        &testSecurityGroupName,
	}

	actual, err := ts.client.CreateSecurityGroup(context.Background(), testZone, &SecurityGroup{
		Description: &testSecurityGroupDescription,
		Name:        &testSecurityGroupName,
	})
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_AddExternalSourceToSecurityGroup() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		added              bool
	)

	ts.mock().
		On(
			"AddExternalSourceToSecurityGroupWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testSecurityGroupID, args.Get(1))
			ts.Require().Equal(
				oapi.AddExternalSourceToSecurityGroupJSONRequestBody{
					Cidr: testSecurityGroupExternalSource,
				},
				args.Get(2),
			)
			added = true
		}).
		Return(
			&oapi.AddExternalSourceToSecurityGroupResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: &oapi.Reference{Id: &testSecurityGroupID},
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: &oapi.Reference{Id: &testSecurityGroupID},
		State:     &testOperationState,
	})

	securityGroup := &SecurityGroup{
		Description: &testSecurityGroupDescription,
		ID:          &testSecurityGroupID,
		Name:        &testSecurityGroupName,
	}

	err := ts.client.AddExternalSourceToSecurityGroup(
		context.Background(),
		testZone,
		securityGroup,
		testSecurityGroupExternalSource)
	ts.Require().NoError(err)
	ts.Require().True(added)
}

func (ts *testSuite) TestClient_DeleteSecurityGroup() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		deleted            = false
	)

	ts.mock().
		On(
			"DeleteSecurityGroupWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testSecurityGroupID, args.Get(1))
			deleted = true
		}).
		Return(
			&oapi.DeleteSecurityGroupResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: &oapi.Reference{Id: &testSecurityGroupID},
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: &oapi.Reference{Id: &testSecurityGroupID},
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.DeleteSecurityGroup(
		context.Background(),
		testZone,
		&SecurityGroup{ID: &testSecurityGroupID},
	))
	ts.Require().True(deleted)
}

func (ts *testSuite) TestClient_FindSecurityGroup() {
	ts.mock().
		On("ListSecurityGroupsWithResponse",
			mock.Anything,                 // ctx
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.ListSecurityGroupsResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &struct {
				SecurityGroups *[]oapi.SecurityGroup `json:"security-groups,omitempty"`
			}{
				SecurityGroups: &[]oapi.SecurityGroup{{
					Id:   &testSecurityGroupID,
					Name: &testSecurityGroupName,
				}},
			},
		}, nil)

	ts.mock().
		On("GetSecurityGroupWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testSecurityGroupID, args.Get(1))
		}).
		Return(&oapi.GetSecurityGroupResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.SecurityGroup{
				Id:   &testSecurityGroupID,
				Name: &testSecurityGroupName,
			},
		}, nil)

	expected := &SecurityGroup{
		ID:   &testSecurityGroupID,
		Name: &testSecurityGroupName,
	}

	actual, err := ts.client.FindSecurityGroup(context.Background(), testZone, *expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)

	actual, err = ts.client.FindSecurityGroup(context.Background(), testZone, *expected.Name)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_GetSecurityGroup() {
	ts.mock().
		On("GetSecurityGroupWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testSecurityGroupID, args.Get(1))
		}).
		Return(&oapi.GetSecurityGroupResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.SecurityGroup{
				Description:     &testSecurityGroupDescription,
				Id:              &testSecurityGroupID,
				Name:            &testSecurityGroupName,
				ExternalSources: &testSecurityGroupExternalSources,
				Rules: &[]oapi.SecurityGroupRule{{
					Description:   &testSecurityGroupRuleDescription,
					EndPort:       func() *int64 { v := int64(testSecurityGroupRuleEndPort); return &v }(),
					FlowDirection: &testSecurityGroupRuleFlowDirection,
					Icmp: &struct {
						Code *int64 `json:"code,omitempty"`
						Type *int64 `json:"type,omitempty"`
					}{
						Code: &testSecurityGroupRuleICMPCode,
						Type: &testSecurityGroupRuleICMPType,
					},
					Id:            &testSecurityGroupRuleID,
					Network:       &testSecurityGroupRuleNetwork,
					Protocol:      &testSecurityGroupRuleProtocol,
					SecurityGroup: &oapi.SecurityGroupResource{Id: testSecurityGroupRuleSecurityGroupID},
					StartPort:     func() *int64 { v := int64(testSecurityGroupRuleStartPort); return &v }(),
				}},
			},
		}, nil)

	expected := &SecurityGroup{
		Description:     &testSecurityGroupDescription,
		ID:              &testSecurityGroupID,
		Name:            &testSecurityGroupName,
		ExternalSources: &testSecurityGroupExternalSources,
		Rules: []*SecurityGroupRule{{
			Description:     &testSecurityGroupRuleDescription,
			EndPort:         &testSecurityGroupRuleEndPort,
			FlowDirection:   (*string)(&testSecurityGroupRuleFlowDirection),
			ICMPCode:        &testSecurityGroupRuleICMPCode,
			ICMPType:        &testSecurityGroupRuleICMPType,
			ID:              &testSecurityGroupRuleID,
			Network:         testSecurityGroupRuleNetworkP,
			Protocol:        (*string)(&testSecurityGroupRuleProtocol),
			SecurityGroupID: &testSecurityGroupRuleSecurityGroupID,
			StartPort:       &testSecurityGroupRuleStartPort,
		}},
	}

	actual, err := ts.client.GetSecurityGroup(context.Background(), testZone, *expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_ListSecurityGroups() {
	ts.mock().
		On("ListSecurityGroupsWithResponse",
			mock.Anything,                 // ctx
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.ListSecurityGroupsResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &struct {
				SecurityGroups *[]oapi.SecurityGroup `json:"security-groups,omitempty"`
			}{
				SecurityGroups: &[]oapi.SecurityGroup{{
					Description: &testSecurityGroupDescription,
					Id:          &testSecurityGroupID,
					Name:        &testSecurityGroupName,
					Rules: &[]oapi.SecurityGroupRule{{
						Description:   &testSecurityGroupRuleDescription,
						EndPort:       func() *int64 { v := int64(testSecurityGroupRuleEndPort); return &v }(),
						FlowDirection: &testSecurityGroupRuleFlowDirection,
						Icmp: &struct {
							Code *int64 `json:"code,omitempty"`
							Type *int64 `json:"type,omitempty"`
						}{
							Code: &testSecurityGroupRuleICMPCode,
							Type: &testSecurityGroupRuleICMPType,
						},
						Id:            &testSecurityGroupRuleID,
						Network:       &testSecurityGroupRuleNetwork,
						Protocol:      &testSecurityGroupRuleProtocol,
						SecurityGroup: &oapi.SecurityGroupResource{Id: testSecurityGroupRuleSecurityGroupID},
						StartPort:     func() *int64 { v := int64(testSecurityGroupRuleStartPort); return &v }(),
					}},
				}},
			},
		}, nil)

	expected := []*SecurityGroup{
		{
			Description: &testSecurityGroupDescription,
			ID:          &testSecurityGroupID,
			Name:        &testSecurityGroupName,
			Rules: []*SecurityGroupRule{{
				Description:     &testSecurityGroupRuleDescription,
				EndPort:         &testSecurityGroupRuleEndPort,
				FlowDirection:   (*string)(&testSecurityGroupRuleFlowDirection),
				ICMPCode:        &testSecurityGroupRuleICMPCode,
				ICMPType:        &testSecurityGroupRuleICMPType,
				ID:              &testSecurityGroupRuleID,
				Network:         testSecurityGroupRuleNetworkP,
				Protocol:        (*string)(&testSecurityGroupRuleProtocol),
				SecurityGroupID: &testSecurityGroupRuleSecurityGroupID,
				StartPort:       &testSecurityGroupRuleStartPort,
			}},
		},
	}

	actual, err := ts.client.ListSecurityGroups(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_RemoveExternalSourceFromSecurityGroup() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		removed            bool
	)

	ts.mock().
		On(
			"RemoveExternalSourceFromSecurityGroupWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testSecurityGroupID, args.Get(1))
			ts.Require().Equal(
				oapi.RemoveExternalSourceFromSecurityGroupJSONRequestBody{
					Cidr: testSecurityGroupExternalSource,
				},
				args.Get(2),
			)
			removed = true
		}).
		Return(
			&oapi.RemoveExternalSourceFromSecurityGroupResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: &oapi.Reference{Id: &testSecurityGroupID},
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: &oapi.Reference{Id: &testSecurityGroupID},
		State:     &testOperationState,
	})

	securityGroup := &SecurityGroup{
		Description: &testSecurityGroupDescription,
		ID:          &testSecurityGroupID,
		Name:        &testSecurityGroupName,
	}

	err := ts.client.RemoveExternalSourceFromSecurityGroup(
		context.Background(),
		testZone,
		securityGroup,
		testSecurityGroupExternalSource)
	ts.Require().NoError(err)
	ts.Require().True(removed)
}
