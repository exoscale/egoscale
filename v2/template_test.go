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
	testTemplateBootMode              = "uefi"
	testTemplateBuild                 = "2020-04-22-ed8fea"
	testTemplateChecksum              = "ed8fea0b3c7c8a62801e414b91e23e74"
	testTemplateCreatedAt, _          = time.Parse(iso8601Format, "2020-05-26T12:09:42Z")
	testTemplateDefaultUser           = "test-user"
	testTemplateDescription           = "Test Template description"
	testTemplateFamily                = "test-family"
	testTemplateID                    = new(clientTestSuite).randomID()
	testTemplateName                  = "test-template"
	testTemplatePasswordEnabled       = true
	testTemplateSize            int64 = 10737418240
	testTemplateSSHKeyEnabled         = true
	testTemplateURL                   = "https://www.exoscale.com/"
	testTemplateVersion               = "1"
	testTemplateVisibility            = "public"
)

func (ts *clientTestSuite) TestClient_ListTemplates() {
	httpmock.RegisterResponder("GET", "/template", func(req *http.Request) (*http.Response, error) {
		ts.Require().Equal(testTemplateVisibility, req.URL.Query().Get("visibility"))
		ts.Require().Equal(testTemplateFamily, req.URL.Query().Get("family"))

		resp, err := httpmock.NewJsonResponse(http.StatusOK,
			struct {
				Templates *[]papi.Template `json:"templates,omitempty"`
			}{
				Templates: &[]papi.Template{{
					BootMode:        (*papi.TemplateBootMode)(&testTemplateBootMode),
					Build:           &testTemplateBuild,
					Checksum:        &testTemplateChecksum,
					CreatedAt:       &testTemplateCreatedAt,
					DefaultUser:     &testTemplateDefaultUser,
					Description:     &testTemplateDescription,
					Family:          &testTemplateFamily,
					Id:              &testTemplateID,
					Name:            &testTemplateName,
					PasswordEnabled: &testTemplatePasswordEnabled,
					Size:            &testTemplateSize,
					SshKeyEnabled:   &testTemplateSSHKeyEnabled,
					Url:             &testTemplateURL,
					Version:         &testTemplateVersion,
					Visibility:      (*papi.TemplateVisibility)(&testTemplateVisibility),
				}},
			})
		if err != nil {
			ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
		}
		return resp, nil
	})

	expected := []*Template{{
		BootMode:        testTemplateBootMode,
		Build:           testTemplateBuild,
		Checksum:        testTemplateChecksum,
		CreatedAt:       testTemplateCreatedAt,
		DefaultUser:     testTemplateDefaultUser,
		Description:     testTemplateDescription,
		Family:          testTemplateFamily,
		ID:              testTemplateID,
		Name:            testTemplateName,
		PasswordEnabled: testTemplatePasswordEnabled,
		SSHKeyEnabled:   testTemplateSSHKeyEnabled,
		Size:            testTemplateSize,
		URL:             testTemplateURL,
		Version:         testTemplateVersion,
		Visibility:      testTemplateVisibility,
	}}

	actual, err := ts.client.ListTemplates(context.Background(), testZone, testTemplateVisibility, testTemplateFamily)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_GetTemplate() {
	ts.mockAPIRequest("GET", fmt.Sprintf("/template/%s", testTemplateID), papi.Template{
		BootMode:        (*papi.TemplateBootMode)(&testTemplateBootMode),
		Build:           &testTemplateBuild,
		Checksum:        &testTemplateChecksum,
		CreatedAt:       &testTemplateCreatedAt,
		DefaultUser:     &testTemplateDefaultUser,
		Description:     &testTemplateDescription,
		Family:          &testTemplateFamily,
		Id:              &testTemplateID,
		Name:            &testTemplateName,
		PasswordEnabled: &testTemplatePasswordEnabled,
		Size:            &testTemplateSize,
		SshKeyEnabled:   &testTemplateSSHKeyEnabled,
		Url:             &testTemplateURL,
		Version:         &testTemplateVersion,
		Visibility:      (*papi.TemplateVisibility)(&testTemplateVisibility),
	})

	expected := &Template{
		BootMode:        testTemplateBootMode,
		Build:           testTemplateBuild,
		Checksum:        testTemplateChecksum,
		CreatedAt:       testTemplateCreatedAt,
		DefaultUser:     testTemplateDefaultUser,
		Description:     testTemplateDescription,
		Family:          testTemplateFamily,
		ID:              testTemplateID,
		Name:            testTemplateName,
		PasswordEnabled: testTemplatePasswordEnabled,
		SSHKeyEnabled:   testTemplateSSHKeyEnabled,
		Size:            testTemplateSize,
		URL:             testTemplateURL,
		Version:         testTemplateVersion,
		Visibility:      testTemplateVisibility,
	}

	actual, err := ts.client.GetTemplate(context.Background(), testZone, expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_DeleteTemplate() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = papi.OperationStateSuccess
		deleted            = false
	)

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("/template/%s", testTemplateID),
		func(req *http.Request) (*http.Response, error) {
			deleted = true

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testTemplateID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testTemplateID},
	})

	ts.Require().NoError(ts.client.DeleteTemplate(context.Background(), testZone, testTemplateID))
	ts.Require().True(deleted)
}
