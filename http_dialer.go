package aws

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// The conn structure represents a 'semi detached' http-client
// It handles redialing & reconnecting on connection errors.
type Conn struct {
	uc *ReusableConn
	c  *httputil.ClientConn
}

// Creates a new connection with the specified dialer function.
func NewConn(d Dialer) *Conn {
	return &Conn{
		uc: NewReusableConnection(d),
		c:  nil,
	}
}

func (self *Conn) dial() (err error) {
	if self.c == nil {
		// Get the underlying connection (or redial)
		err = self.uc.Dial()
		if err == nil {
			self.c = httputil.NewClientConn(self.uc, nil)
		}
	}
	return
}

// Closes the underlying connection
//
// NB: if you re-use the connection after
// this, it will be redialed.
func (self *Conn) Close() (err error) {
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
func (self *Conn) Request(req *http.Request) (resp *http.Response, err error) {
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
		if err == httputil.ErrPersistEOF {
			err = nil
		}
		self.Close()
	}
	return
}

// A generic Dialer to handle both TLS and non TLS http connections.
func URLDialer(u *url.URL, conf *tls.Config) (f func() (c net.Conn, err error)) {
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

	f = func() (c net.Conn, err error) {
		if useTLS {
			return tls.Dial("tcp", host+":"+port, conf)
		}
		return net.Dial("tcp", host+":"+port)
	}
	return
}

// Constructs a basic http.Request based off of a fully-qualified URL
func NewRequest(url_ *url.URL, method string, hdrs http.Header, params url.Values) (req *http.Request) {
	req = &http.Request{
		Method: method,
		URL: &url.URL{
			Path:     url_.Path,
			RawQuery: url_.RawQuery,
		},
		Host:   url_.Host,
		Header: hdrs,
		Form:   params,
	}
	if req.URL.RawQuery != "" {
		req.URL.RawQuery += "&"
	}
	req.URL.RawQuery += params.Encode()
	return
}
