package aws

import (
  "crypto/tls"
  "http"
  "net"
  "os"
)

type Conn struct {
  uc *ReusableConn
  c  *http.ClientConn
}

func NewConn(d Dialer)(*Conn){
  return &Conn {
    uc: NewReusableConnection(d),
    c: nil,
  }
}

func (self *Conn)dial()(err os.Error){
  if self.c == nil {
    // Get the underlying connection (or redial)
    err = self.uc.Dial()
    if err == nil {
      self.c = http.NewClientConn(self.uc, nil)
    }
  }
  return
}

func (self *Conn)Request(req *http.Request)(resp *http.Response, err os.Error){
  err = self.dial()
  if err == nil {
    if req.Form != nil && req.Method == "GET" {
      if req.URL.RawQuery != "" {
        req.URL.RawQuery += "&"
      }
      req.URL.RawQuery += req.Form.Encode()
      req.Form = nil
    }
    // ob, _ := http.DumpRequest(req, true)
    // os.Stdout.Write(ob)

    err = self.c.Write(req)
    if err == nil {
      resp, err = self.c.Read(req)
    }
  }
  return
}


func URLDialer(u *http.URL, conf *tls.Config)(f func()(c net.Conn, err os.Error)){
  host, port, _ := net.SplitHostPort(u.Host)
  if port == "" {
    if u.Scheme == "http" { port = "80" }
    if u.Scheme == "https" { port = "443" }
  }
  if host == "" {
    host = u.Host
  }
  useTLS := (u.Scheme == "https")

  f = func()(c net.Conn, err os.Error){
    if useTLS {
      return tls.Dial("tcp", host + ":" + port, conf)
    }
    return net.Dial("tcp", host + ":" + port)
  }
  return
}
