package errz

import "net/http"

var httpStatusByCode map[Code]int = map[Code]int{
	OK:                 http.StatusOK,
	Cancelled:          499, // See 499 Client Closed Request nginx extension
	Unknown:            http.StatusInternalServerError,
	InvalidArgument:    http.StatusBadRequest,
	DeadlineExceeded:   http.StatusGatewayTimeout, // in practice the context can only be canceled during network requests..
	NotFound:           http.StatusNotFound,
	AlreadyExists:      http.StatusConflict, // there is a conflict of sorts
	PermissionDenied:   http.StatusForbidden,
	ResourceExhausted:  http.StatusTooManyRequests,
	FailedPrecondition: http.StatusPreconditionFailed,           // kind of
	Aborted:            http.StatusConflict,                     // kind of
	OutOfRange:         http.StatusRequestedRangeNotSatisfiable, // kind of
	Unimplemented:      http.StatusNotImplemented,
	Internal:           http.StatusInternalServerError,
	Unavailable:        http.StatusServiceUnavailable,
	DataLoss:           http.StatusInternalServerError,
	Unauthenticated:    http.StatusUnauthorized,
}

// HTTPStatus returns the HTTP status code for the given error code.
func (code Code) HTTPStatus() int {
	status, ok := httpStatusByCode[code]
	if !ok {
		status = http.StatusInternalServerError
	}
	return status
}

// HTTPStatus returns the HTTP status code for the given error.
func HTTPStatus(err error) int {
	switch err := err.(type) {
	case *Error:
		return err.Code.HTTPStatus()
	case nil:
		return http.StatusOK
	}

	return http.StatusInternalServerError
}
