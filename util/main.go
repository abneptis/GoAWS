package main

import (
  "flag"
  "s3_util"
  "os"
  "fmt"
)

func main(){
  flag.Parse()
  if flag.NArg() == 0 {
    fmt.Printf("USAGE: aws [sqs|s3|ec2] ...\n")
    os.Exit(1)
  }
  cmd := flag.Arg(0)
  fmt.Printf("RawArgs: %v\n", os.Args) 	
  os.Args = os.Args[1:]
  switch cmd {
    case "s3": s3_util.Main()
    default: fmt.Printf("Unknown module: %s\n", cmd)
  }
  os.Exit(1)
}
