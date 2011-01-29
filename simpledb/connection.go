package simpledb

import "com.abneptis.oss/aws/awsconn"
import "com.abneptis.oss/aws"
import "com.abneptis.oss/maptools"

import "http"
import "log"
import "os"


type AWSConnection interface {
  WriteRequest(*http.Request)(*http.Response,os.Error)
  Endpoint()(awsconn.Endpoint)
}

type Connection struct {
  conn *http.ClientConn
  ep awsconn.Endpoint
  net string
  local string
}

func NewConnection(ep awsconn.Endpoint, net, local string)(AWSConnection){
  return &Connection{ ep: ep, net: net, local: local}
}

func (self *Connection)Endpoint()(awsconn.Endpoint) {
  return self.ep
}

func (self *Connection)connection()(conn *http.ClientConn, err os.Error){
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
func (self *Connection)abortConnection(err os.Error)(os.Error){
  if self.conn != nil {
    self.conn.Close()
    self.conn = nil
  }
  return err
}

func (self Connection)WriteRequest(req *http.Request)(resp *http.Response, err os.Error){
  c, err := self.connection()
  req.URL.RawQuery = maptools.StringStringJoin(
                     maptools.StringStringEscape(
                     maptools.StringStringsJoin(req.Form, ",", false),
                     aws.Escape, aws.Escape),
                     "=", "&", true)

  if err == nil {
    err = c.Write(req)
    if err != nil {
      return nil, self.abortConnection(err)
    }
  }
  log.Printf("Final URL: %s", req.URL.String())
  if err == nil {
    resp, err = c.Read()
    if err != nil {
      err = self.abortConnection(err)
      // PersistantEOF is NOT an invalid condition for the request, it simply invalidates the ClientConn.
      if err == http.ErrPersistEOF {
        err = nil
      }
    }
  }
  return
}
