package main


//import "com.abneptis.oss/aws"
//import "com.abneptis.oss/aws/sqs"
import "com.abneptis.oss/aws/s3"
//import "http"
import "flag"
import "fmt"
import "log"

func main(){
  flag.Parse()
  id, err := GetAWSIdentity()
  if err != nil {
    log.Fatalf("Unable to get AWS identity: %v\n", err)
  }
  ep, err := GetS3Endpoint()
  if err != nil {
    log.Fatalf("Unable to construct endpoint: %v\n", err)
  }
  out, err := s3.ListBuckets(id, ep)
  if err != nil {
    log.Fatalf("Unable to list buckets: %v\n", err)
  }
  for i := range(out){
    fmt.Printf("%s\n", out[i])
  }
}
