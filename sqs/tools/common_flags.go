package main

import "com.abneptis.oss/aws"
import "com.abneptis.oss/aws/sqs"
import "com.abneptis.oss/aws/auth"

import "flag"
import "http"
import "os"
//import "log"


var secretKeyId = flag.String("secret-key","","Secret key ID")
var accessKeyId = flag.String("access-key","","Access key ID")
var queueName = flag.String("queue","","Name of the queue to use")
var queueURL  = flag.String("queue-url","","Direct queue URL")
var sqsEndpoint = flag.String("endpoint","http://queue.amazonaws.com/","SQS Endpoint")
var proxy = flag.String("proxy","","Proxy to use")

func GetEndpointURL()(u *http.URL, err os.Error){
  return http.ParseURL(*sqsEndpoint)
}

func GetProxyURL()(u *http.URL, err os.Error){
  if proxy != nil && (*proxy) != "" {
    u, err = http.ParseURL(*proxy)
  }
  return
}

func GetEndpoint()(ep *aws.Endpoint, err os.Error){
  purl, err := GetProxyURL()
  if err != nil { return }
  epurl, err := GetEndpointURL()
  if err != nil { return }
  ep = aws.NewEndpoint(epurl, purl)
  return
}

func GetSQSService()(s *sqs.Service, err os.Error){
  ep,err := GetEndpoint()
  if err != nil { return }
  s = sqs.NewService(ep)
  return
}

func GetAWSIdentity()(s auth.Signer, err os.Error){
  if accessKeyId == nil || secretKeyId == nil || *accessKeyId == "" || *secretKeyId == "" {
   return nil, os.NewError("-access-key and -secret-key are both required")
  }
  return auth.NewIdentity("sha256", *accessKeyId, *secretKeyId)
}


func GetQueue()(Q *sqs.Queue, err os.Error){
  proxyURL, err := GetProxyURL()
  if err != nil { return }
  id, err := GetAWSIdentity()
  if err != nil { return }
  if queueURL != nil && *queueURL != ""{
    qrl, err := http.ParseURL(*queueURL)
    if err == nil {
      ep := aws.NewEndpoint(qrl, proxyURL)
      Q = sqs.NewQueueURL(ep)
    }
  } else if queueName != nil && *queueName != "" {
    ep, err := http.ParseURL(*sqsEndpoint)
    if err == nil {
      //log.Printf("Parsed EP url: %v", ep)
      ep := aws.NewEndpoint(ep, proxyURL)
      _sqs := sqs.NewService(ep)
      Q, err = _sqs.CreateQueue(id, *queueName, 90)
      //log.Printf("Q, QUrl, err: %p, %v", Q, err)
    }
  } else {
    err = os.NewError("Either Queue(+Endpoint) or QueueURL are required")
  }
  return
}
