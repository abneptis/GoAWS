package s3

import (
  "aws"
)

import (
  "http"
  "os"
  "xml"
)

func DeleteBucket(id *aws.Signer, c *aws.Conn, name string)(err os.Error){
  var resp *http.Response
  hreq := newRequest("DELETE", name, "", nil, nil) 
  err = signRequest(id, hreq)
   
  if err == nil {
    resp, err = c.Request(hreq)
  }
  if err == nil {
    if resp.StatusCode != http.StatusNoContent {
      err = CodeToError(resp.StatusCode)
    }
  }
  return
}

func CreateBucket(id *aws.Signer, c *aws.Conn, name string)(err os.Error){
  var resp *http.Response
  hreq := newRequest("PUT", name, "", nil, nil) 
  err = signRequest(id, hreq)
   
  if err == nil {
    resp, err = c.Request(hreq)
  }
  if err == nil {
    err = CodeToError(resp.StatusCode)
  }
  return
}


// Returns a list of bucket names known by the endpoint.  Depending on the
// endpoint used, your list may be global or regional in nature.
func ListBuckets(id *aws.Signer, c *aws.Conn)(out []string, err os.Error){
  var resp *http.Response
  hreq := newRequest("GET", "","",nil,nil)
  err = signRequest(id, hreq)
  if err == nil {
    resp, err = c.Request(hreq)
  }
  if err == nil {
    err = CodeToError(resp.StatusCode)
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
