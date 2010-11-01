package s3

import "os"

var ErrorKeyNotFound = os.NewError("Key not found")

type S3Error struct {
  Code string
  Message string
  RequestId string
  HostId string
  BucketName string
}

func (self *S3Error)String()(string){
  return "{S3Error} [" +  self.Code + "]: " + self.Message
}
