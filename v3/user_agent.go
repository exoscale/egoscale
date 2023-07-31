package v3

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
)

// Version string will be embeded into User-Agent header.
// Real value should be set with compiler flag.
const Version = "dev"

// UserAgent is the "User-Agent" HTTP request header added to outgoing HTTP requests.
var UserAgent = fmt.Sprintf("egoscale/%s (%s; %s/%s)",
	Version,
	runtime.Version(),
	runtime.GOOS,
	runtime.GOARCH)

// SetUserAgent is an request editor that adds the "User-Agent" header.
func SetUserAgent(ctx context.Context, req *http.Request) error {
	req.Header.Add("User-Agent", UserAgent)

	return nil
}
