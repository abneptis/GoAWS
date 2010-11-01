package s3

import "os"

var ErrorKeyNotFound = os.NewError("Key not found")
var ErrorAccessDenied = os.NewError("Access denied")

type S3Error struct {
  Code string
  Message string
  RequestId string
  HostId string
  BucketName string
  StringToSignBytes string
}

func (self *S3Error)String()(string){
  return "{S3Error} [" +  self.Code + "]: " + self.Message + self.StringToSignBytes
}
