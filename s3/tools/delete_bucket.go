package main


//import "com.abneptis.oss/aws"
//import "com.abneptis.oss/aws/sqs"
import "com.abneptis.oss/aws/s3"
//import "http"
import "flag"
//import "fmt"
import "log"

func main(){
  flag.Parse()
  id, err := GetAWSIdentity()
  if err != nil {
    log.Exitf("Unable to get AWS identity: %v\n", err)
  }
  ep, err := GetS3Endpoint()
  if err != nil {
    log.Exitf("Unable to construct endpoint: %v\n", err)
  }
  bucketName := flag.Arg(0)
  b := s3.NewBucket(ep, bucketName)
  err = b.Destroy(id)
  if err != nil {
    log.Exitf("Unable to create bucket: %v\n", err)
  }
}
