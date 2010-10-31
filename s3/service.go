package s3

import "com.abneptis.oss/aws/awsconn"
import "com.abneptis.oss/aws/auth"
import "com.abneptis.oss/aws"

import "os"
import "http"

func ListBuckets(id auth.Signer, ep *awsconn.Endpoint)(out []string, err os.Error){
  req := NewRequest("GET", "", "", "", ep, nil, 15)
  req.Set("AWSAccessKeyId", auth.GetSignerIDString(id))
  err = SignS3Request(id, ep.GetURL(), req)
  if err != nil {return}
  hreq, err := req.HTTPRequest(id, ep, "","")
  if err != nil { return }

  bb, _ := http.DumpRequest(hreq, true)
  os.Stdout.Write(bb)

  cconn, err := ep.NewHTTPClientConn("tcp", "", nil)
  if err != nil { return }
  defer cconn.Close()
  resp, err := aws.SendRequest(cconn, hreq)

  bb, _ = http.DumpResponse(resp, true)
  os.Stdout.Write(bb)
  return
}
