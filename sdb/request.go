package sdb

import (
  "http"
)


func newRequest(method string, url *http.URL, hdrs http.Header, params http.Values)(req *http.Request){
  req = &http.Request {
    Method: method,
    URL: &http.URL {
      Path: url.Path,
      RawQuery: url.RawQuery,
    },
    Host: url.Host,
    Header: hdrs,
    Form: params,
  }
  return
}


