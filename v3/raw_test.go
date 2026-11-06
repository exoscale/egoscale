package v3

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"strings"
	"testing"

	"github.com/exoscale/egoscale/v3/credentials"
)

const (
	rawTestKey    = "EXOtest"
	rawTestSecret = "test-secret"
)

// verifySignature checks an EXO2-HMAC-SHA256 Authorization header independently of
// signRequest, following the documented scheme.
func verifySignature(r *http.Request, body []byte) error {
	auth := strings.TrimPrefix(r.Header.Get("Authorization"), "EXO2-HMAC-SHA256 ")
	fields := map[string]string{}
	for _, part := range strings.Split(auth, ",") {
		k, v, _ := strings.Cut(part, "=")
		fields[k] = v
	}
	if fields["credential"] != rawTestKey {
		return errors.New("credential " + fields["credential"])
	}
	var names []string
	if fields["signed-query-args"] != "" {
		names = strings.Split(fields["signed-query-args"], ";")
	}
	if !sort.StringsAreSorted(names) {
		return errors.New("signed query args not sorted")
	}
	values := ""
	for _, n := range names {
		values += r.URL.Query().Get(n)
	}
	msg := strings.Join([]string{r.Method + " " + r.URL.EscapedPath(), string(body), values, "", fields["expires"]}, "\n")
	h := hmac.New(sha256.New, []byte(rawTestSecret))
	h.Write([]byte(msg))
	if want := base64.StdEncoding.EncodeToString(h.Sum(nil)); fields["signature"] != want {
		return errors.New("bad signature")
	}
	return nil
}

type rawRecorded struct {
	method, path, rawQuery, contentType, userAgent, intercepted string
	body                                                        []byte
}

func newRawTestClient(t *testing.T, handler func(w http.ResponseWriter, r *http.Request)) (*Client, *[]rawRecorded) {
	t.Helper()
	var reqs []rawRecorded
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if err := verifySignature(r, body); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = io.WriteString(w, `{"message":"`+err.Error()+`"}`)
			return
		}
		reqs = append(reqs, rawRecorded{r.Method, r.URL.EscapedPath(), r.URL.RawQuery, r.Header.Get("Content-Type"), r.Header.Get("User-Agent"), r.Header.Get("X-Intercepted"), body})
		handler(w, r)
	}))
	t.Cleanup(srv.Close)
	c, err := NewClient(credentials.NewStaticCredentials(rawTestKey, rawTestSecret),
		ClientOptWithEndpoint(Endpoint(srv.URL+"/v2")),
		ClientOptWithHTTPClient(srv.Client()),
		ClientOptWithUserAgent("raw-test"),
		ClientOptWithRequestInterceptors(func(_ context.Context, r *http.Request) error {
			r.Header.Set("X-Intercepted", "yes")
			return nil
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	return c, &reqs
}

func TestDoJSON(t *testing.T) {
	c, reqs := newRawTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			_, _ = io.WriteString(w, `{"instances":[{"id":"i-1","name":"web"}]}`)
		default:
			_, _ = io.WriteString(w, `{"id":"op-1","state":"pending","reference":{"id":"sg-1"}}`)
		}
	})
	ctx := context.Background()

	var list map[string]any
	err := c.DoJSON(ctx, RawRequest{
		Method: "get", Path: "/instance",
		Query: url.Values{"ip-address": {"1.2.3.4"}, "manager-type": {"instance-pool"}, "multi": {"a", "b"}},
	}, &list)
	if err != nil {
		t.Fatal(err)
	}
	if list["instances"].([]any)[0].(map[string]any)["name"] != "web" {
		t.Errorf("decoded %v", list)
	}

	// A generated type, a struct body, an escaped path parameter.
	var op Operation
	if err := c.DoJSON(ctx, RawRequest{Method: http.MethodPost, Path: "/security-group/" + url.PathEscape("a b") + "/rules", Body: map[string]any{"protocol": "tcp"}}, &op); err != nil {
		t.Fatal(err)
	}
	if op.ID != "op-1" || op.State != OperationStatePending || op.Reference.ID != "sg-1" {
		t.Errorf("operation %+v", op)
	}
	// Raw JSON bodies are sent as is; a nil v discards the response.
	if err := c.DoJSON(ctx, RawRequest{Method: http.MethodPut, Path: "/x", Body: []byte(`{"a":1}`)}, nil); err != nil {
		t.Fatal(err)
	}

	got := *reqs
	if len(got) != 3 {
		t.Fatalf("%d requests", len(got))
	}
	if got[0].method != "GET" || got[0].path != "/v2/instance" || got[0].rawQuery != "ip-address=1.2.3.4&manager-type=instance-pool&multi=a&multi=b" ||
		got[0].contentType != "" || len(got[0].body) != 0 || !strings.HasPrefix(got[0].userAgent, "raw-test") || got[0].intercepted != "yes" {
		t.Errorf("GET %+v", got[0])
	}
	if got[1].path != "/v2/security-group/a%20b/rules" || got[1].contentType != "application/json" || string(got[1].body) != `{"protocol":"tcp"}` {
		t.Errorf("POST %+v %s", got[1], got[1].body)
	}
	if got[2].method != "PUT" || string(got[2].body) != `{"a":1}` {
		t.Errorf("PUT %+v %s", got[2], got[2].body)
	}
}

func TestDoErrors(t *testing.T) {
	c, reqs := newRawTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = io.WriteString(w, `{"message":"instance not found"}`)
	})
	ctx := context.Background()

	_, err := c.Do(ctx, RawRequest{Method: http.MethodGet, Path: "/instance/x", OperationID: "get-instance"})
	if !errors.Is(err, ErrNotFound) || !strings.HasPrefix(err.Error(), "get-instance: http response: ") {
		t.Fatalf("error %v", err)
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != 404 || apiErr.Message != "instance not found" ||
		!json.Valid(apiErr.Body()) || !strings.Contains(string(apiErr.Body()), "not found") {
		t.Errorf("APIError %+v %s", apiErr, apiErr.Body())
	}

	before := len(*reqs)
	for _, r := range []RawRequest{
		{Path: "/instance"},
		{Method: "GET", Path: "instance"},
		{Method: "GET", Path: "/instance?ip-address=1.2.3.4"},
		{Method: "POST", Path: "/x", Body: func() {}},
	} {
		if _, err := c.Do(ctx, r); err == nil {
			t.Errorf("%+v accepted", r)
		}
	}
	if len(*reqs) != before {
		t.Errorf("invalid requests were sent")
	}
}
