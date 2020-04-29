package egoscale

import (
	"context"
	"testing"
	"time"
)

func TestClient_Ping(t *testing.T) {
	c := NewClient("https://api-ch-gva-2.exoscale.com/v2.alpha", "KEY", "SECRET")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	if err := c.Ping(ctx); err != nil {
		t.Fatalf("ping failed: %s", err)
	}
}
