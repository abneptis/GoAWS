package main


//import "com.abneptis.oss/goaws"
//import "com.abneptis.oss/goaws/sqs"
//import "http"
import "flag"
import "log"
import "fmt"

func main(){
  flag.Parse()
  q, err := GetQueue()
  if err != nil {
    log.Exitf("Couldn't create queue: %v\n", err)
  }
  fmt.Printf("%v\n", q.URL)
}
