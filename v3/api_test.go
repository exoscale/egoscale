package v3

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
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

func TestDateJSON(t *testing.T) {
	type payload struct {
		Ptr *Date `json:"ptr,omitempty"`
		Val Date  `json:"val,omitempty"`
	}

	date := func(s string) *Date { d := Date(s); return &d }

	tests := []struct {
		name         string
		input        string
		expected     payload
		expectedJSON string
	}{
		{
			name:         "date in a value field",
			input:        `{"val":"2024-02-29"}`,
			expected:     payload{Val: "2024-02-29"},
			expectedJSON: `{"val":"2024-02-29"}`,
		},
		{
			name:         "date in a pointer field",
			input:        `{"ptr":"2026-11-17"}`,
			expected:     payload{Ptr: date("2026-11-17")},
			expectedJSON: `{"ptr":"2026-11-17"}`,
		},
		{
			name:         "unset dates are omitted",
			input:        `{}`,
			expected:     payload{},
			expectedJSON: `{}`,
		},
		{
			name:         "null dates are left unset",
			input:        `{"ptr":null,"val":null}`,
			expected:     payload{},
			expectedJSON: `{}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var actual payload
			if err := json.Unmarshal([]byte(test.input), &actual); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(actual, test.expected) {
				t.Errorf("unmarshal %s = %+v, expected %+v", test.input, actual, test.expected)
			}

			actualJSON, err := json.Marshal(actual)
			if err != nil {
				t.Fatal(err)
			}
			if string(actualJSON) != test.expectedJSON {
				t.Errorf("marshal %+v = %s, expected %s", actual, actualJSON, test.expectedJSON)
			}
		})
	}
}
