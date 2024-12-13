package v3

import (
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
