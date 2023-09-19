// Code generated by v3/generator; DO NOT EDIT.
package recorder

import (
	"context"


	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
v3 "github.com/exoscale/egoscale/v3"
)

type LoadBalancersAPI struct {
    Recordee *v3.LoadBalancersAPI
    Recorder *Recorder
}


func (a *LoadBalancersAPI) List(ctx context.Context) ([]v3.LoadBalancer, error) {
    req := argsToMap()

    resp, err := a.Recordee.List(ctx, )

    writeErr := a.Recorder.RecordCall("LoadBalancersAPI.List", req, resp, err)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

func (a *LoadBalancersAPI) Get(ctx context.Context, id openapi_types.UUID) (*v3.LoadBalancer, error) {
    req := argsToMap(id)

    resp, err := a.Recordee.Get(ctx, id)

    writeErr := a.Recorder.RecordCall("LoadBalancersAPI.Get", req, resp, err)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

func (a *LoadBalancersAPI) Create(ctx context.Context, body v3.CreateLoadBalancerJSONRequestBody) (*v3.Operation, error) {
    req := argsToMap(body)

    resp, err := a.Recordee.Create(ctx, body)

    writeErr := a.Recorder.RecordCall("LoadBalancersAPI.Create", req, resp, err)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

func (a *LoadBalancersAPI) Update(ctx context.Context, id openapi_types.UUID, body v3.UpdateLoadBalancerJSONRequestBody) (*v3.Operation, error) {
    req := argsToMap(id, body)

    resp, err := a.Recordee.Update(ctx, id, body)

    writeErr := a.Recorder.RecordCall("LoadBalancersAPI.Update", req, resp, err)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

func (a *LoadBalancersAPI) Delete(ctx context.Context, id openapi_types.UUID) (*v3.Operation, error) {
    req := argsToMap(id)

    resp, err := a.Recordee.Delete(ctx, id)

    writeErr := a.Recorder.RecordCall("LoadBalancersAPI.Delete", req, resp, err)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

func (a *LoadBalancersAPI) GetService(ctx context.Context, id openapi_types.UUID, serviceId openapi_types.UUID) (*v3.LoadBalancerService, error) {
    req := argsToMap(id, serviceId)

    resp, err := a.Recordee.GetService(ctx, id, serviceId)

    writeErr := a.Recorder.RecordCall("LoadBalancersAPI.GetService", req, resp, err)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

func (a *LoadBalancersAPI) AddService(ctx context.Context, id openapi_types.UUID, body v3.AddServiceToLoadBalancerJSONRequestBody) (*v3.Operation, error) {
    req := argsToMap(id, body)

    resp, err := a.Recordee.AddService(ctx, id, body)

    writeErr := a.Recorder.RecordCall("LoadBalancersAPI.AddService", req, resp, err)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

func (a *LoadBalancersAPI) UpdateService(ctx context.Context, id openapi_types.UUID, serviceId openapi_types.UUID, body v3.UpdateLoadBalancerServiceJSONRequestBody) (*v3.Operation, error) {
    req := argsToMap(id, serviceId, body)

    resp, err := a.Recordee.UpdateService(ctx, id, serviceId, body)

    writeErr := a.Recorder.RecordCall("LoadBalancersAPI.UpdateService", req, resp, err)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

func (a *LoadBalancersAPI) DeleteService(ctx context.Context, id openapi_types.UUID, serviceId openapi_types.UUID) (*v3.Operation, error) {
    req := argsToMap(id, serviceId)

    resp, err := a.Recordee.DeleteService(ctx, id, serviceId)

    writeErr := a.Recorder.RecordCall("LoadBalancersAPI.DeleteService", req, resp, err)
    if writeErr != nil {
       panic(writeErr)
    }

    return resp, err
}

