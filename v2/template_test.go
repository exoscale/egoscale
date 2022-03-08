package v2

import (
	"context"
	"net/http"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/exoscale/egoscale/v2/oapi"
)

var (
	testTemplateBootMode              = "uefi"
	testTemplateBuild                 = "2020-04-22-ed8fea"
	testTemplateChecksum              = "ed8fea0b3c7c8a62801e414b91e23e74"
	testTemplateCreatedAt, _          = time.Parse(iso8601Format, "2020-05-26T12:09:42Z")
	testTemplateDefaultUser           = new(testSuite).randomString(10)
	testTemplateDescription           = new(testSuite).randomString(10)
	testTemplateFamily                = new(testSuite).randomString(10)
	testTemplateID                    = new(testSuite).randomID()
	testTemplateName                  = new(testSuite).randomString(10)
	testTemplatePasswordEnabled       = true
	testTemplateSize            int64 = 10737418240
	testTemplateSSHKeyEnabled         = true
	testTemplateURL                   = "https://example.net/test.qcow2"
	testTemplateVersion               = "1"
	testTemplateVisibility            = "public"
)

func (ts *testSuite) TestClient_CopyTemplate() {
	var (
		dstZone            = "ch-dk-2"
		templateVisibility = "private"
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
	)

	ts.mock().
		On("CopyTemplateWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				oapi.CopyTemplateJSONRequestBody{TargetZone: oapi.Zone{Name: (*oapi.ZoneName)(&dstZone)}},
				args.Get(2),
			)
		}).
		Return(&oapi.CopyTemplateResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.Operation{
				Id:    &testOperationID,
				State: &testOperationState,
				Reference: &struct {
					Command *string `json:"command,omitempty"`
					Id      *string `json:"id,omitempty"` // revive:disable-line
					Link    *string `json:"link,omitempty"`
				}{Id: &testTemplateID},
			},
		}, nil)

	ts.mockGetOperation(&oapi.Operation{
		Id: &testOperationID,
		Reference: &struct {
			Command *string `json:"command,omitempty"`
			Id      *string `json:"id,omitempty"` // revive:disable-line
			Link    *string `json:"link,omitempty"`
		}{Id: &testTemplateID},
		State: &testOperationState,
	})

	ts.mock().
		On("GetTemplateWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testTemplateID, args.Get(1))
		}).
		Return(&oapi.GetTemplateResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.Template{
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
			},
		}, nil)

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
		Zone:            &dstZone,
	}

	actual, err := ts.client.CopyTemplate(context.Background(), testZone, &Template{ID: &testTemplateID}, dstZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_DeleteTemplate() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		deleted            = false
	)

	ts.mock().
		On(
			"DeleteTemplateWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testTemplateID, args.Get(1))
			deleted = true
		}).
		Return(
			&oapi.DeleteTemplateResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id: &testOperationID,
					Reference: &struct {
						Command *string `json:"command,omitempty"`
						Id      *string `json:"id,omitempty"` // revive:disable-line
						Link    *string `json:"link,omitempty"`
					}{Id: &testTemplateID},
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
		}{Id: &testTemplateID},
		State: &testOperationState,
	})

	ts.Require().NoError(ts.client.DeleteTemplate(context.Background(), testZone, &Template{ID: &testTemplateID}))
	ts.Require().True(deleted)
}

func (ts *testSuite) TestClient_GetTemplate() {
	ts.mock().
		On("GetTemplateWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testTemplateID, args.Get(1))
		}).
		Return(&oapi.GetTemplateResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.Template{
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
			},
		}, nil)

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
		Zone:            &testZone,
	}

	actual, err := ts.client.GetTemplate(context.Background(), testZone, *expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_ListTemplates() {
	ts.mock().
		On("ListTemplatesWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // params
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				&oapi.ListTemplatesParams{
					Visibility: (*oapi.ListTemplatesParamsVisibility)(&testTemplateVisibility),
					Family:     &testTemplateFamily,
				},
				args.Get(1),
			)
		}).
		Return(&oapi.ListTemplatesResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &struct {
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
			},
		}, nil)

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
		Zone:            &testZone,
	}}

	actual, err := ts.client.ListTemplates(
		context.Background(),
		testZone,
		ListTemplatesWithVisibility(testTemplateVisibility),
		ListTemplatesWithFamily(testTemplateFamily),
	)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_RegisterTemplate() {
	var (
		templateVisibility = "private"
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
	)

	ts.mock().
		On("RegisterTemplateWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				oapi.RegisterTemplateJSONRequestBody{
					BootMode:        (*oapi.RegisterTemplateJSONBodyBootMode)(&testTemplateBootMode),
					Checksum:        testTemplateChecksum,
					DefaultUser:     &testTemplateDefaultUser,
					Description:     &testTemplateDescription,
					Name:            testTemplateName,
					PasswordEnabled: testTemplatePasswordEnabled,
					SshKeyEnabled:   testTemplateSSHKeyEnabled,
					Url:             testTemplateURL,
				},
				args.Get(1),
			)
		}).
		Return(&oapi.RegisterTemplateResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.Operation{
				Id:    &testOperationID,
				State: &testOperationState,
				Reference: &struct {
					Command *string `json:"command,omitempty"`
					Id      *string `json:"id,omitempty"` // revive:disable-line
					Link    *string `json:"link,omitempty"`
				}{Id: &testTemplateID},
			},
		}, nil)

	ts.mockGetOperation(&oapi.Operation{
		Id: &testOperationID,
		Reference: &struct {
			Command *string `json:"command,omitempty"`
			Id      *string `json:"id,omitempty"` // revive:disable-line
			Link    *string `json:"link,omitempty"`
		}{Id: &testTemplateID},
		State: &testOperationState,
	})

	ts.mock().
		On("GetTemplateWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testTemplateID, args.Get(1))
		}).
		Return(&oapi.GetTemplateResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.Template{
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
			},
		}, nil)

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
		Zone:            &testZone,
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

func (ts *testSuite) TestClient_UpdateTemplate() {
	var (
		testTemplateDescriptionUpdated = testTemplateDescription + "-updated"
		testTemplateNameUpdated        = testTemplateName + "-updated"
		testOperationID                = ts.randomID()
		testOperationState             = oapi.OperationStateSuccess
		updated                        = false
	)

	ts.mock().
		On(
			"UpdateTemplateWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				oapi.UpdateTemplateJSONRequestBody{
					Description: &testTemplateDescriptionUpdated,
					Name:        &testTemplateNameUpdated,
				},
				args.Get(2),
			)
			updated = true
		}).
		Return(
			&oapi.UpdateTemplateResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id: &testOperationID,
					Reference: &struct {
						Command *string `json:"command,omitempty"`
						Id      *string `json:"id,omitempty"` // revive:disable-line
						Link    *string `json:"link,omitempty"`
					}{Id: &testTemplateID},
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
		}{Id: &testTemplateID},
		State: &testOperationState,
	})

	ts.Require().NoError(ts.client.UpdateTemplate(context.Background(), testZone, &Template{
		Description: &testTemplateDescriptionUpdated,
		ID:          &testTemplateID,
		Name:        &testTemplateNameUpdated,
	}))
	ts.Require().True(updated)
}
