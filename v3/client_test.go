package v3

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestGetZoneAPIEndpoint(t *testing.T) {
	c := Client{}
	ctx := context.Background()

	tests := []struct {
		name        string
		zoneName    ZoneName
		want        Endpoint
		wantErr     bool
		errContains string
	}{
		{
			name:     "ch-gva-2",
			zoneName: ZoneNameCHGva2,
			want:     CHGva2,
		},
		{
			name:     "ch-dk-2",
			zoneName: ZoneNameCHDk2,
			want:     CHDk2,
		},
		{
			name:     "de-fra-1",
			zoneName: ZoneNameDEFra1,
			want:     DEFra1,
		},
		{
			name:     "de-muc-1",
			zoneName: ZoneNameDEMuc1,
			want:     DEMuc1,
		},
		{
			name:     "at-vie-1",
			zoneName: ZoneNameATVie1,
			want:     ATVie1,
		},
		{
			name:     "at-vie-2",
			zoneName: ZoneNameATVie2,
			want:     ATVie2,
		},
		{
			name:     "bg-sof-1",
			zoneName: ZoneNameBGSof1,
			want:     BGSof1,
		},
		{
			name:     "hr-zag-1",
			zoneName: ZoneNameHrZag1,
			want:     HrZag1,
		},
		{
			// A hypothetical future zone that doesn't exist in the constants yet
			// should still produce a valid endpoint without needing an API call
			name:     "future zone",
			zoneName: ZoneName("xx-new-1"),
			want:     Endpoint("https://api-xx-new-1.exoscale.com/v2"),
		},
		{
			name:        "empty zone name",
			zoneName:    "",
			wantErr:     true,
			errContains: "zone name is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetZoneAPIEndpoint(ctx, tt.zoneName)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("GetZoneAPIEndpoint(%q) expected error, got nil", tt.zoneName)
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("GetZoneAPIEndpoint(%q) error = %q, want it to contain %q", tt.zoneName, err.Error(), tt.errContains)
				}
				return
			}
			if err != nil {
				t.Fatalf("GetZoneAPIEndpoint(%q) unexpected error: %v", tt.zoneName, err)
			}
			if got != tt.want {
				t.Errorf("GetZoneAPIEndpoint(%q) = %q, want %q", tt.zoneName, got, tt.want)
			}
		})
	}
}

func TestGetZoneName(t *testing.T) {
	c := Client{}
	ctx := context.Background()

	tests := []struct {
		name        string
		endpoint    Endpoint
		want        ZoneName
		wantErr     bool
		errContains string
	}{
		{
			name:     "ch-gva-2",
			endpoint: CHGva2,
			want:     ZoneNameCHGva2,
		},
		{
			name:     "ch-dk-2",
			endpoint: CHDk2,
			want:     ZoneNameCHDk2,
		},
		{
			name:     "de-fra-1",
			endpoint: DEFra1,
			want:     ZoneNameDEFra1,
		},
		{
			name:     "de-muc-1",
			endpoint: DEMuc1,
			want:     ZoneNameDEMuc1,
		},
		{
			name:     "at-vie-1",
			endpoint: ATVie1,
			want:     ZoneNameATVie1,
		},
		{
			name:     "at-vie-2",
			endpoint: ATVie2,
			want:     ZoneNameATVie2,
		},
		{
			name:     "bg-sof-1",
			endpoint: BGSof1,
			want:     ZoneNameBGSof1,
		},
		{
			name:     "hr-zag-1",
			endpoint: HrZag1,
			want:     ZoneNameHrZag1,
		},
		{
			// Round-trip: an endpoint built by GetZoneAPIEndpoint for a future zone
			// should survive a GetZoneName call
			name:     "future zone round-trip",
			endpoint: Endpoint("https://api-xx-new-1.exoscale.com/v2"),
			want:     ZoneName("xx-new-1"),
		},
		{
			name:        "unrecognized endpoint",
			endpoint:    Endpoint("https://example.com/api"),
			wantErr:     true,
			errContains: "not a recognized Exoscale zone endpoint",
		},
		{
			name:        "empty endpoint",
			endpoint:    "",
			wantErr:     true,
			errContains: "not a recognized Exoscale zone endpoint",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetZoneName(ctx, tt.endpoint)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("GetZoneName(%q) expected error, got nil", tt.endpoint)
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("GetZoneName(%q) error = %q, want it to contain %q", tt.endpoint, err.Error(), tt.errContains)
				}
				// Errors about unrecognized endpoints should wrap ErrNotFound
				if !errors.Is(err, ErrNotFound) {
					t.Errorf("GetZoneName(%q) error should wrap ErrNotFound, got: %v", tt.endpoint, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("GetZoneName(%q) unexpected error: %v", tt.endpoint, err)
			}
			if got != tt.want {
				t.Errorf("GetZoneName(%q) = %q, want %q", tt.endpoint, got, tt.want)
			}
		})
	}
}

func TestGetZoneAPIEndpointRoundTrip(t *testing.T) {
	c := Client{}
	ctx := context.Background()

	zones := []ZoneName{
		ZoneNameCHGva2,
		ZoneNameCHDk2,
		ZoneNameDEFra1,
		ZoneNameDEMuc1,
		ZoneNameATVie1,
		ZoneNameATVie2,
		ZoneNameBGSof1,
		ZoneNameHrZag1,
	}

	for _, zone := range zones {
		t.Run(string(zone), func(t *testing.T) {
			endpoint, err := c.GetZoneAPIEndpoint(ctx, zone)
			if err != nil {
				t.Fatalf("GetZoneAPIEndpoint(%q) unexpected error: %v", zone, err)
			}

			roundTripped, err := c.GetZoneName(ctx, endpoint)
			if err != nil {
				t.Fatalf("GetZoneName(%q) unexpected error: %v", endpoint, err)
			}

			if roundTripped != zone {
				t.Errorf("round-trip mismatch: started with %q, got back %q", zone, roundTripped)
			}
		})
	}
}
