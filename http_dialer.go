package aws

import (
	"crypto/tls"
	"http"
	"net"
	"os"
)

// The conn structure represents a 'semi detached' http-client
// It handles redialing & reconnecting on connection errors.
type Conn struct {
	uc *ReusableConn
	c  *http.ClientConn
}

// Creates a new connection with the specified dialer function.
func NewConn(d Dialer) *Conn {
	return &Conn{
		uc: NewReusableConnection(d),
		c:  nil,
	}
}

func (self *Conn) dial() (err os.Error) {
	if self.c == nil {
		// Get the underlying connection (or redial)
		err = self.uc.Dial()
		if err == nil {
			self.c = http.NewClientConn(self.uc, nil)
		}
	}
	return
}


// Closes the underlying connection
//
// NB: if you re-use the connection after
// this, it will be redialed.
func (self *Conn) Close() (err os.Error) {
	if self.c != nil {
		self.c.Close()
		self.c = nil
	}
	if self.uc != nil {
		err = self.uc.Close()
	}
	return
}

// Write a request and read the response;
// This function will also fix-up req.Form for 'GET's
func (self *Conn) Request(req *http.Request) (resp *http.Response, err os.Error) {
	err = self.dial()
	if err == nil {
		if req.Form != nil && req.Method == "GET" {
			if req.URL.RawQuery != "" {
				req.URL.RawQuery += "&"
			}
			req.URL.RawQuery += req.Form.Encode()
			req.Form = nil
		}
		err = self.c.Write(req)
		if err == nil {
			resp, err = self.c.Read(req)
		}
	}
	if err != nil {
		if err == http.ErrPersistEOF {
			err = nil
		}
		self.Close()
	}
	return
}


// A generic Dialer to handle both TLS and non TLS http connections.
func URLDialer(u *http.URL, conf *tls.Config) (f func() (c net.Conn, err os.Error)) {
	host, port, _ := net.SplitHostPort(u.Host)
	if port == "" {
		if u.Scheme == "http" {
			port = "80"
		}
		if u.Scheme == "https" {
			port = "443"
		}
	}
	if host == "" {
		host = u.Host
	}
	useTLS := (u.Scheme == "https")

	f = func() (c net.Conn, err os.Error) {
		if useTLS {
			return tls.Dial("tcp", host+":"+port, conf)
		}
		return net.Dial("tcp", host+":"+port)
	}
	return
}

// Constructs a basic http.Request based off of a fully-qualified URL
func NewRequest(url *http.URL, method string, hdrs http.Header, params http.Values) (req *http.Request ){
	req = &http.Request{
		Method: method,
		URL: &http.URL{
			Path:     url.Path,
			RawQuery: url.RawQuery,
		},
		Host:   url.Host,
		Header: hdrs,
		Form:   params,
	}
  if req.URL.RawQuery != "" { req.URL.RawQuery += "&" }
  req.URL.RawQuery += params.Encode()
  return
}
