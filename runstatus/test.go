package runstatus

import (
	"context"
	"math/rand"
	"os"
	"time"

	egoapi "github.com/exoscale/egoscale/internal/egoscale"
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
	)

	return NewClient(
		context.Background(),
		apiKey,
		apiSecret,
		&ClientOpts{APIEndpoint: apiEndpoint})
}

func pageFixture(name string) (*egoapi.RunstatusPage, func() error, error) {
	if name == "" {
		name = testPrefix + "-" + testRandomString()
	}

	client, err := testClientFromEnv()
	if err != nil {
		return nil, nil, err
	}

	res, err := client.c.CreateRunstatusPage(client.ctx, egoapi.RunstatusPage{
		Name:      name,
		Subdomain: name,
	})
	if err != nil {
		return nil, nil, err
	}

	return res,
		func() error {
			return client.c.DeleteRunstatusPage(client.ctx, *res)
		},
		err
}
