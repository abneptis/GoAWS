package aws

import (
	"os"
	"http"
)

const (
	DEFAULT_SIGNATURE_VERSION = "2"
	DEFAULT_SIGNATURE_METHOD  = "HmacSHA256"
)


var ErrorNotFound os.Error = os.NewError("Not found")
var ErrorUnexpectedResponse os.Error = os.NewError("Unexpected response code")
var ErrorConflicts os.Error = os.NewError("Conflicts with another resources")
var ErrorForbidden os.Error = os.NewError("Access denied")

func CodeToError(i int) (err os.Error) {
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
