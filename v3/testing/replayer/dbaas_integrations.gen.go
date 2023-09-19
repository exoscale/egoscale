// Code generated by v3/generator; DO NOT EDIT.
package replayer

import (
	"context"

v3 "github.com/exoscale/egoscale/v3"
)

type IntegrationsAPI struct {
     Replayer *Replayer
}


func (a *IntegrationsAPI) ListSettings(ctx context.Context, integrationType string, sourceType string, destType string) (*v3.DBaaSIntegrationSettings, error) {
    resp := InitializeReturnType[*v3.DBaaSIntegrationSettings](a.ListSettings)

    var returnErr error
    writeErr := a.Replayer.GetTestCall(&resp, &returnErr)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, returnErr
}

