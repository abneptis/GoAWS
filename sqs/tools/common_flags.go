// Command line tool functions.
//
// I should probably lower-case the functions, since these functions
// generally expect pre-defined flags, and are simply helpers for
// the various tools and exist within package "main" -- but at least
// for akid/skid it could be relevant to other kits.
//
// Todo: evaluate moving those to aws/auth.
package main

import "com.abneptis.oss/aws/awsconn"
import "com.abneptis.oss/aws/sqs"
import "com.abneptis.oss/aws"

import "com.abneptis.oss/cryptools"

import "flag"
import "http"
import "os"
import "log"


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

func GetEndpoint()(ep *awsconn.Endpoint, err os.Error){
  purl, err := GetProxyURL()
  if err != nil { return }
  epurl, err := GetEndpointURL()
  if err != nil { return }
  ep = awsconn.NewEndpoint(epurl, purl)
  return
}

func GetSQSService()(s *sqs.Service, err os.Error){
  ep,err := GetEndpoint()
  if err != nil { return }
  s = sqs.NewService(ep)
  return
}

func GetAWSIdentity()(s cryptools.NamedSigner, err os.Error){
  if accessKeyId == nil || secretKeyId == nil || *accessKeyId == "" || *secretKeyId == "" {
   return nil, os.NewError("-access-key and -secret-key are both required")
  }
  return aws.NewIdentity("sha256", *accessKeyId, *secretKeyId)
}


func GetQueue()(Q *sqs.Queue, err os.Error){
  proxyURL, err := GetProxyURL()
  if err != nil { return }
  id, err := GetAWSIdentity()
  if err != nil { return }
  if queueURL != nil && *queueURL != ""{
    qrl, err := http.ParseURL(*queueURL)
    if err == nil {
      ep := awsconn.NewEndpoint(qrl, proxyURL)
      Q = sqs.NewQueueURL(ep)
    }
  } else if queueName != nil && *queueName != "" {
    ep, err := http.ParseURL(*sqsEndpoint)
    if err == nil {
      log.Printf("Parsed EP url: %v", ep)
      ep := awsconn.NewEndpoint(ep, proxyURL)
      _sqs := sqs.NewService(ep)
      Q, err = _sqs.CreateQueue(id, *queueName, 90)
      if err != nil {
        log.Exitf("Qerr: [%v]", err)
      }
    }
  } else {
    err = os.NewError("Either Queue(+Endpoint) or QueueURL are required")
  }
  return
}
