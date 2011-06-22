package main

import (
  "bytes"
  . "log"
  "flag"
  "http"
  "io"
  "net"
  "path"
  "json"
  "os"
)

type Service struct {
  conf conf
}

func (self Service)ServeHTTP(rw http.ResponseWriter, req *http.Request){
  host, _, err := net.SplitHostPort(req.Host)
  if err != nil {
    host = req.Host
  }
  if req.Method != "GET" {
    http.Error(rw, "Only GET is supported", http.StatusMethodNotAllowed)
    return
  }
  if req.URL.RawQuery != "" {
    http.Error(rw, "Sorry, we do not permit query-string parameters", http.StatusBadRequest)
    return
  }
  if domain, ok := self.conf[host] ; ok {
    buff := bytes.NewBuffer(nil)
    Printf("%s %s %s", req.RemoteAddr, host, req.URL.Path)
    // We could actually get this to be a bit faster if we simply passed in the http.ResponseWriter,
    // but since we want to pass-through headers (content-type, etc) in this example,
    // we have to buffer it off.
    s3_headers, err := domain.Bucket.GetKey(domain.Identity, path.Join(domain.Prefix, req.URL.Path), buff)

    if err == nil {
       reqhdrs := rw.Header()
       for k, v := range(s3_headers){
         reqhdrs[k] = v
       }
      n, err := io.Copy(rw, buff)
      Printf("%s %s %s %d %v", req.RemoteAddr, host, req.URL.Path, n, err)
    } else {
      http.Error(rw, err.String(), http.StatusInternalServerError)
    }
    return
  } else {
    Printf("%s %s %s - Host Unknown", req.RemoteAddr, host, req.URL.Path)
    http.Error(rw, "Invalid host", http.StatusForbidden)
  }
}

var flag_bind_addr *string = flag.String("listen","127.0.0.1:8080","Address/port to listen to")

func main(){
  conf := conf{}
  flag.Parse()
  fp, err := os.Open("config.json")
  if err != nil {
    Fatalf("Couldn't open config file: %v", err)
  }
  err = json.NewDecoder(fp).Decode(&conf)
  if err != nil {
    Fatalf("Couldn't parse config file: %v", err)
  }
  err = http.ListenAndServe(*flag_bind_addr, Service{conf:conf})
  if err != nil {
    Fatalf("Couldn't setup listener: %v", err)
  }
  return
}
