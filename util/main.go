package main

import (
  . "aws/flags" // AWS ID Flags
  "aws/s3/s3_util"
  "aws/sqs/sqs_util"
)

import (
  "flag"
  "os"
  "fmt"
)

func main(){
  flag.Parse()
  if flag.NArg() == 0 {
    fmt.Printf("USAGE: aws [s3] ...\n")
    os.Exit(1)
  }
  cmd := flag.Arg(0)
  // fmt.Printf("RawArgs: %v\n", os.Args) 	
  os.Args = os.Args[1:]
  var err os.Error
  switch cmd {
    case "s3": err = s3_util.Main()
    case "sqs": err = sqs_util.Main()
    default: err = os.NewError("Unknown module:" + cmd)
  }
  if err != nil {
    fmt.Printf("Error: %v\n", err)
    os.Exit(1)
  }
  os.Exit(0)
  UseFlags() // we want the side effects of import...
}
