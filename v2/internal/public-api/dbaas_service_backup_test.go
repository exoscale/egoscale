package publicapi

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDbaasServiceBackup_UnmarshalJSON(t *testing.T) {
	var (
		testBackupTime, _       = time.Parse(iso8601Format, "2020-08-12T11:12:36Z")
		testBackupName          = testBackupTime.Format("2006-01-02_15-04_0.00000000.pghoard")
		testDataSize      int64 = 36259840

		expected = DbaasServiceBackup{
			BackupName: testBackupName,
			BackupTime: testBackupTime,
			DataSize:   testDataSize,
		}

		actual DbaasServiceBackup

		jsonDbaasService = `{
  "backup-name": "` + testBackupName + `",
  "backup-time": "` + testBackupTime.Format(iso8601Format) + `",
  "data-size": ` + fmt.Sprint(testDataSize) + `
}`
	)

	require.NoError(t, json.Unmarshal([]byte(jsonDbaasService), &actual))
	require.Equal(t, expected, actual)
}

func TestDbaasServiceBackup_MarshalJSON(t *testing.T) {
	var (
		testBackupTime, _       = time.Parse(iso8601Format, "2020-08-12T11:12:36Z")
		testBackupName          = testBackupTime.Format("2006-01-02_15-04_0.00000000.pghoard")
		testDataSize      int64 = 36259840

		dbaasServiceBackup = DbaasServiceBackup{
			BackupName: testBackupName,
			BackupTime: testBackupTime,
			DataSize:   testDataSize,
		}

		expected = []byte(`{` +
			`"backup-name":"` + testBackupName + `",` +
			`"backup-time":"` + testBackupTime.Format(iso8601Format) + `",` +
			`"data-size":` + fmt.Sprint(testDataSize) +
			`}`)
	)

	actual, err := json.Marshal(dbaasServiceBackup)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}
