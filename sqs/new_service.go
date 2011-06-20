package sqs

import (
  "aws"
)

import (
  "crypto"
  "http"
  "os"
  "strconv"
  "xml"
)

const (
  DEFAULT_HASH = crypto.SHA256
  MAX_MESSAGE_SIZE = 64*1024
)

type Service struct {
  URL  *http.URL
  conn *aws.Conn
}

func NewService(url *http.URL)(*Service){
  return &Service{
    URL: url,
    conn: aws.NewConn(aws.URLDialer(url, nil)),
  }
}

// Create a queue, returning the Queue object.
func (self *Service)CreateQueue(id *aws.Signer, name string, dvtimeout int)(mq *Queue, err os.Error){
  var resp *http.Response
  parms := http.Values{}
  parms.Set("Action","CreateQueue")
  parms.Set("QueueName",name)
  parms.Set("DefaultVisibilityTimeout",strconv.Itoa(dvtimeout))

  req := newRequest("GET", self.URL, nil, parms)
  err = signRequest(id, req)
  if err == nil {
    resp, err = self.conn.Request(req)
    if err == nil {
      if resp.StatusCode == http.StatusOK {
        xmlresp := createQueueResponse{}
        err = xml.Unmarshal(resp.Body, &xmlresp)
        if err == nil {
          var qrl *http.URL
          qrl, err = http.ParseURL(xmlresp.QueueURL)
          if err == nil {
            mq = NewQueue(qrl)
          } 
        }
      } else {
        err = os.NewError("Unexpected response")
      }
    }
  }

  return
}

type createQueueResponse struct {
  QueueURL string "CreateQueueResult>QueueUrl"
  RequestId string "ResponseMetadata>RequestId"
}

/*
<CreateQueueResponse xmlns="http://queue.amazonaws.com/doc/2009-02-01/">
  <CreateQueueResult>
   <QueueUrl>https://queue.amazonaws.com/930374178234/tQueue</QueueUrl>
  </CreateQueueResult>
  <ResponseMetadata>
   <RequestId>97c6d561-1c7d-4ac7-9ebc-b7f1c87baabf</RequestId>
  </ResponseMetadata>
</CreateQueueResponse>
*/
