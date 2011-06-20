package s3_util

import (
  "aws"
  "aws/s3"
  . "aws/flags"
  . "aws/util/common"
)

import (
  "crypto/tls"
  "flag"
  "fmt"
  "net"
  "os"
  "path"
)

func Main()(err os.Error){
  flag.Parse()

  if flag.NArg() == 0 {
    fmt.Printf("USAGE: s3 [list|get|put] [...]\n")
    os.Exit(1)
  }

  c := aws.NewConn(func()(c net.Conn, err os.Error){
    return tls.Dial("tcp", "s3.amazonaws.com:443", nil)
  })

  signer, err := DefaultSigner()
  if err != nil {
    fmt.Printf("Couldn't extract signer\n")
    os.Exit(1)
  }

  calls := Calls{
    "ls":UserCall{
      F: func()(err os.Error){
        bucket := flag.Arg(1)
        keys := make(chan string)
        go func(){
          for i := range(keys){
            fmt.Printf("%s\n", i)
          }
        }()
        err = s3.ListKeys(signer, c, bucket,"", "", "", keys) 
        close(keys)
        return
      },
    },
    "list_buckets": UserCall{
      F: func()(err os.Error){
        lb, err := s3.ListBuckets(signer, c)
        for b := range(lb) {
          fmt.Println(lb[b])
        }
        return
      },
    },
    "cat": UserCall{
      F: func()(err os.Error){
        bucket := flag.Arg(1)
        key := flag.Arg(2)
        if bucket == "" || key == "" { 
          return os.NewError("Usage: get BUCKET KEY")
        }
        if err == nil {
          _, err = s3.GetKey(signer, c, bucket, key, os.Stdout)
        }
        return
      },
    },
    "key_exists": UserCall{
      F: func()(err os.Error){
        bucket := flag.Arg(1)
        key    := flag.Arg(2)
        if bucket == "" || key == "" {
          fmt.Printf("Usage: key_exists BUCKET KEY\n")
          os.Exit(1)
        }
        err = s3.KeyExists(signer, c, bucket, key)
        return
      },
    },
    "create": UserCall{
      F: func()(err os.Error){
        bucket := flag.Arg(1)
        if bucket == "" { 
          return os.NewError("Usage: create BUCKET")
        }
        err = s3.CreateBucket(signer, c, bucket)
        return
      },
    },
    "delete": UserCall{
      F: func()(err os.Error){
        bucket := flag.Arg(1)
        if bucket == "" || flag.NArg() < 3 { 
          return os.NewError("Usage: delete BUCKET KEY [KEY2...]")
        }
        keys := flag.Args()[2:]
        for key := range(keys){
          if err = s3.DeleteKey(signer, c, bucket, keys[key]) ; err != nil {
            return
          }
        }
        return
      },
    },
    "rm": UserCall{
      F: func()(err os.Error){
        bucket := flag.Arg(1)
        if bucket == "" || flag.NArg() < 3 { 
          return os.NewError("Usage: delete BUCKET KEY [KEY2...]")
        }
        keys := flag.Args()[2:]
        for i := range(keys){
          err = s3.DeleteKey(signer, c, bucket, keys[i])
          if err != nil { break }
        }
        return
      },
    },
    "put": UserCall{
      F: func()(err os.Error){
        bucket := flag.Arg(1)
        prefix := flag.Arg(2)
        if bucket == "" || flag.NArg() < 4 { 
          return os.NewError("Usage: put BUCKET PREFIX FILE [FILE2...]")
        }
        keys := flag.Args()[3:]
        for i := range(keys){
          err = s3.PutLocalFile(signer, c, bucket,
             path.Join(prefix, path.Base(keys[i])),
             keys[i])
          if err != nil { break }
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
