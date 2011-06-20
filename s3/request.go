package s3

import (
  "aws"
)

import (
  "encoding/base64"
  "http"
//  "log"
  "os"
  "path"
  "strconv"
  "time"
)

func newRequest(method, bucket, key string, hdrs http.Header, params http.Values)(req *http.Request){
  req = &http.Request {
    Method: method,
    URL: &http.URL {
      Path: path.Join("/", bucket, key),
    },
    Host: "s3.amazonaws.com",
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
  req.Form.Set("AWSAccessKeyId",id.AccessKey)
  if req.Form.Get("Expires") == "" &&
     req.Header.Get("Date") == "" &&
     req.Header.Get("X-Amz-Date") == "" {
    req.Form.Set("Expires",strconv.Itoa64(time.Seconds() + 30))
  }
  sig, err := id.SignEncoded(DEFAULT_HASH, CanonicalString(req), 
                             base64.StdEncoding)
  if err == nil {
    req.Form.Set("Signature", string(sig))
  }
  // log.Printf("Signed request: '%s';", CanonicalString(req) )
  return
}
