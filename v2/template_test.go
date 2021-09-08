package v2

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/exoscale/egoscale/v2/oapi"
	"github.com/jarcoal/httpmock"
)

var (
	testTemplateBootMode              = "uefi"
	testTemplateBuild                 = "2020-04-22-ed8fea"
	testTemplateChecksum              = "ed8fea0b3c7c8a62801e414b91e23e74"
	testTemplateCreatedAt, _          = time.Parse(iso8601Format, "2020-05-26T12:09:42Z")
	testTemplateDefaultUser           = new(clientTestSuite).randomString(10)
	testTemplateDescription           = new(clientTestSuite).randomString(10)
	testTemplateFamily                = new(clientTestSuite).randomString(10)
	testTemplateID                    = new(clientTestSuite).randomID()
	testTemplateName                  = new(clientTestSuite).randomString(10)
	testTemplatePasswordEnabled       = true
	testTemplateSize            int64 = 10737418240
	testTemplateSSHKeyEnabled         = true
	testTemplateURL                   = "https://example.net/test.qcow2"
	testTemplateVersion               = "1"
	testTemplateVisibility            = "public"
)

func (ts *clientTestSuite) TestClient_DeleteTemplate() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		deleted            = false
	)

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("/template/%s", testTemplateID),
		func(req *http.Request) (*http.Response, error) {
			deleted = true

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testTemplateID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testTemplateID},
	})

	ts.Require().NoError(ts.client.DeleteTemplate(context.Background(), testZone, &Template{ID: &testTemplateID}))
	ts.Require().True(deleted)
}

func (ts *clientTestSuite) TestClient_GetTemplate() {
	ts.mockAPIRequest("GET", fmt.Sprintf("/template/%s", testTemplateID), oapi.Template{
		BootMode:        (*oapi.TemplateBootMode)(&testTemplateBootMode),
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
		Visibility:      (*oapi.TemplateVisibility)(&testTemplateVisibility),
	})

	expected := &Template{
		BootMode:        &testTemplateBootMode,
		Build:           &testTemplateBuild,
		Checksum:        &testTemplateChecksum,
		CreatedAt:       &testTemplateCreatedAt,
		DefaultUser:     &testTemplateDefaultUser,
		Description:     &testTemplateDescription,
		Family:          &testTemplateFamily,
		ID:              &testTemplateID,
		Name:            &testTemplateName,
		PasswordEnabled: &testTemplatePasswordEnabled,
		SSHKeyEnabled:   &testTemplateSSHKeyEnabled,
		Size:            &testTemplateSize,
		URL:             &testTemplateURL,
		Version:         &testTemplateVersion,
		Visibility:      &testTemplateVisibility,
	}

	actual, err := ts.client.GetTemplate(context.Background(), testZone, *expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_ListTemplates() {
	httpmock.RegisterResponder("GET", "/template", func(req *http.Request) (*http.Response, error) {
		ts.Require().Equal(testTemplateVisibility, req.URL.Query().Get("visibility"))
		ts.Require().Equal(testTemplateFamily, req.URL.Query().Get("family"))

		resp, err := httpmock.NewJsonResponse(http.StatusOK,
			struct {
				Templates *[]oapi.Template `json:"templates,omitempty"`
			}{
				Templates: &[]oapi.Template{{
					BootMode:        (*oapi.TemplateBootMode)(&testTemplateBootMode),
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
					Visibility:      (*oapi.TemplateVisibility)(&testTemplateVisibility),
				}},
			})
		if err != nil {
			ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
		}
		return resp, nil
	})

	expected := []*Template{{
		BootMode:        &testTemplateBootMode,
		Build:           &testTemplateBuild,
		Checksum:        &testTemplateChecksum,
		CreatedAt:       &testTemplateCreatedAt,
		DefaultUser:     &testTemplateDefaultUser,
		Description:     &testTemplateDescription,
		Family:          &testTemplateFamily,
		ID:              &testTemplateID,
		Name:            &testTemplateName,
		PasswordEnabled: &testTemplatePasswordEnabled,
		SSHKeyEnabled:   &testTemplateSSHKeyEnabled,
		Size:            &testTemplateSize,
		URL:             &testTemplateURL,
		Version:         &testTemplateVersion,
		Visibility:      &testTemplateVisibility,
	}}

	actual, err := ts.client.ListTemplates(context.Background(), testZone, testTemplateVisibility, testTemplateFamily)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_RegisterTemplate() {
	var (
		templateVisibility = "private"
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
	)

	httpmock.RegisterResponder("POST", "/template",
		func(req *http.Request) (*http.Response, error) {
			var actual oapi.RegisterTemplateJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := oapi.RegisterTemplateJSONRequestBody{
				BootMode:        (*oapi.RegisterTemplateJSONBodyBootMode)(&testTemplateBootMode),
				Checksum:        testTemplateChecksum,
				DefaultUser:     &testTemplateDefaultUser,
				Description:     &testTemplateDescription,
				Name:            testTemplateName,
				PasswordEnabled: testTemplatePasswordEnabled,
				SshKeyEnabled:   testTemplateSSHKeyEnabled,
				Url:             testTemplateURL,
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testTemplateID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testTemplateID},
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/template/%s", testTemplateID), oapi.Template{
		BootMode:        (*oapi.TemplateBootMode)(&testTemplateBootMode),
		Checksum:        &testTemplateChecksum,
		CreatedAt:       &testTemplateCreatedAt,
		DefaultUser:     &testTemplateDefaultUser,
		Description:     &testTemplateDescription,
		Id:              &testTemplateID,
		Name:            &testTemplateName,
		PasswordEnabled: &testTemplatePasswordEnabled,
		Size:            &testTemplateSize,
		SshKeyEnabled:   &testTemplateSSHKeyEnabled,
		Url:             &testTemplateURL,
		Visibility:      (*oapi.TemplateVisibility)(&templateVisibility),
	})

	expected := &Template{
		BootMode:        &testTemplateBootMode,
		Checksum:        &testTemplateChecksum,
		CreatedAt:       &testTemplateCreatedAt,
		DefaultUser:     &testTemplateDefaultUser,
		Description:     &testTemplateDescription,
		ID:              &testTemplateID,
		Name:            &testTemplateName,
		PasswordEnabled: &testTemplatePasswordEnabled,
		SSHKeyEnabled:   &testTemplateSSHKeyEnabled,
		Size:            &testTemplateSize,
		URL:             &testTemplateURL,
		Visibility:      &templateVisibility,
	}

	actual, err := ts.client.RegisterTemplate(context.Background(), testZone, &Template{
		BootMode:        &testTemplateBootMode,
		Checksum:        &testTemplateChecksum,
		DefaultUser:     &testTemplateDefaultUser,
		Description:     &testTemplateDescription,
		Name:            &testTemplateName,
		PasswordEnabled: &testTemplatePasswordEnabled,
		SSHKeyEnabled:   &testTemplateSSHKeyEnabled,
		URL:             &testTemplateURL,
	})
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}
