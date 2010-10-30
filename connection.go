package goaws
/* 
  Copyright (c) 2010, Abneptis LLC.
  See COPYRIGHT and LICENSE for details.
*/
import "com.abneptis.oss/urltools"

import "http"
import "net"
import "os"
import "bufio"

func ClientConnection(nnet, laddr string, url, purl *http.URL, r *bufio.Reader)(cc *http.ClientConn, err os.Error){
  host, err := urltools.ExtractURLHostPort(url)
  if purl != nil {
    host, err = urltools.ExtractURLHostPort(purl)
  }
  if err != nil { return }
  rawc, err := net.Dial(nnet, laddr, host)
  if err != nil { return }
  cc = http.NewClientConn(rawc, r)
  return
}

func SendRequest(cc *http.ClientConn, req *http.Request)(resp *http.Response, err os.Error){
  err = cc.Write(req)
  if err == nil {
    resp, err = cc.Read()
  }
  return
}
