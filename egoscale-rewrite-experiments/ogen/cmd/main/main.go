package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/google/uuid"
	"github.com/hashicorp/go-retryablehttp"

	"github.com/exoscale/egoscale/egoscale-rewrite-experiments/ogen/pkg/oas/api"
	v2api "github.com/exoscale/egoscale/v2/api"
)

var defaultHTTPClient = func() *http.Client {
	rc := retryablehttp.NewClient()
	// silence client by default
	rc.Logger = log.New(io.Discard, "", 0)
	return rc.StandardClient()
}()

type RoundTripper struct {
	Client *http.Client

	SecurityProvider *v2api.SecurityProviderExoscale
}

func (rt *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	err := rt.SecurityProvider.Intercept(context.TODO(), req)
	if err != nil {
		return nil, err
	}

	return rt.Client.Do(req)
}

func main() {
	apiKey := os.Getenv("EXOSCALE_API_KEY")
	apiSecret := os.Getenv("EXOSCALE_API_SECRET")
	apiEndpoint := "https://api-ch-dk-2.exoscale.com/"

	if apiKey == "" || apiSecret == "" {
		panic(fmt.Errorf("missing or incomplete API credentials"))
	}

	spe, err := v2api.NewSecurityProvider(apiKey, apiSecret)
	if err != nil {
		panic(err)
	}

	rt := RoundTripper{
		Client: defaultHTTPClient,

		SecurityProvider: spe,
	}

	httpClient := &http.Client{
		Transport: &rt,
	}

	apiURL, err := url.Parse(apiEndpoint)
	if err != nil {
		panic(fmt.Errorf("unable to initialize API client: %w", err))
	}
	apiURL = apiURL.ResolveReference(&url.URL{Path: v2api.Prefix})

	client, err := api.NewClient(apiURL.String(), api.WithClient(httpClient))
	if err != nil {
		panic(err)
	}

	u, err3 := uuid.Parse("62e3c95d-3ab9-4f3a-b68d-97cb6cc48bec")
	if err3 != nil {
		panic(err3)
	}

	ctx := context.Background()
	template, err2 := client.GetTemplate(ctx, api.GetTemplateParams{
		ID: u,
	})
	if err2 != nil {
		panic(err2)
	}

	fmt.Printf("template: %v\n", template)
}
