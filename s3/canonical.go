package s3

import (
  "http"
  "strings"
)


func CanonicalString(req *http.Request)(string){
  return strings.Join([]string{
    req.Method,
    req.Header.Get("Content-Md5"),
    req.Header.Get("Content-Type"),
    req.Form.Get("Expires"),
    req.URL.Path,
  }, "\n")
}

