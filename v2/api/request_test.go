package api

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testReqEndpointEnv  = "testapi"
	testReqEndpointZone = "xx-test-1"
)

func TestNewReqEndpoint(t *testing.T) {
	re := NewReqEndpoint("", testReqEndpointZone)
	require.Equal(t, defaultReqEndpointEnv, re.env)
	require.Equal(t, testReqEndpointZone, re.zone)

	re = NewReqEndpoint(testReqEndpointEnv, testReqEndpointZone)
	require.Equal(t, testReqEndpointEnv, re.env)
	require.Equal(t, testReqEndpointZone, re.zone)
}

func TestReqEndpoint_Env(t *testing.T) {
	re := NewReqEndpoint(testReqEndpointEnv, testReqEndpointZone)
	require.Equal(t, testReqEndpointEnv, re.Env())
}

func TestReqEndpoint_Zone(t *testing.T) {
	re := NewReqEndpoint(testReqEndpointEnv, testReqEndpointZone)
	require.Equal(t, testReqEndpointZone, re.Zone())
}

func TestReqEndpoint_Host(t *testing.T) {
	testHost := fmt.Sprintf("%s-%s.exoscale.com",
		testReqEndpointEnv,
		testReqEndpointZone)

	re := NewReqEndpoint(testReqEndpointEnv, testReqEndpointZone)
	require.Equal(t, testHost, re.Host())
}

func TestWithEndpoint(t *testing.T) {
	var (
		ctx             = context.Background()
		testReqEndpoint = NewReqEndpoint(testReqEndpointEnv, testReqEndpointZone)
	)

	ctx = WithEndpoint(ctx, testReqEndpoint)
	require.Equal(t, ctx.Value(ReqEndpoint{}), testReqEndpoint)
}

func TestWithZone(t *testing.T) {
	ctx := context.Background()

	ctx = WithZone(ctx, testReqEndpointZone)
	require.Equal(t, ctx.Value(ReqEndpoint{}), ReqEndpoint{
		env:  defaultReqEndpointEnv,
		zone: testReqEndpointZone,
	})
}
