// Code generated by v3/generator; DO NOT EDIT.
package recorder

import (
	"context"

v3 "github.com/exoscale/egoscale/v3"
)

type OrgQuotasAPI struct {
    Recordee *v3.OrgQuotasAPI
}


func (a *OrgQuotasAPI) List(ctx context.Context) ([]v3.Quota, error) {
    return a.Recordee.List(ctx, )
}

func (a *OrgQuotasAPI) Get(ctx context.Context, entity string) (*v3.Quota, error) {
    return a.Recordee.Get(ctx, entity)
}

