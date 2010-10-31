// Amazon AWS Simple Queue Service interface.
//
// This package implements basic SQS functionality
// including queue creation, deletion, enumeration,
// and message pushing/fetching and deletion.
package sqs

import "http"
import "os"
import "strconv"
import "xml"
//import "log"
import "com.abneptis.oss/aws/auth"
//import "com.abneptis.oss/aws"
import "com.abneptis.oss/aws/awsconn"

type Service struct {
  Endpoint *awsconn.Endpoint
}

type createQueueResponse struct {
  CreateQueueResult createQueueResult
}

type createQueueResult struct {
  QueueUrl string
}

func NewService(ep *awsconn.Endpoint)(*Service){
  return &Service{Endpoint: ep}
}

func (self *Service)CreateQueue(id auth.Signer, name string, dvtimeout int)(mq *Queue, err os.Error){
  sqsReq, err := NewSQSRequest(map[string]string{
    "Action": "CreateQueue",
    "QueueName": name,
    "DefaultVisibilityTimeout": strconv.Itoa(dvtimeout),
    "AWSAccessKeyId": string(id.PublicIdentity()),
  })
  if err != nil { return }
  resp, err := SignAndSendSQSRequest(id, "GET", self.Endpoint, &sqsReq)
  if err == nil {
    if resp.StatusCode == 200 {
      parser := xml.NewParser(resp.Body)
      xresp := createQueueResponse{}
      err = parser.Unmarshal(&xresp, nil)
      if err == nil {
        qrl, err := http.ParseURL(xresp.CreateQueueResult.QueueUrl)
        if err == nil {
          ep := awsconn.NewEndpoint(qrl, self.Endpoint.ProxyURL)
          mq = NewQueueURL(ep)
        }
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


func (self *Service)ListQueues(id auth.Signer, prefix string)(out []*Queue, err os.Error){
  sqsReq, err := NewSQSRequest(map[string]string{
    "Action": "ListQueues",
    "AWSAccessKeyId": string(id.PublicIdentity()),
  })
  if err != nil { return }
  if len(prefix) > 0 {
    err = sqsReq.Set("QueueNamePrefix", prefix)
  }
  if err != nil { return }
  resp, err := SignAndSendSQSRequest(id, "GET", self.Endpoint, &sqsReq)
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
          ep := awsconn.NewEndpoint(url, self.Endpoint.ProxyURL)
          out[i] = NewQueueURL(ep)
        }
      }
    } else {
      err = os.NewError("Received unexpected status code" + resp.Status)
    }
  }
  //log.Printf("ListQueue: %v",  out)
  return
}

