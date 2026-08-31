package v3

import (
	"encoding/json"
	"testing"
)

func TestInstancePoolErrorReasonUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		want    InstancePool
	}{
		{
			name:    "absent",
			payload: `{"state":"running"}`,
			want:    InstancePool{State: InstancePoolStateRunning},
		},
		{
			name:    "empty object",
			payload: `{"error_reason":{}}`,
			want:    InstancePool{ErrorReason: &InstancePoolErrorReason{}},
		},
		{
			name:    "all fields with degraded state",
			payload: `{"state":"degraded","error_reason":{"cause":"boom","type":"quota-exceeded","job-id":"abc-123"}}`,
			want: InstancePool{
				State: InstancePoolStateDegraded,
				ErrorReason: &InstancePoolErrorReason{
					Cause: "boom",
					Type:  "quota-exceeded",
					JobID: "abc-123",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var got InstancePool
			if err := json.Unmarshal([]byte(tc.payload), &got); err != nil {
				t.Fatalf("unmarshal error: %v", err)
			}
			if got.State != tc.want.State {
				t.Errorf("State: got %q, want %q", got.State, tc.want.State)
			}
			switch {
			case tc.want.ErrorReason == nil && got.ErrorReason != nil:
				t.Errorf("ErrorReason: got %+v, want nil", got.ErrorReason)
			case tc.want.ErrorReason != nil && got.ErrorReason == nil:
				t.Errorf("ErrorReason: got nil, want %+v", tc.want.ErrorReason)
			case tc.want.ErrorReason != nil && *got.ErrorReason != *tc.want.ErrorReason:
				t.Errorf("ErrorReason: got %+v, want %+v", got.ErrorReason, tc.want.ErrorReason)
			}
		})
	}
}
