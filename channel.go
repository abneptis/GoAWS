package aws

import "com.abneptis.oss/aws/awsconn"
import "com.abneptis.oss/maptools"

import "http"
import "os"


type Channel interface {
  WriteRequest(*http.Request)(*http.Response,os.Error)
  Endpoint()(awsconn.Endpoint)
}

type SignedChannel struct {
  conn *http.ClientConn
  ep *awsconn.Endpoint
  net string
  local string
  reconnect bool
}

func NewConnection(ep *awsconn.Endpoint, net, local string, reconnect bool)(Channel){
  return &SignedChannel{ ep: ep, net: net, local: local, reconnect: reconnect}
}

func (self *SignedChannel)Endpoint()(awsconn.Endpoint) {
  return *self.ep
}

func (self *SignedChannel)connection()(conn *http.ClientConn, err os.Error){
  if self.conn == nil {
    self.conn, err = self.ep.NewHTTPClientConn(self.net, self.local, nil)
  }
  if err == nil {
    conn = self.conn
  }
  return
}

// The (in)(out) error is simply copied.  It's to make
// return abortConnection(err) a trivial action.
func (self *SignedChannel)abortConnection(err os.Error)(os.Error){
  if self.conn != nil {
    self.conn.Close()
    self.conn = nil
  }
  return err
}

func (self *SignedChannel)writeRequest(req *http.Request)(resp *http.Response, err os.Error){
  c, err := self.connection()
  if err == nil {
    updateRawQuery(req)
    err = c.Write(req)
    if err == nil {
      resp, err = c.Read()
    }
    if err != nil {
      if err == http.ErrPersistEOF {
        err = nil
      }
      err = self.abortConnection(err)
    }
  }
  return
}

func updateRawQuery(req *http.Request){
  req.URL.RawQuery = maptools.StringStringJoin(
                     maptools.StringStringEscape(
                     maptools.StringStringsJoin(req.Form, ",", false),
                     Escape, Escape),
                     "=", "&", true)
}

func (self *SignedChannel)WriteRequest(req *http.Request)(resp *http.Response, err os.Error){
  return self.writeRequest(req)
}
