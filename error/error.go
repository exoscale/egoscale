package error

import "errors"

// ErrMissingAPICredentials represent an error due to missing API credentials.
var ErrMissingAPICredentials = errors.New("missing API key/secret")

// ErrMissingAPIEndpoint represent an error due to missing API endpoint.
var ErrMissingAPIEndpoint = errors.New("missing API endpoint")

// ErrResourceNotFound represents an error due to the requested resource not found.
var ErrResourceNotFound = errors.New("resource not found")

// ErrMissingZone represents an error due to missing zone information.
var ErrMissingZone = errors.New("zone not specified")
