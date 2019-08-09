package storage

import (
	"context"
	"math/rand"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	// "github.com/aws/aws-sdk-go/service/s3"
	// "github.com/exoscale/egoscale/internal/egoscale"
	// egoapi "github.com/exoscale/egoscale/internal/egoscale"
)

const (
	testPrefix      = "test-egoscale"
	testDescription = "Created by the egoscale library"
)

func testRandomString() string {
	chars := "1234567890abcdefghijklmnopqrstuvwxyz"

	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 10)
	for i := range b {
		b[i] = chars[rand.Int63()%int64(len(chars))]
	}

	return string(b)
}

func testClientFromEnv() (*Client, error) {
	var (
		apiKey      = os.Getenv("EXOSCALE_API_KEY")
		apiSecret   = os.Getenv("EXOSCALE_API_SECRET")
		apiEndpoint = os.Getenv("EXOSCALE_STORAGE_API_ENDPOINT")
		zone        = os.Getenv("EXOSCALE_STORAGE_ZONE")
	)

	return NewClient(
		context.Background(),
		apiKey,
		apiSecret,
		&ClientOpts{
			APIEndpoint: apiEndpoint,
			Zone:        zone,
		})
}

func bucketFixture(name string) (string, func() error, error) {
	if name == "" {
		name = testPrefix + "-" + testRandomString()
	}

	client, err := testClientFromEnv()
	if err != nil {
		return "", nil, err
	}

	res, err := client.c.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String(name)})
	if err != nil {
		return "", nil, err
	}

	return aws.StringValue(res.Location),
		func() error {
			_, err := client.c.DeleteBucket(&s3.DeleteBucketInput{Bucket: res.Location})
			return err
		},
		err
}
