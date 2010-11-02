// Handles basic connection primatives.
//
// Unlike a simple connection, awsconn carries a second URL
// used for proxy-server data.  It is not a connection, but
// but a set of helpers to the "net" class of functionality.
//
// It is the encouraged way for goaws
// utlities to establish and maintain connection details so
// that proxy configuration data is available to all 
// callers.
package awsconn
/* 
  Copyright (c) 2010, Abneptis LLC.
  See COPYRIGHT and LICENSE for details.
*/

import "com.abneptis.oss/urltools"
import "com.abneptis.oss/maptools"

import "http"
import "net"
import "os"
import "bufio"
import "xml"

// Sends an AWS request. 
// A simple helper exposed as it is often useful.
//
// If req.Form is not empty, but req.URL.RawQuery is,
// req.RawQuery will be filled with the values from
// req.Form
func SendRequest(cc *http.ClientConn, req *http.Request)(resp *http.Response, err os.Error){
  if req.URL.RawQuery == "" && len(req.Form) > 0 {
    req.URL.RawQuery = maptools.StringStringJoin(
                        maptools.StringStringEscape(
                         maptools.StringStringsJoin(req.Form, ",", false),
                         http.URLEscape, http.URLEscape),
                         "=", "&", true)
  }
  //bb, _ := http.DumpRequest(req, true)
  //os.Stderr.Write(bb)
  err = cc.Write(req)
  if err == nil {
    resp, err = cc.Read()
  }
  //bb, _ = http.DumpResponse(resp, true)
  //os.Stderr.Write(bb)
  return
}

// Handles xml unmarshalling to a generic object from
// an http.Response
func ParseResponse(resp *http.Response, o interface{})(err os.Error){
  if resp.Body == nil {
    err = os.NewError("Response body is empty")
  } else {
    parser := xml.NewParser(resp.Body)
    err = parser.Unmarshal(o, nil)
  }
  return
}

// Construct a new Endpoint. 
func NewEndpoint(u, pu *http.URL)(*Endpoint){
  return &Endpoint {
   URL: u,
   ProxyURL: pu,
  }
}

// An endpoint has two URL's, the "URL", or rather, the actual service
// address, and a "ProxyURL", which is used for low-level connections if
// it is present -- else the URL is connected to directly.
type Endpoint struct {
  URL *http.URL
  ProxyURL *http.URL
}

// Return the URL to be used for connection purposes.
// While not expected to be needed by external users,
// this is considered the "correct" way to make that
// determination, and shorter to import and use than write.
func (self *Endpoint)ConnectionURL()(out *http.URL){
  if self.ProxyURL != nil {
    out = self.ProxyURL
  } else {
    out = self.URL
  }
  return
}

// Return the URL to be used for request generation.
// This is far more likely to be useful to end users
// than the above.
func (self *Endpoint)GetURL()(out *http.URL){
  return self.URL
}


// Return a new net.Conn using netname and local as net.Dial does.
// NewConn does not explicitly check to ensure that you are using
// a stream protocol, so if you accept this from a user source, it
// is the callers responsibility to verify.
func (self *Endpoint)NewConn(netname, local string)(rawc net.Conn, err os.Error){
  hps, err := urltools.ExtractURLHostPort(self.ConnectionURL())
  if err == nil {
    rawc, err = net.Dial(netname, local, hps)
  }
  return
}

// Returns a new HTTP connection;  As with NewConn, netname is not
// checked, and the behaviour of an HTTP client over a non stream
// protocol is undefined - but probably interesting to watch.
func (self *Endpoint)NewHTTPClientConn(netname, local string, r *bufio.Reader)(hc *http.ClientConn, err os.Error){
  rawc, err := self.NewConn(netname, local)
  if err == nil {
    hc = http.NewClientConn(rawc, r)
  }
  return
}

// Creates a connection and sends an httpRequest, returning the result and closing the stream.
func (self *Endpoint)SendRequest(req *http.Request)(resp *http.Response, err os.Error){
  cc,err := self.NewHTTPClientConn("tcp", "", nil)
  if err != nil { return }
  defer cc.Close()
  return SendRequest(cc, req)
}

// Attempts to demarshal an XML response into the given response/error types.
func (self *Endpoint)SendParsable(req *http.Request, out interface{}, etype os.Error)(err os.Error){
  resp, err := self.SendRequest(req)
  if err != nil { return }
  if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
   err = ParseResponse(resp, out)
  } else {
   oserr := ParseResponse(resp, etype)
   if oserr != nil {
     err = oserr
   } else {
     err = etype
   }
  }
  return
}


// Create an HTTP request with appropriate details pulled from the local 
// endpoint and other details.
func (self *Endpoint)NewHTTPRequest(method string, path string, params map[string][]string, headers map[string]string)(req *http.Request){
  req = &http.Request {
    Method: method,
    Form: params,
    Header: headers,
    Host: self.URL.Host,
    URL: &http.URL{
      Host: self.URL.Host,
      Scheme: self.URL.Scheme,
      Path: path,
    },
  }
  if req.Header == nil { req.Header = make(map[string]string) }
  if req.Form   == nil { req.Form   = make(map[string][]string) }
  return
}

// Returns the value of the ContentType header (or the empty string).
// Primarily a helper for canonicalization routines.
func ContentType(req *http.Request)(string){
  return req.Header["Content-Type"]
}

// Returns the value of the ContentMd5 header (or the empty string).
// Primarily a helper for canonicalization routines.
func ContentMD5(req *http.Request)(string){
  return req.Header["Content-Md5"]
}

