package v3

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestHandleHTTPErrorResp(t *testing.T) {
	tests := []struct {
		statusCode    int
		body          string
		expectError   bool
		expectedInMsg string
		description   string
	}{
		{
			statusCode:    400,
			body:          `{"message":"something went wrong"}`,
			expectError:   true,
			expectedInMsg: "something went wrong",
			description:   "Legacy message field should be used",
		},
		{
			statusCode:    400,
			body:          `{"type":"#/http-problem-types/bad-request","title":"Bad Request","detail":"required key [id] not found"}`,
			expectError:   true,
			expectedInMsg: "required key [id] not found",
			description:   "RFC 9457 detail field should be used",
		},
		{
			statusCode:    409,
			body:          `{"type":"#/http-problem-types/conflict","title":"Conflict"}`,
			expectError:   true,
			expectedInMsg: "Conflict",
			description:   "RFC 9457 title field as fallback",
		},
		{
			statusCode:    200,
			body:          `{"id":"abc"}`,
			expectError:   false,
			expectedInMsg: "",
			description:   "2xx responses should pass through",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()
			resp := &http.Response{
				StatusCode: test.statusCode,
				Body:       io.NopCloser(strings.NewReader(test.body)),
			}
			err := handleHTTPErrorResp(resp)

			if test.expectError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !strings.Contains(err.Error(), test.expectedInMsg) {
					t.Errorf("expected %q in error, got: %s", test.expectedInMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got: %s", err.Error())
				}
			}
		})
	}
}
