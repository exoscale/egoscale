package v2

import (
	"context"
	"fmt"
	"net/http"
	"time"

	papi "github.com/exoscale/egoscale/v2/internal/public-api"
	"github.com/jarcoal/httpmock"
)

var (
	testSnapshotCreatedAt, _       = time.Parse(iso8601Format, "2020-05-26T12:09:42Z")
	testSnapshotID                 = new(clientTestSuite).randomID()
	testSnapshotInstanceID         = new(clientTestSuite).randomID()
	testSnapshotName               = new(clientTestSuite).randomString(10)
	testSnapshotSize         int64 = 10
	testSnapshotState              = papi.SnapshotStateExported
)

func (ts *clientTestSuite) TestClient_DeleteSnapshot() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = papi.OperationStateSuccess
		deleted            = false
	)

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("/snapshot/%s", testSnapshotID),
		func(req *http.Request) (*http.Response, error) {
			deleted = true

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testSnapshotID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testSnapshotID},
	})

	ts.Require().NoError(ts.client.DeleteSnapshot(context.Background(), testZone, &Snapshot{ID: &testSnapshotID}))
	ts.Require().True(deleted)
}

func (ts *clientTestSuite) TestClient_ExportSnapshot() {
	var (
		testSnapshotExportMD5Sum       = "c9887de796993c2519b463bcd9509e08"
		testSnapshotExportPresignedURL = fmt.Sprintf("https://sos-%s.exo.io/test/%s/%s",
			testZone,
			ts.randomID(),
			testSnapshotID)
		testOperationID    = ts.randomID()
		testOperationState = papi.OperationStateSuccess
		exported           = false
	)

	httpmock.RegisterResponder("POST", fmt.Sprintf("/snapshot/%s:export", testSnapshotID),
		func(req *http.Request) (*http.Response, error) {
			exported = true

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testSnapshotID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testSnapshotID},
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/snapshot/%s", testSnapshotID), papi.Snapshot{
		CreatedAt: &testSnapshotCreatedAt,
		Export: &struct {
			Md5sum       *string `json:"md5sum,omitempty"`
			PresignedUrl *string `json:"presigned-url,omitempty"` // nolint:revive
		}{
			Md5sum:       &testSnapshotExportMD5Sum,
			PresignedUrl: &testSnapshotExportPresignedURL,
		},
		Id:       &testSnapshotID,
		Instance: &papi.Instance{Id: &testSnapshotInstanceID},
		Name:     &testSnapshotName,
		State:    &testSnapshotState,
	})

	expected := &SnapshotExport{
		MD5sum:       &testSnapshotExportMD5Sum,
		PresignedURL: &testSnapshotExportPresignedURL,
	}

	actual, err := ts.client.ExportSnapshot(context.Background(), testZone, &Snapshot{ID: &testSnapshotID})
	ts.Require().NoError(err)
	ts.Require().True(exported)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_ListSnapshots() {
	ts.mockAPIRequest("GET", "/snapshot", struct {
		Snapshots *[]papi.Snapshot `json:"snapshots,omitempty"`
	}{
		Snapshots: &[]papi.Snapshot{{
			CreatedAt: &testSnapshotCreatedAt,
			Id:        &testSnapshotID,
			Instance:  &papi.Instance{Id: &testSnapshotInstanceID},
			Name:      &testSnapshotName,
			Size:      &testSnapshotSize,
			State:     &testSnapshotState,
		}},
	})

	expected := []*Snapshot{{
		CreatedAt:  &testSnapshotCreatedAt,
		ID:         &testSnapshotID,
		InstanceID: &testSnapshotInstanceID,
		Name:       &testSnapshotName,
		Size:       &testSnapshotSize,
		State:      (*string)(&testSnapshotState),
	}}

	actual, err := ts.client.ListSnapshots(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_GetSnapshot() {
	ts.mockAPIRequest("GET", fmt.Sprintf("/snapshot/%s", testSnapshotID), papi.Snapshot{
		CreatedAt: &testSnapshotCreatedAt,
		Id:        &testSnapshotID,
		Instance:  &papi.Instance{Id: &testSnapshotInstanceID},
		Name:      &testSnapshotName,
		Size:      &testSnapshotSize,
		State:     &testSnapshotState,
	})

	expected := &Snapshot{
		CreatedAt:  &testSnapshotCreatedAt,
		ID:         &testSnapshotID,
		InstanceID: &testSnapshotInstanceID,
		Name:       &testSnapshotName,
		Size:       &testSnapshotSize,
		State:      (*string)(&testSnapshotState),
	}

	actual, err := ts.client.GetSnapshot(context.Background(), testZone, *expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}
