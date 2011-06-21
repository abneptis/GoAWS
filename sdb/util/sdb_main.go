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
    "rm":UserCall{
       F: func()(err os.Error){
        if flag.NArg() < 3 {
          return os.NewError("Usage: rm domain_name item [...]")
        }
        d := s.Domain(flag.Arg(1))
        args := flag.Args()[2:]
        for i := range(args){
          err = d.DeleteAttribute(signer, args[i], nil, nil)
          if err != nil { return }
        }
        return
      },
   },
    "get":UserCall{
       F: func()(err os.Error){
        if flag.NArg() < 3 {
          return os.NewError("Usage: get domain_name item [...]") 
        }
        d := s.Domain(flag.Arg(1))
        args := flag.Args()[2:]
        for i := range(args){
          var attrs []sdb.Attribute
          attrs, err = d.GetAttribute(signer, args[i], nil, false)
          if err != nil { return }
          fmt.Printf("Item: %+v", attrs)
        }
        return
      },
   },
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
    "select":UserCall{
      F: func()(err os.Error){
        if flag.NArg() < 3 || flag.NArg() > 4 {
          return os.NewError("Usage: select ('*'|col,col2,...) domain_name [extended expression]")
        }
        colstr := flag.Arg(1) 
        d := s.Domain(flag.Arg(2))
        expr := flag.Arg(3) 
        c := make(chan sdb.Item)
        go func(){
          for i := range(c) {
            fmt.Printf("%s\n", i.Name)
            for ai := range(i.Attribute){
              fmt.Printf("\t%s\t%s\n", i.Attribute[ai].Name, i.Attribute[ai].Value)
            }
          }
        }()
        err = d.Select(signer, colstr, expr, true, c) 
        close(c)

        return
      },
    },

    "domains":UserCall{
      F: func()(err os.Error){
        doms, err := s.ListDomains(signer)
        for i := range(doms){
          fmt.Printf("%s\n", doms[i])
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
