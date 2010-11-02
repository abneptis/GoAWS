package main

import "bytes"
import "io"
import "flag"
import "os"

type nopCloser struct { io.Reader }
func (nopCloser)Close()(n os.Error){return}

var bodyString = flag.String("string","", "String to write to bucket")
var bodyFile   = flag.String("file","", "File to write to bucket")

func getBodyString()(s string, l int64){
  if bodyString != nil && *bodyString != "" {
    s = *bodyString
    l = int64(len(s))
  }
  return
}

func GetRC()(rc io.ReadCloser, l int64, err os.Error){
  s, l := getBodyString()
  if s != "" {
    rc = nopCloser{bytes.NewBufferString(s)}
  } else {
    if bodyFile == nil || *bodyFile == "-" {
      rc = os.Stdin
      l = -1
    } else {
      fp, err := os.Open(*bodyFile, os.O_RDONLY, 0)
      if nil == err {
        fi, err := fp.Stat()
        if err == nil {
          l = fi.Size
        }
        rc = fp
      }
    }
  }
  return
}

