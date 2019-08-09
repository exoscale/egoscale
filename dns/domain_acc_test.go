// +build testacc

package dns

import (
	"testing"

	egoerr "github.com/exoscale/egoscale/error"
	egoapi "github.com/exoscale/egoscale/internal/egoscale"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type domainTestSuite struct {
	suite.Suite
	client *Client
}

func (t *domainTestSuite) SetupTest() {
	var err error

	if t.client, err = testClientFromEnv(); err != nil {
		t.FailNow("unable to initialize API client", err)
	}
}

func (t *domainTestSuite) TestCreateDomain() {
	// IDN Unicode/Punycode conversion courtesy of https://www.punycoder.com/
	var (
		randomString         = testRandomString()
		domainName           = "test-egoscale-" + randomString + ".net"
		unicodeName          = "égzoskèle-" + randomString + ".ch"
		unicodeNamePunycoded = "xn--gzoskle-" + randomString + "-vvbo.ch"
	)

	domain, err := t.client.CreateDomain(domainName)
	if err != nil {
		t.FailNow("domain creation failed", err)
	}
	assert.Greater(t.T(), domain.ID, int64(0))
	assert.Equal(t.T(), domainName, domain.Name)

	if _, err = t.client.c.Request(&egoapi.DeleteDNSDomain{Name: domain.Name}); err != nil {
		t.FailNow("domain deletion failed", err)
	}

	// IDN (Internationalized Domain Name)
	domain, err = t.client.CreateDomain(unicodeName)
	if err != nil {
		t.FailNow("Unicode domain creation failed", err)
	}
	assert.Greater(t.T(), domain.ID, int64(0))
	assert.Equal(t.T(), unicodeNamePunycoded, domain.Name)
	assert.Equal(t.T(), unicodeName, domain.UnicodeName)

	if _, err = t.client.c.Request(&egoapi.DeleteDNSDomain{Name: domain.Name}); err != nil {
		t.FailNow("domain deletion failed", err)
	}
}

func (t *domainTestSuite) TestListDomains() {
	var domainName = "test-egoscale-" + testRandomString() + ".net"

	_, teardown, err := domainFixture(domainName)
	if err != nil {
		t.FailNow("domain fixture setup failed", err)
	}
	defer teardown() // nolint:errcheck

	// We cannot guarantee that there will be only our resources,
	// so we ensure we get at least our fixture domain
	domains, err := t.client.ListDomains()
	if err != nil {
		t.FailNow("domains listing failed", err)
	}
	assert.GreaterOrEqual(t.T(), len(domains), 1)
}

func (t *domainTestSuite) TestGetDomainByID() {
	var domainName = "test-egoscale-" + testRandomString() + ".net"

	res, teardown, err := domainFixture(domainName)
	if err != nil {
		t.FailNow("domain fixture setup failed", err)
	}
	defer teardown() // nolint:errcheck

	domain, err := t.client.GetDomainByID(res.ID)
	if err != nil {
		t.FailNow("domain retrieval by ID failed", err)
	}
	assert.Equal(t.T(), domainName, domain.Name)

	domain, err = t.client.GetDomainByID(1)
	assert.EqualError(t.T(), err, egoerr.ErrResourceNotFound.Error())
	assert.Empty(t.T(), domain)
}

func (t *domainTestSuite) TestGetDomainByName() {
	var domainName = "test-egoscale-" + testRandomString() + ".net"

	_, teardown, err := domainFixture(domainName)
	if err != nil {
		t.FailNow("domain fixture setup failed", err)
	}
	defer teardown() // nolint:errcheck

	domain, err := t.client.GetDomainByName(domainName)
	if err != nil {
		t.FailNow("domain retrieval by name failed", err)
	}
	assert.Equal(t.T(), domainName, domain.Name)

	domain, err = t.client.GetDomainByName("lolnope")
	assert.EqualError(t.T(), err, egoerr.ErrResourceNotFound.Error())
	assert.Empty(t.T(), domain)
}

func (t *domainTestSuite) TestDomainAddRecord() {
	var (
		domainName     = "test-egoscale-" + testRandomString() + ".net"
		recordName     = "test-egoscale"
		recordType     = "MX"
		recordContent  = "mx1.example.net"
		recordPriority = 10
		recordTTL      = 1042
	)

	res, teardown, err := domainFixture(domainName)
	if err != nil {
		t.FailNow("domain fixture setup failed", err)
	}
	defer teardown() // nolint:errcheck
	domain := t.client.domainFromAPI(res)

	record, err := domain.AddRecord(&DomainRecordCreateOpts{
		Name:     recordName,
		Type:     recordType,
		Content:  recordContent,
		Priority: recordPriority,
		TTL:      recordTTL,
	})
	if err != nil {
		t.FailNow("domain record creation failed", err)
	}
	assert.Equal(t.T(), recordName, record.Name)
	assert.Equal(t.T(), recordType, record.Type)
	assert.Equal(t.T(), recordContent, record.Content)
	// assert.Equal(t.T(), recordPriority, record.Priority) // TODO: API bug, uncomment once fixed
	assert.Equal(t.T(), recordTTL, record.TTL)
}

func (t *domainTestSuite) TestDomainRecords() {
	var (
		domainName     = "test-egoscale-" + testRandomString() + ".net"
		recordName     = "test-egoscale"
		recordType     = "MX"
		recordContent  = "mx1.example.net"
		recordPriority = 10
		recordTTL      = 1042
	)

	res, teardown, err := domainFixture(domainName)
	if err != nil {
		t.FailNow("domain fixture setup failed", err)
	}
	defer teardown() // nolint:errcheck

	if _, err = t.client.c.Request(&egoapi.CreateDNSRecord{
		Domain:   res.Name,
		Name:     recordName,
		Type:     recordType,
		Content:  recordContent,
		Priority: recordPriority,
		TTL:      recordTTL,
	}); err != nil {
		t.FailNow("domain record fixture setup failed", err)
	}
	domain := t.client.domainFromAPI(res)

	records, err := domain.Records()
	if err != nil {
		t.FailNow("domain records listing failed", err)
	}
	assert.GreaterOrEqual(t.T(), len(records), 1)

	for _, record := range records {
		if record.Name == "" {
			continue
		}

		assert.Equal(t.T(), recordName, record.Name)
		assert.Equal(t.T(), recordType, record.Type)
		assert.Equal(t.T(), recordContent, record.Content)
		// assert.Equal(t.T(), recordPriority, record.Priority) // TODO: API bug, uncomment once fixed
		assert.Equal(t.T(), recordTTL, record.TTL)
	}
}

func (t *domainTestSuite) TestDomainDelete() {
	var domainName = "test-egoscale-" + testRandomString() + ".net"

	res, _, err := domainFixture(domainName)
	if err != nil {
		t.FailNow("domain fixture setup failed", err)
	}

	domain := t.client.domainFromAPI(res)
	if err = domain.Delete(); err != nil {
		t.FailNow("domain deletion failed", err)
	}
	assert.Equal(t.T(), int64(0), domain.ID)
	assert.Empty(t.T(), domain.Name)

	// We have to list all domains and check if our test domain isn't in the
	// results since there is no way to search for a specific domain via the API
	domains, err := t.client.c.ListWithContext(t.client.ctx, &egoapi.DNSDomain{})
	if err != nil {
		t.FailNow("domains listing failed", err)
	}
	for _, d := range domains {
		assert.NotEqualf(t.T(), d.(*egoapi.DNSDomain).Name, domainName, "domain %q not deleted", domainName)
	}
}

func TestAccDNSDomainTestSuite(t *testing.T) {
	suite.Run(t, new(domainTestSuite))
}
