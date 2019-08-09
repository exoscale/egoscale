// +build testacc

package compute

import (
	"context"
	"os"

	"github.com/stretchr/testify/suite"
)

// Common test environment information
const (
	testPrefix                      = "test-egoscale"
	testDescription                 = "Created by the egoscale library"
	testZoneName                    = "ch-gva-2"
	testZoneID                      = "1128bd56-b4d9-4ac6-a7b9-c715b187ce11"
	testInstanceServiceOfferingName = "Micro"
	testInstanceServiceOfferingID   = "71004023-bb72-4a97-b1e9-bc66dfce9470" // Micro
	testInstanceTemplateName        = "Linux Ubuntu 18.04 LTS 64-bit"
	testInstanceTemplateID          = "095250e3-7c56-441a-a25b-100a3d3f5a6e" // Linux Ubuntu 18.04 LTS 64-bit
	testReverseDNS                  = "egoscale.exoscale.com."
)

func testClientFromEnv() (*Client, error) {
	var (
		apiKey      = os.Getenv("EXOSCALE_API_KEY")
		apiSecret   = os.Getenv("EXOSCALE_API_SECRET")
		apiEndpoint = os.Getenv("EXOSCALE_COMPUTE_API_ENDPOINT")

		tracing bool
	)

	if os.Getenv("EXOSCALE_TRACE") != "" {
		tracing = true
	}

	return NewClient(
		context.Background(),
		apiKey,
		apiSecret,
		&ClientOpts{
			APIEndpoint: apiEndpoint,
			Tracing:     tracing,
		},
	)
}

type accTestSuite struct {
	suite.Suite

	client *Client
}
