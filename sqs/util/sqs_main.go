package sqs_util

import (
  "aws/sqs"
  . "aws/flags"
  . "aws/util/common"
)

import (
  "flag"
  "fmt"
  "http"
  "io"
  "os"
)

func Main()(err os.Error){
  var ep_url_str *string = flag.String("sqs-endpoint","https://queue.amazonaws.com/", "Endpoint to use")
  var flag_default_timeout *int = flag.Int("sqs-queue-timeout", 90, "Queue timeout (create/delete)")
  var flag_pop_timeout     *int = flag.Int("sqs-message-timeout", 90, "Queue timeout (pop/peek)")
  flag.Parse()

  if flag.NArg() == 0 {
    fmt.Printf("USAGE: sqs [list|peek|pop|push] [...]\n")
    os.Exit(1)
  }

  signer, err := DefaultSigner()
  if err != nil { return }
  ep_url, err := http.ParseURL(*ep_url_str)
  if err != nil { return }
  s := sqs.NewService(ep_url)
  

  calls := Calls{
    "create":UserCall{
      F: func()(err os.Error){
        Q, err := s.CreateQueue(signer, flag.Arg(1),  *flag_default_timeout)
        if err == nil {
          fmt.Printf("%s\n", Q.URL)
        }
        return
      },
    },
    "list":UserCall{
      F: func()(err os.Error){
        qs, err := s.ListQueues(signer, "")
        if err == nil {
          for i := range(qs) {
            fmt.Printf("%s\n", qs[i])
          }
        }
        return
      },
    },
    "drop":UserCall{
      F: func()(err os.Error){
        Q, err := s.CreateQueue(signer, flag.Arg(1),  *flag_default_timeout)
        if err == nil {
          err = Q.DeleteQueue(signer)
        }
        return
      },
    },
    "push":UserCall{
      F: func()(err os.Error){
        if flag.NArg() < 2 {
          return os.NewError("Usage: push queuename")
        }
        Q, err := s.CreateQueue(signer, flag.Arg(1),  *flag_default_timeout)
        if err == nil {
          var n int
          lr := io.LimitReader(os.Stdin, sqs.MAX_MESSAGE_SIZE)
          buff := make([]byte, sqs.MAX_MESSAGE_SIZE)
          n, err = io.ReadFull(lr, buff)
          if err == nil || err == io.ErrUnexpectedEOF {
            buff = buff[0:n]
            err = Q.Push(signer, buff)
          }
        }
        return
      },
    },
    "rm":UserCall{
      F: func()(err os.Error){
        if flag.NArg() < 3 {
          return os.NewError("Usage: peek queuename receipthandle")
        }
        Q, err := s.CreateQueue(signer, flag.Arg(1),  *flag_default_timeout)
        if err == nil {
          err = Q.Delete(signer, flag.Arg(2)) 
        }
        return
      },
    },
    "peek":UserCall{
      F: func()(err os.Error){
        if flag.NArg() < 2 {
          return os.NewError("Usage: peek queuename")
        }
        Q, err := s.CreateQueue(signer, flag.Arg(1),  *flag_default_timeout)
        var body []byte
        var id string
        if err == nil {
          body, id, err = Q.Peek(signer, *flag_pop_timeout) 
        }
        if err == nil {
          fmt.Printf("# MessageId %s\n", id)
          os.Stdout.Write(body)
        }
        return
      },
    },
 }

  if call, ok := calls[flag.Arg(0)]; ok {
    err = call.F()
  } else {
    err = os.NewError("Unknown sub-function")
  }

  return
}
