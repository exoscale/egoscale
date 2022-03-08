package v2

import (
	"context"
	"net"
	"net/http"

	"github.com/stretchr/testify/mock"

	"github.com/exoscale/egoscale/v2/oapi"
)

var (
	testSecurityGroupRuleDescription            = new(testSuite).randomString(10)
	testSecurityGroupRuleEndPort         uint16 = 8080
	testSecurityGroupRuleFlowDirection          = oapi.SecurityGroupRuleFlowDirectionIngress
	testSecurityGroupRuleICMPCode        int64  = 0 // nolint:revive
	testSecurityGroupRuleICMPType        int64  = 8
	testSecurityGroupRuleID                     = new(testSuite).randomID()
	testSecurityGroupRuleNetwork                = "1.2.3.0/24"
	_, testSecurityGroupRuleNetworkP, _         = net.ParseCIDR(testSecurityGroupRuleNetwork)
	testSecurityGroupRuleProtocol               = oapi.SecurityGroupRuleProtocolIcmp
	testSecurityGroupRuleSecurityGroupID        = new(testSuite).randomID()
	testSecurityGroupExternalSource             = "8.8.8.8/32"
	testSecurityGroupExternalSources            = []string{"8.8.8.8/32"}
	testSecurityGroupRuleStartPort       uint16 = 8081
)

func (ts *testSuite) TestClient_CreateSecurityGroupRule() {
	var (
		testOperationID                        = ts.randomID()
		testOperationState                     = oapi.OperationStateSuccess
		testAlreadyExistingSecurityGroupRuleID = ts.randomID()
	)

	ts.mock().
		On(
			"AddRuleToSecurityGroupWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testSecurityGroupID, args.Get(1))
			ts.Require().Equal(oapi.AddRuleToSecurityGroupJSONRequestBody{
				Description:   &testSecurityGroupRuleDescription,
				EndPort:       func() *int64 { v := int64(testSecurityGroupRuleEndPort); return &v }(),
				FlowDirection: oapi.AddRuleToSecurityGroupJSONBodyFlowDirection(testSecurityGroupRuleFlowDirection),
				Icmp: &struct {
					Code *int64 `json:"code,omitempty"`
					Type *int64 `json:"type,omitempty"`
				}{
					Code: &testSecurityGroupRuleICMPCode,
					Type: &testSecurityGroupRuleICMPType,
				},
				Network:       &testSecurityGroupRuleNetwork,
				Protocol:      oapi.AddRuleToSecurityGroupJSONBodyProtocol(testSecurityGroupRuleProtocol),
				SecurityGroup: &oapi.SecurityGroupResource{Id: testSecurityGroupRuleSecurityGroupID},
				StartPort:     func() *int64 { v := int64(testSecurityGroupRuleStartPort); return &v }(),
			}, args.Get(2))
		}).
		Return(
			&oapi.AddRuleToSecurityGroupResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id: &testOperationID,
					Reference: &struct {
						Command *string `json:"command,omitempty"`
						Id      *string `json:"id,omitempty"` // revive:disable-line
						Link    *string `json:"link,omitempty"`
					}{Id: &testSecurityGroupID},
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
		}{Id: &testSecurityGroupID},
		State: &testOperationState,
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
				Rules: &[]oapi.SecurityGroupRule{
					{
						Description:   &testSecurityGroupRuleDescription,
						EndPort:       func() *int64 { v := int64(testSecurityGroupRuleEndPort - 1); return &v }(),
						FlowDirection: &testSecurityGroupRuleFlowDirection,
						Icmp: &struct {
							Code *int64 `json:"code,omitempty"`
							Type *int64 `json:"type,omitempty"`
						}{
							Code: &testSecurityGroupRuleICMPCode,
							Type: &testSecurityGroupRuleICMPType,
						},
						Id:            &testAlreadyExistingSecurityGroupRuleID,
						Network:       &testSecurityGroupRuleNetwork,
						Protocol:      &testSecurityGroupRuleProtocol,
						SecurityGroup: &oapi.SecurityGroupResource{Id: testSecurityGroupRuleSecurityGroupID},
						StartPort:     func() *int64 { v := int64(testSecurityGroupRuleStartPort); return &v }(),
					},
				},
			},
		}, nil).
		Once()

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
				Rules: &[]oapi.SecurityGroupRule{
					{
						Description:   &testSecurityGroupRuleDescription,
						EndPort:       func() *int64 { v := int64(testSecurityGroupRuleEndPort - 1); return &v }(),
						FlowDirection: &testSecurityGroupRuleFlowDirection,
						Icmp: &struct {
							Code *int64 `json:"code,omitempty"`
							Type *int64 `json:"type,omitempty"`
						}{
							Code: &testSecurityGroupRuleICMPCode,
							Type: &testSecurityGroupRuleICMPType,
						},
						Id:            &testAlreadyExistingSecurityGroupRuleID,
						Network:       &testSecurityGroupRuleNetwork,
						Protocol:      &testSecurityGroupRuleProtocol,
						SecurityGroup: &oapi.SecurityGroupResource{Id: testSecurityGroupRuleSecurityGroupID},
						StartPort:     func() *int64 { v := int64(testSecurityGroupRuleStartPort); return &v }(),
					},
					{
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
					},
				},
			},
		}, nil)

	securityGroup := &SecurityGroup{
		Description: &testSecurityGroupDescription,
		ID:          &testSecurityGroupID,
		Name:        &testSecurityGroupName,
	}

	expected := SecurityGroupRule{
		Description:     &testSecurityGroupRuleDescription,
		EndPort:         &testSecurityGroupRuleEndPort,
		FlowDirection:   (*string)(&testSecurityGroupRuleFlowDirection),
		ICMPCode:        &testSecurityGroupRuleICMPCode,
		ICMPType:        &testSecurityGroupRuleICMPType,
		ID:              &testSecurityGroupRuleID,
		Network:         func() *net.IPNet { _, v, _ := net.ParseCIDR(testSecurityGroupRuleNetwork); return v }(),
		Protocol:        (*string)(&testSecurityGroupRuleProtocol),
		SecurityGroupID: &testSecurityGroupRuleSecurityGroupID,
		StartPort:       &testSecurityGroupRuleStartPort,
	}

	actual, err := ts.client.CreateSecurityGroupRule(context.Background(), testZone, securityGroup, &expected)
	ts.Require().NoError(err)
	ts.Require().Equal(&expected, actual)
}

func (ts *testSuite) TestClient_DeleteSecurityGroupRule() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		deleted            = false
	)

	ts.mock().
		On(
			"DeleteRuleFromSecurityGroupWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // ruleId
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testSecurityGroupID, args.Get(1))
			ts.Require().Equal(testSecurityGroupRuleID, args.Get(2))
			deleted = true
		}).
		Return(
			&oapi.DeleteRuleFromSecurityGroupResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id: &testOperationID,
					Reference: &struct {
						Command *string `json:"command,omitempty"`
						Id      *string `json:"id,omitempty"` // revive:disable-line
						Link    *string `json:"link,omitempty"`
					}{Id: &testSecurityGroupID},
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
		}{Id: &testSecurityGroupID},
		State: &testOperationState,
	})

	securityGroup := &SecurityGroup{
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
	}

	ts.Require().NoError(ts.client.DeleteSecurityGroupRule(
		context.Background(),
		testZone,
		securityGroup,
		securityGroup.Rules[0],
	))
	ts.Require().True(deleted)
}
