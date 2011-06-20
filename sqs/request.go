package sqs

import (
  "aws"
  "com.abneptis.oss/maptools"
)

import (
  "crypto"
  "encoding/base64"
  "http"
  "log"
  "os"
  "strconv"
  "time"
  "strings"
)


const (
  DEFAULT_VERSION = "2009-02-01"
  DEFAULT_SIGNATURE_VERSION = "2"
  DEFAULT_SIGNATURE_METHOD  = "HmacSHA256"
)

func newRequest(method string, url *http.URL, hdrs http.Header, params http.Values)(req *http.Request){
  req = &http.Request {
    Method: method,
    URL: url,
    Host: url.Host,
    Header: hdrs,
    Form: params,
  }
  return
}

func signRequest(id *aws.Signer, req *http.Request)(err os.Error){
  // log.Printf("Signing request...")
  if req.Form == nil {
    req.Form = http.Values{}
  }
  if req.Form.Get("Version") == "" {  req.Form.Set("Version", DEFAULT_VERSION) }
  if req.Form.Get("SignatureVersion") == "" {  req.Form.Set("SignatureVersion", DEFAULT_SIGNATURE_VERSION) }
  if req.Form.Get("SignatureMethod") == "" {  req.Form.Set("SignatureMethod", DEFAULT_SIGNATURE_METHOD) }
  req.Form.Set("AWSAccessKeyId",id.AccessKey)

  if req.Form.Get("Expires") == "" && req.Form.Get("Timestamp") == "" { 
    req.Form.Set("Expires",strconv.Itoa64(time.Seconds() + 30))
  }
  var sig []byte 
  switch req.Form.Get("SignatureMethod") {
    case "HmacSHA256":
      sig, err = id.SignEncoded(crypto.SHA256, CanonicalString(req), base64.StdEncoding)
    case "HmacSHA1":
      sig, err = id.SignEncoded(crypto.SHA1, CanonicalString(req), base64.StdEncoding)
    default: err = os.NewError("Unknown SignatureMethod")
  }
  if err == nil {
    req.Form.Set("Signature", string(sig))
  }
  log.Printf("Signed request: '%s';", CanonicalString(req) )
  return
}


func CanonicalString(req *http.Request)(cstr string){
  params := maptools.StringStringsJoin(req.Form, ",", true)
  cmap := maptools.StringStringEscape(params, aws.Escape, aws.Escape)
  cstr = strings.Join([]string{req.Method, req.Host, req.URL.Path,
                 maptools.StringStringJoin(cmap, "=", "&", true)}, "\n")
  return
}

