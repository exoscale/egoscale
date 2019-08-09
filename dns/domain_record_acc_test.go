// +build testacc

package dns

import (
	"testing"

	egoapi "github.com/exoscale/egoscale/internal/egoscale"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type domainRecordTestSuite struct {
	suite.Suite
	client *Client
}

func (t *domainRecordTestSuite) SetupTest() {
	var err error

	if t.client, err = testClientFromEnv(); err != nil {
		t.FailNow("unable to initialize API client", err)
	}
}

func (t *domainRecordTestSuite) TestDomainRecordUpdate() {
	var (
		domainName           = "test-egoscale-" + testRandomString() + ".net"
		recordName           = "test-egoscale"
		recordNameEdited     = "test-egoscale-edited"
		recordType           = "MX"
		recordContent        = "mx1.example.net"
		recordContentEdited  = "mx2.example.net"
		recordPriority       = 10
		recordPriorityEdited = 20
		recordTTL            = 1042
		recordTTLEdited      = 1043
	)

	domainRes, teardown, err := domainFixture(domainName)
	if err != nil {
		t.FailNow("domain fixture setup failed", err)
	}
	defer teardown() // nolint:errcheck
	domain := t.client.domainFromAPI(domainRes)

	recordRes, err := t.client.c.Request(&egoapi.CreateDNSRecord{
		Domain:   domain.Name,
		Name:     recordName,
		Type:     recordType,
		Content:  recordContent,
		Priority: recordPriority,
		TTL:      recordTTL,
	})
	if err != nil {
		t.FailNow("domain record fixture creation failed", err)
	}
	record := t.client.domainRecordFromAPI(recordRes.(*egoapi.DNSRecord), domain)

	err = record.Update(&DomainRecordUpdateOpts{
		Name:     recordNameEdited,
		Content:  recordContentEdited,
		Priority: recordPriorityEdited,
		TTL:      recordTTLEdited,
	})
	if err != nil {
		t.FailNow("domain record update failed", err)
	}
	assert.Equal(t.T(), recordNameEdited, record.Name)
	assert.Equal(t.T(), recordContentEdited, record.Content)
	assert.Equal(t.T(), recordPriorityEdited, record.Priority)
	assert.Equal(t.T(), recordTTLEdited, record.TTL)
}

func (t *domainRecordTestSuite) TestDomainRecordDelete() {
	var (
		domainName     = "test-egoscale-" + testRandomString() + ".net"
		recordName     = "test-egoscale"
		recordType     = "MX"
		recordContent  = "mx1.example.net"
		recordPriority = 10
		recordTTL      = 1042
	)

	domainRes, teardown, err := domainFixture(domainName)
	if err != nil {
		t.FailNow("domain fixture setup failed", err)
	}
	defer teardown() // nolint:errcheck
	domain := t.client.domainFromAPI(domainRes)

	recordRes, err := t.client.c.Request(&egoapi.CreateDNSRecord{
		Domain:   domain.Name,
		Name:     recordName,
		Type:     recordType,
		Content:  recordContent,
		Priority: recordPriority,
		TTL:      recordTTL,
	})
	if err != nil {
		t.FailNow("domain record fixture creation failed", err)
	}
	record := t.client.domainRecordFromAPI(recordRes.(*egoapi.DNSRecord), domain)

	if err = record.Delete(); err != nil {
		t.FailNow("domain record deletion failed", err)
	}
	assert.Equal(t.T(), int64(0), record.ID)
	assert.Empty(t.T(), record.Name)
	assert.Empty(t.T(), record.Type)
	assert.Empty(t.T(), record.Content)
	assert.Equal(t.T(), int(0), record.Priority)
	assert.Equal(t.T(), int(0), record.TTL)
	assert.Empty(t.T(), record.Domain)

	// We have to list all records and check if our test record isn't in the
	// results since there is no way to search for a specific record via the API
	records, err := t.client.c.ListWithContext(t.client.ctx, &egoapi.ListDNSRecords{DomainID: domain.ID})
	if err != nil {
		t.FailNow("domain records listing failed", err)
	}
	for _, d := range records {
		assert.NotEqualf(t.T(), d.(*egoapi.DNSRecord).Name, recordName, "domain record %q not deleted", recordName)
	}
}

func TestAccDNSDomainRecordRecordTestSuite(t *testing.T) {
	suite.Run(t, new(domainRecordTestSuite))
}
