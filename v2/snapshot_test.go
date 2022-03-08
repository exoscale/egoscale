package v2

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/exoscale/egoscale/v2/oapi"
)

var (
	testSnapshotCreatedAt, _       = time.Parse(iso8601Format, "2020-05-26T12:09:42Z")
	testSnapshotID                 = new(testSuite).randomID()
	testSnapshotInstanceID         = new(testSuite).randomID()
	testSnapshotName               = new(testSuite).randomString(10)
	testSnapshotSize         int64 = 10
	testSnapshotState              = oapi.SnapshotStateExported
)

func (ts *testSuite) TestClient_DeleteSnapshot() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		deleted            = false
	)

	ts.mock().
		On(
			"DeleteSnapshotWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testSnapshotID, args.Get(1))
			deleted = true
		}).
		Return(
			&oapi.DeleteSnapshotResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id: &testOperationID,
					Reference: &struct {
						Command *string `json:"command,omitempty"`
						Id      *string `json:"id,omitempty"`
						Link    *string `json:"link,omitempty"`
					}{Id: &testSnapshotID},
					State: &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id: &testOperationID,
		Reference: &struct {
			Command *string `json:"command,omitempty"`
			Id      *string `json:"id,omitempty"`
			Link    *string `json:"link,omitempty"`
		}{Id: &testSnapshotID},
		State: &testOperationState,
	})

	ts.Require().NoError(ts.client.DeleteSnapshot(context.Background(), testZone, &Snapshot{ID: &testSnapshotID}))
	ts.Require().True(deleted)
}

func (ts *testSuite) TestClient_ExportSnapshot() {
	var (
		testSnapshotExportMD5Sum       = "c9887de796993c2519b463bcd9509e08"
		testSnapshotExportPresignedURL = fmt.Sprintf("https://sos-%s.exo.io/test/%s/%s",
			testZone,
			ts.randomID(),
			testSnapshotID)
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		exported           = false
	)

	ts.mock().
		On(
			"ExportSnapshotWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testSnapshotID, args.Get(1))
			exported = true
		}).
		Return(
			&oapi.ExportSnapshotResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id: &testOperationID,
					Reference: &struct {
						Command *string `json:"command,omitempty"`
						Id      *string `json:"id,omitempty"`
						Link    *string `json:"link,omitempty"`
					}{Id: &testSnapshotID},
					State: &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id: &testOperationID,
		Reference: &struct {
			Command *string `json:"command,omitempty"`
			Id      *string `json:"id,omitempty"`
			Link    *string `json:"link,omitempty"`
		}{Id: &testSnapshotID},
		State: &testOperationState,
	})

	ts.mock().
		On("GetSnapshotWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.GetSnapshotResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.Snapshot{
				CreatedAt: &testSnapshotCreatedAt,
				Export: &struct {
					Md5sum       *string `json:"md5sum,omitempty"`
					PresignedUrl *string `json:"presigned-url,omitempty"` // nolint:revive
				}{
					Md5sum:       &testSnapshotExportMD5Sum,
					PresignedUrl: &testSnapshotExportPresignedURL,
				},
				Id:       &testSnapshotID,
				Instance: &oapi.Instance{Id: &testSnapshotInstanceID},
				Name:     &testSnapshotName,
				State:    &testSnapshotState,
			},
		}, nil)

	expected := &SnapshotExport{
		MD5sum:       &testSnapshotExportMD5Sum,
		PresignedURL: &testSnapshotExportPresignedURL,
	}

	actual, err := ts.client.ExportSnapshot(context.Background(), testZone, &Snapshot{ID: &testSnapshotID})
	ts.Require().NoError(err)
	ts.Require().True(exported)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_ListSnapshots() {
	ts.mock().
		On("ListSnapshotsWithResponse",
			mock.Anything,                 // ctx
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.ListSnapshotsResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &struct {
				Snapshots *[]oapi.Snapshot `json:"snapshots,omitempty"`
			}{
				Snapshots: &[]oapi.Snapshot{{
					CreatedAt: &testSnapshotCreatedAt,
					Id:        &testSnapshotID,
					Instance:  &oapi.Instance{Id: &testSnapshotInstanceID},
					Name:      &testSnapshotName,
					Size:      &testSnapshotSize,
					State:     &testSnapshotState,
				}},
			},
		}, nil)

	expected := []*Snapshot{{
		CreatedAt:  &testSnapshotCreatedAt,
		ID:         &testSnapshotID,
		InstanceID: &testSnapshotInstanceID,
		Name:       &testSnapshotName,
		Size:       &testSnapshotSize,
		State:      (*string)(&testSnapshotState),
		Zone:       &testZone,
	}}

	actual, err := ts.client.ListSnapshots(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_GetSnapshot() {
	ts.mock().
		On("GetSnapshotWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testSnapshotID, args.Get(1))
		}).
		Return(&oapi.GetSnapshotResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.Snapshot{
				CreatedAt: &testSnapshotCreatedAt,
				Id:        &testSnapshotID,
				Instance:  &oapi.Instance{Id: &testSnapshotInstanceID},
				Name:      &testSnapshotName,
				Size:      &testSnapshotSize,
				State:     &testSnapshotState,
			},
		}, nil)

	expected := &Snapshot{
		CreatedAt:  &testSnapshotCreatedAt,
		ID:         &testSnapshotID,
		InstanceID: &testSnapshotInstanceID,
		Name:       &testSnapshotName,
		Size:       &testSnapshotSize,
		State:      (*string)(&testSnapshotState),
		Zone:       &testZone,
	}

	actual, err := ts.client.GetSnapshot(context.Background(), testZone, *expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}
