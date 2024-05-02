package metadata

import (
	"context"
	"io"
	"net/http"
	"net/url"
)

type Endpoint string

const (
	AvailabilityZone Endpoint = "availability-zone"
	CloudIdentifier  Endpoint = "cloud-identifier"
	InstanceID       Endpoint = "instance-id"
	LocalHostname    Endpoint = "local-hostname"
	LocalIpv4        Endpoint = "local-ipv4"
	PublicHostname   Endpoint = "public-hostname"
	PublicIpv4       Endpoint = "public-ipv4"
	ServiceOffering  Endpoint = "service-offering"
	VMID             Endpoint = "vm-id"
)

const (
	URL         = "http://metadata.exoscale.com/latest/"
	MetaDataURL = URL + "meta-data"
	UserDataURL = URL + "user-data"
)

func UserData(ctx context.Context) (string, error) {
	return httpGet(ctx, UserDataURL)
}

func Get(ctx context.Context, endpoint Endpoint) (string, error) {
	url, err := url.JoinPath(MetaDataURL, string(endpoint))
	if err != nil {
		return "", err
	}

	return httpGet(ctx, url)
}

func httpGet(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
