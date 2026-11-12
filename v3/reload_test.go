package v3

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/exoscale/egoscale/v3/credentials"
)

const (
	testOldKey    = "EXOaaaaaaaaaaaaaaaaaaaaaaaa"
	testOldSecret = "secret-old_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	testNewKey    = "EXObbbbbbbbbbbbbbbbbbbbbbbb"
	testNewSecret = "secret-new_bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
)

// fakeAPI is an API server accepting a single API key, and recording the requests it gets.
type fakeAPI struct {
	*httptest.Server

	mu        sync.Mutex
	validKey  string
	forbidden string // message of the 403 sent to the valid key, if not empty
	requests  []fakeRequest
}

type fakeRequest struct {
	method, path, apiKey, body string
}

func newFakeAPI(t *testing.T, validKey string) *fakeAPI {
	t.Helper()

	api := &fakeAPI{validKey: validKey}
	api.Server = httptest.NewServer(http.HandlerFunc(api.handle))
	t.Cleanup(api.Close)

	return api
}

func (a *fakeAPI) handle(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	key, _ := signedAPIKey(r)

	a.mu.Lock()
	a.requests = append(a.requests, fakeRequest{method: r.Method, path: r.URL.Path, apiKey: key, body: string(body)})
	validKey, forbidden := a.validKey, a.forbidden
	a.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	switch {
	case key != validKey:
		w.WriteHeader(http.StatusForbidden)
		_, _ = io.WriteString(w, `{"message":"Invalid key or request signature"}`)
	case forbidden != "":
		w.WriteHeader(http.StatusForbidden)
		_, _ = io.WriteString(w, `{"message":"`+forbidden+`"}`)
	case strings.HasSuffix(r.URL.Path, "/quota"):
		_, _ = io.WriteString(w, `{"quotas":[]}`)
	default:
		_, _ = io.WriteString(w, `{"id":"7f2b2a3e-4f42-4d6c-9b3a-111111111111","state":"success"}`)
	}
}

func (a *fakeAPI) received() []fakeRequest {
	a.mu.Lock()
	defer a.mu.Unlock()

	return append([]fakeRequest(nil), a.requests...)
}

// rotatingSource is a credentials.Source whose value the test changes.
type rotatingSource struct {
	value atomic.Value
}

func newRotatingSource(key, secret string) *rotatingSource {
	s := &rotatingSource{}
	s.set(key, secret)

	return s
}

func (s *rotatingSource) set(key, secret string) {
	s.value.Store(credentials.Value{APIKey: key, APISecret: secret})
}

func (s *rotatingSource) read() (credentials.Value, error) {
	return s.value.Load().(credentials.Value), nil
}

func newTestClient(t *testing.T, api *fakeAPI, creds *credentials.Credentials) *Client {
	t.Helper()

	client, err := NewClient(creds,
		ClientOptWithHTTPClient(api.Client()),
		ClientOptWithEndpoint(Endpoint(api.URL)))
	if err != nil {
		t.Fatal(err)
	}

	return client
}

func newReloadableTestCredentials(t *testing.T, source *rotatingSource) *credentials.Credentials {
	t.Helper()

	p, err := credentials.NewReloadableProvider(source.read)
	if err != nil {
		t.Fatal(err)
	}

	return credentials.NewCredentials(p)
}

func checkKeys(t *testing.T, requests []fakeRequest, want ...string) {
	t.Helper()

	if len(requests) != len(want) {
		t.Fatalf("got %d requests, want %d: %+v", len(requests), len(want), requests)
	}
	for i, r := range requests {
		if r.apiKey != want[i] {
			t.Errorf("request %d signed with %q, want %q", i, r.apiKey, want[i])
		}
	}
}

// TestClientUsesReloadedCredentials checks that a client, and its clones, sign the next request
// with the new credentials once they are reloaded.
func TestClientUsesReloadedCredentials(t *testing.T) {
	api := newFakeAPI(t, testOldKey)
	source := newRotatingSource(testOldKey, testOldSecret)
	creds := newReloadableTestCredentials(t, source)
	client := newTestClient(t, api, creds)
	clone := client.WithUserAgent("test")

	if _, err := client.ListQuotas(t.Context()); err != nil {
		t.Fatal(err)
	}

	api.mu.Lock()
	api.validKey = testNewKey
	api.mu.Unlock()
	source.set(testNewKey, testNewSecret)
	if err := creds.Reload(t.Context()); err != nil {
		t.Fatal(err)
	}

	if _, err := client.ListQuotas(t.Context()); err != nil {
		t.Fatal(err)
	}
	if _, err := clone.ListQuotas(t.Context()); err != nil {
		t.Fatal(err)
	}

	checkKeys(t, api.received(), testOldKey, testNewKey, testNewKey)
}

// TestClientRetriesWithReloadedCredentials checks that a request rejected for its API key is
// signed again and retried once, body included, when the credentials changed in the meantime.
func TestClientRetriesWithReloadedCredentials(t *testing.T) {
	t.Run("without logs", func(t *testing.T) { testClientRetriesWithReloadedCredentials(t, false) })
	// The debug logs read the bodies of the requests and responses.
	t.Run("with debug logs", func(t *testing.T) { testClientRetriesWithReloadedCredentials(t, true) })
}

func testClientRetriesWithReloadedCredentials(t *testing.T, debugLogs bool) {
	api := newFakeAPI(t, testNewKey)
	source := newRotatingSource(testOldKey, testOldSecret)
	client := newTestClient(t, api, newReloadableTestCredentials(t, source))
	var logs strings.Builder
	if debugLogs {
		client = client.WithLogHandler(slog.NewTextHandler(&logs, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	// The key was rotated, but the credentials weren't reloaded yet.
	source.set(testNewKey, testNewSecret)

	if _, err := client.RegisterSSHKey(t.Context(), RegisterSSHKeyRequest{Name: "k", PublicKey: "ssh-ed25519 AAAA"}); err != nil {
		t.Fatal(err)
	}

	requests := api.received()
	checkKeys(t, requests, testOldKey, testNewKey)
	if requests[1].method != http.MethodPost || requests[1].body != requests[0].body || requests[1].body == "" {
		t.Errorf("retried request %+v differs from the first one %+v", requests[1], requests[0])
	}

	// The next requests use the new key directly.
	if _, err := client.ListQuotas(t.Context()); err != nil {
		t.Fatal(err)
	}
	checkKeys(t, api.received()[2:], testNewKey)

	if debugLogs && !strings.Contains(logs.String(), "retrying with the reloaded credentials") {
		t.Errorf("retry not logged:\n%s", logs.String())
	}
}

// TestClientDoesntRetryWithSameCredentials checks that a request rejected for its API key is not
// retried when the credentials didn't change, and that the API error is returned.
func TestClientDoesntRetryWithSameCredentials(t *testing.T) {
	for name, creds := range map[string]func(t *testing.T) *credentials.Credentials{
		"reloadable": func(t *testing.T) *credentials.Credentials {
			return newReloadableTestCredentials(t, newRotatingSource(testOldKey, testOldSecret))
		},
		"static": func(*testing.T) *credentials.Credentials {
			return credentials.NewStaticCredentials(testOldKey, testOldSecret)
		},
	} {
		t.Run(name, func(t *testing.T) {
			api := newFakeAPI(t, testNewKey)
			client := newTestClient(t, api, creds(t))

			_, err := client.ListQuotas(t.Context())
			var apiErr *APIError
			if !errors.Is(err, ErrForbidden) || !errors.As(err, &apiErr) || apiErr.Message != "Invalid key or request signature" {
				t.Fatalf("got error %v, want the API error", err)
			}
			checkKeys(t, api.received(), testOldKey)
		})
	}
}

// TestClientDoesntRetryOtherForbidden checks that a 403 not about the API key isn't retried.
func TestClientDoesntRetryOtherForbidden(t *testing.T) {
	api := newFakeAPI(t, testOldKey)
	api.forbidden = "Forbidden: missing permission"
	source := newRotatingSource(testOldKey, testOldSecret)
	client := newTestClient(t, api, newReloadableTestCredentials(t, source))
	source.set(testNewKey, testNewSecret)

	if _, err := client.ListQuotas(t.Context()); !errors.Is(err, ErrForbidden) {
		t.Fatalf("got error %v, want %v", err, ErrForbidden)
	}
	checkKeys(t, api.received(), testOldKey)
}
