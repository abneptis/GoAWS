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
	AccessKey       string
	secretAccessKey []byte
}

func NewSigner(akid, sak string) *Signer {
	return &Signer{akid, bytes.NewBufferString(sak).Bytes()}
}

// the core function of the Signer, generates the raw hmac of he bytes.
func (self *Signer) SignBytes(h crypto.Hash, buff []byte) (sig []byte, err os.Error) {
	hh := hmac.New(func() hash.Hash {
		return h.New()
	}, self.secretAccessKey)
	_, err = hh.Write(buff)
	if err == nil {
		sig = hh.Sum()
	}
	return
}

// Same as SignBytes, but with strings.
func (self *Signer) SignString(h crypto.Hash, s string) (os string, err os.Error) {
	ob, err := self.SignBytes(h, bytes.NewBufferString(s).Bytes())
	if err == nil {
		os = string(ob)
	}
	return
}

// SignBytes, but will base64 encode based on the specified encoder.
func (self *Signer) SignEncoded(h crypto.Hash, s string, e *base64.Encoding) (out []byte, err os.Error) {
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
func (self *Signer) SignRequestV2(req *http.Request, canon func(*http.Request)(string, os.Error), api_ver string, exp int64) (err os.Error) {
	// log.Printf("Signing request...")

  qstring, err := http.ParseQuery(req.URL.RawQuery)
  if err != nil { return }
  qstring["SignatureVersion"] =  []string{DEFAULT_SIGNATURE_VERSION}
  if _, ok := qstring["SignatureMethod"]; !ok || len(qstring["SignatureMethod"]) == 0 {
	  qstring["SignatureMethod"]  =  []string{DEFAULT_SIGNATURE_METHOD}
  }
	qstring["Version"]  =  []string{api_ver}

  if exp > 0 {
    qstring["Expires"] = []string{strconv.Itoa64(time.Seconds()+exp)}
  } else {
    qstring["Timestamp"] = []string{time.UTC().Format(ISO8601TimestampFormat)}
  }
  qstring["Signature"] = nil, false
  qstring["AWSAccessKeyId"] = []string{ self.AccessKey}

	var sig []byte
  req.URL.RawQuery = http.Values(qstring).Encode()
  can, err := canon(req)
  if err != nil { return }
  //log.Printf("String-to-sign: '%s'", can)

	switch qstring["SignatureMethod"][0] {
	case "HmacSHA256":
		sig, err = self.SignEncoded(crypto.SHA256, can, base64.StdEncoding)
	case "HmacSHA1":
		sig, err = self.SignEncoded(crypto.SHA1, can, base64.StdEncoding)
	default:
		err = os.NewError("Unknown SignatureMethod:" + req.Form.Get("SignatureMethod"))
	}

	if err == nil {
    req.URL.RawQuery += "&" + http.Values{"Signature": []string{string(sig)}}.Encode()
	}

	return
}

// Used exclusively by S3 to the best of my knowledge...
func (self *Signer) SignRequestV1(req *http.Request, canon func(*http.Request) (string, os.Error), exp int64) (err os.Error) {
  qstring, err := http.ParseQuery(req.URL.RawQuery)

  if err != nil { return }

  if exp > 0 {
    qstring["Expires"] = []string{strconv.Itoa64(time.Seconds()+exp)}
  } else {
    qstring["Timestamp"] = []string{time.UTC().Format(ISO8601TimestampFormat)}
  }
  qstring["Signature"] = nil, false
  qstring["AWSAccessKeyId"] = []string{ self.AccessKey}


  req.URL.RawQuery = http.Values(qstring).Encode()


  can, err := canon(req)
  if err != nil { return }

	var sig []byte
	sig, err = self.SignEncoded(crypto.SHA1, can, base64.StdEncoding)

	if err == nil {
    req.URL.RawQuery += "&" + http.Values{"Signature": []string{string(sig)}}.Encode()
	}

	return
}

// Generates the canonical string-to-sign for (most) AWS services.
// You shouldn't need to use this directly.
func Canonicalize(req *http.Request) (out string, err os.Error) {
  fv, err := http.ParseQuery(req.URL.RawQuery)
  if err == nil { 
	  out = strings.Join([]string{req.Method, req.Host, req.URL.Path, SortedEscape(fv)}, "\n")
  }
  return
}

// Generates the canonical string-to-sign for S3 services.
// You shouldn't need to use this directly unless you're pre-signing URL's.
func CanonicalizeS3(req *http.Request) (out string, err os.Error ){
  fv, err := http.ParseQuery(req.URL.RawQuery)
  if err == nil || len(fv["Expires"]) != 1 { 
	  out = strings.Join([]string{req.Method, req.Header.Get("Content-Md5"), req.Header.Get("Content-Type"), fv["Expires"][0], req.URL.Path}, "\n")
  }
  return
}
