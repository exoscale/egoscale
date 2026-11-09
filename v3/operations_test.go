package v3

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/exoscale/egoscale/v3/credentials"
)

func TestGetSubnetInstanceWithEmptyIPv4(t *testing.T) {
	const (
		vpcID      UUID = "11111111-1111-1111-1111-111111111111"
		subnetID   UUID = "22222222-2222-2222-2222-222222222222"
		instanceID UUID = "33333333-3333-3333-3333-333333333333"
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wantPath := fmt.Sprintf("/vpc/%s/subnet/%s", vpcID, subnetID)
		if r.Method != http.MethodGet || r.URL.Path != wantPath {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"id": %q,
			"name": "my-subnet",
			"ipv4-block": "10.0.0.0/24",
			"instances": [
				{"id": %q, "name": "my-instance", "ipv4": ""}
			]
		}`, subnetID, instanceID)
	}))
	defer server.Close()

	client, err := NewClient(
		credentials.NewStaticCredentials("EXOtest", "secret"),
		ClientOptWithEndpoint(Endpoint(server.URL)),
	)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	subnet, err := client.GetSubnet(context.Background(), vpcID, subnetID)
	if err != nil {
		t.Fatalf("GetSubnet: unexpected error: %v", err)
	}

	if len(subnet.Instances) != 1 {
		t.Fatalf("expected 1 instance, got %d", len(subnet.Instances))
	}

	instance := subnet.Instances[0]
	if instance.ID != instanceID {
		t.Errorf("instance ID = %q, want %q", instance.ID, instanceID)
	}
	if instance.Name != "my-instance" {
		t.Errorf("instance name = %q, want %q", instance.Name, "my-instance")
	}
	if instance.Ipv4 != nil {
		t.Errorf("instance ipv4 = %v, want nil", instance.Ipv4)
	}
}
