package aws

import (
  "http"
  "os"
)

type Conn struct {
  uc *ReusableConn
  c  *http.ClientConn
}

func NewConn(d Dialer)(*Conn){
  return &Conn {
    uc: NewReusableConnection(d),
  }
}

func (self Conn)dial()(err os.Error){
  if self.c == nil {
    // Get the underlying connection (or redial)
    err = self.uc.Dial()
    if err == nil {
      self.c = http.NewClientConn(self.uc, nil)
    }
  }
  return
}

func (self Conn)Request(req *http.Request)(resp *http.Response, err os.Error){
  err = self.dial()
  if err == nil {
    err = self.c.Write(req)
    if err == nil {
      resp, err = self.c.Read(req)
    }
  }
  return
}

