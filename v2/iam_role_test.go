package v2

import (
	"context"
	"net/http"

	"github.com/exoscale/egoscale/v2/oapi"
	"github.com/stretchr/testify/mock"
)

var (
	testIAMRoleID          = new(testSuite).randomString(10)
	testIAMRoleName        = new(testSuite).randomString(10)
	testIAMRoleDescription = new(testSuite).randomString(10)
	testIAMRoleEdiable     = true
)

func (ts *testSuite) TestClient_GetIAMRole() {
	ts.mock().
		On("GetIamRoleWithResponse",
			mock.Anything, // ctx
			testIAMRoleID,
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.GetIamRoleResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.IamRole{
				Id:          &testIAMRoleID,
				Name:        &testIAMRoleName,
				Description: &testIAMRoleDescription,
				Editable:    &testIAMRoleEdiable,
				Labels: &oapi.Labels{
					AdditionalProperties: map[string]string{
						"foo": "bar",
					},
				},
				Permissions: &[]oapi.IamRolePermissions{
					"bypass-governance-retention",
				},
				Policy: &oapi.IamPolicy{
					DefaultServiceStrategy: oapi.IamPolicyDefaultServiceStrategyAllow,
					Services: oapi.IamPolicy_Services{
						AdditionalProperties: map[string]oapi.IamServicePolicy{
							"test": oapi.IamServicePolicy{
								Type: (*oapi.IamServicePolicyType)(&testIAMPolicyServiceType),
								Rules: &[]oapi.IamServicePolicyRule{
									{
										Action:     (*oapi.IamServicePolicyRuleAction)(&testIAMPolicyServiceRuleAction),
										Expression: &testIAMPolicyServiceRuleExpression,
									},
								},
							},
						},
					},
				},
			},
		}, nil)

	expected := &IAMRole{
		ID:          &testIAMRoleID,
		Name:        &testIAMRoleName,
		Description: &testIAMRoleDescription,
		Editable:    &testIAMRoleEdiable,
		Labels: map[string]string{
			"foo": "bar",
		},
		Permissions: []string{
			"bypass-governance-retention",
		},
		Policy: &IAMPolicy{
			DefaultServiceStrategy: string(oapi.IamPolicyDefaultServiceStrategyAllow),
			Services: map[string]IAMPolicyService{
				"test": IAMPolicyService{
					Type: &testIAMPolicyServiceType,
					Rules: []IAMPolicyServiceRule{
						IAMPolicyServiceRule{
							Action:     (*string)(&testIAMPolicyServiceRuleAction),
							Expression: &testIAMPolicyServiceRuleExpression,
						},
					},
				},
			},
		},
	}

	actual, err := ts.client.GetIAMRole(context.Background(), testZone, testIAMRoleID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_ListIAMRoles() {
	ts.mock().
		On("ListIamRolesWithResponse",
			mock.Anything,                 // ctx
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.ListIamRolesResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &struct {
				IamRoles *[]oapi.IamRole `json:"iam-roles,omitempty"`
			}{
				IamRoles: &[]oapi.IamRole{
					{
						Id:          &testIAMRoleID,
						Name:        &testIAMRoleName,
						Description: &testIAMRoleDescription,
						Editable:    &testIAMRoleEdiable,
						Labels: &oapi.Labels{
							AdditionalProperties: map[string]string{
								"foo": "bar",
							},
						},
						Permissions: &[]oapi.IamRolePermissions{
							"bypass-governance-retention",
						},
						Policy: &oapi.IamPolicy{
							DefaultServiceStrategy: oapi.IamPolicyDefaultServiceStrategyAllow,
							Services: oapi.IamPolicy_Services{
								AdditionalProperties: map[string]oapi.IamServicePolicy{
									"test": oapi.IamServicePolicy{
										Type: (*oapi.IamServicePolicyType)(&testIAMPolicyServiceType),
										Rules: &[]oapi.IamServicePolicyRule{
											{
												Action:     (*oapi.IamServicePolicyRuleAction)(&testIAMPolicyServiceRuleAction),
												Expression: &testIAMPolicyServiceRuleExpression,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}, nil)

	expected := []*IAMRole{
		&IAMRole{
			ID:          &testIAMRoleID,
			Name:        &testIAMRoleName,
			Description: &testIAMRoleDescription,
			Editable:    &testIAMRoleEdiable,
			Labels: map[string]string{
				"foo": "bar",
			},
			Permissions: []string{
				"bypass-governance-retention",
			},
			Policy: &IAMPolicy{
				DefaultServiceStrategy: string(oapi.IamPolicyDefaultServiceStrategyAllow),
				Services: map[string]IAMPolicyService{
					"test": IAMPolicyService{
						Type: &testIAMPolicyServiceType,
						Rules: []IAMPolicyServiceRule{
							IAMPolicyServiceRule{
								Action:     (*string)(&testIAMPolicyServiceRuleAction),
								Expression: &testIAMPolicyServiceRuleExpression,
							},
						},
					},
				},
			},
		},
	}

	actual, err := ts.client.ListIAMRoles(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_CreateIAMRole() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
	)

	ts.mock().
		On(
			"CreateIamRoleWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				oapi.CreateIamRoleJSONRequestBody{
					Name:        testIAMRoleName,
					Description: &testIAMRoleDescription,
					Editable:    &testIAMRoleEdiable,
					Labels: &oapi.Labels{
						AdditionalProperties: map[string]string{
							"foo": "bar",
						},
					},
					Permissions: &[]oapi.CreateIamRoleJSONBodyPermissions{
						"bypass-governance-retention",
					},
					Policy: &oapi.IamPolicy{
						DefaultServiceStrategy: oapi.IamPolicyDefaultServiceStrategyAllow,
						Services: oapi.IamPolicy_Services{
							AdditionalProperties: map[string]oapi.IamServicePolicy{
								"test": oapi.IamServicePolicy{
									Type: (*oapi.IamServicePolicyType)(&testIAMPolicyServiceType),
									Rules: &[]oapi.IamServicePolicyRule{
										{
											Action:     (*oapi.IamServicePolicyRuleAction)(&testIAMPolicyServiceRuleAction),
											Expression: &testIAMPolicyServiceRuleExpression,
										},
									},
								},
							},
						},
					},
				},
				args.Get(1),
			)
		}).
		Return(
			&oapi.CreateIamRoleResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testIAMRoleID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testIAMRoleID, nil),
		State:     &testOperationState,
	})

	ts.mock().
		On("GetIamRoleWithResponse",
			mock.Anything,                 // ctx
			testIAMRoleID,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.GetIamRoleResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.IamRole{
				Id:          &testIAMRoleID,
				Name:        &testIAMRoleName,
				Description: &testIAMRoleDescription,
				Editable:    &testIAMRoleEdiable,
				Labels: &oapi.Labels{
					AdditionalProperties: map[string]string{
						"foo": "bar",
					},
				},
				Permissions: &[]oapi.IamRolePermissions{
					"bypass-governance-retention",
				},
				Policy: &oapi.IamPolicy{
					DefaultServiceStrategy: oapi.IamPolicyDefaultServiceStrategyAllow,
					Services: oapi.IamPolicy_Services{
						AdditionalProperties: map[string]oapi.IamServicePolicy{
							"test": oapi.IamServicePolicy{
								Type: (*oapi.IamServicePolicyType)(&testIAMPolicyServiceType),
								Rules: &[]oapi.IamServicePolicyRule{
									{
										Action:     (*oapi.IamServicePolicyRuleAction)(&testIAMPolicyServiceRuleAction),
										Expression: &testIAMPolicyServiceRuleExpression,
									},
								},
							},
						},
					},
				},
			},
		}, nil)

	expected := &IAMRole{
		ID:          &testIAMRoleID,
		Name:        &testIAMRoleName,
		Description: &testIAMRoleDescription,
		Editable:    &testIAMRoleEdiable,
		Labels: map[string]string{
			"foo": "bar",
		},
		Permissions: []string{
			"bypass-governance-retention",
		},
		Policy: &IAMPolicy{
			DefaultServiceStrategy: string(oapi.IamPolicyDefaultServiceStrategyAllow),
			Services: map[string]IAMPolicyService{
				"test": IAMPolicyService{
					Type: &testIAMPolicyServiceType,
					Rules: []IAMPolicyServiceRule{
						IAMPolicyServiceRule{
							Action:     (*string)(&testIAMPolicyServiceRuleAction),
							Expression: &testIAMPolicyServiceRuleExpression,
						},
					},
				},
			},
		},
	}

	actual, err := ts.client.CreateIAMRole(context.Background(), testZone, &IAMRole{
		Name:        &testIAMRoleName,
		Description: &testIAMRoleDescription,
		Editable:    &testIAMRoleEdiable,
		Labels: map[string]string{
			"foo": "bar",
		},
		Permissions: []string{
			"bypass-governance-retention",
		},
		Policy: &IAMPolicy{
			DefaultServiceStrategy: string(oapi.IamPolicyDefaultServiceStrategyAllow),
			Services: map[string]IAMPolicyService{
				"test": IAMPolicyService{
					Type: &testIAMPolicyServiceType,
					Rules: []IAMPolicyServiceRule{
						IAMPolicyServiceRule{
							Action:     (*string)(&testIAMPolicyServiceRuleAction),
							Expression: &testIAMPolicyServiceRuleExpression,
						},
					},
				},
			},
		},
	})
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_DeleteIAMRole() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		deleted            = false
	)

	ts.mock().
		On(
			"DeleteIamRoleWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testIAMRoleID, args.Get(1))
			deleted = true
		}).
		Return(
			&oapi.DeleteIamRoleResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testIAMRoleID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testIAMRoleID, nil),
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.DeleteIAMRole(
		context.Background(),
		testZone,
		&IAMRole{ID: &testIAMRoleID},
	))
	ts.Require().True(deleted)
}

func (ts *testSuite) TestClient_UpdateIAMRole() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
	)

	ts.mock().
		On(
			"UpdateIamRoleWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				testIAMRoleID,
				args.Get(1),
			)
			ts.Require().Equal(
				oapi.UpdateIamRoleJSONRequestBody{
					Description: &testIAMRoleDescription,
					Labels: &oapi.Labels{
						AdditionalProperties: map[string]string{
							"foo": "bar",
						},
					},
					Permissions: &[]oapi.UpdateIamRoleJSONBodyPermissions{
						"bypass-governance-retention",
					},
				},
				args.Get(2),
			)
		}).
		Return(
			&oapi.UpdateIamRoleResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testIAMRoleID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testIAMRoleID, nil),
		State:     &testOperationState,
	})

	err := ts.client.UpdateIAMRole(context.Background(), testZone, &IAMRole{
		ID:          &testIAMRoleID,
		Description: &testIAMRoleDescription,
		Labels: map[string]string{
			"foo": "bar",
		},
		Permissions: []string{
			"bypass-governance-retention",
		},
	})
	ts.Require().NoError(err)
}

func (ts *testSuite) TestClient_UpdateIAMRolePolicy() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
	)

	ts.mock().
		On(
			"UpdateIamRolePolicyWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				testIAMRoleID,
				args.Get(1),
			)
			ts.Require().Equal(
				oapi.UpdateIamRolePolicyJSONRequestBody{
					DefaultServiceStrategy: oapi.IamPolicyDefaultServiceStrategyAllow,
					Services: oapi.IamPolicy_Services{
						AdditionalProperties: map[string]oapi.IamServicePolicy{
							"test": oapi.IamServicePolicy{
								Type: (*oapi.IamServicePolicyType)(&testIAMPolicyServiceType),
								Rules: &[]oapi.IamServicePolicyRule{
									{
										Action:     (*oapi.IamServicePolicyRuleAction)(&testIAMPolicyServiceRuleAction),
										Expression: &testIAMPolicyServiceRuleExpression,
									},
								},
							},
						},
					},
				},
				args.Get(2),
			)
		}).
		Return(
			&oapi.UpdateIamRolePolicyResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testIAMRoleID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testIAMRoleID, nil),
		State:     &testOperationState,
	})

	err := ts.client.UpdateIAMRolePolicy(context.Background(), testZone, &IAMRole{
		ID: &testIAMRoleID,
		Policy: &IAMPolicy{
			DefaultServiceStrategy: string(oapi.IamPolicyDefaultServiceStrategyAllow),
			Services: map[string]IAMPolicyService{
				"test": IAMPolicyService{
					Type: &testIAMPolicyServiceType,
					Rules: []IAMPolicyServiceRule{
						IAMPolicyServiceRule{
							Action:     (*string)(&testIAMPolicyServiceRuleAction),
							Expression: &testIAMPolicyServiceRuleExpression,
						},
					},
				},
			},
		},
	})
	ts.Require().NoError(err)
}
