package v2

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/jarcoal/httpmock"

	papi "github.com/exoscale/egoscale/v2/internal/public-api"
)

var (
	testSecurityGroupDescription                = new(clientTestSuite).randomString(10)
	testSecurityGroupID                         = new(clientTestSuite).randomID()
	testSecurityGroupName                       = new(clientTestSuite).randomString(10)
	testSecurityGroupRuleDescription            = new(clientTestSuite).randomString(10)
	testSecurityGroupRuleEndPort         uint16 = 8080
	testSecurityGroupRuleFlowDirection          = papi.SecurityGroupRuleFlowDirectionIngress
	testSecurityGroupRuleICMPCode        int64  = 0 // nolint:revive
	testSecurityGroupRuleICMPType        int64  = 8
	testSecurityGroupRuleID                     = new(clientTestSuite).randomID()
	testSecurityGroupRuleNetwork                = "1.2.3.0/24"
	_, testSecurityGroupRuleNetworkP, _         = net.ParseCIDR(testSecurityGroupRuleNetwork)
	testSecurityGroupRuleProtocol               = papi.SecurityGroupRuleProtocolIcmp
	testSecurityGroupRuleSecurityGroupID        = new(clientTestSuite).randomID()
	testSecurityGroupRuleStartPort       uint16 = 8081
)

func (ts *clientTestSuite) TestClient_CreateSecurityGroup() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = papi.OperationStateSuccess
	)

	httpmock.RegisterResponder("POST", "/security-group",
		func(req *http.Request) (*http.Response, error) {
			var actual papi.CreateSecurityGroupJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.CreateSecurityGroupJSONRequestBody{
				Description: &testSecurityGroupDescription,
				Name:        testSecurityGroupName,
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
		Reference: &papi.Reference{Id: &testSecurityGroupID},
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/security-group/%s", testSecurityGroupID), papi.SecurityGroup{
		Description: &testSecurityGroupDescription,
		Id:          &testSecurityGroupID,
		Name:        &testSecurityGroupName,
	})

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

func (ts *clientTestSuite) TestClient_CreateSecurityGroupRule() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = papi.OperationStateSuccess
	)

	httpmock.RegisterResponder("POST", fmt.Sprintf("/security-group/%s/rules", testSecurityGroupID),
		func(req *http.Request) (*http.Response, error) {
			var actual papi.AddRuleToSecurityGroupJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.AddRuleToSecurityGroupJSONRequestBody{
				Description:   &testSecurityGroupRuleDescription,
				EndPort:       func() *int64 { v := int64(testSecurityGroupRuleEndPort); return &v }(),
				FlowDirection: papi.AddRuleToSecurityGroupJSONBodyFlowDirection(testSecurityGroupRuleFlowDirection),
				Icmp: &struct {
					Code *int64 `json:"code,omitempty"`
					Type *int64 `json:"type,omitempty"`
				}{
					Code: &testSecurityGroupRuleICMPCode,
					Type: &testSecurityGroupRuleICMPType,
				},
				Network:       &testSecurityGroupRuleNetwork,
				Protocol:      papi.AddRuleToSecurityGroupJSONBodyProtocol(testSecurityGroupRuleProtocol),
				SecurityGroup: &papi.SecurityGroupResource{Id: testSecurityGroupRuleSecurityGroupID},
				StartPort:     func() *int64 { v := int64(testSecurityGroupRuleStartPort); return &v }(),
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
		Reference: &papi.Reference{Id: &testSecurityGroupID},
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/security-group/%s", testSecurityGroupID), papi.SecurityGroup{
		Description: &testSecurityGroupDescription,
		Id:          &testSecurityGroupID,
		Name:        &testSecurityGroupName,
		Rules: &[]papi.SecurityGroupRule{{
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
			SecurityGroup: &papi.SecurityGroupResource{Id: testSecurityGroupRuleSecurityGroupID},
			StartPort:     func() *int64 { v := int64(testSecurityGroupRuleStartPort); return &v }(),
		}},
	})

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

func (ts *clientTestSuite) TestClient_DeleteSecurityGroup() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = papi.OperationStateSuccess
		deleted            = false
	)

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("/security-group/%s", testSecurityGroupID),
		func(req *http.Request) (*http.Response, error) {
			deleted = true

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testSecurityGroupID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testSecurityGroupID},
	})

	ts.Require().NoError(ts.client.DeleteSecurityGroup(
		context.Background(),
		testZone,
		&SecurityGroup{ID: &testSecurityGroupID},
	))
	ts.Require().True(deleted)
}

func (ts *clientTestSuite) TestClient_DeleteSecurityGroupRule() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = papi.OperationStateSuccess
		deleted            = false
	)

	httpmock.RegisterResponder("DELETE",
		fmt.Sprintf("=~^/security-group/%s/rules/.*", testSecurityGroupID),
		func(req *http.Request) (*http.Response, error) {
			deleted = true

			ts.Require().Equal(
				fmt.Sprintf("/security-group/%s/rules/%s", testSecurityGroupID, testSecurityGroupRuleID),
				req.URL.String())

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testSecurityGroupID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testSecurityGroupID},
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

func (ts *clientTestSuite) TestClient_FindSecurityGroup() {
	ts.mockAPIRequest("GET", "/security-group", struct {
		SecurityGroups *[]papi.SecurityGroup `json:"security-groups,omitempty"`
	}{
		SecurityGroups: &[]papi.SecurityGroup{{
			Id:   &testSecurityGroupID,
			Name: &testSecurityGroupName,
		}},
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/security-group/%s", testSecurityGroupID), papi.SecurityGroup{
		Id:   &testSecurityGroupID,
		Name: &testSecurityGroupName,
	})

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

func (ts *clientTestSuite) TestClient_GetSecurityGroup() {
	ts.mockAPIRequest("GET", fmt.Sprintf("/security-group/%s", testSecurityGroupID), papi.SecurityGroup{
		Description: &testSecurityGroupDescription,
		Id:          &testSecurityGroupID,
		Name:        &testSecurityGroupName,
		Rules: &[]papi.SecurityGroupRule{{
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
			SecurityGroup: &papi.SecurityGroupResource{Id: testSecurityGroupRuleSecurityGroupID},
			StartPort:     func() *int64 { v := int64(testSecurityGroupRuleStartPort); return &v }(),
		}},
	})

	expected := &SecurityGroup{
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

	actual, err := ts.client.GetSecurityGroup(context.Background(), testZone, *expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_ListSecurityGroups() {
	ts.mockAPIRequest("GET", "/security-group", struct {
		SecurityGroups *[]papi.SecurityGroup `json:"security-groups,omitempty"`
	}{
		SecurityGroups: &[]papi.SecurityGroup{{
			Description: &testSecurityGroupDescription,
			Id:          &testSecurityGroupID,
			Name:        &testSecurityGroupName,
			Rules: &[]papi.SecurityGroupRule{{
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
				SecurityGroup: &papi.SecurityGroupResource{Id: testSecurityGroupRuleSecurityGroupID},
				StartPort:     func() *int64 { v := int64(testSecurityGroupRuleStartPort); return &v }(),
			}},
		}},
	})

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
