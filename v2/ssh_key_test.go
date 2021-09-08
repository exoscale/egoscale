package v2

import (
	"context"
	"fmt"
	"net/http"

	"github.com/exoscale/egoscale/v2/oapi"
	"github.com/jarcoal/httpmock"
)

var (
	testSSHKeyFingerprint = new(clientTestSuite).randomString(10)
	testSSHKeyName        = new(clientTestSuite).randomString(10)
	testSSHKeyPublicKey   = new(clientTestSuite).randomString(10)
)

func (ts *clientTestSuite) TestClient_DeleteSSHKey() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		deleted            = false
	)

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("/ssh-key/%s", testSSHKeyName),
		func(req *http.Request) (*http.Response, error) {
			deleted = true

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testSSHKeyName},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testSSHKeyName},
	})

	ts.Require().NoError(ts.client.DeleteSSHKey(context.Background(), testZone, &SSHKey{Name: &testSSHKeyName}))
	ts.Require().True(deleted)
}

func (ts *clientTestSuite) TestClient_GetSSHKey() {
	ts.mockAPIRequest("GET", fmt.Sprintf("/ssh-key/%s", testSSHKeyName), oapi.SshKey{
		Fingerprint: &testSSHKeyFingerprint,
		Name:        &testSSHKeyName,
	})

	expected := &SSHKey{
		Fingerprint: &testSSHKeyFingerprint,
		Name:        &testSSHKeyName,
	}

	actual, err := ts.client.GetSSHKey(context.Background(), testZone, *expected.Name)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_ListSSHKeys() {
	httpmock.RegisterResponder("GET", "/ssh-key", func(req *http.Request) (*http.Response, error) {
		resp, err := httpmock.NewJsonResponse(http.StatusOK,
			struct {
				SSHKeys *[]oapi.SshKey `json:"ssh-keys,omitempty"`
			}{
				SSHKeys: &[]oapi.SshKey{{
					Fingerprint: &testSSHKeyFingerprint,
					Name:        &testSSHKeyName,
				}},
			})
		if err != nil {
			ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
		}
		return resp, nil
	})

	expected := []*SSHKey{{
		Fingerprint: &testSSHKeyFingerprint,
		Name:        &testSSHKeyName,
	}}

	actual, err := ts.client.ListSSHKeys(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_RegisterSSHKey() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
	)

	httpmock.RegisterResponder("POST", "/ssh-key",
		func(req *http.Request) (*http.Response, error) {
			var actual oapi.RegisterSshKeyJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := oapi.RegisterSshKeyJSONRequestBody{
				Name:      testSSHKeyName,
				PublicKey: testSSHKeyPublicKey,
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testSSHKeyName},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/ssh-key/%s", testSSHKeyName), oapi.SshKey{
		Fingerprint: &testSSHKeyFingerprint,
		Name:        &testSSHKeyName,
	})

	expected := &SSHKey{
		Fingerprint: &testSSHKeyFingerprint,
		Name:        &testSSHKeyName,
	}

	actual, err := ts.client.RegisterSSHKey(context.Background(), testZone, testSSHKeyName, testSSHKeyPublicKey)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}
