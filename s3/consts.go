package s3

import (
  "crypto"
  "http"
  "os"
)


// Do S3 endpoints actually play any role?
const (
  USWEST_HOST = "us-west-1.s3.amazonaws.com"
  USEAST_HOST = "s3.amazonaws.com"
  APSOUTHEAST_HOST = "ap-southeast-1.s3.amazonaws.com"
  EUWEST_HOST = "eu-west-1.s3.amazonaws.com"
)

const (
  DEFAULT_HASH = crypto.SHA1
)


var ErrorUnknownKey os.Error = os.NewError("Unknown Key")
var ErrorUnexpectedResponse os.Error = os.NewError("Unexpected response code")
var ErrorConflicts os.Error = os.NewError("Conflicts with anothers resources")
var ErrorPermissions os.Error = os.NewError("Access denied")

func CodeToError(i int)(err os.Error){
  switch i {
    case http.StatusOK:
    case http.StatusNotFound: err = ErrorUnknownKey
    case http.StatusConflict: err = ErrorConflicts
    case http.StatusForbidden: err = ErrorPermissions
    default: err = ErrorUnexpectedResponse
  }
  return
}


