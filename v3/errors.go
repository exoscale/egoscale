package v3

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrBadRequest                    = fmt.Errorf(http.StatusText(http.StatusBadRequest))
	ErrUnauthorized                  = fmt.Errorf(http.StatusText(http.StatusUnauthorized))
	ErrPaymentRequired               = fmt.Errorf(http.StatusText(http.StatusPaymentRequired))
	ErrForbidden                     = fmt.Errorf(http.StatusText(http.StatusForbidden))
	ErrNotFound                      = fmt.Errorf(http.StatusText(http.StatusNotFound))
	ErrMethodNotAllowed              = fmt.Errorf(http.StatusText(http.StatusMethodNotAllowed))
	ErrNotAcceptable                 = fmt.Errorf(http.StatusText(http.StatusNotAcceptable))
	ErrProxyAuthRequired             = fmt.Errorf(http.StatusText(http.StatusProxyAuthRequired))
	ErrRequestTimeout                = fmt.Errorf(http.StatusText(http.StatusRequestTimeout))
	ErrConflict                      = fmt.Errorf(http.StatusText(http.StatusConflict))
	ErrGone                          = fmt.Errorf(http.StatusText(http.StatusGone))
	ErrLengthRequired                = fmt.Errorf(http.StatusText(http.StatusLengthRequired))
	ErrPreconditionFailed            = fmt.Errorf(http.StatusText(http.StatusPreconditionFailed))
	ErrRequestEntityTooLarge         = fmt.Errorf(http.StatusText(http.StatusRequestEntityTooLarge))
	ErrRequestURITooLong             = fmt.Errorf(http.StatusText(http.StatusRequestURITooLong))
	ErrUnsupportedMediaType          = fmt.Errorf(http.StatusText(http.StatusUnsupportedMediaType))
	ErrRequestedRangeNotSatisfiable  = fmt.Errorf(http.StatusText(http.StatusRequestedRangeNotSatisfiable))
	ErrExpectationFailed             = fmt.Errorf(http.StatusText(http.StatusExpectationFailed))
	ErrTeapot                        = fmt.Errorf(http.StatusText(http.StatusTeapot))
	ErrMisdirectedRequest            = fmt.Errorf(http.StatusText(http.StatusMisdirectedRequest))
	ErrUnprocessableEntity           = fmt.Errorf(http.StatusText(http.StatusUnprocessableEntity))
	ErrLocked                        = fmt.Errorf(http.StatusText(http.StatusLocked))
	ErrFailedDependency              = fmt.Errorf(http.StatusText(http.StatusFailedDependency))
	ErrTooEarly                      = fmt.Errorf(http.StatusText(http.StatusTooEarly))
	ErrUpgradeRequired               = fmt.Errorf(http.StatusText(http.StatusUpgradeRequired))
	ErrPreconditionRequired          = fmt.Errorf(http.StatusText(http.StatusPreconditionRequired))
	ErrTooManyRequests               = fmt.Errorf(http.StatusText(http.StatusTooManyRequests))
	ErrRequestHeaderFieldsTooLarge   = fmt.Errorf(http.StatusText(http.StatusRequestHeaderFieldsTooLarge))
	ErrUnavailableForLegalReasons    = fmt.Errorf(http.StatusText(http.StatusUnavailableForLegalReasons))
	ErrInternalServerError           = fmt.Errorf(http.StatusText(http.StatusInternalServerError))
	ErrNotImplemented                = fmt.Errorf(http.StatusText(http.StatusNotImplemented))
	ErrBadGateway                    = fmt.Errorf(http.StatusText(http.StatusBadGateway))
	ErrServiceUnavailable            = fmt.Errorf(http.StatusText(http.StatusServiceUnavailable))
	ErrGatewayTimeout                = fmt.Errorf(http.StatusText(http.StatusGatewayTimeout))
	ErrHTTPVersionNotSupported       = fmt.Errorf(http.StatusText(http.StatusHTTPVersionNotSupported))
	ErrVariantAlsoNegotiates         = fmt.Errorf(http.StatusText(http.StatusVariantAlsoNegotiates))
	ErrInsufficientStorage           = fmt.Errorf(http.StatusText(http.StatusInsufficientStorage))
	ErrLoopDetected                  = fmt.Errorf(http.StatusText(http.StatusLoopDetected))
	ErrNotExtended                   = fmt.Errorf(http.StatusText(http.StatusNotExtended))
	ErrNetworkAuthenticationRequired = fmt.Errorf(http.StatusText(http.StatusNetworkAuthenticationRequired))
)

var httpStatusCodeErrors = map[int]error{
	http.StatusBadRequest:                    ErrBadRequest,
	http.StatusUnauthorized:                  ErrUnauthorized,
	http.StatusPaymentRequired:               ErrPaymentRequired,
	http.StatusForbidden:                     ErrForbidden,
	http.StatusNotFound:                      ErrNotFound,
	http.StatusMethodNotAllowed:              ErrMethodNotAllowed,
	http.StatusNotAcceptable:                 ErrNotAcceptable,
	http.StatusProxyAuthRequired:             ErrProxyAuthRequired,
	http.StatusRequestTimeout:                ErrRequestTimeout,
	http.StatusConflict:                      ErrConflict,
	http.StatusGone:                          ErrGone,
	http.StatusLengthRequired:                ErrLengthRequired,
	http.StatusPreconditionFailed:            ErrPreconditionFailed,
	http.StatusRequestEntityTooLarge:         ErrRequestEntityTooLarge,
	http.StatusRequestURITooLong:             ErrRequestURITooLong,
	http.StatusUnsupportedMediaType:          ErrUnsupportedMediaType,
	http.StatusRequestedRangeNotSatisfiable:  ErrRequestedRangeNotSatisfiable,
	http.StatusExpectationFailed:             ErrExpectationFailed,
	http.StatusTeapot:                        ErrTeapot,
	http.StatusMisdirectedRequest:            ErrMisdirectedRequest,
	http.StatusUnprocessableEntity:           ErrUnprocessableEntity,
	http.StatusLocked:                        ErrLocked,
	http.StatusFailedDependency:              ErrFailedDependency,
	http.StatusTooEarly:                      ErrTooEarly,
	http.StatusUpgradeRequired:               ErrUpgradeRequired,
	http.StatusPreconditionRequired:          ErrPreconditionRequired,
	http.StatusTooManyRequests:               ErrTooManyRequests,
	http.StatusRequestHeaderFieldsTooLarge:   ErrRequestHeaderFieldsTooLarge,
	http.StatusUnavailableForLegalReasons:    ErrUnavailableForLegalReasons,
	http.StatusInternalServerError:           ErrInternalServerError,
	http.StatusNotImplemented:                ErrNotImplemented,
	http.StatusBadGateway:                    ErrBadGateway,
	http.StatusServiceUnavailable:            ErrServiceUnavailable,
	http.StatusGatewayTimeout:                ErrGatewayTimeout,
	http.StatusHTTPVersionNotSupported:       ErrHTTPVersionNotSupported,
	http.StatusVariantAlsoNegotiates:         ErrVariantAlsoNegotiates,
	http.StatusInsufficientStorage:           ErrInsufficientStorage,
	http.StatusLoopDetected:                  ErrLoopDetected,
	http.StatusNotExtended:                   ErrNotExtended,
	http.StatusNetworkAuthenticationRequired: ErrNetworkAuthenticationRequired,
}

func handleHTTPErrorResp(resp *http.Response) error {
	if resp.StatusCode >= 400 && resp.StatusCode <= 599 {
		var res struct {
			Message string `json:"message"`
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading response body: %s", err)
		}

		if json.Valid(data) {
			if err = json.Unmarshal(data, &res); err != nil {
				return fmt.Errorf("error unmarshaling response: %s", err)
			}
		} else {
			res.Message = string(data)
		}

		err, ok := httpStatusCodeErrors[resp.StatusCode]
		if ok {
			if err == ErrNotFound {
				return ErrNotFound
			}
			return fmt.Errorf("%w: %s", err, res.Message)
		}

		return fmt.Errorf("unmapped HTTP error: status code %d, message: %s", resp.StatusCode, res.Message)
	}

	return nil
}
