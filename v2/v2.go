// Package v2 is the new Exoscale client API binding.
// Reference: https://openapi-v2.exoscale.com/
package v2

import (
	"context"
)

type getter interface {
	get(ctx context.Context, client *Client, zone, id string) (interface{}, error)
}
