package compute

import "github.com/exoscale/egoscale/v3/oapi"

// Compute provides access to [Exoscale Compute] API resources.
//
// [Exoscale Compute]: https://community.exoscale.com/documentation/compute/
type Compute struct {
	oapiClient *oapi.ClientWithResponses
}

// NewCompute initializes Compute with provided oapi Client.
func NewCompute(c *oapi.ClientWithResponses) *Compute {
	return &Compute{c}
}

//func (a *Compute) Instances() *Instaces {
//return NewInstances(a.oapiClient)
//}