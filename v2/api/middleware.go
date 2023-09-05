package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
)

type Middleware interface {
	http.RoundTripper
}

// ErrorHandlerMiddleware is an Exoscale API HTTP client middleware that
// returns concrete Go errors according to API response errors.
type ErrorHandlerMiddleware struct {
	next http.RoundTripper
}

func NewAPIErrorHandlerMiddleware(next http.RoundTripper) Middleware {
	if next == nil {
		next = http.DefaultTransport
	}

	return &ErrorHandlerMiddleware{next: next}
}

func (m *ErrorHandlerMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := m.next.RoundTrip(req)
	if err != nil {
		// If the request returned a Go error don't bother analyzing the response
		// body, as there probably won't be any (e.g. connection timeout/refused).
		return resp, err
	}

	if resp.StatusCode >= 400 && resp.StatusCode <= 599 {
		var res struct {
			Message string `json:"message"`
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %s", err)
		}

		if json.Valid(data) {
			if err = json.Unmarshal(data, &res); err != nil {
				return nil, fmt.Errorf("error unmarshaling response: %s", err)
			}
		} else {
			res.Message = string(data)
		}

		switch {
		case resp.StatusCode == http.StatusNotFound:
			return nil, ErrNotFound

		case resp.StatusCode >= 400 && resp.StatusCode < 500:
			return nil, fmt.Errorf("%w: %s", ErrInvalidRequest, res.Message)

		case resp.StatusCode >= 500:
			return nil, fmt.Errorf("%w: %s", ErrAPIError, res.Message)
		}
	}

	return resp, err
}

// TraceMiddleware is a client HTTP middleware that dumps HTTP requests and responses content.
type TraceMiddleware struct {
	next http.RoundTripper
}

func NewTraceMiddleware(next http.RoundTripper) Middleware {
	if next == nil {
		next = http.DefaultTransport
	}

	return &TraceMiddleware{next: next}
}

func (t *TraceMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
	if dump, err := httputil.DumpRequest(req, true); err == nil {
		fmt.Fprintf(os.Stderr, ">>> %s\n", dump)
	}

	fmt.Fprintln(os.Stderr, "----------------------------------------------------------------------")

	resp, err := t.next.RoundTrip(req)

	if resp != nil {
		if dump, err := httputil.DumpResponse(resp, true); err == nil {
			fmt.Fprintf(os.Stderr, "<<< %s\n", dump)
		}
	}

	return resp, err
}

// RecordMiddleware is a client HTTP middleware that dumps HTTP requests and responses content.
type RecordMiddleware struct {
	next http.RoundTripper
}

func NewRecordMiddleware(next http.RoundTripper) Middleware {
	if next == nil {
		next = http.DefaultTransport
	}

	return &RecordMiddleware{next: next}
}

func cloneReadCloser(rc io.ReadCloser) (io.ReadCloser, io.ReadCloser, error) {
	if rc == nil {
		return nil, nil, nil
	}

	// Read the data from the original ReadCloser
	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, nil, err
	}

	// Create a new ReadCloser from a bytes.Buffer containing the data
	buffer := bytes.NewBuffer(data)
	newReadCloser := ioutil.NopCloser(buffer)

	// Create a new ReadCloser from the original data
	originalReadCloser := ioutil.NopCloser(bytes.NewReader(data))

	return newReadCloser, originalReadCloser, nil
}

func dumpContent(r io.ReadCloser) map[string]interface{} {
	// if r == nil {
	// 	fmt.Printf("%s: (empty)\n", prefix)

	// 	return
	// }

	content := make(map[string]interface{})
	if r == nil {
		return content
	}

	err := json.NewDecoder(r).Decode(&content)
	if err != nil {
		fmt.Println(err.Error())
	}

	return content

	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	respJSON, err := json.MarshalIndent(content, " ", "    ")
	// 	if err != nil {
	// 		fmt.Println(err.Error())
	// 	} else {
	// 		fmt.Println(prefix+":", string(respJSON))
	// 	}
	// }
}

func (t *RecordMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
	n, o, err := cloneReadCloser(req.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	req.Body = o
	reqContent := dumpContent(n)

	resp, err := t.next.RoundTrip(req)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		n, o, err := cloneReadCloser(resp.Body)
		if err != nil {
			fmt.Println(err.Error())
		}
		resp.Body = o
		respContent := dumpContent(n)

		err = WriteTestdata(reqContent, respContent, resp.StatusCode)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	return resp, err
}
