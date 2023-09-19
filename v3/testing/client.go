package testing

import (
	"flag"
	"fmt"
	"os"
	"path"
	"regexp"
	"sync"
	"testing"

	v3 "github.com/exoscale/egoscale/v3"
	"github.com/exoscale/egoscale/v3/testing/recorder"
	"github.com/exoscale/egoscale/v3/testing/replayer"
)

type ClientIface interface {
	IAM() IAMAPIIface
	// DBaaS() *v3.DBaaSAPI
	// Compute() *v3.ComputeAPI
	// DNS() *v3.DNSAPI
	// Global() *v3.GlobalAPI
}

type TestClient struct {
	Client   *v3.Client
	Recorder *recorder.Recorder
	Replayer *replayer.Replayer
}

// IAM provides access to IAM resources on Exoscale platform.
func (tc *TestClient) IAM() IAMAPIIface {
	if tc.Recorder != nil && tc.Replayer != nil {
		panic("can't record and replay at the same time")
	}

	if tc.Recorder != nil {
		return &IAMAPIRecorder{
			Recorder: tc.Recorder,
			client:   tc,
		}
	} else {
		return &IAMAPIReplayer{
			Replayer: tc.Replayer,
		}
	}
}

// // DBaaS provides access to DBaaS resources on Exoscale platform.
// func (c *TestClient) DBaaS() *DBaaSAPI {
// 	return &v3.DBaaSAPI{c}
// }

// // Compute provides access to Compute resources on Exoscale platform.
// func (c *TestClient) Compute() *ComputeAPI {
// 	return &v3.ComputeAPI{c}
// }

// // DNS provides access to DNS resources on Exoscale platform.
// func (c *TestClient) DNS() *DNSAPI {
// 	return &v3.DNSAPI{c}
// }

// // Global provides access to global resources on Exoscale platform.
// func (c *TestClient) Global() *GlobalAPI {
// 	return &v3.GlobalAPI{c}
// }

var (
	recordCalls     bool
	parseFlags      sync.Once
	testdataDirName = "testdata"
)

func init() {
	flag.BoolVar(&recordCalls, "record-calls", false, "record calls to egoscale during tests")
}

func getRecordingFlag() bool {
	parseFlags.Do(func() {
		flag.Parse()
	})

	return recordCalls
}

func createTestdataDir() error {
	if _, err := os.Stat(testdataDirName); os.IsNotExist(err) {
		err := os.Mkdir(testdataDirName, os.FileMode(0755))
		if err != nil {
			return fmt.Errorf("failed to create testdata directory: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check existence of testdata directory: %w", err)
	}

	return nil
}

func normalizeFilename(s string) string {
	// Replace whitespace with underscores
	re := regexp.MustCompile(`\s+`)
	s = re.ReplaceAllString(s, "_")

	// Replace invalid characters with underscores
	invalidChars := regexp.MustCompile(`[\/:*?"<>|]`)
	s = invalidChars.ReplaceAllString(s, "_")

	// Truncate to 255 characters as that's a limit on some filesystems
	if len(s) > 255 {
		s = s[:255]
	}

	return s
}

func NewClient(t *testing.T, initializer func() (*v3.Client, error)) (ClientIface, error) {
	var record bool = getRecordingFlag()

	recordingPath := path.Join(testdataDirName, normalizeFilename(t.Name()+".json"))
	if record {
		c, err := initializer()
		if err != nil {
			return nil, err
		}

		if err := createTestdataDir(); err != nil {
			return nil, err
		}

		return &TestClient{
			Client: c,
			Recorder: &recorder.Recorder{
				Filename: recordingPath,
			},
		}, nil
	}

	return &TestClient{
		Replayer: &replayer.Replayer{
			Filename: recordingPath,
		},
	}, nil
}
