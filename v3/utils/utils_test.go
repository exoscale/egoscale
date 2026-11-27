package utils_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/exoscale/egoscale/v3/utils"
)

func TestPrintHeaders(t *testing.T) {
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
		{
			input:  nil,
			output: "",
		},
		{
			input:  http.Header{},
			output: "",
		},
	}

	for _, testcase := range testcases {
		if got := utils.PrintHeaders(testcase.input); got != testcase.output {
			t.Errorf("expected %q, got %q", testcase.output, got)
		}
	}
}

func TestIsLoggableContentType(t *testing.T) {
	testcases := []struct {
		input  http.Header
		output bool
	}{
		{
			input:  http.Header{"Content-Type": []string{utils.LoggableContentTypes[0]}},
			output: true,
		},
		{
			input:  http.Header{"Content-Type": []string{utils.LoggableContentTypes[len(utils.LoggableContentTypes)-1]}},
			output: true,
		},
		{
			input:  http.Header{"Content-Type": []string{strings.ToUpper(utils.LoggableContentTypes[0]) + `; charset="utf-8"`}},
			output: true,
		},
		{
			input:  nil,
			output: false,
		},
		{
			input:  nil,
			output: false,
		},
	}

	for _, testcase := range testcases {
		if got := utils.IsLoggableContentType(testcase.input); got != testcase.output {
			t.Errorf("expected %v, got %v", testcase.output, got)
		}
	}
}
