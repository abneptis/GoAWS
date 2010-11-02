package sqs

//import "com.abneptis.oss/aws"
import "com.abneptis.oss/aws/awsconn"
import "com.abneptis.oss/aws/auth"
import "com.abneptis.oss/uuid"
import "com.abneptis.oss/maptools"
import "bytes"
import "os"
import "path"
import "strconv"
//import "xml"
import "time"
import "http"

var MaxNumberOfMessages = 10

type Queue struct {
  Name string
  Endpoint *awsconn.Endpoint
}


// Creates a new queue from a name and endpoint.
func NewQueue(n string, ep *awsconn.Endpoint)(*Queue){
  return &Queue{Name:n, Endpoint: ep}
}

// Creates a queue by parsing the endpoint URL.
// This could potentially break down the road if
// amazon ever changes their SQS URI format substantially.
func NewQueueURL(ep *awsconn.Endpoint)(mq *Queue){
  _, name := path.Split(ep.URL.Path)
  mq = NewQueue(name, ep)
  return
}

// Deletes an SQS queue.
func (self *Queue)Delete(id auth.Signer)(err os.Error){
  if self == nil || self.Endpoint == nil { return os.NewError("Undefined endpoint!") }
  sqsReq, err := self.signedRequest(id, map[string]string{
    "Action": "DeleteQueue",
  })
  if err != nil { return }

  xresp := &deleteQueueResponse{}
  xerr := &Error{}
  err = self.Endpoint.SendParsable(sqsReq, xresp, xerr)

  if err != nil { return }
  return
}

// Pushes the raw bytes into SQS, returning the message UUID.
//
// NB, we don't do any verification of the MD5
func (self *Queue)Push(id auth.Signer, body []byte)(msgid *uuid.UUID, err os.Error){
  sqsReq, err := self.signedRequest(id, map[string]string{
    "Action": "SendMessage",
    "MessageBody": string(body),
  })
  if err != nil { return }

  xresp := &sendMessageResponse{}
  xerr := &errorResponse{}
  err = self.Endpoint.SendParsable(sqsReq, xresp, xerr)
  if err != nil { return }
  msgid, err = uuid.Parse(xresp.SendMessageResult.MessageId)
  return
}

// Helper function that Push()'s a string instead of []bytes.
func (self *Queue)PushString(id auth.Signer, body string)(*uuid.UUID, os.Error){
  buff := bytes.NewBufferString(body)
  return self.Push(id, buff.Bytes())
}

// Fetches messages from SQS.  Note SQS may raise an error if you
// set "lim" higher than 10 (lim <0 will default to MaxNumberOfMessages.).
// Timeout is 
func (self *Queue)FetchMessages(id auth.Signer, lim, timeout int)(m []*Message, err os.Error){
  if lim <= 0 { lim = MaxNumberOfMessages }
  parms := map[string]string{
    "Action": "ReceiveMessage",
    "AttributeName": "All",
    "MaxNumberOfMessages": strconv.Itoa(lim),
  }
  if timeout >= 0 { parms["Timeout"] = strconv.Itoa(timeout) }
  sqsReq, err := self.signedRequest(id, parms)
  if err != nil { return }
  xresp := &receiveMessageResponse{}
  xerr := &errorResponse{}
  err = self.Endpoint.SendParsable(sqsReq, xresp, xerr)
  if err != nil { return }
  m = make([]*Message, len(xresp.ReceiveMessageResult.Message))
  for mi := range(xresp.ReceiveMessageResult.Message) {
    m[mi], err = xresp.ReceiveMessageResult.Message[mi].Message()
    if err != nil { break }
  }
  return
}


func (self *Queue)DeleteMessage(id auth.Signer, m *Message)(err os.Error){
  sqsReq, err := self.signedRequest(id, map[string]string{
    "Action": "DeleteMessage",
    "ReceiptHandle": m.ReceiptHandle,
  })
  if err != nil { return }

  xresp := &deleteMessageResponse{}
  xerr := &Error{}
  err = self.Endpoint.SendParsable(sqsReq, xresp, xerr)

  if err != nil { return }
  return
}

func (self *Queue)signedRequest(id auth.Signer, params map[string]string)(req *http.Request, err os.Error){
  req = self.Endpoint.NewHTTPRequest("GET", self.Endpoint.GetURL().Path, maptools.StringStringToStringStrings(params), nil)
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
