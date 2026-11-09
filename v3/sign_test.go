package v3

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/exoscale/egoscale/v3/credentials"
)

const (
	testAPIKey    = "EXOabcdef0123456789abcdef01"
	testAPISecret = "secret-0123456789"
)

func testSignature(message string) string {
	h := hmac.New(sha256.New, []byte(testAPISecret))
	_, _ = h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func TestSignHTTPRequest(t *testing.T) {
	expiration := time.Date(2026, 11, 9, 12, 0, 0, 0, time.UTC)
	ts := strconv.FormatInt(expiration.Unix(), 10)

	tests := []struct {
		name       string
		method     string
		url        string
		body       string
		message    string // the signed message, written out
		signedArgs string
	}{
		{
			name:    "no query no body",
			method:  "GET",
			url:     "https://api.example.net/v2/instance",
			message: "GET /v2/instance\n\n\n\n" + ts,
		},
		{
			name:       "single-valued query parameters, sorted by name",
			method:     "GET",
			url:        "https://api.example.net/v2/instance?manager-type=instance-pool&ip-address=1.2.3.4",
			message:    "GET /v2/instance\n\n1.2.3.4instance-pool\n\n" + ts,
			signedArgs: "ip-address;manager-type",
		},
		{
			name:       "repeated parameters are not signed, empty values are",
			method:     "GET",
			url:        "https://api.example.net/v2/x?a=1&a=2&b=3&empty=",
			message:    "GET /v2/x\n\n3\n\n" + ts,
			signedArgs: "b;empty",
		},
		{
			name:    "escaped path",
			method:  "GET",
			url:     "https://api.example.net/v2/dns-domain/a%20b%2Fc/record",
			message: "GET /v2/dns-domain/a%20b%2Fc/record\n\n\n\n" + ts,
		},
		{
			name:       "body and query",
			method:     "DELETE",
			url:        "https://api.example.net/v2/x/1?force=true",
			body:       `{"a":[1,2]}`,
			message:    "DELETE /v2/x/1\n{\"a\":[1,2]}\ntrue\n\n" + ts,
			signedArgs: "force",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body io.Reader
			if tt.body != "" {
				body = strings.NewReader(tt.body)
			}
			req, err := http.NewRequest(tt.method, tt.url, body)
			if err != nil {
				t.Fatal(err)
			}

			if err := SignHTTPRequest(req, testAPIKey, testAPISecret, expiration); err != nil {
				t.Fatal(err)
			}

			want := "EXO2-HMAC-SHA256 credential=" + testAPIKey
			if tt.signedArgs != "" {
				want += ",signed-query-args=" + tt.signedArgs
			}
			want += ",expires=" + ts + ",signature=" + testSignature(tt.message)
			if got := req.Header.Get("Authorization"); got != want {
				t.Errorf("Authorization\n got %s\nwant %s", got, want)
			}

			if req.Body != nil {
				sent, err := io.ReadAll(req.Body)
				if err != nil {
					t.Fatal(err)
				}
				if string(sent) != tt.body {
					t.Errorf("body after signing is %q, want %q", sent, tt.body)
				}
			}
		})
	}
}

// TestClientSignsWithSignHTTPRequest checks that the Client operations sign with SignHTTPRequest.
func TestClientSignsWithSignHTTPRequest(t *testing.T) {
	var sent *http.Request
	var sentBody []byte
	hc := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		sent = r
		sentBody, _ = io.ReadAll(r.Body)
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": {"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{}`)),
			Request:    r,
		}, nil
	})}
	client, err := NewClient(credentials.NewStaticCredentials(testAPIKey, testAPISecret),
		ClientOptWithHTTPClient(hc),
		ClientOptWithEndpoint("https://api.example.net/v2"))
	if err != nil {
		t.Fatal(err)
	}

	if _, err := client.CreateSecurityGroup(context.Background(), CreateSecurityGroupRequest{Name: "web"}); err != nil {
		t.Fatal(err)
	}

	got := sent.Header.Get("Authorization")
	m := regexp.MustCompile(`,expires=([0-9]+),`).FindStringSubmatch(got)
	if m == nil {
		t.Fatalf("Authorization %q has no expires", got)
	}
	ts, _ := strconv.ParseInt(m[1], 10, 64)

	req, _ := http.NewRequest(sent.Method, sent.URL.String(), bytes.NewReader(sentBody))
	if err := SignHTTPRequest(req, testAPIKey, testAPISecret, time.Unix(ts, 0)); err != nil {
		t.Fatal(err)
	}
	if want := req.Header.Get("Authorization"); got != want {
		t.Errorf("Authorization\n got %s\nwant %s", got, want)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
