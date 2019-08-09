package dns

import (
	"context"
	"math/rand"
	"os"
	"time"

	"github.com/exoscale/egoscale/internal/egoscale"
	egoapi "github.com/exoscale/egoscale/internal/egoscale"
)

const (
	testPrefix      = "test-egoscale"
	testDescription = "Created by the egoscale library"
	testDomainName  = "example.net"
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
		apiEndpoint = os.Getenv("EXOSCALE_DNS_API_ENDPOINT")
	)

	return NewClient(
		context.Background(),
		apiKey,
		apiSecret,
		&ClientOpts{APIEndpoint: apiEndpoint})
}

func domainFixture(name string) (*egoapi.DNSDomain, func() error, error) {
	client, err := testClientFromEnv()
	if err != nil {
		return nil, nil, err
	}

	if name == "" {
		name = testDomainName
	}

	res, err := client.c.Request(&egoapi.CreateDNSDomain{Name: name})
	if err != nil {
		return nil, nil, err
	}

	return res.(*egoscale.DNSDomain),
		func() error {
			_, err := client.c.Request(&egoapi.DeleteDNSDomain{Name: name})
			return err
		},
		err
}
