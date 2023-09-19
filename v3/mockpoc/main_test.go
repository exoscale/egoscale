package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	v3 "github.com/exoscale/egoscale/v3"
	v3testing "github.com/exoscale/egoscale/v3/testing"
)

func TestMock(t *testing.T) {
	testKeyName := "egomock-test-key"

	ctx := context.Background()

	// recording
	client, err := v3testing.NewClient(t, func() (*v3.Client, error) {
		zc, err := v3.DefaultClient(v3.ClientOptWithCredentialsFromEnv())

		return &zc.Client, err
	})
	assert.NoError(t, err)

	// accKeysClient := tc.IAM().AccessKey()
	// accKeys := recorder.AccessKeyAPI{
	// 	Recordee: accKeysClient,
	// }
	// recording
	accKeys := client.IAM().AccessKey()

	// replaying
	// accKeys := replayer.AccessKeyAPI{}
	// replaying
	createKeyResp, err := accKeys.Create(ctx, v3.CreateAccessKeyJSONRequestBody{
		Name: v3.FromString(testKeyName),
	})
	assert.NoError(t, err)

	getKeyResp, err := accKeys.Get(ctx, *createKeyResp.Key)
	assert.NoError(t, err)

	assert.Equal(t, testKeyName, *getKeyResp.Name)

	listKeysResp, err := accKeys.List(ctx)
	assert.NoError(t, err)

	for _, key := range listKeysResp {
		if *key.Name == testKeyName {
			revocation, err := accKeys.Revoke(ctx, *key.Key)
			assert.NoError(t, err)

			assert.Equal(t, v3.OperationStateSuccess, *revocation.State)
		}
	}
}
