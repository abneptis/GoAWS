// Interface types and functions for Amazon's Simple Storage Service
//
// S3 uses an Bucket/Key tuple to store data in a user-defined format.
//
// The interface exposed will use readers/writers, and will likely need
// to be wrapped in user serializers/deserializers.
package s3

import "com.abneptis.oss/aws/awsconn"
import "com.abneptis.oss/aws/auth"
//import "com.abneptis.oss/aws"

import "os"

//<ListAllMyBucketsResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
// <Owner><ID>19f48e038756359c402c774f40ea9b193668d906b8836c783823b9fd33b270ef</ID><DisplayName>amazon</DisplayName></Owner>
// <Buckets>
//  <Bucket><Name>records.abneptis.com</Name><CreationDate>2010-10-31T18:51:43.000Z</CreationDate></Bucket>
//  <Bucket><Name>records.abneptis.net</Name><CreationDate>2010-10-31T18:51:56.000Z</CreationDate></Bucket>
// </Buckets>
//</ListAllMyBucketsResult>

type listBucketsResult struct {
  Owner bucketOwner
  Buckets bucketResults
}

type bucketOwner struct {
  ID string
  DisplayName string
}

type bucketResults struct {
  Bucket []bucketRecord
}

type bucketRecord struct {
  CreationDate string
  Name string
}

// Returns a list of bucket names known by the endpoint.  Depending on the 
// endpoint used, your list may be global or regional in nature.
func ListBuckets(id auth.Signer, ep *awsconn.Endpoint)(out []string, err os.Error){
  qr, err := NewQueryRequest("GET", ep, "", "", "", "", "", nil,
           map[string]string{"AWSAccessKeyId": auth.GetSignerIDString(id)}, 15)
  if err != nil { return }
  resp, err := qr.Send(id, ep)
  if err != nil { return }
  if resp.StatusCode == 404 {
    err = ErrorKeyNotFound
  }
  result := &listBucketsResult{}
  err = awsconn.ParseResponse(resp, result)
  if err != nil { return }

  out = make([]string, len(result.Buckets.Bucket))
  for i := range(result.Buckets.Bucket){
    out[i] = result.Buckets.Bucket[i].Name
  }
  return
}
