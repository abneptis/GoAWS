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

func (self *Bucket)Create(id auth.Signer)(err os.Error){
  req := NewRequest("PUT", self.Name, "", "", self.Endpoint, nil, 15)
  req.Set("AWSAccessKeyId", auth.GetSignerIDString(id))
  err = SignS3Request(id, self.Endpoint.GetURL(), req)
  if err != nil {return}
  hreq, err := req.HTTPRequest(id, self.Endpoint, self.Name,"")
  if err != nil { return }

  bb, _ := http.DumpRequest(hreq, true)
  os.Stdout.Write(bb)

  cconn, err := self.Endpoint.NewHTTPClientConn("tcp", "", nil)
  if err != nil { return }
  defer cconn.Close()
  resp, err := awsconn.SendRequest(cconn, hreq)

  bb, _ = http.DumpResponse(resp, true)
  os.Stdout.Write(bb)
  return
}
