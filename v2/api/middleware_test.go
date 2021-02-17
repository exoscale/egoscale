package api

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestErrorHandlerMiddleware_RoundTrip(t *testing.T) {
	client := http.DefaultClient

	httpmock.ActivateNonDefault(client)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, "/",
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNotFound, ""), nil
		})

	httpmock.RegisterResponder(http.MethodPost, "/",
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusBadRequest, `{"message":"not this way"}`), nil
		})

	httpmock.RegisterResponder(http.MethodPost, "/broken",
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusInternalServerError, "API is broken"), nil
		})

	httpmock.RegisterResponder(http.MethodGet, "/ok",
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, "test"), nil
		})

	client.Transport = NewAPIErrorHandlerMiddleware(client.Transport)

	// Test for ErrNotFound when receiving a http.StatusNotFound status
	req, err := http.NewRequest(http.MethodGet, "http://example.net/", nil)
	require.NoError(t, err)
	resp, err := client.Do(req)
	require.True(t, errors.Is(err, ErrNotFound))
	require.Nil(t, resp)

	// Test for ErrInvalidRequest when receiving a http.StatusBadRequest status
	req, err = http.NewRequest(http.MethodPost, "http://example.net/", nil)
	require.NoError(t, err)
	resp, err = client.Do(req)
	require.True(t, errors.Is(err, ErrInvalidRequest))
	require.True(t, strings.Contains(err.Error(), "not this way"))
	require.Nil(t, resp)

	// Test for ErrAPIError when receiving a http.StatusInternalServerError status
	req, err = http.NewRequest(http.MethodPost, "http://example.net/broken", nil)
	require.NoError(t, err)
	resp, err = client.Do(req)
	require.True(t, errors.Is(err, ErrAPIError))
	require.Nil(t, resp)

	// Test for successful request
	req, err = http.NewRequest(http.MethodGet, "http://example.net/ok", nil)
	require.NoError(t, err)
	resp, err = client.Do(req)
	require.NoError(t, err)
	actual, _ := ioutil.ReadAll(resp.Body)
	require.Equal(t, []byte("test"), actual)
}
