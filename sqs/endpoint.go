package sqs

import "http"
import "os"
import "strconv"
import "xml"
//import "log"
import "com.abneptis.oss/goaws/auth"

type Endpoint struct {
  ProxyURL *http.URL
  URL *http.URL
}

func NewEndpoint(u, pu *http.URL)(*Endpoint){
  return &Endpoint {
   URL: u,
   ProxyURL: pu,
  }
}

type createQueueResponse struct {
  CreateQueueResult createQueueResult
}

type createQueueResult struct {
  QueueUrl string
}

func (self *Endpoint)CreateQueue(id auth.Signer, name string, dvtimeout int)(mq *Queue, err os.Error){
  sqsReq, err := NewSQSRequest(map[string]string{
    "Action": "CreateQueue",
    "QueueName": name,
    "DefaultVisibilityTimeout": strconv.Itoa(dvtimeout),
    "AWSAccessKeyId": string(id.PublicIdentity()),
  })
  if err != nil { return }
  resp, err := SignAndSendSQSRequest(id, "GET", self.URL, self.ProxyURL, &sqsReq)
  if err == nil {
    if resp.StatusCode == 200 {
      parser := xml.NewParser(resp.Body)
      xresp := createQueueResponse{}
      err = parser.Unmarshal(&xresp, nil)
      if err == nil {
        mq, err = NewQueueString(name, xresp.CreateQueueResult.QueueUrl, self.ProxyURL)
      }
    } else {
      err = os.NewError("Received unexpected status code" + resp.Status)
    }
  }
  return
}


type listQueuesResponse struct {
  ListQueuesResult listQueuesResult
}

type listQueuesResult struct {
  QueueUrl []string
}


func (self *Endpoint)ListQueues(id auth.Signer, prefix string)(out []*Queue, err os.Error){
  sqsReq, err := NewSQSRequest(map[string]string{
    "Action": "ListQueues",
    "AWSAccessKeyId": string(id.PublicIdentity()),
  })
  if err != nil { return }
  if len(prefix) > 0 {
    err = sqsReq.Set("QueueNamePrefix", prefix)
  }
  if err != nil { return }
  resp, err := SignAndSendSQSRequest(id, "GET", self.URL, self.ProxyURL, &sqsReq)
  if err == nil {
    if resp.StatusCode == 200 {
      //bb, err := http.DumpResponse(resp, true)
      //os.Stderr.Write(bb)
      parser := xml.NewParser(resp.Body)
      xresp := listQueuesResponse{}
      err = parser.Unmarshal(&xresp, nil)
      if err == nil {
        out = make([]*Queue, len(xresp.ListQueuesResult.QueueUrl))
        for i := range(xresp.ListQueuesResult.QueueUrl){
          url, err := http.ParseURL(xresp.ListQueuesResult.QueueUrl[i])
          if err != nil { break }
          out[i] = NewQueueURL(url, self.ProxyURL)
        }
      }
    } else {
      err = os.NewError("Received unexpected status code" + resp.Status)
    }
  }
  //log.Printf("ListQueue: %v",  out)
  return
}

