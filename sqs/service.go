// Amazon AWS Simple Queue Service interface.
//
// This package implements basic SQS functionality
// including queue creation, deletion, enumeration,
// and message pushing/fetching and deletion.
package sqs

import "com.abneptis.oss/aws/auth"
import "com.abneptis.oss/aws/awsconn"
import "com.abneptis.oss/maptools"
import "com.abneptis.oss/cryptools/signer"

import "encoding/base64"
import "http"
//import "fmt"
import "os"
import "strconv"
import "strings"
import "time"

var DefaultSQSVersion = "2009-02-01"
//var DefaultSignatureVersion = "2"
//var DefaultSignatureMethod = "HmacSHA256"



type Service struct {
  Endpoint *awsconn.Endpoint
}

func NewService(ep *awsconn.Endpoint)(*Service){
  return &Service{Endpoint: ep}
}

func CanonicalizeRequest(req *http.Request)(cstr string){
  params := maptools.StringStringsJoin(req.Form, ",", true)
  cmap := maptools.StringStringEscape(params, sqsEscape, sqsEscape)
  cstr = strings.Join([]string{req.Method, req.Host, req.URL.Path,
                 maptools.StringStringJoin(cmap, "=", "&", true)}, "\n")
  return
}

func SignRequest(id auth.Signer, req *http.Request)(err os.Error){
  cstr := CanonicalizeRequest(req)
  //fmt.Printf("Canon String:\n==========\n{%s}\n=========\n", cstr)
  sig, err := signer.SignString64(id, base64.StdEncoding, cstr)
  if err == nil {
    req.Form["Signature"] = []string{sig}
  }
  return
}

func (self *Service)signedRequest(id auth.Signer, path string, params map[string]string)(req *http.Request, err os.Error){
  req = self.Endpoint.NewHTTPRequest("GET", path, maptools.StringStringToStringStrings(params), nil)
  req.Form["AWSAccessKeyId"] = []string{auth.GetSignerIDString(id)}
  if len(req.Form["Version"]) == 0 {
    req.Form["Version"] = []string{DefaultSQSVersion}
  }
  if len(req.Form["SignatureMethod"]) == 0 {
    req.Form["SignatureMethod"] = []string{DefaultSignatureMethod}
  }
  if len(req.Form["SignatureVersion"]) == 0 {
    req.Form["SignatureVersion"] = []string{DefaultSignatureVersion}
  }
  if len(req.Form["Expires"]) == 0 && len(req.Form["Timestamp"]) == 0{
    req.Form["Timestamp"] = []string{ strconv.Itoa64(time.Seconds()) }
  }
  err = SignRequest(id, req)
  return
}



func (self *Service)CreateQueue(id auth.Signer, name string, dvtimeout int)(mq *Queue, err os.Error){
  sqsReq, err := self.signedRequest(id, "/", map[string]string{
    "Action": "CreateQueue",
    "QueueName": name,
    "DefaultVisibilityTimeout": strconv.Itoa(dvtimeout),
  })
  if err != nil { return }

  xresp := &createQueueResponse{}
  xerr := &errorResponse{}
  err = self.Endpoint.SendParsable(sqsReq, xresp, xerr)

  if err != nil { return }
  qrl, err := http.ParseURL(xresp.CreateQueueResult.QueueUrl)
  if err != nil { return }
  ep := awsconn.NewEndpoint(qrl, self.Endpoint.ProxyURL)
  mq = NewQueueURL(ep)
  return
}

func (self *Service)ListQueues(id auth.Signer, prefix string)(out []*Queue, err os.Error){
  sqsReq, err := self.signedRequest(id, "/", map[string]string{
    "Action": "ListQueues",
  })
  if err != nil { return }
  xresp := &listQueuesResponse{}
  xerr  := &errorResponse{}
  err = self.Endpoint.SendParsable(sqsReq, xresp, xerr)
  if err != nil { return }
  out = make([]*Queue, len(xresp.ListQueuesResult.QueueUrl))
  for i := range(xresp.ListQueuesResult.QueueUrl){
    url, err := http.ParseURL(xresp.ListQueuesResult.QueueUrl[i])
    if err != nil { break }
    ep := awsconn.NewEndpoint(url, self.Endpoint.ProxyURL)
    out[i] = NewQueueURL(ep)
  }
  //log.Printf("ListQueue: %v",  out)
  return
}

