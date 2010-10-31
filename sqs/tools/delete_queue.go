package main


//import "com.abneptis.oss/aws"
//import "com.abneptis.oss/aws/sqs"
//import "http"
import "flag"
import "log"

func main(){
  flag.Parse()
  q, err := GetQueue()
  if err != nil {
    log.Exitf("Couldn't get queue: %v\n", err)
  }
  id, err := GetAWSIdentity()
  if err != nil {
    log.Exitf("Unable to get AWS identity: %v\n", err)
  }
  err = q.Delete(id)
  if err != nil {
    log.Exitf("Failed to delete queue: %v\n", err)
  }
}
