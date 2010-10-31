package aws
/* 
  Copyright (c) 2010, Abneptis LLC.
  See COPYRIGHT and LICENSE for details.
*/

import "http"
import "os"

func SendRequest(cc *http.ClientConn, req *http.Request)(resp *http.Response, err os.Error){
  err = cc.Write(req)
  if err == nil {
    resp, err = cc.Read()
  }
  return
}
