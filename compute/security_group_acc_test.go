// +build testacc

package compute

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"testing"

	egoerr "github.com/exoscale/egoscale/error"
	egoapi "github.com/exoscale/egoscale/internal/egoscale"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type securityGroupFixture struct {
	c   *Client
	req *egoapi.CreateSecurityGroup
	res *egoapi.SecurityGroup
}

func newSecurityGroupFixture(c *Client, opts ...securityGroupFixtureOpt) *securityGroupFixture {
	var fixture = &securityGroupFixture{
		c:   c,
		req: &egoapi.CreateSecurityGroup{},
	}

	// Fixture default options
	for _, opt := range []securityGroupFixtureOpt{
		securityGroupFixtureOptName(testPrefix + "-" + testRandomString()),
		securityGroupFixtureOptDescription(testDescription),
	} {
		opt(fixture)
	}

	for _, opt := range opts {
		opt(fixture)
	}

	return fixture
}

func (f *securityGroupFixture) setup() (*securityGroupFixture, error) { // nolint:unused,deadcode
	res, err := f.c.c.RequestWithContext(f.c.ctx, f.req)
	if err != nil {
		return nil, f.c.csError(err)
	}
	f.res = res.(*egoapi.SecurityGroup)

	return f, nil
}

func (f *securityGroupFixture) teardown() error { // nolint:unused,deadcode
	_, err := f.c.c.RequestWithContext(f.c.ctx, &egoapi.DeleteSecurityGroup{ID: f.res.ID})
	return f.c.csError(err)
}

func (t *accTestSuite) withSecurityGroupFixture(f func(*securityGroupFixture), opts ...securityGroupFixtureOpt) {
	securityGroupFixture, err := newSecurityGroupFixture(t.client, opts...).setup()
	if err != nil {
		t.FailNow("Security Group fixture setup failed", err)
	}

	f(securityGroupFixture)
}

type securityGroupFixtureOpt func(*securityGroupFixture)

func securityGroupFixtureOptName(name string) securityGroupFixtureOpt { // nolint:unused,deadcode
	return func(f *securityGroupFixture) { f.req.Name = name }
}

func securityGroupFixtureOptDescription(description string) securityGroupFixtureOpt { // nolint:unused,deadcode
	return func(f *securityGroupFixture) { f.req.Description = description }
}

type securityGroupTestSuite struct {
	suite.Suite
	*accTestSuite
}

func (t *securityGroupTestSuite) SetupTest() {
	client, err := testClientFromEnv()
	if err != nil {
		t.FailNow("unable to initialize API client", err)
	}

	t.accTestSuite = &accTestSuite{
		Suite:  t.Suite,
		client: client,
	}
}

func (t *securityGroupTestSuite) TestCreateSecurityGroup() {
	var securityGroupName = testPrefix + "-" + testRandomString()

	securityGroup, err := t.client.CreateSecurityGroup(
		securityGroupName,
		&SecurityGroupCreateOpts{Description: testDescription},
	)
	if err != nil {
		t.FailNow("Security Group creation failed", err)
	}
	assert.NotEmpty(t.T(), securityGroup.ID)

	actualSecurityGroup := egoapi.SecurityGroup{}
	if err := json.Unmarshal(securityGroup.Raw(), &actualSecurityGroup); err != nil {
		t.FailNow("unable to unmarshal raw resource", err)
	}

	assert.Equal(t.T(), securityGroupName, actualSecurityGroup.Name)
	assert.Equal(t.T(), securityGroupName, securityGroup.Name)
	assert.Equal(t.T(), testDescription, actualSecurityGroup.Description)
	assert.Equal(t.T(), testDescription, securityGroup.Description)

	if _, err = t.client.c.RequestWithContext(t.client.ctx, &egoapi.DeleteSecurityGroup{
		ID: egoapi.MustParseUUID(securityGroup.ID)}); err != nil {
		t.FailNow("Security Group deletion failed", err)
	}
}

func (t *securityGroupTestSuite) TestListSecurityGroups() {
	t.withSecurityGroupFixture(func(securityGroupFixture *securityGroupFixture) {
		defer securityGroupFixture.teardown() // nolint:errcheck

		securityGroups, err := t.client.ListSecurityGroups()
		if err != nil {
			t.FailNow("Security Groups listing failed", err)
		}

		// We cannot guarantee that there will be only our resources in the
		// testing environment, so we ensure we get at least our fixture SG + the default SG
		assert.GreaterOrEqual(t.T(), len(securityGroups), 2)
	})
}

func (t *securityGroupTestSuite) TestGetSecurityGroupByID() {
	t.withSecurityGroupFixture(func(securityGroupFixture *securityGroupFixture) {
		defer securityGroupFixture.teardown() // nolint:errcheck

		securityGroup, err := t.client.GetSecurityGroupByID(securityGroupFixture.res.ID.String())
		if err != nil {
			t.FailNow("Security Group retrieval by ID failed", err)
		}
		assert.Equal(t.T(), securityGroupFixture.res.ID.String(), securityGroup.ID)

		securityGroup, err = t.client.GetSecurityGroupByID("00000000-0000-0000-0000-000000000000")
		assert.EqualError(t.T(), err, egoerr.ErrResourceNotFound.Error())
		assert.Empty(t.T(), securityGroup)
	})
}

func (t *securityGroupTestSuite) TestGetSecurityGroupByName() {
	var securityGroupName = testPrefix + "-" + testRandomString()

	t.withSecurityGroupFixture(func(securityGroupFixture *securityGroupFixture) {
		defer securityGroupFixture.teardown() // nolint:errcheck

		securityGroup, err := t.client.GetSecurityGroupByName(securityGroupFixture.res.Name)
		if err != nil {
			t.FailNow("Security Group retrieval by name failed", err)
		}
		assert.Equal(t.T(), securityGroupFixture.res.ID.String(), securityGroup.ID)

		securityGroup, err = t.client.GetSecurityGroupByName("lolnope")
		assert.EqualError(t.T(), err, egoerr.ErrResourceNotFound.Error())
		assert.Empty(t.T(), securityGroup)
	}, securityGroupFixtureOptName(securityGroupName))
}

func (t *securityGroupTestSuite) TestSecurityGroupAddRule() {
	t.withSecurityGroupFixture(func(securityGroupFixture *securityGroupFixture) {
		defer securityGroupFixture.teardown() // nolint:errcheck

		var (
			securityGroup       = t.client.securityGroupFromAPI(securityGroupFixture.res)
			testDescription     = "test-egoscale"
			testNetworkCIDR     = "1.1.1.1/32"
			testPortIngress     = "80-81"
			testPortEgress      = "53"
			testProtocolIngress = "tcp"
			testProtocolEgress  = "udp"
		)

		_, testCIDR, _ := net.ParseCIDR(testNetworkCIDR)

		for _, rule := range []SecurityGroupRule{
			SecurityGroupRule{
				Type:        "ingress",
				Description: testDescription,
				NetworkCIDR: testCIDR,
				Port:        testPortIngress,
				Protocol:    testProtocolIngress,
			},
			SecurityGroupRule{
				Type:          "egress",
				Description:   testDescription,
				SecurityGroup: &SecurityGroup{Name: "default"},
				Port:          testPortEgress,
				Protocol:      testProtocolEgress,
			},
		} {
			rule := rule
			if err := securityGroup.AddRule(&rule); err != nil {
				t.FailNow("Security Group rule adding failed", err)
			}
		}

		res, err := t.client.c.ListWithContext(t.client.ctx, &egoapi.SecurityGroup{Name: securityGroup.Name})
		if err != nil {
			t.FailNow("Security Group rules listing failed", err)
		}
		for _, item := range res {
			sg := item.(*egoapi.SecurityGroup)
			assert.Len(t.T(), sg.IngressRule, 1)
			assert.Len(t.T(), sg.EgressRule, 1)
			assert.Equal(t.T(), testDescription, sg.IngressRule[0].Description)
			assert.Equal(t.T(), testNetworkCIDR, sg.IngressRule[0].CIDR.String())
			assert.Equal(t.T(), strings.Split(testPortIngress, "-")[0], fmt.Sprint(sg.IngressRule[0].StartPort))
			assert.Equal(t.T(), strings.Split(testPortIngress, "-")[1], fmt.Sprint(sg.IngressRule[0].EndPort))
			assert.Equal(t.T(), testProtocolIngress, sg.IngressRule[0].Protocol)
			assert.Equal(t.T(), testDescription, sg.EgressRule[0].Description)
			assert.Equal(t.T(), "default", sg.EgressRule[0].SecurityGroupName)
			assert.Equal(t.T(), testPortEgress, fmt.Sprint(sg.EgressRule[0].StartPort))
			assert.Equal(t.T(), testPortEgress, fmt.Sprint(sg.EgressRule[0].EndPort))
			assert.Equal(t.T(), testProtocolEgress, sg.EgressRule[0].Protocol)
		}
	})
}

func (t *securityGroupTestSuite) TestSecurityGroupIngressRules() {
	t.withSecurityGroupFixture(func(securityGroupFixture *securityGroupFixture) {
		defer securityGroupFixture.teardown() // nolint:errcheck

		res, err := t.client.c.ListWithContext(t.client.ctx, &egoapi.SecurityGroup{Name: "default"})
		if err != nil {
			t.FailNow("unable to retrieve the default Security Group", err)
		}
		securityGroupDefault := res[0].(*egoapi.SecurityGroup)

		securityGroup := t.client.securityGroupFromAPI(securityGroupFixture.res)

		for _, rule := range []*egoapi.AuthorizeSecurityGroupIngress{
			&egoapi.AuthorizeSecurityGroupIngress{
				SecurityGroupName: securityGroup.Name,
				Description:       "test-egoscale",
				CIDRList:          []egoapi.CIDR{*egoapi.MustParseCIDR("0.0.0.0/0")},
				StartPort:         8000,
				EndPort:           9000,
				Protocol:          "tcp",
			},
			&egoapi.AuthorizeSecurityGroupIngress{
				SecurityGroupName:     securityGroup.Name,
				Description:           "test-egoscale",
				UserSecurityGroupList: []egoapi.UserSecurityGroup{securityGroupDefault.UserSecurityGroup()},
				Protocol:              "icmp",
				IcmpType:              8,
				IcmpCode:              0,
			},
		} {
			if _, err := t.client.c.RequestWithContext(t.client.ctx, rule); err != nil {
				t.FailNow("unable to add a test rule to the fixture Security group", err)
			}
		}

		rules, err := securityGroup.IngressRules()
		if err != nil {
			t.FailNow("Security Group ingress rules listing failed", err)
		}
		assert.Len(t.T(), rules, 2)
		assert.NotEmpty(t.T(), rules[0].ID)
		assert.Equal(t.T(), "ingress", rules[0].Type)
		assert.Equal(t.T(), "test-egoscale", rules[0].Description)
		assert.Equal(t.T(), "default", rules[0].SecurityGroup.Name)
		assert.Equal(t.T(), "icmp", rules[0].Protocol)
		assert.Equal(t.T(), uint8(8), rules[0].ICMPType)
		assert.Equal(t.T(), uint8(0), rules[0].ICMPCode)
		assert.Equal(t.T(), "0.0.0.0/0", rules[1].NetworkCIDR.String())
		assert.Equal(t.T(), "tcp", rules[1].Protocol)
		assert.Equal(t.T(), "8000-9000", rules[1].Port)
	})
}

func (t *securityGroupTestSuite) TestSecurityGroupEgressRules() {
	t.withSecurityGroupFixture(func(securityGroupFixture *securityGroupFixture) {
		defer securityGroupFixture.teardown() // nolint:errcheck

		securityGroup := t.client.securityGroupFromAPI(securityGroupFixture.res)

		for _, rule := range []*egoapi.AuthorizeSecurityGroupEgress{
			&egoapi.AuthorizeSecurityGroupEgress{
				SecurityGroupName: securityGroup.Name,
				Description:       "DNS",
				CIDRList:          []egoapi.CIDR{*egoapi.MustParseCIDR("0.0.0.0/0")},
				StartPort:         53,
				EndPort:           53,
				Protocol:          "tcp",
			},
			&egoapi.AuthorizeSecurityGroupEgress{
				SecurityGroupName: securityGroup.Name,
				Description:       "DNS",
				CIDRList:          []egoapi.CIDR{*egoapi.MustParseCIDR("0.0.0.0/0")},
				StartPort:         53,
				EndPort:           53,
				Protocol:          "udp",
			},
		} {
			if _, err := t.client.c.RequestWithContext(t.client.ctx, rule); err != nil {
				t.FailNow("unable to add a test rule to the fixture Security group", err)
			}
		}

		rules, err := securityGroup.EgressRules()
		if err != nil {
			t.FailNow("Security Group egress rules listing failed", err)
		}
		assert.Len(t.T(), rules, 2)
		assert.NotEmpty(t.T(), rules[0].ID)
		assert.Equal(t.T(), "egress", rules[0].Type)
		assert.Equal(t.T(), "DNS", rules[0].Description)
		assert.Equal(t.T(), "0.0.0.0/0", rules[0].NetworkCIDR.String())
		assert.Equal(t.T(), "53", rules[0].Port)
		assert.Equal(t.T(), "tcp", rules[0].Protocol)
		assert.Equal(t.T(), "udp", rules[1].Protocol)
	})
}

func (t *securityGroupTestSuite) TestSecurityGroupDelete() {
	t.withSecurityGroupFixture(func(securityGroupFixture *securityGroupFixture) {
		securityGroup := t.client.securityGroupFromAPI(securityGroupFixture.res)
		if err := securityGroup.Delete(); err != nil {
			t.FailNow("Security Group deletion failed", err)
		}
		assert.Empty(t.T(), securityGroup.ID)
		assert.Empty(t.T(), securityGroup.Name)
		assert.Empty(t.T(), securityGroup.Description)

		r, _ := t.client.c.ListWithContext(t.client.ctx, &egoapi.SecurityGroup{ID: securityGroupFixture.res.ID})
		assert.Len(t.T(), r, 0)
	})
}

func (t *securityGroupTestSuite) TestSecurityGroupRuleDelete() {
	t.withSecurityGroupFixture(func(securityGroupFixture *securityGroupFixture) {
		defer securityGroupFixture.teardown() // nolint:errcheck

		securityGroup := t.client.securityGroupFromAPI(securityGroupFixture.res)

		res, err := t.client.c.RequestWithContext(t.client.ctx, &egoapi.AuthorizeSecurityGroupIngress{
			SecurityGroupName: securityGroup.Name,
			CIDRList:          []egoapi.CIDR{*egoapi.MustParseCIDR("0.0.0.0/0")},
			StartPort:         22,
			EndPort:           22,
			Protocol:          "tcp",
		})
		if err != nil {
			t.FailNow("unable to add a test rule to the fixture Security group", err)
		}

		rule, err := t.client.securityGroupRuleFromAPI(&(res.(*egoapi.SecurityGroup).IngressRule[0]))
		if err != nil {
			t.FailNow("Security Group rule retrieval failed", err)
		}

		if err := rule.Delete(); err != nil {
			t.FailNow("Security Group rule deletion failed", err)
		}

		res3, err := t.client.c.ListWithContext(t.client.ctx, &egoapi.SecurityGroup{Name: securityGroup.Name})
		if err != nil {
			t.FailNow("Security Group rules listing failed", err)
		}
		for _, item := range res3 {
			sg := item.(*egoapi.SecurityGroup)
			assert.Len(t.T(), sg.IngressRule, 0)
		}
	})
}

func TestAccComputeSecurityGroupTestSuite(t *testing.T) {
	suite.Run(t, new(securityGroupTestSuite))
}
