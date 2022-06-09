package v2

import (
	"context"
	"net/http"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/exoscale/egoscale/v2/oapi"
)

var (
	testDnsDomainCreatedAt, _ = time.Parse(iso8601Format, "2020-05-26T12:09:42Z")
	testDnsDomainId           = new(testSuite).randomID()
	testDnsDomainUnicodeName  = "test.domain"

	testDnsDomainRecordContent            = new(testSuite).randomString(10)
	testDnsDomainRecordCreatedAt, _       = time.Parse(iso8601Format, "2021-04-27T12:09:42Z")
	testDnsDomainRecordId                 = new(testSuite).randomID()
	testDnsDomainRecordName               = new(testSuite).randomString(10)
	testDnsDomainRecordPriority     int64 = 5
	testDnsDomainRecordTtl          int64 = 1800
	testDnsDomainRecordType               = "CNAME"
	testDnsDomainRecordUpdatedAt, _       = time.Parse(iso8601Format, "2021-06-27T12:09:42Z")
)

func (ts *testSuite) TestClient_ListDnsDomains() {
	domains := struct {
		DnsDomains *[]oapi.DnsDomain `json:"dns-domains,omitempty"`
	}{
		DnsDomains: &[]oapi.DnsDomain{
			oapi.DnsDomain{
				CreatedAt:   &testDnsDomainCreatedAt,
				Id:          &testDnsDomainId,
				UnicodeName: &testDnsDomainUnicodeName,
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

	expected := []DnsDomain{
		{
			CreatedAt:   &testDnsDomainCreatedAt,
			ID:          &testDnsDomainId,
			UnicodeName: &testDnsDomainUnicodeName,
		},
	}

	actual, err := ts.client.ListDnsDomains(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_GetDnsDomain() {
	ts.mock().
		On(
			"GetDnsDomainWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				testDnsDomainId,
				args.Get(1),
			)
		}).
		Return(
			&oapi.GetDnsDomainResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.DnsDomain{
					CreatedAt:   &testDnsDomainCreatedAt,
					Id:          &testDnsDomainId,
					UnicodeName: &testDnsDomainUnicodeName,
				},
			},
			nil,
		)

	expected := &DnsDomain{
		CreatedAt:   &testDnsDomainCreatedAt,
		ID:          &testDnsDomainId,
		UnicodeName: &testDnsDomainUnicodeName,
	}

	actual, err := ts.client.GetDnsDomain(context.Background(), testZone, testDnsDomainId)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_CreateDnsDomain() {
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
					UnicodeName: &testDnsDomainUnicodeName,
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
		Reference: oapi.NewReference(nil, &testDnsDomainId, nil),
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
				CreatedAt:   &testDnsDomainCreatedAt,
				Id:          &testDnsDomainId,
				UnicodeName: &testDnsDomainUnicodeName,
			},
		}, nil)

	expected := &DnsDomain{
		CreatedAt:   &testDnsDomainCreatedAt,
		ID:          &testDnsDomainId,
		UnicodeName: &testDnsDomainUnicodeName,
	}

	actual, err := ts.client.CreateDnsDomain(context.Background(), testZone, &DnsDomain{
		UnicodeName: &testDnsDomainUnicodeName,
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
				testDnsDomainId,
				args.Get(1),
			)
		}).
		Return(
			&oapi.DeleteDnsDomainResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testDnsDomainId, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testDnsDomainId, nil),
		State:     &testOperationState,
	})

	err := ts.client.DeleteDnsDomain(context.Background(), testZone, &DnsDomain{ID: &testDnsDomainId})
	ts.Require().NoError(err)
}

func (ts *testSuite) TestClient_GetDnsDomainZoneFile() {
	var expected string = `$ORIGIN example.com.
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
				testDnsDomainId,
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

	actual, err := ts.client.GetDnsDomainZoneFile(context.Background(), testZone, testDnsDomainId)
	ts.Require().NoError(err)
	ts.Require().Equal(actual, []byte(expected))
}

func (ts *testSuite) TestClient_ListDnsDomainRecords() {
	tp := oapi.DnsDomainRecordType(testDnsDomainRecordType)
	records := struct {
		DnsDomainRecords *[]oapi.DnsDomainRecord `json:"dns-domain-records,omitempty"`
	}{
		&[]oapi.DnsDomainRecord{
			oapi.DnsDomainRecord{
				Content:   &testDnsDomainRecordContent,
				CreatedAt: &testDnsDomainRecordCreatedAt,
				Id:        &testDnsDomainRecordId,
				Name:      &testDnsDomainRecordName,
				Priority:  &testDnsDomainRecordPriority,
				Ttl:       &testDnsDomainRecordTtl,
				Type:      &tp,
				UpdatedAt: &testDnsDomainRecordUpdatedAt,
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

	expected := []DnsDomainRecord{
		{
			Content:   &testDnsDomainRecordContent,
			CreatedAt: &testDnsDomainRecordCreatedAt,
			ID:        &testDnsDomainRecordId,
			Name:      &testDnsDomainRecordName,
			Priority:  &testDnsDomainRecordPriority,
			Ttl:       &testDnsDomainRecordTtl,
			Type:      &testDnsDomainRecordType,
			UpdatedAt: &testDnsDomainRecordUpdatedAt,
		},
	}

	actual, err := ts.client.ListDnsDomainRecords(context.Background(), testZone, testDnsDomainId)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_GetDnsDomainRecord() {
	tp := oapi.DnsDomainRecordType(testDnsDomainRecordType)
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
				testDnsDomainId,
				args.Get(1),
			)
			ts.Require().Equal(
				testDnsDomainRecordId,
				args.Get(2),
			)
		}).
		Return(
			&oapi.GetDnsDomainRecordResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.DnsDomainRecord{
					Content:   &testDnsDomainRecordContent,
					CreatedAt: &testDnsDomainRecordCreatedAt,
					Id:        &testDnsDomainRecordId,
					Name:      &testDnsDomainRecordName,
					Priority:  &testDnsDomainRecordPriority,
					Ttl:       &testDnsDomainRecordTtl,
					Type:      &tp,
					UpdatedAt: &testDnsDomainRecordUpdatedAt,
				},
			},
			nil,
		)

	expected := &DnsDomainRecord{
		Content:   &testDnsDomainRecordContent,
		CreatedAt: &testDnsDomainRecordCreatedAt,
		ID:        &testDnsDomainRecordId,
		Name:      &testDnsDomainRecordName,
		Priority:  &testDnsDomainRecordPriority,
		Ttl:       &testDnsDomainRecordTtl,
		Type:      &testDnsDomainRecordType,
		UpdatedAt: &testDnsDomainRecordUpdatedAt,
	}

	actual, err := ts.client.GetDnsDomainRecord(context.Background(), testZone, testDnsDomainId, testDnsDomainRecordId)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_CreateDnsDomainRecord() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
	)
	tp := oapi.DnsDomainRecordType(testDnsDomainRecordType)

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
				testDnsDomainId,
				args.Get(1),
			)
			ts.Require().Equal(
				oapi.CreateDnsDomainRecordJSONRequestBody{
					Content:  testDnsDomainRecordContent,
					Name:     testDnsDomainRecordName,
					Priority: &testDnsDomainRecordPriority,
					Ttl:      &testDnsDomainRecordTtl,
					Type:     oapi.CreateDnsDomainRecordJSONBodyType(testDnsDomainRecordType),
				},
				args.Get(2),
			)
		}).
		Return(
			&oapi.CreateDnsDomainRecordResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testDnsDomainRecordId, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testDnsDomainRecordId, nil),
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
				Content:   &testDnsDomainRecordContent,
				CreatedAt: &testDnsDomainRecordCreatedAt,
				Id:        &testDnsDomainRecordId,
				Name:      &testDnsDomainRecordName,
				Priority:  &testDnsDomainRecordPriority,
				Ttl:       &testDnsDomainRecordTtl,
				Type:      &tp,
				UpdatedAt: &testDnsDomainRecordUpdatedAt,
			},
		}, nil)

	expected := &DnsDomainRecord{
		Content:   &testDnsDomainRecordContent,
		CreatedAt: &testDnsDomainRecordCreatedAt,
		ID:        &testDnsDomainRecordId,
		Name:      &testDnsDomainRecordName,
		Priority:  &testDnsDomainRecordPriority,
		Ttl:       &testDnsDomainRecordTtl,
		Type:      &testDnsDomainRecordType,
		UpdatedAt: &testDnsDomainRecordUpdatedAt,
	}

	actual, err := ts.client.CreateDnsDomainRecord(context.Background(), testZone, testDnsDomainId, &DnsDomainRecord{
		Content:  &testDnsDomainRecordContent,
		Name:     &testDnsDomainRecordName,
		Priority: &testDnsDomainRecordPriority,
		Ttl:      &testDnsDomainRecordTtl,
		Type:     &testDnsDomainRecordType,
	})
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_DeleteDnsDomainRecord() {
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
				testDnsDomainId,
				args.Get(1),
			)
			ts.Require().Equal(
				testDnsDomainRecordId,
				args.Get(2),
			)
		}).
		Return(
			&oapi.DeleteDnsDomainRecordResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testDnsDomainRecordId, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testDnsDomainRecordId, nil),
		State:     &testOperationState,
	})

	err := ts.client.DeleteDnsDomainRecord(context.Background(), testZone, testDnsDomainId, &DnsDomainRecord{ID: &testDnsDomainRecordId})
	ts.Require().NoError(err)
}

func (ts *testSuite) TestClient_UpdateDnsDomainRecord() {
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
				testDnsDomainId,
				args.Get(1),
			)
			ts.Require().Equal(
				testDnsDomainRecordId,
				args.Get(2),
			)
			ts.Require().Equal(
				oapi.UpdateDnsDomainRecordJSONRequestBody{
					Content:  &testDnsDomainRecordContent,
					Name:     &testDnsDomainRecordName,
					Priority: &testDnsDomainRecordPriority,
					Ttl:      &testDnsDomainRecordTtl,
				},
				args.Get(3),
			)
		}).
		Return(
			&oapi.UpdateDnsDomainRecordResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testDnsDomainRecordId, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testDnsDomainRecordId, nil),
		State:     &testOperationState,
	})

	err := ts.client.UpdateDnsDomainRecord(context.Background(), testZone, testDnsDomainId, &DnsDomainRecord{
		ID:       &testDnsDomainRecordId,
		Content:  &testDnsDomainRecordContent,
		Name:     &testDnsDomainRecordName,
		Priority: &testDnsDomainRecordPriority,
		Ttl:      &testDnsDomainRecordTtl,
		Type:     &testDnsDomainRecordType,
	})
	ts.Require().NoError(err)
}
