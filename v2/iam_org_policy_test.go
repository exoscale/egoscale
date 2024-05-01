package v2

import (
	"context"
	"net/http"

	"github.com/exoscale/egoscale/v2/oapi"
	"github.com/stretchr/testify/mock"
)

var (
	testIAMPolicyID                    = new(testSuite).randomString(10)
	testIAMPolicyServiceType           = new(testSuite).randomString(10)
	testIAMPolicyServiceRuleAction     = "allow"
	testIAMPolicyServiceRuleExpression = new(testSuite).randomString(10)
)

func (ts *testSuite) TestClient_GetIAMOrgPolicy() {
	ts.mock().
		On("GetIamOrganizationPolicyWithResponse",
			mock.Anything,                 // ctx
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.GetIamOrganizationPolicyResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.IamPolicy{
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
		}, nil)

	expected := &IAMPolicy{
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
	}

	actual, err := ts.client.GetIAMOrgPolicy(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_UpdateIAMOrgPolicy() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
	)

	ts.mock().
		On(
			"UpdateIamOrganizationPolicyWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				oapi.UpdateIamOrganizationPolicyJSONRequestBody{
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
				args.Get(1),
			)
		}).
		Return(
			&oapi.UpdateIamOrganizationPolicyResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testIAMPolicyID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testIAMPolicyID, nil),
		State:     &testOperationState,
	})

	err := ts.client.UpdateIAMOrgPolicy(context.Background(), testZone, &IAMPolicy{
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
	})
	ts.Require().NoError(err)
}
