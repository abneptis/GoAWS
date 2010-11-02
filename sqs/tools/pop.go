package main

//import "com.abneptis.oss/aws"
//import "com.abneptis.oss/aws/sqs"
import "http"
import "json"
import "flag"
import "log"
import "os"
import "fmt"

var max = flag.Int("max-messages", 1, "Maximum number of messages to receive")
var timeout = flag.Int("fetch-timeout", -1, "MessageVisibilityTimeout (-1 for queue default)")
var outform = flag.String("outform", "text", "Output format (*text|shell|json|raw)")

func main(){
  flag.Parse()
  q, err := GetQueue()
  if err != nil {
    log.Exitf("Couldn't create queue: %v\n", err)
  }
  id, err := GetAWSIdentity()
  if err != nil {
    log.Exitf("Couldn't get identity: %v\n", err)
  }
  msgs, err := q.FetchMessages(id, *max, *timeout)
  if err != nil {
    log.Exitf("Couldn't pop from queue: %v\n", err)
  }
  log.Printf("#[%d messages received]", len(msgs))
  _form := *outform
  for mi := range(msgs) {
    err = q.DeleteMessage(id, msgs[mi])
    if err != nil {
      log.Printf("Couldn't delete message, not displaying (%v)", err)
      continue
    }
    switch _form {
      case "text":
        fmt.Printf("MessageID\t%s\n", msgs[mi].MessageId.String())
        fmt.Printf("ReceiptHandle\t%s\n", msgs[mi].ReceiptHandle)
        fmt.Printf("MessageBody\t%s\n", string(msgs[mi].Body))
      case "shell":
        // TODO: Escape these properly.
        fmt.Printf("export SQSMessage%dId=\"%s\";\n", mi, msgs[mi].MessageId.String())
        fmt.Printf("export SQSMessage%dReceipt=\"%s\";\n", mi, msgs[mi].ReceiptHandle)
        fmt.Printf("export SQSMessage%dBody=\"%s\";\n", mi, http.URLEscape(string(msgs[mi].Body)))
      case "raw":
        os.Stdout.Write(msgs[mi].Body)
    }
  }
  if _form == "json" {
    enc := json.NewEncoder(os.Stdout)
    err := enc.Encode(msgs)
    if err != nil { log.Exitf("Error decoding messages: %v", err) }
  }
}

