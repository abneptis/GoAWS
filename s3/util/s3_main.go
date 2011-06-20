package s3_util

import (
  "com.abneptis.oss/aws"
  "com.abneptis.oss/aws/s3"
)

import (
  "flag"
  "fmt"
  "net"
  "os"
  "crypto/tls"
)

func Main(){
  flag.Parse()
  if flag.NArg() == 0 {
    fmt.Printf("USAGE: s3 [list|get|put] [...]\n")
    os.Exit(1)
  }
  c := aws.NewConn(func()(c net.Conn, err os.Error){
    return tls.Dial("tcp", "s3.amazonaws.com", nil)
  })
  switch flag.Arg(0) {
    case "listbuckets":
      s3.ListBuckets(nil, c)
      fmt.Printf("C == %+v\n", c)
    default:
      fmt.Printf("s3: Unknown subcommand %s\n", flag.Arg(0))
      os.Exit(1)
  }
}
