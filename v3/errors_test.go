package v3

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

func makeResp(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func TestHandleHTTPErrorResp(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		body         string
		wantSentinel error
		wantAPIErr   *APIError // nil means only check sentinel + nil-ness
	}{
		{
			name:       "2xx returns nil",
			statusCode: 200,
			body:       `{"message":"ok"}`,
		},
		{
			// RFC 9457 with detail + string errors array
			name:         "400 detail with string errors",
			statusCode:   400,
			body:         `{"type":"#/http-problem-types/bad-request","title":"Bad Request","detail":"One or more vLLM parameters are invalid.","errors":["Parameter 'foo-bar' is not allowed for vLLM 0.19.1"]}`,
			wantSentinel: ErrBadRequest,
			wantAPIErr: &APIError{
				Title:  "Bad Request",
				Detail: "One or more vLLM parameters are invalid.",
				Errors: []APIErrorEntry{{Detail: "Parameter 'foo-bar' is not allowed for vLLM 0.19.1"}},
			},
		},
		{
			// RFC 9457 with detail only, no errors field
			name:         "400 detail only",
			statusCode:   400,
			body:         `{"type":"#/http-problem-types/bad-request","title":"Bad Request","detail":"GPU type must be valid and count must be between 1 and 8."}`,
			wantSentinel: ErrBadRequest,
			wantAPIErr: &APIError{
				Title:  "Bad Request",
				Detail: "GPU type must be valid and count must be between 1 and 8.",
			},
		},
		{
			// RFC 9457 with object errors array carrying location + detail
			name:         "422 object errors with location",
			statusCode:   422,
			body:         `{"type":"#/http-problem-types/request-invalid-body","errors":[{"path":"...","pointer":"...","location":"$['inference-engine-version']","detail":"does not have a value in the enumeration [\"0.12.0\", \"0.19.1\"]"}],"title":"Invalid Request Body"}`,
			wantSentinel: ErrUnprocessableEntity,
			wantAPIErr: &APIError{
				Title: "Invalid Request Body",
				Errors: []APIErrorEntry{
					{
						Location: "$['inference-engine-version']",
						Detail:   `does not have a value in the enumeration ["0.12.0", "0.19.1"]`,
					},
				},
			},
		},
		{
			// RFC 9457 with detail only, no errors field
			name:         "404 detail",
			statusCode:   404,
			body:         `{"type":"#/http-problem-types/not-found","title":"Not Found","detail":"The referenced model was not found."}`,
			wantSentinel: ErrNotFound,
			wantAPIErr: &APIError{
				Title:  "Not Found",
				Detail: "The referenced model was not found.",
			},
		},
		{
			// simple {"message":"..."} format
			name:         "400 simple message field",
			statusCode:   400,
			body:         `{"message":"something went wrong"}`,
			wantSentinel: ErrBadRequest,
			wantAPIErr:   &APIError{Message: "something went wrong"},
		},
		{
			// simple {"error":"..."} format
			name:         "500 simple error field",
			statusCode:   500,
			body:         `{"error":"internal failure"}`,
			wantSentinel: ErrInternalServerError,
			wantAPIErr:   &APIError{Message: "internal failure"},
		},
		{
			// non-JSON body
			name:         "503 plain text body",
			statusCode:   503,
			body:         "service unavailable",
			wantSentinel: ErrServiceUnavailable,
			wantAPIErr:   &APIError{Message: "service unavailable"},
		},
		{
			// unmapped status code still returns an error
			name:       "599 unmapped status code",
			statusCode: 599,
			body:       `{"message":"custom error"}`,
			wantAPIErr: &APIError{Message: "custom error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleHTTPErrorResp(makeResp(tt.statusCode, tt.body))

			if tt.wantSentinel == nil && tt.wantAPIErr == nil {
				if err != nil {
					t.Fatalf("expected nil, got %v", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("expected error, got nil")
			}

			if tt.wantSentinel != nil && !errors.Is(err, tt.wantSentinel) {
				t.Errorf("errors.Is(%v) = false, want true", tt.wantSentinel)
			}

			apiErr, ok := err.(*APIError)
			if !ok {
				t.Fatalf("expected *APIError, got %T", err)
			}

			want := tt.wantAPIErr
			if want.Message != "" && apiErr.Message != want.Message {
				t.Errorf("Message = %q, want %q", apiErr.Message, want.Message)
			}
			if want.Title != "" && apiErr.Title != want.Title {
				t.Errorf("Title = %q, want %q", apiErr.Title, want.Title)
			}
			if want.Detail != "" && apiErr.Detail != want.Detail {
				t.Errorf("Detail = %q, want %q", apiErr.Detail, want.Detail)
			}
			if len(want.Errors) > 0 {
				if len(apiErr.Errors) != len(want.Errors) {
					t.Fatalf("len(Errors) = %d, want %d", len(apiErr.Errors), len(want.Errors))
				}
				for i, we := range want.Errors {
					got := apiErr.Errors[i]
					if we.Location != "" && got.Location != we.Location {
						t.Errorf("Errors[%d].Location = %q, want %q", i, got.Location, we.Location)
					}
					if we.Detail != "" && got.Detail != we.Detail {
						t.Errorf("Errors[%d].Detail = %q, want %q", i, got.Detail, we.Detail)
					}
				}
			}
		})
	}
}

// TestHandleHTTPErrorRespSentinelWrapping verifies that callers can use errors.Is
// to identify the HTTP status class after wrapping, matching the pattern in cli/cmd/root.go.
func TestHandleHTTPErrorRespSentinelWrapping(t *testing.T) {
	cases := []struct {
		statusCode int
		sentinel   error
	}{
		{400, ErrBadRequest},
		{401, ErrUnauthorized},
		{403, ErrForbidden},
		{404, ErrNotFound},
		{409, ErrConflict},
		{422, ErrUnprocessableEntity},
		{429, ErrTooManyRequests},
		{500, ErrInternalServerError},
		{502, ErrBadGateway},
		{503, ErrServiceUnavailable},
		{504, ErrGatewayTimeout},
	}

	for _, c := range cases {
		err := handleHTTPErrorResp(makeResp(c.statusCode, `{"message":"err"}`))
		if err == nil {
			t.Errorf("status %d: expected error, got nil", c.statusCode)
			continue
		}
		wrapped := fmt.Errorf("create instance: %w", err)
		if !errors.Is(wrapped, c.sentinel) {
			t.Errorf("status %d: errors.Is(wrapped, sentinel) = false", c.statusCode)
		}
	}
}
