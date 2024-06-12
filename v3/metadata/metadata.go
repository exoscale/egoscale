// This package provides functions to interact with the Exoscale metadata server
// and retrieve user-data (Cloudinit or Ignition data).
package metadata

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	diskfs "github.com/diskfs/go-diskfs"
)

// Endpoint represents different types of metadata
// available on the Exoscale server.
type Endpoint string

// These constants define the various types of
// Exoscale metadata you can retrieve.
// Use the Get function to access specific metadata.
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

	CdRomPath = "/dev/disk/by-label/cidata"
)

// UserData retrieves the user-data associated with the current instance from the Exoscale server.
// This data is typically used for Cloudinit/Ignition configuration.
func UserData(ctx context.Context) (string, error) {
	return httpGet(ctx, UserDataURL)
}

// Get retrieves the value for a specific type of Exoscale metadata.
// Provide the desired Endpoint constant as an argument.
func Get(ctx context.Context, endpoint Endpoint) (string, error) {
	url, err := url.JoinPath(MetaDataURL, string(endpoint))
	if err != nil {
		return "", err
	}

	return httpGet(ctx, url)
}

// FromCdRom retrieves metadata for Exoscale Private Instance,
// from the attached CD-ROM(iso9660) device file system.
// Important note: Run this code as privileged user.
// Not Windows compatible.
func FromCdRom(endpoint Endpoint) (string, error) {
	iso, err := diskfs.Open(CdRomPath, diskfs.WithOpenMode(diskfs.ReadOnly))
	if err != nil {
		return "", fmt.Errorf("disk open: %w", err)
	}
	defer iso.File.Close()

	return getFileMetaDataValue(iso.File, string(endpoint))
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

func getFileMetaDataValue(file *os.File, endpoint string) (string, error) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		if strings.Contains(parts[0], endpoint) {
			return strings.TrimSpace(parts[1]), nil
		}
	}

	return "", fmt.Errorf("endpoint '%s' not found in file", endpoint)
}
