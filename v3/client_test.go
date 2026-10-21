package v3_test

import (
	"net/http"
	"testing"

	"github.com/exoscale/egoscale/v3"
)

func TestOnelineHeaders(t *testing.T) {
	testcases := []struct {
		input  http.Header
		output string
	}{
		{
			input:  http.Header{"key1": {"val1", "val2"}},
			output: "key1:val1,val2",
		},
		{
			input:  http.Header{"key2 key3": {"val1", ""}, "key4": {""}},
			output: "key2 key3:val1, key4:",
		},
	}

	for _, testcase := range testcases {
		if got := v3.OnelineHeaders(testcase.input); got != testcase.output {
			t.Errorf("expected %q, got %q", testcase.output, got)
		}
	}
}
