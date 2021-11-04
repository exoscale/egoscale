package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type testHandler struct {
	resStatus int
	resText   string
}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(h.resStatus)
	_, _ = w.Write([]byte(h.resText))
}

func TestErrorHandlerMiddleware_RoundTrip(t *testing.T) {
	tests := []struct {
		name     string
		handler  *testHandler
		testFunc func(t *testing.T, res *http.Response, err error)
	}{
		{
			name:    "ErrNotFound",
			handler: &testHandler{resStatus: http.StatusNotFound},
			testFunc: func(t *testing.T, res *http.Response, err error) {
				require.ErrorIs(t, err, ErrNotFound)
				require.Nil(t, res)
			},
		},
		{
			name:    "ErrInvalidRequest",
			handler: &testHandler{resStatus: http.StatusBadRequest},
			testFunc: func(t *testing.T, res *http.Response, err error) {
				require.ErrorIs(t, err, ErrInvalidRequest)
				require.Nil(t, res)
			},
		},
		{
			name:    "ErrAPIError",
			handler: &testHandler{resStatus: http.StatusInternalServerError},
			testFunc: func(t *testing.T, res *http.Response, err error) {
				require.ErrorIs(t, err, ErrAPIError)
				require.Nil(t, res)
			},
		},
		{
			name:    "OK",
			handler: &testHandler{resStatus: http.StatusOK, resText: "test"},
			testFunc: func(t *testing.T, res *http.Response, err error) {
				require.NoError(t, err)
				actual, _ := ioutil.ReadAll(res.Body)
				require.Equal(t, []byte("test"), actual)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testServer := httptest.NewServer(test.handler)
			defer testServer.Close()

			testClient := testServer.Client()
			testClient.Transport = &ErrorHandlerMiddleware{next: testClient.Transport}

			res, err := testClient.Get(testServer.URL)
			test.testFunc(t, res, err)
		})
	}
}
