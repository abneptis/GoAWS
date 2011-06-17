package main

import "com.abneptis.oss/aws"
import "flag"
import "log"
import "fmt"

func main(){
  aws.AwsIDFlags()
  flag.Parse()
  q, err := GetQueue()
  if err != nil {
    log.Fatalf("Couldn't create queue: %v\n", err)
  }
  id, err := GetAWSIdentity()
  if err != nil {
    log.Fatalf("Couldn't get identity: %v\n", err)
  }
  msg, err := GetMessage()
  if err != nil {
    log.Fatalf("Couldn't read message: %v", err)
  }
  fmt.Printf("Read Message [%d]: %v\n", len(msg), msg)
  mid, err := q.Push(id, msg)
  if err != nil {
    log.Fatalf("Couldn't push to queue: %v\n", err)
  }
  fmt.Printf("%s\n",mid.String())
}

