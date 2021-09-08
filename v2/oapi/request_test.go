package oapi

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMultiRequestsEditor(t *testing.T) {
	var (
		testRequestEditorFn = func(k, v string) RequestEditorFn {
			return func(_ context.Context, req *http.Request) error {
				req.Header.Add(k, v)
				return nil
			}
		}
		req, _ = http.NewRequest("GET", "/test", nil)
	)

	multiRequestsEditorFn := MultiRequestsEditor(
		testRequestEditorFn("H1", "v1"),
		testRequestEditorFn("H2", "v2"),
	)

	err := multiRequestsEditorFn(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, http.Header{
		"H1": []string{"v1"},
		"H2": []string{"v2"},
	}, req.Header)
}
