package main

import (
	. "log"
	"flag"
	"http"
	"net"
	"path"
	"json"
	"os"
	"strconv"
)

type Service struct {
	conf conf
}

func (self Service) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
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
	if domain, ok := self.conf[host]; ok {
		Printf("%s %s %s", req.RemoteAddr, host, req.URL.Path)
		key_path := path.Join(domain.Prefix, req.URL.Path)
		resp, err := domain.Bucket.HeadKey(domain.Identity, key_path)
		if err != nil {
			resp, err = domain.Bucket.HeadKey(domain.Identity, key_path)
		}
		// Error or no, we copy in the headers if we got a valid response 
		if resp != nil {
			for k, v := range resp.Header {
				rw.Header()[k] = v
			}
			// We're only copying the body in for 200's.
			if resp.ContentLength > 0 && err == nil {
				rw.Header().Set("Content-Length", strconv.Itoa64(resp.ContentLength))
			}
		}
		outcode := http.StatusInternalServerError
		if resp != nil {
			outcode = resp.StatusCode
		}
		rw.WriteHeader(outcode)
		if err == nil {
			_, err = domain.Bucket.GetKey(domain.Identity, key_path, rw)
		}
		Printf("%d: %s %s %s %d %v", outcode, req.RemoteAddr, host, req.URL.Path, resp.ContentLength, err)
		return
	} else {
		Printf("%s %s %s - Host Unknown", req.RemoteAddr, host, req.URL.Path)
		http.Error(rw, "Invalid host", http.StatusForbidden)
	}
}

var flag_bind_addr *string = flag.String("listen", "127.0.0.1:8080", "Address/port to listen to")

func main() {
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
	err = http.ListenAndServe(*flag_bind_addr, Service{conf: conf})
	if err != nil {
		Fatalf("Couldn't setup listener: %v", err)
	}
	return
}
