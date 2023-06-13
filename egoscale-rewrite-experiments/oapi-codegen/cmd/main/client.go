package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/sauterp/egoscale-rewrite-experiments/oapi-codegen/gen/oapi"

	"github.com/exoscale/egoscale/v2/api"
)

var defaultHTTPClient = func() *http.Client {
	rc := retryablehttp.NewClient()
	// silence client by default
	rc.Logger = log.New(io.Discard, "", 0)
	return rc.StandardClient()
}()

func MultiRequestsEditor(fns ...oapi.RequestEditorFn) oapi.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		for _, fn := range fns {
			if err := fn(ctx, req); err != nil {
				return err
			}
		}

		return nil
	}
}

func newClient() (*oapi.ClientWithResponses, error) {
	apiKey := os.Getenv("EXOSCALE_API_KEY")
	apiSecret := os.Getenv("EXOSCALE_API_SECRET")
	apiEndpoint := "https://api-ch-dk-2.exoscale.com/"

	if apiKey == "" || apiSecret == "" {
		return nil, fmt.Errorf("missing or incomplete API credentials")
	}

	httpClient := defaultHTTPClient

	apiSecurityProvider, err := api.NewSecurityProvider(apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize API security provider: %w", err)
	}

	apiURL, err := url.Parse(apiEndpoint)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize API client: %w", err)
	}
	apiURL = apiURL.ResolveReference(&url.URL{Path: api.Prefix})

	oapiOpts := []oapi.ClientOption{
		oapi.WithHTTPClient(httpClient),
		oapi.WithRequestEditorFn(
			MultiRequestsEditor(
				apiSecurityProvider.Intercept,
			),
		),
	}

	c, err := oapi.NewClientWithResponses(apiURL.String(), oapiOpts...)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize API client: %w", err)
	}

	return c, nil
}
