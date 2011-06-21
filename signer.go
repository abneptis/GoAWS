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

type Signer struct {
  AccessKey string
  secretAccessKey []byte
}

func NewSigner(akid, sak string)(*Signer){
  return &Signer{akid,bytes.NewBufferString(sak).Bytes()}
}

func (self Signer)SignBytes(h crypto.Hash, buff []byte)(sig []byte, err os.Error){
  hh := hmac.New(func()(hash.Hash){
     return h.New()
  }, self.secretAccessKey)
  _, err = hh.Write(buff)
  if err == nil {
    sig = hh.Sum()
  }
  return
}

func (self Signer)SignString(h crypto.Hash, s string)(os string, err os.Error){
  ob, err := self.SignBytes(h, bytes.NewBufferString(s).Bytes())
  if err == nil {
    os = string(ob)
  }
  return
}

func (self Signer)SignEncoded(h crypto.Hash, s string, e *base64.Encoding)(out []byte, err os.Error){
  ob, err := self.SignBytes(h, bytes.NewBufferString(s).Bytes())
  if err == nil { 
    out = make([]byte, e.EncodedLen(len(ob)))
    e.Encode(out, ob)
  }
  return
}

// The V2 denotes amazon signing version 2, not version 2 of this particular function...
func (self *Signer)SignRequestV2(req *http.Request, canon func(*http.Request)(string), api_ver string, exp int64)(err os.Error){
  // log.Printf("Signing request...")

  if req.Form == nil { req.Form = http.Values{} }

  // Setup some defaults
  if req.Form.Get("SignatureVersion") == "" {  req.Form.Set("SignatureVersion", DEFAULT_SIGNATURE_VERSION) }
  if req.Form.Get("SignatureMethod") == "" {  req.Form.Set("SignatureMethod", DEFAULT_SIGNATURE_METHOD) }

  req.Form.Set("Version", api_ver)
  req.Form.Set("AWSAccessKeyId",self.AccessKey)
  req.Form.Del("Signature")
  req.Form.Set("Expires",strconv.Itoa64(time.Seconds() + exp))

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
    if req.Method == "GET" {
      if req.URL.RawQuery != "" {
        req.URL.RawQuery += "&"
      }
      req.URL.RawQuery += req.Form.Encode()
    }
  }

  return
}


func Canonicalize(req *http.Request)(out string){
  return strings.Join([]string{req.Method, req.Host, req.URL.Path, SortedEscape(req.Form)}, "\n")
}

func CanonicalizeS3(req *http.Request)(string){
  return strings.Join([]string{ req.Method, req.Header.Get("Content-Md5"), req.Header.Get("Content-Type"), req.Form.Get("Expires"), req.URL.Path, }, "\n")
}


