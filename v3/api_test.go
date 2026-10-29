package v3

import (
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/exoscale/egoscale/v3/credentials"
)

func TestPollInterval(t *testing.T) {
	tests := []struct {
		runTime     time.Duration
		expectedMin time.Duration
		expectedMax time.Duration
		description string
	}{
		{
			runTime:     10 * time.Second,
			expectedMin: 3 * time.Second,
			expectedMax: 3 * time.Second,
			description: "Polling at 10 seconds should return 3 seconds",
		},
		{
			runTime:     30 * time.Second,
			expectedMin: 3 * time.Second,
			expectedMax: 3 * time.Second,
			description: "Polling at 30 seconds should still return 3 seconds",
		},
		{
			runTime:     60 * time.Second,
			expectedMin: 3 * time.Second,
			expectedMax: 7 * time.Second, // Expected range after 30s should increase linearly
			description: "Polling at 60 seconds should return a value greater than 3 seconds but less than 7 seconds",
		},
		{
			runTime:     300 * time.Second,
			expectedMin: 3 * time.Second,
			expectedMax: 24 * time.Second, // Interval keeps increasing linearly
			description: "Polling at 5 minutes should return a value in the correct range (up to 24 seconds)",
		},
		{
			runTime:     900 * time.Second, // 15 minutes
			expectedMin: 60 * time.Second,
			expectedMax: 60 * time.Second,
			description: "Polling at 15 minutes should return exactly 60 seconds",
		},
		{
			runTime:     1200 * time.Second, // 20 minutes
			expectedMin: 60 * time.Second,
			expectedMax: 60 * time.Second,
			description: "Polling beyond 15 minutes should cap at 60 seconds",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			interval := pollInterval(test.runTime)
			if interval < test.expectedMin || interval > test.expectedMax {
				t.Errorf("pollInterval(%v) = %v, expected between %v and %v",
					test.runTime, interval, test.expectedMin, test.expectedMax)
			}
		})
	}
}

// TestUpdateLoadBalancerServiceRequestLabelSelectorStates pins the three meanings of
// the optional label-selector field on update: a nil pointer leaves it untouched, an
// allocated selector with an empty map removes it, a populated map replaces it. This
// only holds because the outer field keeps the JSON "omitempty" tag while the inner
// map does not, which is what lets `{"labels":{}}` survive marshalling.
func TestUpdateLoadBalancerServiceRequestLabelSelectorStates(t *testing.T) {
	tests := []struct {
		name          string
		labelSelector *LoadBalancerServiceLabelSelector
		wantSelector  string // expected "label-selector" value, "" means the key must be absent
	}{
		{
			name:          "nil pointer leaves the selector untouched",
			labelSelector: nil,
			wantSelector:  "",
		},
		{
			name:          "allocated selector with empty labels removes it",
			labelSelector: &LoadBalancerServiceLabelSelector{Labels: Labels{}},
			wantSelector:  `{"labels":{}}`,
		},
		{
			name:          "non-empty labels replace the selector wholesale",
			labelSelector: &LoadBalancerServiceLabelSelector{Labels: Labels{"role": "web"}},
			wantSelector:  `{"labels":{"role":"web"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := prepareJSONBody(UpdateLoadBalancerServiceRequest{
				LabelSelector: tt.labelSelector,
			})
			if err != nil {
				t.Fatalf("prepareJSONBody: %v", err)
			}

			buf, err := io.ReadAll(body)
			if err != nil {
				t.Fatalf("read request body: %v", err)
			}

			fields := map[string]json.RawMessage{}
			if err := json.Unmarshal(buf, &fields); err != nil {
				t.Fatalf("unmarshal request body %q: %v", buf, err)
			}

			got, present := fields["label-selector"]
			switch {
			case tt.wantSelector == "" && present:
				t.Fatalf(`"label-selector" = %s, want the key to be absent (body: %s)`, got, buf)
			case tt.wantSelector != "" && !present:
				t.Fatalf(`"label-selector" key missing, want %s (body: %s)`, tt.wantSelector, buf)
			case tt.wantSelector != "" && string(got) != tt.wantSelector:
				t.Errorf(`"label-selector" = %s, want %s`, got, tt.wantSelector)
			}
		})
	}
}

// TestValidateSelectorOnlyLoadBalancerService verifies that a service without any
// instance pool still validates, "instance-pool" not being required anymore. The
// cross-field rules are enforced by the API only, not by the client.
func TestValidateSelectorOnlyLoadBalancerService(t *testing.T) {
	client, err := NewClient(credentials.NewStaticCredentials("key", "secret"))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	req := AddServiceToLoadBalancerRequest{
		Name:          "web",
		Protocol:      AddServiceToLoadBalancerRequestProtocolTCP,
		Strategy:      AddServiceToLoadBalancerRequestStrategyRoundRobin,
		Port:          443,
		TargetPort:    8443,
		Healthcheck:   &LoadBalancerServiceHealthcheck{Mode: LoadBalancerServiceHealthcheckModeTCP},
		LabelSelector: &LoadBalancerServiceLabelSelector{Labels: Labels{"role": "web"}},
	}

	if err := client.Validate(req); err != nil {
		t.Errorf("Validate: unexpected error: %v", err)
	}
}
