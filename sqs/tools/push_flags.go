package main

//import "com.abneptis.oss/aws"
import "com.abneptis.oss/aws/sqs"
//import "http"
import "bytes"
import "flag"
import "os"
import "io"


var messageFile   = flag.String("message-file","-", "File to read message from")
var messageString = flag.String("message","", "Raw message to use")

func GetMessageReadCloser()(rc io.ReadCloser, err os.Error){
  if messageFile == nil || *messageFile == "-" {
    rc = os.Stdin
  } else {
    rc, err = os.Open(*messageFile)
  }
  return
}

func GetMessage()(out []byte, err os.Error){
  if messageString != nil && *messageString != "" {
    buff := bytes.NewBufferString(*messageString)
    out = buff.Bytes()
  } else {
    rc, err := GetMessageReadCloser()
    out = make([]byte, sqs.MaxMessageSize)
    ipos := 0
    for n, err := rc.Read(out[ipos:sqs.MaxMessageSize]);
        err != os.EOF && ipos < sqs.MaxMessageSize ;
        n, err = rc.Read(out[ipos:sqs.MaxMessageSize]){
      if n > 0 {
        ipos += n
      }
      if err != os.EOF && err !=nil{
        break
      }
    }
    out = out[0:ipos]
    if err == os.EOF {
      err = nil
    } else {
      if err == nil { err = os.NewError("No end-of-file found") }
    }
  }
  return
}

