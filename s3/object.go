package s3

import "io"
import "os"

type Object struct {
  Key string
  Body io.ReadCloser
}

func (self *Object)Close()(os.Error){
  return self.Body.Close()
}

func (self *Object)Read(b []byte)(int, os.Error){
  return self.Body.Read(b)
}
