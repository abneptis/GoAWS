package s3

import (
  "aws"
)

import (
  "http"
  "os"
  "path"
  "xml"
)

type Service struct {
  URL  *http.URL
  conn *aws.Conn
}

func NewService(url *http.URL)(s *Service){
  s = &Service {
    URL: url,
  }
  if s.URL == nil { s.URL, _  = http.ParseURL(USEAST_HOST) }
  s.conn = aws.NewConn(aws.URLDialer(s.URL, nil))
  return
}

func (self *Service)Bucket(name string)(*Bucket){
  return NewBucket(self.URL, name, self.conn)
}

func s3Path(bucket, key string)(string){
  return path.Join("/", bucket, key)
}

func (self *Service)bucket_url(bucket string)(*http.URL){
  return &http.URL {
    Host: self.URL.Host,
    Path: path.Join(self.URL.Path, s3Path(bucket,"")),
    Scheme: self.URL.Scheme,
  }
}

func (self *Service)DeleteBucket(id *aws.Signer, name string)(err os.Error){
  var resp *http.Response
  hreq := aws.NewRequest(self.bucket_url(name), "DELETE", nil, nil)
  err   = id.SignRequestV1(hreq, aws.CanonicalizeS3, 15)
   
  if err == nil {
    resp, err = self.conn.Request(hreq)
  }

  if err == nil {
    if resp.StatusCode != http.StatusNoContent {
      err = aws.CodeToError(resp.StatusCode)
    }
  }

  return
}

func (self *Service)CreateBucket(id *aws.Signer, name string)(err os.Error){
  var resp *http.Response
  hreq := aws.NewRequest(self.bucket_url(name), "PUT", nil, nil)
  err   = id.SignRequestV1(hreq, aws.CanonicalizeS3, 15)
   
  if err == nil {
    resp, err = self.conn.Request(hreq)
  }
  if err == nil {
    err = aws.CodeToError(resp.StatusCode)
  }
  return
}


// Returns a list of bucket names known by the endpoint.  Depending on the
// endpoint used, your list may be global or regional in nature.
func (self *Service)ListBuckets(id *aws.Signer)(out []string, err os.Error){
  var resp *http.Response
  hreq := aws.NewRequest(self.bucket_url(""), "GET", nil, nil)
  err   = id.SignRequestV1(hreq, aws.CanonicalizeS3, 15)
  if err == nil {
    resp, err = self.conn.Request(hreq)
  }
  if err == nil {
    err = aws.CodeToError(resp.StatusCode)
  }
  if err == nil {
    result := listAllMyBucketsResult{}
    err = xml.Unmarshal(resp.Body, &result)
    if err == nil {
      out = result.Buckets
    }
  }
  return
}

type listAllMyBucketsResult struct {
  Owner          owner
  Buckets        []string "Buckets>Bucket>Name"
}
