package sdb_util

import (
  "aws/sdb"
  . "aws/flags"
  . "aws/util/common"
)

import (
  "flag"
  "fmt"
  "http"
  "os"
)

func Main()(err os.Error){
  var flag_sdb_endpoint *string = flag.String("sdb-endpoint","https://sdb.amazonaws.com/", "Endpoint to use")
  flag.Parse()

  if flag.NArg() == 0 {
    fmt.Printf("USAGE: sdb [list|create|destroy] [...]\n")
    os.Exit(1)
  }

  signer, err := DefaultSigner()
  if err != nil { return }
  ep_url, err := http.ParseURL(*flag_sdb_endpoint)
  if err != nil { return }
  s := sdb.NewService(ep_url)
  

  calls := Calls{
    "create":UserCall{
      F: func()(err os.Error){
        if flag.NArg() != 2 {
          return os.NewError("Usage: create domain_name")
        }
        err = s.CreateDomain(signer, flag.Arg(1))
        return
      },
    },
    "drop":UserCall{
      F: func()(err os.Error){
        if flag.NArg() != 2 {
          return os.NewError("Usage: drop domain_name")
        }
        err = s.DestroyDomain(signer, flag.Arg(1))
        return
      },
    },
    "list":UserCall{
      F: func()(err os.Error){
        doms, err := s.ListDomains(signer)
        for i := range(doms){
          fmt.Printf("%s\n", doms[i])
        }
        return
      },
    },
    "rm":UserCall{
      F: func()(err os.Error){
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
