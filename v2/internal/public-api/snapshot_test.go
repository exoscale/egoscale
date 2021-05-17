package publicapi

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSnapshot_UnmarshalJSON(t *testing.T) {
	var (
		testCreatedAt, _ = time.Parse(iso8601Format, "2020-08-12T11:12:36Z")
		testID           = testRandomID(t)
		testInstanceID   = testRandomID(t)
		testName         = "test_ROOT-846459_20200706132534"
		testExportMD5Sum = "c9887de796993c2519b463bcd9509e08"
		testExportURL    = fmt.Sprintf("https://sos-ch-gva-2.exo.io/test/%s/%s", testRandomID(t), testID)
		testState        = SnapshotStateExported

		expected = Snapshot{
			CreatedAt: &testCreatedAt,
			Export: &struct {
				Md5sum       *string `json:"md5sum,omitempty"`
				PresignedUrl *string `json:"presigned-url,omitempty"` // nolint:revive
			}{
				Md5sum:       &testExportMD5Sum,
				PresignedUrl: &testExportURL,
			},
			Id:       &testID,
			Instance: &Instance{Id: &testInstanceID},
			Name:     &testName,
			State:    &testState,
		}

		actual Snapshot

		jsonSnapshot = `{
  "created-at": "` + testCreatedAt.Format(iso8601Format) + `",
  "export": {"md5sum": "` + testExportMD5Sum + `", "presigned-url": "` + testExportURL + `"},
  "id": "` + testID + `",
  "instance": {"id": "` + testInstanceID + `"},
  "name": "` + testName + `",
  "state": "` + string(testState) + `"
}`
	)

	require.NoError(t, json.Unmarshal([]byte(jsonSnapshot), &actual))
	require.Equal(t, expected, actual)
}

func TestSnapshot_MarshalJSON(t *testing.T) {
	var (
		testCreatedAt, _ = time.Parse(iso8601Format, "2020-08-12T11:12:36Z")
		testID           = testRandomID(t)
		testInstanceID   = testRandomID(t)
		testName         = "test_ROOT-846459_20200706132534"
		testExportMD5Sum = "c9887de796993c2519b463bcd9509e08"
		testExportURL    = fmt.Sprintf("https://sos-ch-gva-2.exo.io/test/%s/%s", testRandomID(t), testID)
		testState        = SnapshotStateExported

		snapshot = Snapshot{
			CreatedAt: &testCreatedAt,
			Export: &struct {
				Md5sum       *string `json:"md5sum,omitempty"`
				PresignedUrl *string `json:"presigned-url,omitempty"` // nolint:revive
			}{
				Md5sum:       &testExportMD5Sum,
				PresignedUrl: &testExportURL,
			},
			Id:       &testID,
			Instance: &Instance{Id: &testInstanceID},
			Name:     &testName,
			State:    &testState,
		}

		expected = []byte(`{` +
			`"created-at":"` + testCreatedAt.Format(iso8601Format) + `",` +
			`"export":{"md5sum":"` + testExportMD5Sum + `","presigned-url":"` + testExportURL + `"},` +
			`"id":"` + testID + `",` +
			`"instance":{"id":"` + testInstanceID + `"},` +
			`"name":"` + testName + `",` +
			`"state":"` + string(testState) + `"` +
			`}`)
	)

	actual, err := json.Marshal(snapshot)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}
