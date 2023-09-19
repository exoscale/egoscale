package testing

import (
	"testing"

	v3 "github.com/exoscale/egoscale/v3"
)

type ClientIface interface {
	IAM() IAMAPIIface
	// DBaaS() *v3.DBaaSAPI
	// Compute() *v3.ComputeAPI
	// DNS() *v3.DNSAPI
	// Global() *v3.GlobalAPI
}

type TestClient struct {
	Client *v3.Client
}

// IAM provides access to IAM resources on Exoscale platform.
func (c *TestClient) IAM() IAMAPIIface {
	return &IAMAPIRecorder{c}
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

func NewClient(t *testing.T, initializer func() (*v3.Client, error)) (ClientIface, error) {
	c, err := initializer()
	if err != nil {
		return nil, err
	}

	return &TestClient{
		Client: c,
	}, nil
}
