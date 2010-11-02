package s3

import "com.abneptis.oss/aws/awsconn"
import "com.abneptis.oss/aws/auth"
import "com.abneptis.oss/cryptools/signer"
import "com.abneptis.oss/urltools"
import "com.abneptis.oss/maptools"

import "encoding/base64"
import "http"
import "os"
import "strconv"
import "strings"
import "time"


func s3EscapeTest(i byte)(out bool){
  switch i {
    case 'a','b','c','d','e','f','g','h','i','j','k','l','m',
         'A','B','C','D','E','F','G','H','I','J','K','L','M',
         'n','o','p','q','r','s','t','u','v','w','x','y','z',
         'N','O','P','Q','R','S','T','U','V','W','X','Y','Z',
         '0','1','2','3','4','5','6','7','8','9','-':
      out = false
    default:
      out = true
  }
  return
}

func s3Escape(in string)(out string){
  return urltools.Escape(in, s3EscapeTest, urltools.PercentUpper)
}



// Constructs a query request object;  This complex type implements
// both the Canonicalization and Signing methods necessary to access
// S3.
//
// For the most part, users are expected to use the exposed bucket/endpoint/obj
// functions
//
// If Expiration is < 1 year (in seconds), the expiration is assumed to
// mean seconds-from-now (local clock based).
func NewQueryRequest(id auth.Signer, endpoint *awsconn.Endpoint,
                     method, bucket, key, ctype, cmd5 string,
                     params, hdrs map[string]string)(req *http.Request, err os.Error){
  req = &http.Request {
    Method: method,
    Host: endpoint.GetURL().Host,
    URL: &http.URL {
      Scheme: endpoint.GetURL().Scheme,
      Host: endpoint.GetURL().Host,
      Path: endpoint.GetURL().Path,
    },
    Header: hdrs,
    Form: maptools.StringStringToStringStrings(params),
  }
  if req.Header == nil { req.Header = make(map[string]string) }
  if ctype != "" { req.Header["Content-Type"] = ctype }
  if cmd5  != "" { req.Header["Content-Md5"] = cmd5}
  if bucket != "" {
    req.URL.Path = "/" + bucket
    if key != "" {
      req.URL.Path += "/" + key
    }
  }
  req.Form["AWSAccessKeyId"] = []string{auth.GetSignerIDString(id)}
  if len(req.Form["Expires"]) == 0 {
    req.Form["Expires"] = []string{strconv.Itoa64(time.Seconds() + 30)}
  }
  if len(req.Form["Signature"]) == 0 {
    sig, err := signer.SignString64(id, base64.StdEncoding, CanonicalString(req))
    if err != nil { return }
    req.Form["Signature"] = []string{sig}
  }
  return
}

func CanonicalString(req *http.Request)(string){
  return strings.Join([]string{
    req.Method,
    awsconn.ContentMD5(req),
    awsconn.ContentType(req),
    req.Form["Expires"][0], "" + req.URL.Path,
  }, "\n")
}
