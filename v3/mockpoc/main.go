package main

import (
	"context"
	"fmt"
	"log"

	v3 "github.com/exoscale/egoscale/v3"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	client, err := v3.DefaultClient(v3.ClientOptWithCredentialsFromEnv())
	check(err)

	testKeyName := "egomock-test-key"

	ctx := context.Background()
	accKeys := client.IAM().AccessKey()
	iamKey, err := accKeys.Create(ctx, v3.CreateAccessKeyJSONRequestBody{
		Name: v3.FromString(testKeyName),
	})
	check(err)

	fmt.Printf("iamKey.Name: %v\n", *iamKey.Name)
	fmt.Printf("iamKey.Secret: %v\n", *iamKey.Secret)

	keys, err := accKeys.List(ctx)
	check(err)

	for _, key := range keys {
		if *key.Name == testKeyName {
			fmt.Printf("removing Access Key: %s\n", *key.Key)

			revocation, err := accKeys.Revoke(ctx, *key.Key)
			check(err)

			fmt.Printf("o: %v\n", revocation)
		}
	}
}
