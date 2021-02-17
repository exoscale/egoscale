package publicapi

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTemplate_UnmarshalJSON(t *testing.T) {
	var (
		testID                    = "c19542b7-d269-4bd4-bf7c-2cae36d066d3"
		testName                  = "Linux Ubuntu 20.04 LTS 64-bit"
		testCreatedAt, _          = time.Parse(iso8601Format, "2020-08-12T11:12:36Z")
		testDescription           = "Linux Ubuntu 20.04 LTS 64-bit 2020-08-11-e15f6a"
		testBootMode              = "uefi"
		testBuild                 = "2020-08-12-e15f6a"
		testChecksum              = "e15f6a918f0a645b14238a01568d8cb7"
		testDefaultUser           = "ubuntu"
		testFamily                = "20.04 lts"
		testPasswordEnabled       = true
		testSize            int64 = 10737418240
		testSSHKeyEnabled         = true
		testURL                   = "https://ubuntu.com/"
		testVersion               = "20.04 lts"
		testVisibility            = "public"

		expected = Template{
			BootMode:        &testBootMode,
			Build:           &testBuild,
			Checksum:        &testChecksum,
			CreatedAt:       &testCreatedAt,
			DefaultUser:     &testDefaultUser,
			Description:     &testDescription,
			Family:          &testFamily,
			Id:              &testID,
			Name:            &testName,
			PasswordEnabled: &testPasswordEnabled,
			Size:            &testSize,
			SshKeyEnabled:   &testSSHKeyEnabled,
			Url:             &testURL,
			Version:         &testVersion,
			Visibility:      &testVisibility,
		}

		actual Template

		jsonTemplate = `{
  "boot-mode": "` + testBootMode + `",
  "build": "` + testBuild + `",
  "checksum": "` + testChecksum + `",
  "created-at": "` + testCreatedAt.Format(iso8601Format) + `",
  "default-user": "` + testDefaultUser + `",
  "description": "` + testDescription + `",
  "family": "` + testFamily + `",
  "id": "` + testID + `",
  "name": "` + testName + `",
  "password-enabled": ` + fmt.Sprint(testPasswordEnabled) + `,
  "size": ` + fmt.Sprint(testSize) + `,
  "ssh-key-enabled": ` + fmt.Sprint(testSSHKeyEnabled) + `,
  "version": "` + testVersion + `",
  "visibility": "` + testVisibility + `",
  "url": "` + testURL + `"
}`
	)

	require.NoError(t, json.Unmarshal([]byte(jsonTemplate), &actual))
	require.Equal(t, expected, actual)
}

func TestTemplate_MarshalJSON(t *testing.T) {
	var (
		testID                    = "c19542b7-d269-4bd4-bf7c-2cae36d066d3"
		testName                  = "Linux Ubuntu 20.04 LTS 64-bit"
		testCreatedAt, _          = time.Parse(iso8601Format, "2020-08-12T11:12:36Z")
		testDescription           = "Linux Ubuntu 20.04 LTS 64-bit 2020-08-11-e15f6a"
		testBootMode              = "uefi"
		testBuild                 = "2020-08-12-e15f6a"
		testChecksum              = "e15f6a918f0a645b14238a01568d8cb7"
		testDefaultUser           = "ubuntu"
		testFamily                = "20.04 lts"
		testPasswordEnabled       = true
		testSize            int64 = 10737418240
		testSSHKeyEnabled         = true
		testVersion               = "20.04 lts"
		testVisibility            = "public"
		testURL                   = "https://ubuntu.com/"

		template = Template{
			BootMode:        &testBootMode,
			Build:           &testBuild,
			Checksum:        &testChecksum,
			CreatedAt:       &testCreatedAt,
			DefaultUser:     &testDefaultUser,
			Description:     &testDescription,
			Family:          &testFamily,
			Id:              &testID,
			Name:            &testName,
			PasswordEnabled: &testPasswordEnabled,
			Size:            &testSize,
			SshKeyEnabled:   &testSSHKeyEnabled,
			Version:         &testVersion,
			Visibility:      &testVisibility,
			Url:             &testURL,
		}

		expected = []byte(`{` +
			`"boot-mode":"` + testBootMode + `",` +
			`"build":"` + testBuild + `",` +
			`"checksum":"` + testChecksum + `",` +
			`"created-at":"` + testCreatedAt.Format(iso8601Format) + `",` +
			`"default-user":"` + testDefaultUser + `",` +
			`"description":"` + testDescription + `",` +
			`"family":"` + testFamily + `",` +
			`"id":"` + testID + `",` +
			`"name":"` + testName + `",` +
			`"password-enabled":` + fmt.Sprint(testPasswordEnabled) + `,` +
			`"size":` + fmt.Sprint(testSize) + `,` +
			`"ssh-key-enabled":` + fmt.Sprint(testSSHKeyEnabled) + `,` +
			`"url":"` + testURL + `",` +
			`"version":"` + testVersion + `",` +
			`"visibility":"` + testVisibility + `"` +
			`}`)
	)

	actual, err := json.Marshal(template)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}
