package s3

import "com.abneptis.oss/aws/awsconn"
import "com.abneptis.oss/aws/auth"
//import "com.abneptis.oss/aws"

import "os"
import "http"

//<ListAllMyBucketsResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/"><Owner><ID>19f48e038756359c402c774f40ea9b193668d906b8836c783823b9fd33b270ef</ID><DisplayName>amazon</DisplayName></Owner><Buckets><Bucket><Name>records.abneptis.com</Name><CreationDate>2010-10-31T18:51:43.000Z</CreationDate></Bucket><Bucket><Name>records.abneptis.net</Name><CreationDate>2010-10-31T18:51:56.000Z</CreationDate></Bucket></Buckets></ListAllMyBucketsResult>

type listBucketsResult struct {
  Owner bucketOwner
  Buckets []bucketResults
}

type bucketOwner struct {
  ID string
  DisplayName string
}

type bucketResults struct {
  Bucket bucketRecord
}

type bucketRecord struct {
  CreationDate string
  Name string
}

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
  resp, err := awsconn.SendRequest(cconn, hreq)

  bb, _ = http.DumpResponse(resp, true)
  os.Stdout.Write(bb)
  return
}
