package egoscale

import (
	"context"
	"testing"
)

func TestClient_Ping(t *testing.T) {
	c := NewClient("https://api-ch-gva-2.exoscale.com/v2.alpha", "KEY", "SECRET")
	if err := c.Ping(context.Background()); err != nil {
		t.Fatalf("ping failed: %s", err)
	}
}
