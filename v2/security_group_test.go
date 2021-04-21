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
	testSecurityGroupDescription         = "Test Security Group description"
	testSecurityGroupID                  = new(clientTestSuite).randomID()
	testSecurityGroupName                = "test-security-group"
	testSecurityGroupRuleDescription     = "Test rule"
	testSecurityGroupRuleEndPort         = 8080
	testSecurityGroupRuleFlowDirection   = "ingress"
	testSecurityGroupRuleICMPCode        = 0
	testSecurityGroupRuleICMPType        = 8
	testSecurityGroupRuleID              = new(clientTestSuite).randomID()
	testSecurityGroupRuleNetwork         = "1.2.3.0/24"
	testSecurityGroupRuleProtocol        = "icmp"
	testSecurityGroupRuleSecurityGroupID = new(clientTestSuite).randomID()
	testSecurityGroupRuleStartPort       = 8081
)

func (ts *clientTestSuite) TestSecurityGroup_get() {
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
				Code: func() *int64 { v := int64(testSecurityGroupRuleICMPCode); return &v }(),
				Type: func() *int64 { v := int64(testSecurityGroupRuleICMPType); return &v }(),
			},
			Id:            &testSecurityGroupRuleID,
			Network:       &testSecurityGroupRuleNetwork,
			Protocol:      &testSecurityGroupRuleProtocol,
			SecurityGroup: &papi.SecurityGroupResource{Id: &testSecurityGroupRuleSecurityGroupID},
			StartPort:     func() *int64 { v := int64(testSecurityGroupRuleStartPort); return &v }(),
		}},
	})

	expected := &SecurityGroup{
		Description: testSecurityGroupDescription,
		ID:          testSecurityGroupID,
		Name:        testSecurityGroupName,
		Rules: []*SecurityGroupRule{{
			Description:     testSecurityGroupRuleDescription,
			EndPort:         uint16(testSecurityGroupRuleEndPort),
			FlowDirection:   testSecurityGroupRuleFlowDirection,
			ICMPCode:        uint8(testSecurityGroupRuleICMPCode),
			ICMPType:        uint8(testSecurityGroupRuleICMPType),
			ID:              testSecurityGroupRuleID,
			Network:         func() *net.IPNet { _, v, _ := net.ParseCIDR(testSecurityGroupRuleNetwork); return v }(),
			Protocol:        testSecurityGroupRuleProtocol,
			SecurityGroupID: testSecurityGroupRuleSecurityGroupID,
			StartPort:       uint16(testSecurityGroupRuleStartPort),
		}},

		c:    ts.client,
		zone: testZone,
	}

	actual, err := new(SecurityGroup).get(context.Background(), ts.client, testZone, expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestSecurityGroup_AddRule() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = "success"
	)

	httpmock.RegisterResponder("POST", fmt.Sprintf("/security-group/%s/rules", testSecurityGroupID),
		func(req *http.Request) (*http.Response, error) {
			var actual papi.AddRuleToSecurityGroupJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.AddRuleToSecurityGroupJSONRequestBody{
				Description:   &testSecurityGroupRuleDescription,
				EndPort:       func() *int64 { v := int64(testSecurityGroupRuleEndPort); return &v }(),
				FlowDirection: testSecurityGroupRuleFlowDirection,
				Icmp: &struct {
					Code *int64 `json:"code,omitempty"`
					Type *int64 `json:"type,omitempty"`
				}{
					Code: func() *int64 { v := int64(testSecurityGroupRuleICMPCode); return &v }(),
					Type: func() *int64 { v := int64(testSecurityGroupRuleICMPType); return &v }(),
				},
				Network:       &testSecurityGroupRuleNetwork,
				Protocol:      testSecurityGroupRuleProtocol,
				SecurityGroup: &papi.SecurityGroupResource{Id: &testSecurityGroupRuleSecurityGroupID},
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
				Code: func() *int64 { v := int64(testSecurityGroupRuleICMPCode); return &v }(),
				Type: func() *int64 { v := int64(testSecurityGroupRuleICMPType); return &v }(),
			},
			Id:            &testSecurityGroupRuleID,
			Network:       &testSecurityGroupRuleNetwork,
			Protocol:      &testSecurityGroupRuleProtocol,
			SecurityGroup: &papi.SecurityGroupResource{Id: &testSecurityGroupRuleSecurityGroupID},
			StartPort:     func() *int64 { v := int64(testSecurityGroupRuleStartPort); return &v }(),
		}},
	})

	securityGroup := &SecurityGroup{
		Description: testSecurityGroupDescription,
		ID:          testSecurityGroupID,
		Name:        testSecurityGroupName,

		c: ts.client,
	}

	expected := SecurityGroupRule{
		Description:     testSecurityGroupRuleDescription,
		EndPort:         uint16(testSecurityGroupRuleEndPort),
		FlowDirection:   testSecurityGroupRuleFlowDirection,
		ICMPCode:        uint8(testSecurityGroupRuleICMPCode),
		ICMPType:        uint8(testSecurityGroupRuleICMPType),
		ID:              testSecurityGroupRuleID,
		Network:         func() *net.IPNet { _, v, _ := net.ParseCIDR(testSecurityGroupRuleNetwork); return v }(),
		Protocol:        testSecurityGroupRuleProtocol,
		SecurityGroupID: testSecurityGroupRuleSecurityGroupID,
		StartPort:       uint16(testSecurityGroupRuleStartPort),
	}

	actual, err := securityGroup.AddRule(context.Background(), &expected)
	ts.Require().NoError(err)
	ts.Require().Equal(&expected, actual)
}

func (ts *clientTestSuite) TestSecurityGroup_DeleteRule() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = "success"
	)

	httpmock.RegisterResponder("DELETE",
		fmt.Sprintf("=~^/security-group/%s/rules/.*", testSecurityGroupID),
		func(req *http.Request) (*http.Response, error) {
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
		Description: testSecurityGroupDescription,
		ID:          testSecurityGroupID,
		Name:        testSecurityGroupName,
		Rules: []*SecurityGroupRule{{
			Description:     testSecurityGroupRuleDescription,
			EndPort:         uint16(testSecurityGroupRuleEndPort),
			FlowDirection:   testSecurityGroupRuleFlowDirection,
			ICMPCode:        uint8(testSecurityGroupRuleICMPCode),
			ICMPType:        uint8(testSecurityGroupRuleICMPType),
			ID:              testSecurityGroupRuleID,
			Network:         func() *net.IPNet { _, v, _ := net.ParseCIDR(testSecurityGroupRuleNetwork); return v }(),
			Protocol:        testSecurityGroupRuleProtocol,
			SecurityGroupID: testSecurityGroupRuleSecurityGroupID,
			StartPort:       uint16(testSecurityGroupRuleStartPort),
		}},

		c: ts.client,
	}

	ts.Require().NoError(securityGroup.DeleteRule(context.Background(), securityGroup.Rules[0]))
}

func (ts *clientTestSuite) TestClient_CreateSecurityGroup() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = "success"
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
		Description: testSecurityGroupDescription,
		ID:          testSecurityGroupID,
		Name:        testSecurityGroupName,

		c:    ts.client,
		zone: testZone,
	}

	actual, err := ts.client.CreateSecurityGroup(context.Background(), testZone, &SecurityGroup{
		Description: testSecurityGroupDescription,
		Name:        testSecurityGroupName,
	})
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
					Code: func() *int64 { v := int64(testSecurityGroupRuleICMPCode); return &v }(),
					Type: func() *int64 { v := int64(testSecurityGroupRuleICMPType); return &v }(),
				},
				Id:            &testSecurityGroupRuleID,
				Network:       &testSecurityGroupRuleNetwork,
				Protocol:      &testSecurityGroupRuleProtocol,
				SecurityGroup: &papi.SecurityGroupResource{Id: &testSecurityGroupRuleSecurityGroupID},
				StartPort:     func() *int64 { v := int64(testSecurityGroupRuleStartPort); return &v }(),
			}},
		}},
	})

	expected := []*SecurityGroup{
		{
			Description: testSecurityGroupDescription,
			ID:          testSecurityGroupID,
			Name:        testSecurityGroupName,
			Rules: []*SecurityGroupRule{{
				Description:     testSecurityGroupRuleDescription,
				EndPort:         uint16(testSecurityGroupRuleEndPort),
				FlowDirection:   testSecurityGroupRuleFlowDirection,
				ICMPCode:        uint8(testSecurityGroupRuleICMPCode),
				ICMPType:        uint8(testSecurityGroupRuleICMPType),
				ID:              testSecurityGroupRuleID,
				Network:         func() *net.IPNet { _, v, _ := net.ParseCIDR(testSecurityGroupRuleNetwork); return v }(),
				Protocol:        testSecurityGroupRuleProtocol,
				SecurityGroupID: testSecurityGroupRuleSecurityGroupID,
				StartPort:       uint16(testSecurityGroupRuleStartPort),
			}},

			c:    ts.client,
			zone: testZone,
		},
	}

	actual, err := ts.client.ListSecurityGroups(context.Background(), testZone)
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
				Code: func() *int64 { v := int64(testSecurityGroupRuleICMPCode); return &v }(),
				Type: func() *int64 { v := int64(testSecurityGroupRuleICMPType); return &v }(),
			},
			Id:            &testSecurityGroupRuleID,
			Network:       &testSecurityGroupRuleNetwork,
			Protocol:      &testSecurityGroupRuleProtocol,
			SecurityGroup: &papi.SecurityGroupResource{Id: &testSecurityGroupRuleSecurityGroupID},
			StartPort:     func() *int64 { v := int64(testSecurityGroupRuleStartPort); return &v }(),
		}},
	})

	expected := &SecurityGroup{
		Description: testSecurityGroupDescription,
		ID:          testSecurityGroupID,
		Name:        testSecurityGroupName,
		Rules: []*SecurityGroupRule{{
			Description:     testSecurityGroupRuleDescription,
			EndPort:         uint16(testSecurityGroupRuleEndPort),
			FlowDirection:   testSecurityGroupRuleFlowDirection,
			ICMPCode:        uint8(testSecurityGroupRuleICMPCode),
			ICMPType:        uint8(testSecurityGroupRuleICMPType),
			ID:              testSecurityGroupRuleID,
			Network:         func() *net.IPNet { _, v, _ := net.ParseCIDR(testSecurityGroupRuleNetwork); return v }(),
			Protocol:        testSecurityGroupRuleProtocol,
			SecurityGroupID: testSecurityGroupRuleSecurityGroupID,
			StartPort:       uint16(testSecurityGroupRuleStartPort),
		}},

		c:    ts.client,
		zone: testZone,
	}

	actual, err := ts.client.GetSecurityGroup(context.Background(), testZone, expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_DeleteSecurityGroup() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = "success"
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

	ts.Require().NoError(ts.client.DeleteSecurityGroup(context.Background(), testZone, testSecurityGroupID))
	ts.Require().True(deleted)
}
