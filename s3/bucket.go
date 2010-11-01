package s3

import "com.abneptis.oss/aws/awsconn"
import "com.abneptis.oss/aws/auth"
//import "com.abneptis.oss/aws"

import "os"
import "http"
import "strconv"

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
  return qr.Send(id, self.Endpoint)
}

// Create a new bucket in S3.  Note that namespaces for S3 (unlike SQS)
// are global in nature, so you may not conflict with another users bucket-name.
//
// Frequent good choices are dns names (forward or backwards: com.abneptis/foo or
// abneptis.com/foo should be equally unique) or GUIDs.
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

// Destroys an S3 bucket.  It is NOT an error to delete a bucket with
// contents.
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

// Get an s3.Object with a ReadCloser for the body.
func (self *Bucket)GetKey(id auth.Signer, key string)(obj *Object, err os.Error){
  qr, err := NewQueryRequest("GET", self.Endpoint, self.Name, key, "", "", "", nil,
           map[string]string{"AWSAccessKeyId": auth.GetSignerIDString(id)}, 15)
  if err != nil { return }
  resp, err := self.send(id, qr)
  if err != nil { return }
  switch resp.StatusCode {
    case 403:
      err = ErrorAccessDenied
    case 404:
      err = ErrorKeyNotFound
    case 200:
      obj = &Object{Key: key, Body: resp.Body}
    default:
      err = os.NewError("Unhandled response code: " + resp.Status )
  }
  return
}
/*
<ListBucketResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
 <Name>records.abneptis.net</Name>
 <Prefix></Prefix>
 <Marker></Marker>
 <MaxKeys>1000</MaxKeys>
 <IsTruncated>false</IsTruncated>
 <Contents>
  <Key>accounts/9v_eecx0HuSb_-hkT0Qp_es0_rt_KSYEZXZsLi8ZF1uXkjSXGI_DnT1DM1_IcJG09FKQyiFHI-GjvZQ18RKSJA==/vP9j8c-NCFc0WLyMi12jDrIyhJW4jJXR2wAfhoHa6aN_i0N363D0euzBhB9CLfXlgeTk98Drx0gTk7JSvgQ8tQ==</Key>
  <LastModified>2010-10-17T23:51:53.000Z</LastModified>
  <ETag>&quot;c3d9e26b5b9ec3b7933ad622a716d25c&quot;</ETag>
  <Size>76</Size>
  <Owner><ID>19f48e038756359c402c774f40ea9b193668d906b8836c783823b9fd33b270ef</ID><DisplayName>amazon</DisplayName></Owner>
  <StorageClass>STANDARD</StorageClass>
 </Contents>
 <Contents>
  <Key>test-key</Key>
  <LastModified>2010-10-17T22:32:09.000Z</LastModified>
  <ETag>&quot;08ff08d3b2981eb6c611a385ffa4f865&quot;</ETag>
  <Size>11</Size>
  <Owner><ID>19f48e038756359c402c774f40ea9b193668d906b8836c783823b9fd33b270ef</ID><DisplayName>amazon</DisplayName></Owner>
  <StorageClass>STANDARD</StorageClass>
 </Contents>
</ListBucketResult>
*/

type listBucketResult struct {
  Name string
  Prefix string
  Marker string
  MaxKeys int
  IsTruncated bool
  Contents []bucketResult
}

type bucketResult struct {
  Key string
  LastModified string
  Size int
  StorageClass string
  Owner bucketOwner
}

// Returns a list of Object pointers with the Name field set.
//
// Users should be aware that there is no Body in the objects returned
// by ListKeys.
func (self *Bucket)ListKeys(id auth.Signer, delim, marker, prefix string, max int)(out []*Object, err os.Error){
  qr, err := NewQueryRequest("GET", self.Endpoint, self.Name, "", "", "", "", nil,
           map[string]string{"AWSAccessKeyId": auth.GetSignerIDString(id)}, 15)
  if err != nil { return }
  if delim != "" {
    err = qr.Parameters.Set("delimiter", delim)
    if err != nil { return }
  }
  if marker != "" {
    err = qr.Parameters.Set("marker", marker)
    if err != nil { return }
  }
  if prefix != "" {
    err = qr.Parameters.Set("prefix", marker)
    if err != nil { return }
  }
  qr.Parameters.Set("max-keys", strconv.Itoa(max))
  resp, err := self.send(id, qr)
  if err != nil { return }
  switch resp.StatusCode {
    case 404:
      err = ErrorKeyNotFound
    case 200:
      obj := &listBucketResult{}
      err = awsconn.ParseResponse(resp, obj)
      out = make([]*Object, len(obj.Contents))
      if err == nil {
        for i := range(obj.Contents){
          out[i] = &Object{Key: obj.Contents[i].Key}
        }
      }
    default:
      s3err := &S3Error{}
      err = awsconn.ParseResponse(resp, s3err)
      if err == nil { err = s3err }
  }
  return
}
