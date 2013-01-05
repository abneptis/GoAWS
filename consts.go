package aws

import (
	"errors"
	"net/http"
)

const (
	DEFAULT_SIGNATURE_VERSION = "2"
	DEFAULT_SIGNATURE_METHOD  = "HmacSHA256"
)

var ErrorNotFound = errors.New("Not found")
var ErrorUnexpectedResponse = errors.New("Unexpected response code")
var ErrorConflicts = errors.New("Conflicts with another resources")
var ErrorForbidden = errors.New("Access denied")

func CodeToError(i int) (err error) {
	switch i {
	case http.StatusOK:
	case http.StatusNotFound:
		err = ErrorNotFound
	case http.StatusConflict:
		err = ErrorConflicts
	case http.StatusForbidden:
		err = ErrorForbidden
	default:
		err = ErrorUnexpectedResponse
	}
	return
}
