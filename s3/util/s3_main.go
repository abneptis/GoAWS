package s3_util

import (
  "aws/s3"
  . "aws/flags"
  . "aws/util/common"
)

import (
  "flag"
  "fmt"
  "http"
  "os"
  "path"
)

func Main()(err os.Error){
  s3_ep := flag.String("s3-endpoint", "https://" + s3.USEAST_HOST + "/", "S3 Endpoint to use")
  flag.Parse()
  var svc *s3.Service

  s3_url, err := http.ParseURL(*s3_ep)

  if err != nil { return }

  if flag.NArg() == 0 {
    return os.NewError("USAGE: s3 [list|get|put] [...]")
  }
  svc = s3.NewService(s3_url)

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
        err = svc.Bucket(bucket).ListKeys(signer, "", "", "", keys) 
        close(keys)
        return
      },
    },
    "buckets": UserCall{
      F: func()(err os.Error){
        lb, err := svc.ListBuckets(signer)
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
          _, err = svc.Bucket(bucket).GetKey(signer, key, os.Stdout)
        }
        return
      },
    },
    "exists": UserCall{
      F: func()(err os.Error){
        bucket := flag.Arg(1)
        key    := flag.Arg(2)
        if bucket == "" || key == "" {
          fmt.Printf("Usage: exists BUCKET KEY\n")
          os.Exit(1)
        }
        err = svc.Bucket(bucket).Exists(signer, key)
        return
      },
    },
    "drop": UserCall{
      F: func()(err os.Error){
        bucket := flag.Arg(1)
        if bucket == "" { return os.NewError("Usage: drop BUCKET") }
        err = svc.DeleteBucket(signer, bucket)
        return
      },
    },
    "create": UserCall{
      F: func()(err os.Error){
        bucket := flag.Arg(1)
        if bucket == "" { return os.NewError("Usage: create BUCKET") }
        err = svc.CreateBucket(signer, bucket)
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
          err = svc.Bucket(bucket).Delete(signer, keys[i])
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
          err = svc.Bucket(bucket).PutLocalFile(signer, path.Join(prefix, path.Base(keys[i])), keys[i])
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
