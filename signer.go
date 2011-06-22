package aws

import (
  "bytes"
  "crypto"
  "crypto/hmac"
  "encoding/base64"
  "hash"
  "http"
  "os"
  "strings"
  "strconv"
  "time"
)

// A signer simply holds the access & secret access keys
// necessary for aws, and proivides helper functions
// to assist in generating an appropriate signature.
type Signer struct {
  AccessKey string
  secretAccessKey []byte
}

func NewSigner(akid, sak string)(*Signer){
  return &Signer{akid,bytes.NewBufferString(sak).Bytes()}
}

// the core function of the Signer, generates the raw hmac of he bytes.
func (self *Signer)SignBytes(h crypto.Hash, buff []byte)(sig []byte, err os.Error){
  hh := hmac.New(func()(hash.Hash){
     return h.New()
  }, self.secretAccessKey)
  _, err = hh.Write(buff)
  if err == nil {
    sig = hh.Sum()
  }
  return
}

// Same as SignBytes, but with strings.
func (self *Signer)SignString(h crypto.Hash, s string)(os string, err os.Error){
  ob, err := self.SignBytes(h, bytes.NewBufferString(s).Bytes())
  if err == nil {
    os = string(ob)
  }
  return
}

// SignBytes, but will base64 encode based on the specified encoder.
func (self *Signer)SignEncoded(h crypto.Hash, s string, e *base64.Encoding)(out []byte, err os.Error){
  ob, err := self.SignBytes(h, bytes.NewBufferString(s).Bytes())
  if err == nil { 
    out = make([]byte, e.EncodedLen(len(ob)))
    e.Encode(out, ob)
  }
  return
}

// The V2 denotes amazon signing version 2, not version 2 of this particular function...
// V2 is used by all services but S3;
// Note, some services vary in their exact implementation of escaping w/r/t signatures,
// so it is recommended you use this function.
//
// Final note: if exp is set to 0, a Timestamp will be used, otherwise an expiration.
func (self *Signer)SignRequestV2(req *http.Request, canon func(*http.Request)(string), api_ver string, exp int64)(err os.Error){
  // log.Printf("Signing request...")

  if req.Form == nil { req.Form = http.Values{} }

  // Setup some defaults
  if req.Form.Get("SignatureVersion") == "" {  req.Form.Set("SignatureVersion", DEFAULT_SIGNATURE_VERSION) }
  if req.Form.Get("SignatureMethod") == "" {  req.Form.Set("SignatureMethod", DEFAULT_SIGNATURE_METHOD) }

  req.Form.Set("Version", api_ver)
  req.Form.Set("AWSAccessKeyId",self.AccessKey)
  req.Form.Del("Signature")
  if req.Form.Get("Timestamp") == "" && req.Form.Get("Expires") == "" {
    if exp > 0 {
      req.Form.Set("Expires",strconv.Itoa64(time.Seconds() + exp))
    } else {
      req.Form.Set("Expires",time.UTC().Format(ISO8601TimestampFormat))
    }
  }

  var sig []byte 
  switch req.Form.Get("SignatureMethod") {
    case "HmacSHA256": sig, err = self.SignEncoded(crypto.SHA256, canon(req), base64.StdEncoding)
    case "HmacSHA1": sig, err = self.SignEncoded(crypto.SHA1, canon(req), base64.StdEncoding)
    default: err = os.NewError("Unknown SignatureMethod:" + req.Form.Get("SignatureMethod"))
  }

  if err == nil {
    req.Form.Set("Signature", string(sig))
    if req.Method == "GET" {
      if req.URL.RawQuery != "" {
        req.URL.RawQuery += "&"
      }
      req.URL.RawQuery += req.Form.Encode()
    }
  }

  return
}

// Used exclusively by S3 to the best of my knowledge...
func (self *Signer)SignRequestV1(req *http.Request, canon func(*http.Request)(string), exp int64)(err os.Error){
  if req.Form == nil { req.Form = http.Values{} }

  req.Form.Set("AWSAccessKeyId",self.AccessKey)
  req.Form.Del("Signature")
  req.Form.Set("Expires",strconv.Itoa64(time.Seconds() + exp))

  var sig []byte 
  sig, err = self.SignEncoded(crypto.SHA1, canon(req), base64.StdEncoding)

  if err == nil {
    req.Form.Set("Signature", string(sig))
    if req.Method == "GET" || req.Method == "PUT" || req.Method == "DELETE" {
      if req.URL.RawQuery != "" {
        req.URL.RawQuery += "&"
      }
      req.URL.RawQuery += req.Form.Encode()
    }
  }

  return
}

// Generates the canonical string-to-sign for (most) AWS services.
// You shouldn't need to use this directly.
func Canonicalize(req *http.Request)(out string){
  return strings.Join([]string{req.Method, req.Host, req.URL.Path, SortedEscape(req.Form)}, "\n")
}

// Generates the canonical string-to-sign for S3 services.
// You shouldn't need to use this directly unless you're pre-signing URL's.
func CanonicalizeS3(req *http.Request)(string){
  return strings.Join([]string{ req.Method, req.Header.Get("Content-Md5"), req.Header.Get("Content-Type"), req.Form.Get("Expires"), req.URL.Path, }, "\n")
}


