package goaws
/* 
  Copyright (c) 2010, Abneptis LLC.
  See COPYRIGHT and LICENSE for details.
*/

import "http"
import "net"
import "os"
import "strings"
import "bufio"

func ExtractURLHostPort(url *http.URL)(hostport string, err os.Error){
  if url == nil || url.Host == "" { return "", os.NewError("Invalid URL (or empty Host)") }
  portIdx := strings.LastIndex(url.Host, ":")
  var host,port string
  if portIdx >= 0 {
    host = url.Host[0:portIdx]
    port = url.Host[portIdx+1:]
  } else {
    host = url.Host
  }
  if port == "" {
    switch url.Scheme {
      case "http": port = "80"
      case "https": port = "443"
      case "ftp": port = "21"
      case "smtp": port = "25"
      default:
        err = os.NewError("No port specified, and unknown scheme: " + url.Scheme)
    }
  }
  hostport = host + ":" + port
  return
}

func ClientConnection(nnet, laddr string, url, purl *http.URL, r *bufio.Reader)(cc *http.ClientConn, err os.Error){
  host, err := ExtractURLHostPort(url)
  if purl != nil {
    host, err = ExtractURLHostPort(purl)
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
