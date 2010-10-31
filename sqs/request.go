package sqs

import "com.abneptis.oss/cryptools/signer"
import "com.abneptis.oss/aws"
import "com.abneptis.oss/aws/auth"
import "com.abneptis.oss/maptools"
import "com.abneptis.oss/urltools"

import "encoding/base64"
import "http"
import "fmt"
import "os"
import "strings"
import "time"

var DefaultVersion = "2009-02-01"
var DefaultSignatureVersion = "2"
var DefaultSignatureMethod = "HmacSHA256"

type Request *http.Request

func MakeHTTPRequest(url *http.URL, method string, params map[string]string)(req *http.Request){
  req = &http.Request {
    Host: url.Host,
    URL: url,
    Method: method,
    Form: maptools.StringStringToStringStrings(params),
  }
  cmap := maptools.StringStringEscape(params, sqsEscape, sqsEscape)

  req.URL.RawQuery = maptools.StringStringJoin(cmap, "=", "&", true)

  return
}


func NewSQSRequest(params map[string]string)(out aws.RequestMap, err os.Error){
  out = aws.RequestMap{
    Values: map[string]string{},
    Allowed: map[string]bool{
     "Action":true,
     "AWSAccessKeyId": true,
     "SignatureMethod": true,
     "SignatureVersion": true,
     "Signature": true,
     "Version": true,
     // One or the other of the following is required.
     "Expires": false,
     "Timestamp": false,
     // CreateQueue params.
     "DefaultVisibilityTimeout": false,
     "QueueName": false,
     // ListQueues
     "QueueNamePrefix": false,
     // PushMessage
     "MessageBody": false,
     // ReceiveMessage
     "MaxNumberOfMessages": false,
     "AttributeName": false,
     "VisibilityTimeout": false,
    },
  }
  for k,v := range(params){
    err = out.Set(k,v)
    if err != nil { break}
  }
  if err != nil { return }
  if ! out.IsSet("Version") {
    out.Set("Version", DefaultVersion)
  }
  if ! out.IsSet("SignatureMethod") {
    out.Set("SignatureMethod", DefaultSignatureMethod)
  }
  if ! out.IsSet("SignatureVersion") {
    out.Set("SignatureVersion", DefaultSignatureVersion)
  }
  if ! out.IsSet("Expires") && ! out.IsSet("Timestamp") {
    t := time.LocalTime()
    out.Set("Timestamp", t.Format(aws.SQSTimestampFormat))
  }
  return
}

func sqsEscapeTest(i byte)(out bool){
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

func sqsEscape(in string)(out string){
  return urltools.Escape(in, sqsEscapeTest, urltools.PercentUpper)
}



func SignSQSRequest(id auth.Signer, m string, u *http.URL, in *aws.RequestMap)(err os.Error){
  canonMap := maptools.StringStringEscape(in.Values, sqsEscape, sqsEscape)
  host := strings.Split(u.Host, ":", 2)
  canonString := fmt.Sprintf("%s\n%s\n%s\n%s", m, host[0], u.Path,
                 maptools.StringStringJoin(canonMap, "=", "&", true))
  //fmt.Printf("CanonString: [%s]\n", canonString)
  sig, err := signer.SignString64(id, base64.StdEncoding, canonString)
  if err == nil {
    err = in.Set("Signature", sig)
  }
  return
}


func SignAndSendSQSRequest(id auth.Signer, method string, ep *aws.Endpoint, in *aws.RequestMap)(resp *http.Response, err os.Error){
  err = SignSQSRequest(id, method, ep.URL, in)
  if err != nil { return }
  hreq := MakeHTTPRequest(ep.URL, method, in.Values)
  hreq.Close = true
  //bb, _ := http.DumpRequest(hreq, true)
  //os.Stderr.Write(bb)
  cc, err := ep.NewHTTPClientConn("tcp", "", nil)
  if err != nil && err != http.ErrPersistEOF { return }
  resp, err = aws.SendRequest(cc, hreq)
  return
}
