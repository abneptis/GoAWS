package cobranding
// package com.abneptis.oss/aws/fps/cobranding

import "com.abneptis.oss/cryptools"
import "com.abneptis.oss/maptools"
import "com.abneptis.oss/aws/awsconn"
import "com.abneptis.oss/aws"

import "http"
import "os"
import "encoding/base64"
import "strings"


type Request http.Request


func calculateCobrandingURL(id cryptools.NamedSigner, endpoint *awsconn.Endpoint, params map[string]string)(hreq *http.Request, err os.Error){

  if params == nil { params = make(map[string]string)}
  req := &Request {
    //Method: method,
    Host: endpoint.GetURL().Host,
    URL: &http.URL {
      Scheme: endpoint.GetURL().Scheme,
      Host: endpoint.GetURL().Host,
      Path: endpoint.GetURL().Path,
    },
  }
  if req.URL.Path == "" { req.URL.Path = "/" }
  if _, ok := params["Version"]; ! ok { params["Version"] = DEFAULT_VERSION }
  if _, ok := params["SignatureVersion"]; ! ok { params["SignatureVersion"] = DEFAULT_SIGNATURE_VERSION  }
  if _, ok := params["SignatureMethod"]; ! ok { params["SignatureMethod"] = DEFAULT_SIGNATURE_METHOD }
  // Timestamp is only mandatory on certain contexts.
  //if _, ok := params["Timestamp"]; ! ok { params["Timestamp"] = time.UTC().Format("2006-01-02T15:04:05-07:00")}
  // ListDomains, Select don't require a DomainName attribute
  params["callerKey"] = id.SignerName()

  req.Form =  maptools.StringStringToStringStrings(params)
  req.Method = "GET"

  if _, ok := params["Signature"]; !ok {
    err = req.Sign(id)
  }
  if err == nil {
    hreq = (*http.Request)(req)
  }
  return
}


func (self *Request)Sign(id cryptools.NamedSigner)(err os.Error){
  params := maptools.StringStringsJoin(self.Form, ",", true)
  cmap := maptools.StringStringEscape(params, aws.Escape, aws.Escape)
  cstr := strings.Join([]string{self.Method, self.Host, self.URL.Path,
                 maptools.StringStringJoin(cmap, "=", "&", true)}, "\n")
  sig, err := cryptools.SignString64(id, base64.StdEncoding, cryptools.SignableString(cstr))
  if err == nil {
    self.Form["Signature"] = []string{sig}
  }
  return
}
