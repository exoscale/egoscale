package v3

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/exoscale/egoscale/v3/credentials"
)

func TestWithLogHandlerDoesNotLeakToOriginal(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"quotas":[]}`))
	}))
	defer srv.Close()

	client, err := NewClient(credentials.NewStaticCredentials("EXOxxx", "secret"),
		ClientOptWithEndpoint(Endpoint(srv.URL)))
	if err != nil {
		t.Fatal(err)
	}

	var logs bytes.Buffer
	_ = client.WithLogHandler(slog.NewTextHandler(&logs, &slog.HandlerOptions{Level: slog.LevelDebug}))

	if _, err := client.ListQuotas(t.Context()); err != nil {
		t.Fatal(err)
	}
	if logs.Len() > 0 {
		t.Errorf("original client logged through the clone's handler:\n%s", logs.String())
	}
}
