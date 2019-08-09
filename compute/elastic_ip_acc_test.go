// +build testacc

package compute

import (
	"encoding/json"
	"testing"
	"time"

	egoerr "github.com/exoscale/egoscale/error"
	egoapi "github.com/exoscale/egoscale/internal/egoscale"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type elasticIPFixture struct {
	c   *Client
	req *egoapi.AssociateIPAddress
	res *egoapi.IPAddress
}

func newElasticIPFixture(c *Client, opts ...elasticIPFixtureOpt) *elasticIPFixture {
	var fixture = &elasticIPFixture{
		c:   c,
		req: &egoapi.AssociateIPAddress{ZoneID: egoapi.MustParseUUID(testZoneID)},
	}

	// Fixture default options
	for _, opt := range []elasticIPFixtureOpt{
		elasticIPFixtureOptZone(testZoneID),
	} {
		opt(fixture)
	}

	for _, opt := range opts {
		opt(fixture)
	}

	return fixture
}

func (f *elasticIPFixture) setup() (*elasticIPFixture, error) { // nolint:unused,deadcode
	res, err := f.c.c.RequestWithContext(f.c.ctx, f.req)
	if err != nil {
		return nil, f.c.csError(err)
	}
	f.res = res.(*egoapi.IPAddress)

	return f, nil
}

func (f *elasticIPFixture) teardown() error { // nolint:unused,deadcode
	_, err := f.c.c.RequestWithContext(f.c.ctx, &egoapi.DisassociateIPAddress{ID: f.res.ID})
	return f.c.csError(err)
}

type elasticIPFixtureOpt func(*elasticIPFixture)

func elasticIPFixtureOptZone(id string) elasticIPFixtureOpt { // nolint:unused,deadcode
	return func(f *elasticIPFixture) { f.req.ZoneID = egoapi.MustParseUUID(id) }
}

func elasticIPFixtureOptHealthcheckMode(mode string) elasticIPFixtureOpt { // nolint:unused,deadcode
	return func(f *elasticIPFixture) { f.req.HealthcheckMode = mode }
}

func elasticIPFixtureOptHealthcheckPort(port uint16) elasticIPFixtureOpt { // nolint:unused,deadcode
	return func(f *elasticIPFixture) { f.req.HealthcheckPort = int64(port) }
}

func elasticIPFixtureOptHealthcheckPath(path string) elasticIPFixtureOpt { // nolint:unused,deadcode
	return func(f *elasticIPFixture) { f.req.HealthcheckPath = path }
}

func elasticIPFixtureOptHealthcheckInterval(interval time.Duration) elasticIPFixtureOpt { // nolint:unused,deadcode
	return func(f *elasticIPFixture) { f.req.HealthcheckInterval = int64(interval.Seconds()) }
}

func elasticIPFixtureOptHealthcheckTimeout(timeout time.Duration) elasticIPFixtureOpt { // nolint:unused,deadcode
	return func(f *elasticIPFixture) { f.req.HealthcheckTimeout = int64(timeout.Seconds()) }
}

func elasticIPFixtureOptHealthcheckStrikesOK(n int) elasticIPFixtureOpt { // nolint:unused,deadcode
	return func(f *elasticIPFixture) { f.req.HealthcheckStrikesOk = int64(n) }
}

func elasticIPFixtureOptHealthcheckStrikesFail(n int) elasticIPFixtureOpt { // nolint:unused,deadcode
	return func(f *elasticIPFixture) { f.req.HealthcheckStrikesFail = int64(n) }
}

func (t *accTestSuite) withElasticIPFixture(f func(*elasticIPFixture), opts ...elasticIPFixtureOpt) {
	elasticIPFixture, err := newElasticIPFixture(t.client, opts...).setup()
	if err != nil {
		t.FailNow("Elastic IP fixture setup failed", err)
	}

	f(elasticIPFixture)
}

type elasticIPTestSuite struct {
	suite.Suite
	*accTestSuite
}

func (t *elasticIPTestSuite) SetupTest() {
	client, err := testClientFromEnv()
	if err != nil {
		t.FailNow("unable to initialize API client", err)
	}

	t.accTestSuite = &accTestSuite{
		Suite:  t.Suite,
		client: client,
	}
}

func (t *elasticIPTestSuite) TestCreateElasticIP() {
	t.withZoneFixture(func(zoneFixture *zoneFixture) {
		var (
			zone                          = t.client.zoneFromAPI(zoneFixture.res)
			healthcheckMode               = "http"
			healthcheckPort        uint16 = 80
			healthcheckPath               = "/health"
			healthcheckInterval           = 5 * time.Second
			healthcheckTimeout            = 3 * time.Second
			healthcheckStrikesOK          = 2
			healthcheckStrikesFail        = 1
		)

		elasticIP, err := t.client.CreateElasticIP(
			zone,
			&ElasticIPCreateOpts{
				HealthcheckMode:        healthcheckMode,
				HealthcheckPort:        healthcheckPort,
				HealthcheckPath:        healthcheckPath,
				HealthcheckInterval:    healthcheckInterval,
				HealthcheckTimeout:     healthcheckTimeout,
				HealthcheckStrikesOK:   healthcheckStrikesOK,
				HealthcheckStrikesFail: healthcheckStrikesFail,
			})
		if err != nil {
			t.FailNow("Elastic IP creation failed", err)
		}
		assert.NotEmpty(t.T(), elasticIP.ID)

		actualElasticIP := egoapi.IPAddress{}
		if err := json.Unmarshal(elasticIP.Raw(), &actualElasticIP); err != nil {
			t.FailNow("unable to unmarshal raw resource", err)
		}

		assert.Equal(t.T(), zone.ID, elasticIP.Zone.ID)
		assert.Equal(t.T(), actualElasticIP.IPAddress, elasticIP.Address)
		assert.Equal(t.T(), actualElasticIP.Healthcheck.Mode, healthcheckMode)
		assert.Equal(t.T(), elasticIP.HealthcheckMode, healthcheckMode)
		assert.Equal(t.T(), uint16(actualElasticIP.Healthcheck.Port), healthcheckPort)
		assert.Equal(t.T(), elasticIP.HealthcheckPort, healthcheckPort)
		assert.Equal(t.T(), actualElasticIP.Healthcheck.Path, healthcheckPath)
		assert.Equal(t.T(), elasticIP.HealthcheckPath, healthcheckPath)

		if _, err = t.client.c.RequestWithContext(t.client.ctx, &egoapi.DisassociateIPAddress{
			ID: egoapi.MustParseUUID(elasticIP.ID),
		}); err != nil {
			t.FailNow("Elastic IP deletion failed", err)
		}
	})
}

func (t *elasticIPTestSuite) TestListElasticIPs() {
	t.withZoneFixture(func(zoneFixture *zoneFixture) {

		t.withElasticIPFixture(func(elasticIPFixture *elasticIPFixture) {
			defer elasticIPFixture.teardown() // nolint:errcheck

			zone := t.client.zoneFromAPI(zoneFixture.res)

			elasticIPs, err := t.client.ListElasticIPs(zone)
			if err != nil {
				t.FailNow("Elastic IPs listing failed", err)
			}

			// We cannot guarantee that there will be only our resources in the
			// testing environment, so we ensure we get at least our fixture EIP
			assert.GreaterOrEqual(t.T(), len(elasticIPs), 1)
		})
	})
}

func (t *elasticIPTestSuite) TestGetElasticIPByID() {
	t.withZoneFixture(func(zoneFixture *zoneFixture) {

		t.withElasticIPFixture(func(elasticIPFixture *elasticIPFixture) {
			defer elasticIPFixture.teardown() // nolint:errcheck

			zone := t.client.zoneFromAPI(zoneFixture.res)
			elasticIP, err := t.client.GetElasticIPByID(zone, elasticIPFixture.res.ID.String())
			if err != nil {
				t.FailNow("Elastic IP retrieval failed", err)
			}
			assert.Equal(t.T(), elasticIPFixture.res.ID.String(), elasticIP.ID)

			elasticIP, err = t.client.GetElasticIPByID(zone, "00000000-0000-0000-0000-000000000000")
			assert.EqualError(t.T(), err, egoerr.ErrResourceNotFound.Error())
			assert.Empty(t.T(), elasticIP)
		})
	})
}

func (t *elasticIPTestSuite) TestGetElasticIPByAddress() {
	t.withZoneFixture(func(zoneFixture *zoneFixture) {

		t.withElasticIPFixture(func(elasticIPFixture *elasticIPFixture) {
			defer elasticIPFixture.teardown() // nolint:errcheck

			zone := t.client.zoneFromAPI(zoneFixture.res)
			elasticIP, err := t.client.GetElasticIPByAddress(zone, elasticIPFixture.res.IPAddress.String())
			if err != nil {
				t.FailNow("Elastic IP retrieval failed", err)
			}
			assert.Equal(t.T(), elasticIPFixture.res.ID.String(), elasticIP.ID)

			elasticIP, err = t.client.GetElasticIPByAddress(zone, "0.0.0.0")
			assert.EqualError(t.T(), err, egoerr.ErrResourceNotFound.Error())
			assert.Empty(t.T(), elasticIP)
		})
	})
}

func (t *elasticIPTestSuite) TestElasticIPReverseDNS() {
	t.withElasticIPFixture(func(elasticIPFixture *elasticIPFixture) {
		defer elasticIPFixture.teardown() // nolint:errcheck

		elasticIP, err := t.client.elasticIPFromAPI(elasticIPFixture.res)
		if err != nil {
			t.FailNow("Elastic IP fixture setup failed", err)
		}

		if _, err := t.client.c.RequestWithContext(t.client.ctx, &egoapi.UpdateReverseDNSForPublicIPAddress{
			ID:         egoapi.MustParseUUID(elasticIP.ID),
			DomainName: testReverseDNS,
		}); err != nil {
			t.FailNow("Elastic IP reverse DNS record update failed", err)
		}

		reverseDNS, err := elasticIP.ReverseDNS()
		if err != nil {
			t.FailNow("Elastic IP reverse DNS retrieval failed", err)
		}

		assert.Equal(t.T(), testReverseDNS, reverseDNS)
	})
}

func (t *elasticIPTestSuite) TestElasticIPSetReverseDNS() {
	t.withElasticIPFixture(func(elasticIPFixture *elasticIPFixture) {
		defer elasticIPFixture.teardown() // nolint:errcheck

		elasticIP, err := t.client.elasticIPFromAPI(elasticIPFixture.res)
		if err != nil {
			t.FailNow("Elastic IP fixture setup failed", err)
		}

		if err := elasticIP.SetReverseDNS(testReverseDNS); err != nil {
			t.FailNow("Elastic IP reverse DNS setting failed", err)
		}

		res2, err := t.client.c.RequestWithContext(t.client.ctx, &egoapi.QueryReverseDNSForPublicIPAddress{
			ID: egoapi.MustParseUUID(elasticIP.ID),
		})
		if err != nil {
			t.FailNow("Elastic IP reverse DNS record retrieval failed", err)
		}
		actualElasticIP := res2.(*egoapi.IPAddress)

		assert.Len(t.T(), actualElasticIP.ReverseDNS, 1)
		assert.Equal(t.T(), testReverseDNS, actualElasticIP.ReverseDNS[0].DomainName)
	})
}

func (t *elasticIPTestSuite) TestElasticIPUnsetReverseDNS() {
	t.withElasticIPFixture(func(elasticIPFixture *elasticIPFixture) {
		defer elasticIPFixture.teardown() // nolint:errcheck

		elasticIP, err := t.client.elasticIPFromAPI(elasticIPFixture.res)
		if err != nil {
			t.FailNow("Elastic IP fixture setup failed", err)
		}

		if _, err := t.client.c.RequestWithContext(t.client.ctx, &egoapi.UpdateReverseDNSForPublicIPAddress{
			ID:         egoapi.MustParseUUID(elasticIP.ID),
			DomainName: testReverseDNS,
		}); err != nil {
			t.FailNow("Elastic IP reverse DNS record update failed", err)
		}

		if err := elasticIP.UnsetReverseDNS(); err != nil {
			t.FailNow("Elastic IP reverse DNS unsetting failed", err)
		}

		res2, err := t.client.c.RequestWithContext(t.client.ctx, &egoapi.QueryReverseDNSForPublicIPAddress{
			ID: egoapi.MustParseUUID(elasticIP.ID),
		})
		if err != nil {
			t.FailNow("unable to retrieve Elastic IP reverse DNS", err)
		}
		actualElasticIP := res2.(*egoapi.IPAddress)

		assert.Empty(t.T(), actualElasticIP.ReverseDNS)
	})
}

func (t *elasticIPTestSuite) TestElasticIPInstances() {
	t.withElasticIPFixture(func(elasticIPFixture *elasticIPFixture) {
		defer elasticIPFixture.teardown() // nolint:errcheck

		t.withInstanceFixture(func(instanceFixture *instanceFixture) {
			defer instanceFixture.teardown() // nolint:errcheck

			elasticIP, err := t.client.elasticIPFromAPI(elasticIPFixture.res)
			if err != nil {
				t.FailNow("Elastic IP fixture setup failed", err)
			}
			instance, err := t.client.instanceFromAPI(instanceFixture.res)
			if err != nil {
				t.FailNow("Compute instance fixture setup failed", err)
			}

			if _, err := t.client.c.RequestWithContext(t.client.ctx, &egoapi.AddIPToNic{
				NicID:     instanceFixture.res.DefaultNic().ID,
				IPAddress: elasticIPFixture.res.IPAddress,
			}); err != nil {
				t.FailNow("unable to attach Elastic IP to Compute instance fixture", err)
			}

			elasticIPs, err := elasticIP.Instances()
			if err != nil {
				t.FailNow("Elastic IP attached Compute instances retrieval failed", err)
			}

			assert.Len(t.T(), elasticIPs, 1)
			assert.Equal(t.T(), instance.ID, elasticIPs[0].ID)
		})
	})
}

func (t *elasticIPTestSuite) TestElasticIPAttachInstance() {
	t.withElasticIPFixture(func(elasticIPFixture *elasticIPFixture) {
		defer elasticIPFixture.teardown() // nolint:errcheck

		t.withInstanceFixture(func(instanceFixture *instanceFixture) {
			defer instanceFixture.teardown() // nolint:errcheck

			elasticIP, err := t.client.elasticIPFromAPI(elasticIPFixture.res)
			if err != nil {
				t.FailNow("Elastic IP fixture setup failed", err)
			}
			instance, err := t.client.instanceFromAPI(instanceFixture.res)
			if err != nil {
				t.FailNow("Compute instance fixture setup failed", err)
			}

			if err := elasticIP.AttachInstance(instance); err != nil {
				t.FailNow("Compute instance Elastic IP attachment failed", err)
			}

			res, err := t.client.c.RequestWithContext(t.client.ctx, &egoapi.ListNics{
				NicID: instanceFixture.res.DefaultNic().ID,
			})
			if err != nil {
				t.FailNow("unable to retrieve Compute instance fixture NIC", err)
			}
			instanceNICs := res.(*egoapi.ListNicsResponse).Nic
			assert.Len(t.T(), instanceNICs, 1)
			assert.Len(t.T(), instanceNICs[0].SecondaryIP, 1)
			assert.True(t.T(), instanceNICs[0].SecondaryIP[0].IPAddress.Equal(elasticIPFixture.res.IPAddress))
		})
	})
}

func (t *elasticIPTestSuite) TestElasticIPDetachInstance() {
	t.withElasticIPFixture(func(elasticIPFixture *elasticIPFixture) {
		defer elasticIPFixture.teardown() // nolint:errcheck

		t.withInstanceFixture(func(instanceFixture *instanceFixture) {
			defer instanceFixture.teardown() // nolint:errcheck

			elasticIP, err := t.client.elasticIPFromAPI(elasticIPFixture.res)
			if err != nil {
				t.FailNow("Elastic IP fixture setup failed", err)
			}
			instance, err := t.client.instanceFromAPI(instanceFixture.res)
			if err != nil {
				t.FailNow("Compute instance fixture setup failed", err)
			}

			if _, err := t.client.c.RequestWithContext(t.client.ctx, &egoapi.AddIPToNic{
				NicID:     instanceFixture.res.DefaultNic().ID,
				IPAddress: elasticIPFixture.res.IPAddress,
			}); err != nil {
				t.FailNow("unable to attach Elastic IP to Compute instance fixture", err)
			}

			if err := elasticIP.DetachInstance(instance); err != nil {
				t.FailNow("Compute instance Elastic IP attachment failed", err)
			}

			res, err := t.client.c.RequestWithContext(t.client.ctx, &egoapi.ListNics{
				NicID: instanceFixture.res.DefaultNic().ID,
			})
			if err != nil {
				t.FailNow("unable to retrieve Compute instance fixture NIC", err)
			}
			instanceNICs := res.(*egoapi.ListNicsResponse).Nic
			assert.Len(t.T(), instanceNICs, 1)
			assert.Len(t.T(), instanceNICs[0].SecondaryIP, 0)
		})
	})
}

func (t *elasticIPTestSuite) TestElasticIPUpdate() {
	t.withZoneFixture(func(zoneFixture *zoneFixture) {

		t.withElasticIPFixture(func(elasticIPFixture *elasticIPFixture) {
			defer elasticIPFixture.teardown() // nolint:errcheck

			var (
				zone                          = t.client.zoneFromAPI(zoneFixture.res)
				healthcheckMode               = "http"
				healthcheckPort        uint16 = 80
				healthcheckPath               = "/health"
				healthcheckInterval           = 5 * time.Second
				healthcheckTimeout            = 3 * time.Second
				healthcheckStrikesOK          = 2
				healthcheckStrikesFail        = 1
			)

			elasticIP, err := t.client.elasticIPFromAPI(elasticIPFixture.res)
			if err != nil {
				t.FailNow("Elastic IP fixture setup failed", err)
			}

			if err := elasticIP.Update(&ElasticIPUpdateOpts{
				HealthcheckMode:        healthcheckMode,
				HealthcheckPort:        healthcheckPort,
				HealthcheckPath:        healthcheckPath,
				HealthcheckInterval:    healthcheckInterval,
				HealthcheckTimeout:     healthcheckTimeout,
				HealthcheckStrikesOK:   healthcheckStrikesOK,
				HealthcheckStrikesFail: healthcheckStrikesFail,
			}); err != nil {
				t.FailNow("Elastic IP update failed", err)
			}

			res, err := t.client.c.ListWithContext(t.client.ctx, &egoapi.IPAddress{
				ID:     egoapi.MustParseUUID(elasticIP.ID),
				ZoneID: egoapi.MustParseUUID(zone.ID),
			})
			if err != nil || len(res) == 0 {
				t.FailNow("Elastic IP retrieval failed", err)
			}
			actualElasticIP := res[0].(*egoapi.IPAddress)

			assert.NotNil(t.T(), actualElasticIP.Healthcheck)
			assert.Equal(t.T(), healthcheckMode, actualElasticIP.Healthcheck.Mode)
			assert.Equal(t.T(), healthcheckMode, elasticIP.HealthcheckMode)
			assert.Equal(t.T(), healthcheckPort, uint16(actualElasticIP.Healthcheck.Port))
			assert.Equal(t.T(), healthcheckPort, elasticIP.HealthcheckPort)
			assert.Equal(t.T(), healthcheckPath, actualElasticIP.Healthcheck.Path)
			assert.Equal(t.T(), healthcheckPath, elasticIP.HealthcheckPath)
			assert.Equal(t.T(), healthcheckInterval, time.Duration(actualElasticIP.Healthcheck.Interval)*time.Second)
			assert.Equal(t.T(), healthcheckInterval, elasticIP.HealthcheckInterval)
			assert.Equal(t.T(), healthcheckTimeout, time.Duration(actualElasticIP.Healthcheck.Timeout)*time.Second)
			assert.Equal(t.T(), healthcheckTimeout, elasticIP.HealthcheckTimeout)
			assert.Equal(t.T(), healthcheckStrikesOK, int(actualElasticIP.Healthcheck.StrikesOk))
			assert.Equal(t.T(), healthcheckStrikesOK, elasticIP.HealthcheckStrikesOK)
			assert.Equal(t.T(), healthcheckStrikesFail, int(actualElasticIP.Healthcheck.StrikesFail))
			assert.Equal(t.T(), healthcheckStrikesFail, elasticIP.HealthcheckStrikesFail)
		})
	})
}

func (t *elasticIPTestSuite) TestElasticIPDelete() {
	t.withZoneFixture(func(zoneFixture *zoneFixture) {

		t.withElasticIPFixture(func(elasticIPFixture *elasticIPFixture) {
			elasticIP, err := t.client.elasticIPFromAPI(elasticIPFixture.res)
			if err != nil {
				t.FailNow("Elastic IP fixture setup failed", err)
			}

			if err = elasticIP.Delete(); err != nil {
				t.FailNow("Elastic IP deletion failed", err)
			}
			assert.Empty(t.T(), elasticIP.ID)
			assert.Empty(t.T(), elasticIP.Address)
			assert.Empty(t.T(), elasticIP.Zone)
			assert.Empty(t.T(), elasticIP.HealthcheckMode)
			assert.Empty(t.T(), elasticIP.HealthcheckPort)
			assert.Empty(t.T(), elasticIP.HealthcheckPath)
			assert.Empty(t.T(), elasticIP.HealthcheckInterval)
			assert.Empty(t.T(), elasticIP.HealthcheckTimeout)
			assert.Empty(t.T(), elasticIP.HealthcheckStrikesOK)
			assert.Empty(t.T(), elasticIP.HealthcheckStrikesFail)

			r, _ := t.client.c.ListWithContext(t.client.ctx, &egoapi.IPAddress{ID: elasticIPFixture.res.ID})
			assert.Len(t.T(), r, 0)
		})
	})
}

func TestAccComputeElasticIPTestSuite(t *testing.T) {
	suite.Run(t, new(elasticIPTestSuite))
}
