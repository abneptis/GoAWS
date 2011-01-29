package s3

import "com.abneptis.oss/aws/awsconn"
import "com.abneptis.oss/aws/auth"
import "com.abneptis.oss/cryptools/hextools"
//import "com.abneptis.oss/aws"

import "bytes"
import "crypto/md5"
import "io"
import "json"
import "os"
import "strconv"

type Bucket struct {
  Endpoint *awsconn.Endpoint
  Name     string
}

func NewBucket(ep *awsconn.Endpoint, name string)(*Bucket){
  return &Bucket{Endpoint: ep, Name: name}
}

// Create a new bucket in S3.  Note that namespaces for S3 (unlike SQS)
// are global in nature, so you may not conflict with another users bucket-name.
//
// Frequent good choices are dns names (forward or backwards: com.abneptis/foo or
// abneptis.com/foo should be equally unique) or GUIDs.
type createBucketResponse struct {
  CreateBucketResponse bucketList
}

type bucketList struct {
  Bucket []string
}

func (self *Bucket)Create(id auth.Signer)(err os.Error){
  hreq, err := NewQueryRequest(id, self.Endpoint, "PUT", self.Name,"","","", nil, nil)
  if err != nil { return }
  resp, err := self.Endpoint.SendRequest(hreq)
  if err != nil { return }
  if resp.StatusCode != 200 {
    err = os.NewError("Unable to create: " + resp.Status)
  }
  return
}


// Destroys an S3 bucket.  It is NOT an error to delete a bucket with
// contents.
func (self *Bucket)Destroy(id auth.Signer)(err os.Error){
  hreq, err := NewQueryRequest(id, self.Endpoint, "DELETE", self.Name,"","","", nil, nil)
  if err != nil { return }
  resp, err := self.Endpoint.SendRequest(hreq)
  if err != nil { return }
  if resp.StatusCode != 204 {
    err = os.NewError(resp.Status)
  }
  return
}

// Get an s3.Object with a ReadCloser for the body.
func (self *Bucket)GetKey(id auth.Signer, key string)(obj *Object, err os.Error){
  hreq, err := NewQueryRequest(id, self.Endpoint, "GET", self.Name,key,"","", nil, nil)
  if err != nil { return }
  cc, err := self.Endpoint.NewHTTPClientConn("tcp","", nil)
  if err != nil { return }
  defer cc.Close()
  resp, err := awsconn.SendRequest(cc, hreq)
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

// Delete an S3 Key.
func (self *Bucket)DeleteKey(id auth.Signer, key string)(err os.Error){
  hreq, err := NewQueryRequest(id, self.Endpoint, "DELETE", self.Name,key,"","", nil, nil)
  if err != nil { return }
  resp, err := self.Endpoint.SendRequest(hreq)
  if err != nil { return }
  switch resp.StatusCode {
    case 403:
      err = ErrorAccessDenied
    case 404:
      err = ErrorKeyNotFound
    case 204:
    default:
      err = os.NewError("Unhandled response code: " + resp.Status )
  }
  return
}

// Write S3 Key.
func (self *Bucket)PutKey(id auth.Signer, key, ctype, cmd5 string, llen int64, rc io.ReadCloser)(err os.Error){
  hreq, err := NewQueryRequest(id, self.Endpoint, "PUT", self.Name,key,ctype,cmd5, nil, nil)
  hreq.ContentLength = llen
  hreq.Body = rc
  if err != nil { return }
  resp, err := self.Endpoint.SendRequest(hreq)
  if err != nil { return }
  switch resp.StatusCode {
    case 403:
      err = ErrorAccessDenied
    case 200:
    default:
      err = os.NewError("Unhandled response code: " + resp.Status )
  }
  return
}


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
  hreq, err := NewQueryRequest(id, self.Endpoint, "GET", self.Name,"","","", nil, nil)
  if err != nil { return }
  if delim != "" {
    hreq.Form["delimiter"] = []string{delim}
  }
  if marker != "" {
    hreq.Form["marker"] = []string{marker}
  }
  if prefix != "" {
    hreq.Form["prefix"] = []string{prefix}
  }
  hreq.Form["max-keys"] =  []string{strconv.Itoa(max)}
  etype := &errorResponse{}
  obj := &listBucketResult{}
  err = self.Endpoint.SendParsable(hreq, obj, etype)
  if err != nil { return }
  out = make([]*Object, len(obj.Contents))
  for i := range(obj.Contents){
    out[i] = &Object{Key: obj.Contents[i].Key}
  }
  return
}

// Returns true if a key exists in the bucket.
// Note(1): This is not a "cheap" operation.  Cheaper than get-object, yes, 
// but it still does a full remote operation to test.
// Note(2): Due to amazons eventual consistancy, you should not rely solely
// on the result of this function as deterministic that a key has not been
// set.  It is quite possible that a key will be set and an immediate check
// of the same will fail, so retry logic should be built into the caller 
// (stack)
func (self *Bucket)Exists(id auth.Signer, key string)(exists bool, err os.Error){
  hreq, err := NewQueryRequest(id, self.Endpoint, "HEAD", self.Name,key,"","", nil, nil)
  if err != nil { return }
  cc, err := self.Endpoint.NewHTTPClientConn("tcp","", nil)
  if err != nil { return }
  defer cc.Close()
  resp, err := awsconn.SendRequest(cc, hreq)
  if err != nil { return }
  switch resp.StatusCode {
    case 404:
      exists = false
    case 200:
      exists = true
    default:
      err = os.NewError("Unexpected response: " + resp.Status )
  }
  return
}

type nopCloser struct { io.Reader }
func (nopCloser)Close()(n os.Error){return}

func (self *Bucket)PutJSON(id auth.Signer, key, ctype string, obj interface{})(err os.Error){
  rawb, err := json.Marshal(obj)
  if err != nil { return }
  mdh := md5.New()
  _, err = mdh.Write(rawb)
  if err != nil { return }
  sum := mdh.Sum()
  hout, err := hextools.HexString(false, true, string(sum))
  if err != nil { return }

  err = self.PutKey(id, key, ctype, hout, int64(len(rawb)), nopCloser{bytes.NewBuffer(rawb)})
  return
}
