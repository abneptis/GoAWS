package elb

import "os"
import "http"

import "com.abneptis.oss/aws/awsconn"

func GetEndpoint(region string, proxy *http.URL)(ep *awsconn.Endpoint, err os.Error){
  url, err := http.ParseURL("http://elasticloadbalancing." + region + ".amazonaws.com")
  if err == nil {
    ep = awsconn.NewEndpoint(url, proxy)
  }
  return
}

