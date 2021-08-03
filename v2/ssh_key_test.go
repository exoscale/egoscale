package v2

import (
	"context"
	"fmt"
	"net/http"

	papi "github.com/exoscale/egoscale/v2/internal/public-api"
	"github.com/jarcoal/httpmock"
)

var (
	testSSHKeyFingerprint = new(clientTestSuite).randomString(10)
	testSSHKeyName        = new(clientTestSuite).randomString(10)
	testSSHKeyPublicKey   = new(clientTestSuite).randomString(10)
)

func (ts *clientTestSuite) TestClient_RegisterSSHKey() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = papi.OperationStateSuccess
	)

	httpmock.RegisterResponder("POST", "/ssh-key",
		func(req *http.Request) (*http.Response, error) {
			var actual papi.RegisterSshKeyJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.RegisterSshKeyJSONRequestBody{
				Name:      testSSHKeyName,
				PublicKey: testSSHKeyPublicKey,
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testSSHKeyName},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testSSHKeyName},
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/ssh-key/%s", testSSHKeyName), papi.SshKey{
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

func (ts *clientTestSuite) TestClient_ListSSHKeys() {
	httpmock.RegisterResponder("GET", "/ssh-key", func(req *http.Request) (*http.Response, error) {
		resp, err := httpmock.NewJsonResponse(http.StatusOK,
			struct {
				SSHKeys *[]papi.SshKey `json:"ssh-keys,omitempty"`
			}{
				SSHKeys: &[]papi.SshKey{{
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

func (ts *clientTestSuite) TestClient_GetSSHKey() {
	ts.mockAPIRequest("GET", fmt.Sprintf("/ssh-key/%s", testSSHKeyName), papi.SshKey{
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

func (ts *clientTestSuite) TestClient_DeleteSSHKey() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = papi.OperationStateSuccess
		deleted            = false
	)

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("/ssh-key/%s", testSSHKeyName),
		func(req *http.Request) (*http.Response, error) {
			deleted = true

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testSSHKeyName},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testSSHKeyName},
	})

	ts.Require().NoError(ts.client.DeleteSSHKey(context.Background(), testZone, testSSHKeyName))
	ts.Require().True(deleted)
}
