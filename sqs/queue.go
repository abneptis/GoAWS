package sqs

import "com.abneptis.oss/aws"
import "com.abneptis.oss/aws/auth"
import "com.abneptis.oss/uuid"
import "bytes"
import "os"
import "path"
import "strconv"
import "xml"

var MaxNumberOfMessages = 10

type Queue struct {
  Name string
  Endpoint *aws.Endpoint
}

func NewQueue(n string, ep *aws.Endpoint)(*Queue){
  return &Queue{Name:n, Endpoint: ep}
}

func NewQueueURL(ep *aws.Endpoint)(mq *Queue){
  _, name := path.Split(ep.URL.Path)
  mq = NewQueue(name, ep)
  return
}

func (self *Queue)Delete(id auth.Signer)(err os.Error){
  sqsReq, err := NewSQSRequest(map[string]string{
    "Action": "DeleteQueue",
    "AWSAccessKeyId": string(id.PublicIdentity()),
  })
  if err != nil { return }
  resp, err := SignAndSendSQSRequest(id, "GET", self.Endpoint, &sqsReq)
  if err == nil {
    if resp.StatusCode == 200 {
      parser := xml.NewParser(resp.Body)
      xresp := createQueueResponse{}
      err = parser.Unmarshal(&xresp, nil)
    } else {
      err = os.NewError("Received unexpected status: " + resp.Status)
    }
  }
  return
}

type sendMessageResponse struct {
  SendMessageResult  sendMessageResult
}

type sendMessageResult struct {
  MD5OfMessageBody string
  MessageId string
}


// NB, we don't do any verification of the MD5
func (self *Queue)Push(id auth.Signer, body []byte)(msgid *uuid.UUID, err os.Error){
   sqsReq, err := NewSQSRequest(map[string]string{
    "Action": "SendMessage",
    "AWSAccessKeyId": string(id.PublicIdentity()),
    "MessageBody": string(body),
  })
  if err != nil { return }
  resp, err := SignAndSendSQSRequest(id, "GET", self.Endpoint, &sqsReq)
  //bb, _ := http.DumpResponse(resp, true)
  //os.Stdout.Write(bb)
  if err == nil {
    if resp.StatusCode == 200 {
      parser := xml.NewParser(resp.Body)
      xresp := sendMessageResponse{}
      err = parser.Unmarshal(&xresp, nil)
      if err == nil { msgid, err = uuid.Parse(xresp.SendMessageResult.MessageId) }
    } else {
      err = os.NewError("Received unexpected status: " + resp.Status)
    }
  }
  return
}

func (self *Queue)PushString(id auth.Signer, body string)(*uuid.UUID, os.Error){
  buff := bytes.NewBufferString(body)
  return self.Push(id, buff.Bytes())
}

type receiveMessageResponse struct {
  ReceiveMessageResult receiveMessageResult
}

type receiveMessageResult struct {
  Message []*rawMessage
}

func (self *Queue)FetchMessages(id auth.Signer, lim, timeout int)(m []*Message, err os.Error){
  if lim <= 0 { lim = MaxNumberOfMessages }
   sqsReq, err := NewSQSRequest(map[string]string{
    "Action": "ReceiveMessage",
    "AWSAccessKeyId": string(id.PublicIdentity()),
    "AttributeName": "All",
    "MaxNumberOfMessages": strconv.Itoa(lim),
  })
  if err != nil { return }
  if timeout >= 0 {
    err = sqsReq.Set("VisibilityTimeout", strconv.Itoa(timeout))
    if err != nil { return }
  }
  resp, err := SignAndSendSQSRequest(id, "GET", self.Endpoint, &sqsReq)
  //bb, _ := http.DumpResponse(resp, true)
  //os.Stdout.Write(bb)
  if err == nil {
    if resp.StatusCode == 200 {
      parser := xml.NewParser(resp.Body)
      xresp := receiveMessageResponse{}
      err = parser.Unmarshal(&xresp, nil)
      if err == nil {
        m = make([]*Message, len(xresp.ReceiveMessageResult.Message))
        for mi := range(xresp.ReceiveMessageResult.Message) {
          m[mi], err = xresp.ReceiveMessageResult.Message[mi].Message()
          if err != nil { break }
        }
      }
    } else {
      err = os.NewError("Received unexpected status: " + resp.Status)
    }
  }
  return
}

