package s3

import "com.abneptis.oss/aws/awsconn"
import "com.abneptis.oss/aws/auth"
import "com.abneptis.oss/aws"
import "com.abneptis.oss/cryptools/signer"
import "com.abneptis.oss/maptools"
import "com.abneptis.oss/urltools"

import "encoding/base64"
//import "fmt"
import "http"
import "os"
import "strings"
import "strconv"
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



type QueryRequest struct {
  Method string
  Endpoint *awsconn.Endpoint
  AmzHeaders   map[string]string
  Parameters   *aws.RequestMap
  Flag  string
  Bucket string
  Key    string
  ContentType string
  ContentMD5  string
}

func queryParamsMap()(*aws.RequestMap){
  return &aws.RequestMap{
    Allowed: map[string]bool{
      "Signature": true,
      "AWSAccessKeyId": true,
      "Expires": true,
    },
    Values: map[string]string{},
  }
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
func NewQueryRequest(method string, endpoint *awsconn.Endpoint,
                     bucket, key, flag, ctype, cmd5 string,
                    amz, params map[string]string,
                    expiration int64)(qr *QueryRequest, err os.Error){
  qr = &QueryRequest{
    Method: method,
    Endpoint: endpoint,
    Flag: flag, Bucket: bucket, Key: key,
    ContentType: ctype, ContentMD5: cmd5,
    Parameters: queryParamsMap(),
  }

  if amz == nil {
    qr.AmzHeaders = make(map[string]string)
  } else {
    qr.AmzHeaders = amz
  }

  qr.Parameters = queryParamsMap()
  for k, v := range(amz){
    err = qr.Parameters.Set(k,v)
  }
  if expiration < (365*24*60*60) {
    expiration += time.Seconds()
  }
  qr.Parameters.Set("Expires", strconv.Itoa64(expiration))
  return
}

func (self *QueryRequest)canonHeaders()(out string){
  for k, v := range(self.AmzHeaders) {
    out += strings.ToLower(k) + "=" + v + "\n"
  }
  return
}

func (self *QueryRequest)canonResource()(out string){
  // TODO: Escape key(?)
  // TODO: Handle Flag(? is that part of canon?)
  out += "/" + self.Bucket
  if self.Key != "" {
    out += "/" + self.Key
  }
  return
}

func (self *QueryRequest)CanonicalString()(out string){
  exp, _ := self.Parameters.Get("Expires")
  out = strings.Join([]string{self.Method, self.ContentMD5, self.ContentType,
                     exp,
                     self.canonHeaders() + self.canonResource()}, "\n")
  return
}

func (self *QueryRequest)Sign(id auth.Signer)(err os.Error){
  err = self.Parameters.Set("AWSAccessKeyId", string(id.PublicIdentity()))
  if err != nil { return }
  cs := self.CanonicalString()
  //fmt.Printf("CanonString: [%s]\n", cs)
  sig, err := signer.SignString64(id, base64.StdEncoding, cs)
  if err == nil {
    err = self.Parameters.Set("Signature", sig)
  }
  return
}

func (self *QueryRequest)IsVhosted()(bool){
  if self.Endpoint.GetURL().Host == self.Bucket ||
     self.Endpoint.GetURL().Host == self.Bucket + ".s3.amazonaws.com" {
    return true
  }
  return false
}

func (self *QueryRequest)HTTPRequest()(req *http.Request){
  req = &http.Request {
    Method: self.Method,
    URL: &http.URL {
      Scheme: self.Endpoint.GetURL().Scheme,
      Host: self.Endpoint.GetURL().Host,
      Path: "/",
    },
    Header: map[string]string{},
  }
  if self.Bucket != "" {
    req.URL.Path += self.Bucket
    if self.Key != "" {
      req.URL.Path += self.Key
    }
  }
  if self.ContentType != "" {
    req.Header["Content-Type"] = self.ContentType
  }
  if self.ContentMD5 != "" {
    // is Case correct?
    req.Header["Content-Md5"] = self.ContentMD5
  }
  if self.Flag != "" {
    req.URL.RawQuery = self.Flag + "&"
  }
  if len(req.URL.RawQuery) > 0 { req.URL.RawQuery += "&" }
  cmap := maptools.StringStringEscape(self.Parameters.Values, s3Escape, s3Escape)
  req.URL.RawQuery += maptools.StringStringJoin(cmap, "=", "&", true)
  return
}


// Sign and send a low-level qury-request;  this includes formulating the
// appropriate http request, creating a connection, sending the data, and
// returning the ultimate http response.
func (self *QueryRequest)Send(id auth.Signer, ep *awsconn.Endpoint)(resp *http.Response, err os.Error){
  self.Sign(id)
  if err != nil {return}
  hreq := self.HTTPRequest()
  cconn, err := ep.NewHTTPClientConn("tcp", "", nil)
  if err != nil { return }
  defer cconn.Close()
  resp, err = awsconn.SendRequest(cconn, hreq)
  return
}

