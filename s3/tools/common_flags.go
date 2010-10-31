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
import "com.abneptis.oss/aws/auth"
//import "com.abneptis.oss/aws/s3"

import "flag"
import "http"
import "os"
//import "log"


var secretKeyId = flag.String("secret-key","","Secret key ID")
var accessKeyId = flag.String("access-key","","Access key ID")
var s3Endpoint = flag.String("endpoint","http://s3.amazonaws.com/","S3 Endpoint")
var proxy = flag.String("proxy","","Proxy to use")

func GetProxyURL()(u *http.URL, err os.Error){
  if proxy != nil && (*proxy) != "" {
    u, err = http.ParseURL(*proxy)
  }
  return
}

func GetEndpointURL()(u *http.URL, err os.Error){
  return http.ParseURL(*s3Endpoint)
}

func GetS3Endpoint()(ep *awsconn.Endpoint, err os.Error){
  purl, err := GetProxyURL()
  if err != nil { return }
  epurl, err := GetEndpointURL()
  if err != nil { return }
  ep = awsconn.NewEndpoint(epurl, purl)
  return
}

func GetAWSIdentity()(s auth.Signer, err os.Error){
  if accessKeyId == nil || secretKeyId == nil || *accessKeyId == "" || *secretKeyId == "" {
   return nil, os.NewError("-access-key and -secret-key are both required")
  }
  return auth.NewIdentity("sha1", *accessKeyId, *secretKeyId)
}

