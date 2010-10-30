package main

import "com.abneptis.oss/goaws"
import "com.abneptis.oss/goaws/sqs"
import "http"
import "flag"
import "log"
import "fmt"

func main(){
  url, _ := http.ParseURL("http://queue.amazonaws.com/")
  //qs := sqs.NewEndpoint(url, nil)
  //qName := flag.Arg(0)
  //id, err := goaws.NewIdentity("sha256", "AKIAJVBFI6WMZSSBBFRQ",
  //                  "m6ILj9yFxPlbJfITm64KoSHzweRXfZzNArP80OkD")
  //if err != nil {
  // log.Exitf("Couldn't create identity: %v\n", err)
  //}
  //q, err := qs.CreateQueue(id, qName, 90)
  //if err != nil {
  //  log.Exitf("Couldn't find/create queue: %v\n", err)
  // }
  
}
