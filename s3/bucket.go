package s3

import "com.abneptis.oss/aws/awsconn"
import "com.abneptis.oss/aws/auth"
//import "com.abneptis.oss/aws"

import "os"
import "http"

type Bucket struct {
  Endpoint *awsconn.Endpoint
  Name     string
}

func NewBucket(ep *awsconn.Endpoint, name string)(*Bucket){
  return &Bucket{Endpoint: ep, Name: name}
}

// Sign and send a low-level qury-request;  this includes formulating the
// appropriate http request, creating a connection, sending the data, and
// returning the ultimate http response.
func (self *Bucket)send(id auth.Signer, qr *QueryRequest)(resp *http.Response, err os.Error){
  qr.Sign(id)
  if err != nil {return}
  hreq := qr.HTTPRequest()
  cconn, err := self.Endpoint.NewHTTPClientConn("tcp", "", nil)
  if err != nil { return }
  defer cconn.Close()
  resp, err = awsconn.SendRequest(cconn, hreq)
  return
}

func (self *Bucket)Create(id auth.Signer)(err os.Error){
  qr, err := NewQueryRequest("PUT", self.Endpoint, self.Name, "", "", "", "", nil,
           map[string]string{"AWSAccessKeyId": auth.GetSignerIDString(id)}, 15)
  if err != nil { return }
  resp, err := self.send(id, qr)
  if resp.StatusCode != 200 {
    s3err := &S3Error{}
    err = awsconn.ParseResponse(resp, s3err)
    if err == nil { err = s3err }
    return
  }

  return
}

func (self *Bucket)Destroy(id auth.Signer)(err os.Error){
  qr, err := NewQueryRequest("DELETE", self.Endpoint, self.Name, "", "", "", "", nil,
           map[string]string{"AWSAccessKeyId": auth.GetSignerIDString(id)}, 15)
  if err != nil { return }
  resp, err := self.send(id, qr)
  if resp.StatusCode != 204 {
    s3err := &S3Error{}
    err = awsconn.ParseResponse(resp, s3err)
    if err == nil { err = s3err }
    return
  }
  return
}

func (self *Bucket)GetKey(id auth.Signer, key string)(obj *Object, err os.Error){
  qr, err := NewQueryRequest("GET", self.Endpoint, self.Name, key, "", "", "", nil,
           map[string]string{"AWSAccessKeyId": auth.GetSignerIDString(id)}, 15)
  if err != nil { return }
  resp, err := self.send(id, qr)
  if err != nil { return }
  if resp.StatusCode == 404 {
    err = ErrorKeyNotFound
  }
  return
}
