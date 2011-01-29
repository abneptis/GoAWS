package simpledb

import "com.abneptis.oss/aws/auth"

import "http"
import "os"

type Handler struct {
  conn AWSConnection
  signer auth.Signer
}

func NewHandler(c AWSConnection, a auth.Signer)(*Handler){
  return &Handler{ conn: c, signer: a}
}


func (self *Handler)CreateDomain(dn string)(response SimpledbResponse, err os.Error){
  req, err := newQuery(self.signer, self.conn.Endpoint(), dn, "CreateDomain", nil)
  if err == nil {
    var resp *http.Response
    resp, err = self.conn.WriteRequest(req)
    if err == nil {
      response, err = ((*Response)(resp)).ParseResponse()
    }
  }
  return
}

func (self *Handler)DeleteDomain(dn string)(response SimpledbResponse, err os.Error){
  req, err := newQuery(self.signer, self.conn.Endpoint(), dn, "DeleteDomain", nil)
  if err == nil {
    var resp *http.Response
    resp, err = self.conn.WriteRequest(req)
    if err == nil {
      response, err = ((*Response)(resp)).ParseResponse()
    }
  }
  return
}
