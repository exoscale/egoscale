package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	v3 "github.com/exoscale/egoscale/v3"
)

func TestMock(t *testing.T) {
	client, err := v3.DefaultClient(v3.ClientOptWithCredentialsFromEnv())
	assert.NoError(t, err)

	testKeyName := "egomock-test-key"

	ctx := context.Background()
	accKeys := client.IAM().AccessKey()
	iamKey, err := accKeys.Create(ctx, v3.CreateAccessKeyJSONRequestBody{
		Name: v3.FromString(testKeyName),
	})
	assert.NoError(t, err)

	fmt.Printf("iamKey.Name: %v\n", *iamKey.Name)
	fmt.Printf("iamKey.Secret: %v\n", *iamKey.Secret)

	keys, err := accKeys.List(ctx)
	assert.NoError(t, err)

	for _, key := range keys {
		if *key.Name == testKeyName {
			fmt.Printf("removing Access Key: %s\n", *key.Key)

			revocation, err := accKeys.Revoke(ctx, *key.Key)
			assert.NoError(t, err)

			fmt.Printf("revocation: %v\n", revocation)
		}
	}
}
