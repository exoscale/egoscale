package v3

import (
	"encoding/json"
	"testing"
)

func TestInstancePoolUnmarshalErrorReasonAndState(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		check   func(t *testing.T, p InstancePool)
	}{
		{
			name:    "no error_reason field",
			payload: `{"name":"pool1","state":"running"}`,
			check: func(t *testing.T, p InstancePool) {
				if p.ErrorReason != nil {
					t.Fatalf("expected ErrorReason nil, got %+v", p.ErrorReason)
				}
				if p.State != InstancePoolStateRunning {
					t.Fatalf("expected running, got %q", p.State)
				}
			},
		},
		{
			name:    "state degraded",
			payload: `{"name":"pool1","state":"degraded"}`,
			check: func(t *testing.T, p InstancePool) {
				if p.State != InstancePoolStateDegraded {
					t.Fatalf("expected degraded, got %q", p.State)
				}
			},
		},
		{
			name:    "empty error_reason object",
			payload: `{"error_reason":{}}`,
			check: func(t *testing.T, p InstancePool) {
				if p.ErrorReason == nil {
					t.Fatal("expected non-nil ErrorReason")
				}
				if p.ErrorReason.Cause != "" || p.ErrorReason.Type != "" || p.ErrorReason.JobID != "" {
					t.Fatalf("expected zero-value fields, got %+v", p.ErrorReason)
				}
			},
		},
		{
			name:    "only cause",
			payload: `{"error_reason":{"cause":"boom"}}`,
			check: func(t *testing.T, p InstancePool) {
				if p.ErrorReason == nil || p.ErrorReason.Cause != "boom" {
					t.Fatalf("expected cause=boom, got %+v", p.ErrorReason)
				}
				if p.ErrorReason.Type != "" || p.ErrorReason.JobID != "" {
					t.Fatalf("expected other fields empty, got %+v", p.ErrorReason)
				}
			},
		},
		{
			name:    "only type",
			payload: `{"error_reason":{"type":"quota-exceeded"}}`,
			check: func(t *testing.T, p InstancePool) {
				if p.ErrorReason == nil || p.ErrorReason.Type != "quota-exceeded" {
					t.Fatalf("expected type=quota-exceeded, got %+v", p.ErrorReason)
				}
			},
		},
		{
			name:    "only job-id",
			payload: `{"error_reason":{"job-id":"abc-123"}}`,
			check: func(t *testing.T, p InstancePool) {
				if p.ErrorReason == nil || p.ErrorReason.JobID != "abc-123" {
					t.Fatalf("expected job-id=abc-123, got %+v", p.ErrorReason)
				}
			},
		},
		{
			name:    "all fields populated with degraded state",
			payload: `{"state":"degraded","error_reason":{"cause":"boom","type":"quota-exceeded","job-id":"abc-123"}}`,
			check: func(t *testing.T, p InstancePool) {
				if p.State != InstancePoolStateDegraded {
					t.Fatalf("expected degraded, got %q", p.State)
				}
				if p.ErrorReason == nil {
					t.Fatal("expected non-nil ErrorReason")
				}
				if p.ErrorReason.Cause != "boom" || p.ErrorReason.Type != "quota-exceeded" || p.ErrorReason.JobID != "abc-123" {
					t.Fatalf("unexpected ErrorReason: %+v", p.ErrorReason)
				}
			},
		},
		{
			name:    "unknown extra field is ignored",
			payload: `{"error_reason":{"cause":"boom","unknown-field":"x"}}`,
			check: func(t *testing.T, p InstancePool) {
				if p.ErrorReason == nil || p.ErrorReason.Cause != "boom" {
					t.Fatalf("expected cause=boom, got %+v", p.ErrorReason)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var p InstancePool
			if err := json.Unmarshal([]byte(tc.payload), &p); err != nil {
				t.Fatalf("unmarshal error: %v", err)
			}
			tc.check(t, p)
		})
	}
}

func TestInstancePoolErrorReasonRoundTrip(t *testing.T) {
	original := InstancePool{
		Name:  "pool1",
		State: InstancePoolStateDegraded,
		ErrorReason: &InstancePoolErrorReason{
			Cause: "boom",
			Type:  "quota-exceeded",
			JobID: "abc-123",
		},
	}

	buf, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded InstancePool
	if err := json.Unmarshal(buf, &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if decoded.State != InstancePoolStateDegraded {
		t.Fatalf("state mismatch: %q", decoded.State)
	}
	if decoded.ErrorReason == nil {
		t.Fatal("expected non-nil ErrorReason after round-trip")
	}
	if *decoded.ErrorReason != *original.ErrorReason {
		t.Fatalf("mismatch: got %+v, want %+v", decoded.ErrorReason, original.ErrorReason)
	}
}

func TestInstancePoolErrorReasonOmitEmpty(t *testing.T) {
	p := InstancePool{Name: "pool1"}
	buf, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var raw map[string]any
	if err := json.Unmarshal(buf, &raw); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if _, ok := raw["error_reason"]; ok {
		t.Fatalf("expected error_reason omitted when nil, got %s", string(buf))
	}
}
