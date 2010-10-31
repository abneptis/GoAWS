package s3

import "com.abneptis.oss/aws"
import "com.abneptis.oss/aws/auth"
import "com.abneptis.oss/aws/awsconn"
import "com.abneptis.oss/maptools"
import "com.abneptis.oss/urltools"
import "com.abneptis.oss/cryptools/signer"

import "os"
import "http"
import "path"
import "time"
import "encoding/base64"
import "strconv"
import "strings"
//import "fmt"

type Request struct {
  Method string
  Params *aws.RequestMap
  Flag string // acl, torrent, etc.
  Key string
  Bucket string
  Endpoint *awsconn.Endpoint
  ContentMD5 string
  ContentType string
}

func NewRequest(method, bucket, key, flag string, ep *awsconn.Endpoint, parms map[string]string, expires int64)(req *Request){
  req = &Request{
    Params: &aws.RequestMap{
     Values: make(map[string]string),
     Allowed: map[string]bool{
       "AWSAccessKeyId": true,
       "Signature": true,
       "Expires": true,
     },
    },
    Bucket: bucket,
    Key: key,
    Flag: flag,
    Method:method,
  }
  if expires < 60*60*365 {
    expires += time.Seconds()
  }
  req.Set("Expires", strconv.Itoa64(expires))
  return
}

func (self *Request)Set(k, v string)(os.Error){
  return self.Params.Set(k,v)
}

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

func SignS3Request(id auth.Signer, url *http.URL, in *Request)(err os.Error){
  exp, nil := in.Params.Get("Expires")
  canonString, _ := in.CanonicalQueryString("","", exp)
  //fmt.Printf("CanonString: [%s]\n", canonString)
  sig, err := signer.SignString64(id, base64.StdEncoding, canonString)
  if err == nil {
    err = in.Set("Signature", sig)
  }
  return
}

func (self *Request)CanonicalResource()(string){
  return path.Join("/", self.Bucket, self.Key)
}

func (self *Request)CanonicalQueryString(md5, ctype, expiry string)(out string, err os.Error){
  out = strings.Join([]string{self.Method, self.ContentMD5, self.ContentType, expiry, self.CanonicalResource()}, "\n")
  return
}

func (self *Request)HTTPRequest(id auth.Signer, ep *awsconn.Endpoint, bucket, key string)(req *http.Request, err os.Error){
  _url := ep.GetURL()
  if _url == nil { err = os.NewError("Invalid endpoint URL"); return }
  req = &http.Request{
    Method: self.Method,
    URL: &http.URL{
      Scheme: _url.Scheme,
      Host: _url.Host,
    },
    Host: _url.Host,
  }
  if key != "" {
    req.URL.Path = path.Join("/", bucket)
  } else {
    req.URL.Path = path.Join("/", bucket, key)
  }
  //err = SignS3Request(id, req.URL, self)
  //if err != nil {return}
  cmap := maptools.StringStringEscape(self.Params.Values, s3Escape, s3Escape)
  req.URL.RawQuery = maptools.StringStringJoin(cmap, "=", "&", true)

  return
}
