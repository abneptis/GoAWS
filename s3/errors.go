package s3

import "os"

var ErrorKeyNotFound = os.NewError("Key not found")
var ErrorAccessDenied = os.NewError("Access denied")

type errorResponse struct {
  Error Error
}

func (self errorResponse)String()(string){return self.Error.String()}


// S3s error format is somewhat open
// so fields can be added here if needed.
//
// We make no effort to respond to errors correctly,
// merely report them with their appropriate details. 
type Error struct {
  Code string
  Message string
  RequestId string
  HostId string
  BucketName string
  StringToSignBytes string
}

func (self *Error)String()(string){
  return "{S3Error} [" +  self.Code + "]: " + self.Message
}
