package sqs

import "com.abneptis.oss/goaws"
//import "log"
import "http"
import "os"
import "path"
import "xml"

type Queue struct {
  Name string
  URL  *http.URL
  ProxyURL  *http.URL
}

func NewQueue(n string, u *http.URL, pu *http.URL)(*Queue){
  return &Queue{Name:n, URL: u, ProxyURL: pu}
}

func NewQueueURL(u *http.URL, pu *http.URL)(mq *Queue){
  mq =  &Queue{URL: u, ProxyURL: pu}
  _, mq.Name = path.Split(u.Path)
  return
}

func NewQueueString(n string, u string, pu *http.URL)(q *Queue, err os.Error){
  _u, err := http.ParseURL(u)
  if err == nil {
    q = NewQueue(n, _u, pu)
  }
  return
}

func (self *Queue)Delete(id goaws.Signer)(err os.Error){
  sqsReq, err := NewSQSRequest(map[string]string{
    "Action": "DeleteQueue",
    "AWSAccessKeyId": string(id.PublicIdentity()),
  })
  if err != nil { return }
  resp, err := SignAndSendSQSRequest(id, "GET", self.URL, self.ProxyURL, &sqsReq)
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

