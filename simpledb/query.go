package simpledb

import "com.abneptis.oss/aws/awsconn"
import "com.abneptis.oss/aws/auth"
import "com.abneptis.oss/aws"
import "com.abneptis.oss/maptools"
import "com.abneptis.oss/cryptools/signer"

import "encoding/base64"
import "http"
import "os"
import "strings"
import "time"


type Request http.Request

func (self *Request)Sign(id auth.Signer)(err os.Error){
  params := maptools.StringStringsJoin(self.Form, ",", true)
  cmap := maptools.StringStringEscape(params, aws.Escape, aws.Escape)
  cstr := strings.Join([]string{self.Method, self.Host, self.URL.Path,
                 maptools.StringStringJoin(cmap, "=", "&", true)}, "\n")
  sig, err := signer.SignString64(id, base64.StdEncoding, cstr)
  if err == nil {
    self.Form["Signature"] = []string{sig}
    //log.Printf("CanonString: {\n%s\n}", cstr)
    //log.Printf("Signature: %s", string(sig))
  }
  return
}


func newQuery(id auth.Signer, endpoint awsconn.Endpoint,
              domain, action string,
              params map[string]string)(hreq *http.Request, err os.Error){

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
  if _, ok := params["Timestamp"]; ! ok { params["Timestamp"] = time.UTC().Format("2006-01-02T15:04:05-07:00")}
  // ListDomains, Select don't require a DomainName attribute
  if domain != "" { params["DomainName"] = domain }
  params["Action"] = action
  params["AWSAccessKeyId"] = string(id.PublicIdentity())

  req.Form =  maptools.StringStringToStringStrings(params)
  if len(http.EncodeQuery(req.Form)) > 512{
    req.Method = "POST"
  } else {
    req.Method = "GET"
  }

  if _, ok := params["Signature"]; !ok {
    err = req.Sign(id)
  }
  if err == nil {
    hreq = (*http.Request)(req)
  }
  return
}

