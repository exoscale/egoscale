package v2

import (
	"context"
	"net/http"

	"github.com/stretchr/testify/mock"

	"github.com/exoscale/egoscale/v2/oapi"
)

var (
	testSSHKeyFingerprint = new(testSuite).randomString(10)
	testSSHKeyName        = new(testSuite).randomString(10)
	testSSHKeyPublicKey   = new(testSuite).randomString(10)
)

func (ts *testSuite) TestClient_DeleteSSHKey() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		deleted            = false
	)

	ts.mock().
		On(
			"DeleteSshKeyWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // name
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testSSHKeyName, args.Get(1))
			deleted = true
		}).
		Return(
			&oapi.DeleteSshKeyResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id: &testOperationID,
					Reference: &struct {
						Command *string `json:"command,omitempty"`
						Id      *string `json:"id,omitempty"` // revive:disable-line
						Link    *string `json:"link,omitempty"`
					}{Id: &testSSHKeyName},
					State: &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id: &testOperationID,
		Reference: &struct {
			Command *string `json:"command,omitempty"`
			Id      *string `json:"id,omitempty"` // revive:disable-line
			Link    *string `json:"link,omitempty"`
		}{Id: &testSSHKeyName},
		State: &testOperationState,
	})

	ts.Require().NoError(ts.client.DeleteSSHKey(context.Background(), testZone, &SSHKey{Name: &testSSHKeyName}))
	ts.Require().True(deleted)
}

func (ts *testSuite) TestClient_GetSSHKey() {
	ts.mock().
		On("GetSshKeyWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // name
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testSSHKeyName, args.Get(1))
		}).
		Return(&oapi.GetSshKeyResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.SshKey{
				Fingerprint: &testSSHKeyFingerprint,
				Name:        &testSSHKeyName,
			},
		}, nil)

	expected := &SSHKey{
		Fingerprint: &testSSHKeyFingerprint,
		Name:        &testSSHKeyName,
	}

	actual, err := ts.client.GetSSHKey(context.Background(), testZone, *expected.Name)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_ListSSHKeys() {
	ts.mock().
		On("ListSshKeysWithResponse",
			mock.Anything,                 // ctx
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.ListSshKeysResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &struct {
				SshKeys *[]oapi.SshKey `json:"ssh-keys,omitempty"` // nolint:revive
			}{
				SshKeys: &[]oapi.SshKey{{
					Fingerprint: &testSSHKeyFingerprint,
					Name:        &testSSHKeyName,
				}},
			},
		}, nil)

	expected := []*SSHKey{{
		Fingerprint: &testSSHKeyFingerprint,
		Name:        &testSSHKeyName,
	}}

	actual, err := ts.client.ListSSHKeys(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_RegisterSSHKey() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
	)

	ts.mock().
		On(
			"RegisterSshKeyWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				oapi.RegisterSshKeyJSONRequestBody{
					Name:      testSSHKeyName,
					PublicKey: testSSHKeyPublicKey,
				},
				args.Get(1),
			)
		}).
		Return(
			&oapi.RegisterSshKeyResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id: &testOperationID,
					Reference: &struct {
						Command *string `json:"command,omitempty"`
						Id      *string `json:"id,omitempty"` // revive:disable-line
						Link    *string `json:"link,omitempty"`
					}{Id: &testSSHKeyName},
					State: &testOperationState,
				},
			},
			nil,
		)

	ts.mock().
		On("GetSshKeyWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // name
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.GetSshKeyResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.SshKey{
				Fingerprint: &testSSHKeyFingerprint,
				Name:        &testSSHKeyName,
			},
		}, nil)

	expected := &SSHKey{
		Fingerprint: &testSSHKeyFingerprint,
		Name:        &testSSHKeyName,
	}

	actual, err := ts.client.RegisterSSHKey(context.Background(), testZone, testSSHKeyName, testSSHKeyPublicKey)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}
