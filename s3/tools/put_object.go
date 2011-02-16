package main


//import "com.abneptis.oss/aws"
//import "com.abneptis.oss/aws/sqs"
import "com.abneptis.oss/aws/s3"
//import "http"
import "flag"
//import "fmt"
import "log"
//import "io"
//import "os"

func main(){
  flag.Parse()
  id, err := GetAWSIdentity()
  if err != nil {
    log.Fatalf("Unable to get AWS identity: %v", err)
  }
  ep, err := GetS3Endpoint()
  if err != nil {
    log.Fatalf("Unable to construct endpoint: %v", err)
  }
  bucket := s3.NewBucket(ep, flag.Arg(0))
  rc, ll, err := GetRC()
  if err != nil {
    log.Fatalf("Error reading data for key: %v",err)
  }
  //log.Fatalf("Error reading data for key: %v/%v/%v",rc,ll,err)
  err = bucket.PutKey(id, flag.Arg(1), "","", ll, rc)
  if err != nil {
    log.Fatalf("Couldn't put key: %v", err)
  }
  //io.Copy(os.Stdout,o.Body)
}
