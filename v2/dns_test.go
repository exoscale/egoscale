package v2

import (
	"context"
	"net/http"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/exoscale/egoscale/v2/oapi"
)

var (
	testDNSDomainCreatedAt, _ = time.Parse(iso8601Format, "2020-05-26T12:09:42Z")
	testDNSDomainID           = new(testSuite).randomID()
	testDNSDomainUnicodeName  = "test.domain"

	testDNSDomainRecordContent            = new(testSuite).randomString(10)
	testDNSDomainRecordCreatedAt, _       = time.Parse(iso8601Format, "2021-04-27T12:09:42Z")
	testDNSDomainRecordID                 = new(testSuite).randomID()
	testDNSDomainRecordName               = new(testSuite).randomString(10)
	testDNSDomainRecordPriority     int64 = 5
	testDNSDomainRecordTTL          int64 = 1800
	testDNSDomainRecordType               = "CNAME"
	testDNSDomainRecordUpdatedAt, _       = time.Parse(iso8601Format, "2021-06-27T12:09:42Z")
)

func (ts *testSuite) TestClient_ListDnsDomains() {
	domains := struct {
		DnsDomains *[]oapi.DnsDomain `json:"dns-domains,omitempty"` //nolint: revive
	}{
		&[]oapi.DnsDomain{
			oapi.DnsDomain{
				CreatedAt:   &testDNSDomainCreatedAt,
				Id:          &testDNSDomainID,
				UnicodeName: &testDNSDomainUnicodeName,
			},
		},
	}
	ts.mock().
		On(
			"ListDnsDomainsWithResponse",
			mock.Anything,                 // ctx
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(
			&oapi.ListDnsDomainsResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200:      &domains,
			},
			nil,
		)

	expected := []DNSDomain{
		{
			CreatedAt:   &testDNSDomainCreatedAt,
			ID:          &testDNSDomainID,
			UnicodeName: &testDNSDomainUnicodeName,
		},
	}

	actual, err := ts.client.ListDNSDomains(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_GetDNSDomain() {
	ts.mock().
		On(
			"GetDnsDomainWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				testDNSDomainID,
				args.Get(1),
			)
		}).
		Return(
			&oapi.GetDnsDomainResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.DnsDomain{
					CreatedAt:   &testDNSDomainCreatedAt,
					Id:          &testDNSDomainID,
					UnicodeName: &testDNSDomainUnicodeName,
				},
			},
			nil,
		)

	expected := &DNSDomain{
		CreatedAt:   &testDNSDomainCreatedAt,
		ID:          &testDNSDomainID,
		UnicodeName: &testDNSDomainUnicodeName,
	}

	actual, err := ts.client.GetDNSDomain(context.Background(), testZone, testDNSDomainID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_CreateDNSDomain() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
	)

	ts.mock().
		On(
			"CreateDnsDomainWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				oapi.CreateDnsDomainJSONRequestBody{
					UnicodeName: &testDNSDomainUnicodeName,
				},
				args.Get(1),
			)
		}).
		Return(
			&oapi.CreateDnsDomainResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.DnsDomain{
					Id: &testOperationID,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testDNSDomainID, nil),
		State:     &testOperationState,
	})

	ts.mock().
		On("GetDnsDomainWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.GetDnsDomainResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.DnsDomain{
				CreatedAt:   &testDNSDomainCreatedAt,
				Id:          &testDNSDomainID,
				UnicodeName: &testDNSDomainUnicodeName,
			},
		}, nil)

	expected := &DNSDomain{
		CreatedAt:   &testDNSDomainCreatedAt,
		ID:          &testDNSDomainID,
		UnicodeName: &testDNSDomainUnicodeName,
	}

	actual, err := ts.client.CreateDNSDomain(context.Background(), testZone, &DNSDomain{
		UnicodeName: &testDNSDomainUnicodeName,
	})
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_DeleteDnsDomain() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
	)

	ts.mock().
		On(
			"DeleteDnsDomainWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				testDNSDomainID,
				args.Get(1),
			)
		}).
		Return(
			&oapi.DeleteDnsDomainResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testDNSDomainID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testDNSDomainID, nil),
		State:     &testOperationState,
	})

	err := ts.client.DeleteDNSDomain(context.Background(), testZone, &DNSDomain{ID: &testDNSDomainID})
	ts.Require().NoError(err)
}

func (ts *testSuite) TestClient_GetDNSDomainZoneFile() {
	var expected = `$ORIGIN example.com.
$TTL 1h
example.com. 3600 IN SOA ns1.exoscale.ch. admin.dnsimple.com. 1654788358 86400 7200 604800 300
example.com. 3600 IN NS ns1.exoscale.ch.
example.com. 3600 IN NS ns1.exoscale.com.
example.com. 3600 IN NS ns1.exoscale.io.
example.com. 3600 IN NS ns1.exoscale.net.`

	ts.mock().
		On(
			"GetDnsDomainZoneFileWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				testDNSDomainID,
				args.Get(1),
			)
		}).
		Return(
			&oapi.GetDnsDomainZoneFileResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				Body:         []byte(expected),
			},
			nil,
		)

	actual, err := ts.client.GetDNSDomainZoneFile(context.Background(), testZone, testDNSDomainID)
	ts.Require().NoError(err)
	ts.Require().Equal(actual, []byte(expected))
}

func (ts *testSuite) TestClient_ListDNSDomainRecords() {
	tp := oapi.DnsDomainRecordType(testDNSDomainRecordType)
	records := struct {
		DnsDomainRecords *[]oapi.DnsDomainRecord `json:"dns-domain-records,omitempty"` //nolint: revive
	}{
		&[]oapi.DnsDomainRecord{
			oapi.DnsDomainRecord{
				Content:   &testDNSDomainRecordContent,
				CreatedAt: &testDNSDomainRecordCreatedAt,
				Id:        &testDNSDomainRecordID,
				Name:      &testDNSDomainRecordName,
				Priority:  &testDNSDomainRecordPriority,
				Ttl:       &testDNSDomainRecordTTL,
				Type:      &tp,
				UpdatedAt: &testDNSDomainRecordUpdatedAt,
			},
		},
	}
	ts.mock().
		On(
			"ListDnsDomainRecordsWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // domainId
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(
			&oapi.ListDnsDomainRecordsResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200:      &records,
			},
			nil,
		)

	expected := []DNSDomainRecord{
		{
			Content:   &testDNSDomainRecordContent,
			CreatedAt: &testDNSDomainRecordCreatedAt,
			ID:        &testDNSDomainRecordID,
			Name:      &testDNSDomainRecordName,
			Priority:  &testDNSDomainRecordPriority,
			TTL:       &testDNSDomainRecordTTL,
			Type:      &testDNSDomainRecordType,
			UpdatedAt: &testDNSDomainRecordUpdatedAt,
		},
	}

	actual, err := ts.client.ListDNSDomainRecords(context.Background(), testZone, testDNSDomainID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_GetDNSDomainRecord() {
	tp := oapi.DnsDomainRecordType(testDNSDomainRecordType)
	ts.mock().
		On(
			"GetDnsDomainRecordWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // domainId
			mock.Anything,                 // recordId
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				testDNSDomainID,
				args.Get(1),
			)
			ts.Require().Equal(
				testDNSDomainRecordID,
				args.Get(2),
			)
		}).
		Return(
			&oapi.GetDnsDomainRecordResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.DnsDomainRecord{
					Content:   &testDNSDomainRecordContent,
					CreatedAt: &testDNSDomainRecordCreatedAt,
					Id:        &testDNSDomainRecordID,
					Name:      &testDNSDomainRecordName,
					Priority:  &testDNSDomainRecordPriority,
					Ttl:       &testDNSDomainRecordTTL,
					Type:      &tp,
					UpdatedAt: &testDNSDomainRecordUpdatedAt,
				},
			},
			nil,
		)

	expected := &DNSDomainRecord{
		Content:   &testDNSDomainRecordContent,
		CreatedAt: &testDNSDomainRecordCreatedAt,
		ID:        &testDNSDomainRecordID,
		Name:      &testDNSDomainRecordName,
		Priority:  &testDNSDomainRecordPriority,
		TTL:       &testDNSDomainRecordTTL,
		Type:      &testDNSDomainRecordType,
		UpdatedAt: &testDNSDomainRecordUpdatedAt,
	}

	actual, err := ts.client.GetDNSDomainRecord(context.Background(), testZone, testDNSDomainID, testDNSDomainRecordID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_CreateDNSDomainRecord() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
	)
	tp := oapi.DnsDomainRecordType(testDNSDomainRecordType)

	ts.mock().
		On(
			"CreateDnsDomainRecordWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // domainId
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				testDNSDomainID,
				args.Get(1),
			)
			ts.Require().Equal(
				oapi.CreateDnsDomainRecordJSONRequestBody{
					Content:  testDNSDomainRecordContent,
					Name:     testDNSDomainRecordName,
					Priority: &testDNSDomainRecordPriority,
					Ttl:      &testDNSDomainRecordTTL,
					Type:     oapi.CreateDnsDomainRecordJSONBodyType(testDNSDomainRecordType),
				},
				args.Get(2),
			)
		}).
		Return(
			&oapi.CreateDnsDomainRecordResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testDNSDomainRecordID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testDNSDomainRecordID, nil),
		State:     &testOperationState,
	})

	ts.mock().
		On("GetDnsDomainRecordWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // domainId
			mock.Anything,                 // recordId
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.GetDnsDomainRecordResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.DnsDomainRecord{
				Content:   &testDNSDomainRecordContent,
				CreatedAt: &testDNSDomainRecordCreatedAt,
				Id:        &testDNSDomainRecordID,
				Name:      &testDNSDomainRecordName,
				Priority:  &testDNSDomainRecordPriority,
				Ttl:       &testDNSDomainRecordTTL,
				Type:      &tp,
				UpdatedAt: &testDNSDomainRecordUpdatedAt,
			},
		}, nil)

	expected := &DNSDomainRecord{
		Content:   &testDNSDomainRecordContent,
		CreatedAt: &testDNSDomainRecordCreatedAt,
		ID:        &testDNSDomainRecordID,
		Name:      &testDNSDomainRecordName,
		Priority:  &testDNSDomainRecordPriority,
		TTL:       &testDNSDomainRecordTTL,
		Type:      &testDNSDomainRecordType,
		UpdatedAt: &testDNSDomainRecordUpdatedAt,
	}

	actual, err := ts.client.CreateDNSDomainRecord(context.Background(), testZone, testDNSDomainID, &DNSDomainRecord{
		Content:  &testDNSDomainRecordContent,
		Name:     &testDNSDomainRecordName,
		Priority: &testDNSDomainRecordPriority,
		TTL:      &testDNSDomainRecordTTL,
		Type:     &testDNSDomainRecordType,
	})
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_DeleteDNSDomainRecord() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
	)

	ts.mock().
		On(
			"DeleteDnsDomainRecordWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // domainId
			mock.Anything,                 // recordId
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				testDNSDomainID,
				args.Get(1),
			)
			ts.Require().Equal(
				testDNSDomainRecordID,
				args.Get(2),
			)
		}).
		Return(
			&oapi.DeleteDnsDomainRecordResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testDNSDomainRecordID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testDNSDomainRecordID, nil),
		State:     &testOperationState,
	})

	err := ts.client.DeleteDNSDomainRecord(context.Background(), testZone, testDNSDomainID, &DNSDomainRecord{ID: &testDNSDomainRecordID})
	ts.Require().NoError(err)
}

func (ts *testSuite) TestClient_UpdateDNSDomainRecord() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
	)

	ts.mock().
		On(
			"UpdateDnsDomainRecordWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // domainId
			mock.Anything,                 // recordId
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				testDNSDomainID,
				args.Get(1),
			)
			ts.Require().Equal(
				testDNSDomainRecordID,
				args.Get(2),
			)
			ts.Require().Equal(
				oapi.UpdateDnsDomainRecordJSONRequestBody{
					Content:  &testDNSDomainRecordContent,
					Name:     &testDNSDomainRecordName,
					Priority: &testDNSDomainRecordPriority,
					Ttl:      &testDNSDomainRecordTTL,
				},
				args.Get(3),
			)
		}).
		Return(
			&oapi.UpdateDnsDomainRecordResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testDNSDomainRecordID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testDNSDomainRecordID, nil),
		State:     &testOperationState,
	})

	err := ts.client.UpdateDNSDomainRecord(context.Background(), testZone, testDNSDomainID, &DNSDomainRecord{
		ID:       &testDNSDomainRecordID,
		Content:  &testDNSDomainRecordContent,
		Name:     &testDNSDomainRecordName,
		Priority: &testDNSDomainRecordPriority,
		TTL:      &testDNSDomainRecordTTL,
		Type:     &testDNSDomainRecordType,
	})
	ts.Require().NoError(err)
}
