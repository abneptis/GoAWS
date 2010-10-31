package aws
/* 
  Copyright (c) 2010, Abneptis LLC.
  See COPYRIGHT and LICENSE for details.
*/

import "http"
import "os"

/*
  Simple wrapper that does an http.ClientConn Write/Read and returns
  the response/errors upstream as appropriate.
*/
func SendRequest(cc *http.ClientConn, req *http.Request)(resp *http.Response, err os.Error){
  err = cc.Write(req)
  if err == nil {
    resp, err = cc.Read()
  }
  return
}
