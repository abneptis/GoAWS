package main


//import "com.abneptis.oss/aws"
//import "com.abneptis.oss/aws/sqs"
//import "http"
import "flag"
import "fmt"
import "log"

func main(){
  flag.Parse()
  id, err := GetAWSIdentity()
  if err != nil {
    log.Exitf("Unable to get AWS identity: %v\n", err)
  }
  ep, err := GetSQSService()
  if err != nil {
    log.Exitf("Unable to construct endpoint: %v\n", err)
  }
  qs, err := ep.ListQueues(id, "")
  if err != nil {
    log.Exitf("Unable to list queues: %v\n", err)
  }
  for i := range(qs){
    fmt.Printf("%s\t%s\n", qs[i].Name, qs[i].Endpoint.URL.String())
  }
}
