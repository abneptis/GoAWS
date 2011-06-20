package s3

import (
  "aws"
  "http"
  "io"
  "os"
  "xml"
)

func PutLocalFile(id *aws.Signer, c *aws.Conn, bucket, key, file string)(err os.Error){
  fp, err := os.Open(file)
  if err == nil {
    defer fp.Close()
    err = PutFile(id,c,bucket,key,fp)
  }
  return
}

func PutFile(id *aws.Signer, c *aws.Conn, bucket, key string, fp *os.File)(err os.Error){
  if key == "" || key == "/" || bucket == "" {
    return os.NewError("Bucket and key are required (and cannot be '/')")
  }
  var resp *http.Response
  if fp == nil {
    return os.NewError("invalid file descriptor")
  }
  fi, err := fp.Stat()
  if err == nil {
    fsize := fi.Size
    hdr := http.Header{}
    hreq := newRequest("PUT", bucket, key, hdr, nil)
    hreq.ContentLength = fsize
    hreq.Body = fp
    err = signRequest(id, hreq)
    if err == nil {
      resp, err = c.Request(hreq)
      if err == nil { err = CodeToError(resp.StatusCode) }
    }
  }
  return
}

func DeleteKey(id *aws.Signer, c *aws.Conn, bucket, key string)(err os.Error){
  var resp *http.Response
  if bucket == "" || key == "" {
    err = os.NewError("Bucket/Key both required")
    return
  }
  hreq := newRequest("DELETE", bucket, key, nil, nil)
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

func GetKey(id *aws.Signer, c *aws.Conn, bucket, key string, w io.Writer)(http.Header, err os.Error){
  var resp *http.Response
  hreq := newRequest("GET", bucket, key, nil, nil)
  err = signRequest(id, hreq)
  if err == nil {
    resp, err = c.Request(hreq)
  }
  if err == nil {
    err = CodeToError(resp.StatusCode)
  } 

  if resp != nil {
    _, err2 := io.Copy(w, resp.Body)
    if err == nil { err = err2 }
  }
  return
}

func KeyExists(id *aws.Signer, c *aws.Conn, bucket, key string)(err os.Error){
  var resp *http.Response
  hreq := newRequest("HEAD", bucket, key, nil, nil)
  err = signRequest(id, hreq)
  if err == nil {
    resp, err = c.Request(hreq)
  }
  if err == nil {
    err = CodeToError(resp.StatusCode)
  } 
  return
}

type owner struct {
  ID          string
  DisplayName string
}

func ListKeys(id *aws.Signer, c *aws.Conn, bucket string,
              prefix, delim, marker string, out chan<- string)(err os.Error){
  var done bool
  var resp *http.Response
  result := listBucketResult{}
  result.Prefix = prefix
  result.Marker = marker
  for err == nil && ! done {
    form := http.Values{}
    form.Set("prefix", result.Prefix)
    form.Set("marker", result.Marker)
    form.Set("delimeter", delim)
    hreq := newRequest("GET", bucket,"/",nil,form)
    err = signRequest(id, hreq)
    if err == nil {
      resp, err = c.Request(hreq)
    }
    if err == nil {
      err = CodeToError(resp.StatusCode)
    } 
    if err == nil {
      err = xml.Unmarshal(resp.Body, &result)
      if err == nil {
        for i := range(result.Contents) {
          out <- result.Contents[i].Key
        }
        done = ! result.IsTruncated
      }
    }
  }
  return
}

type listBucketResult struct {
  Name string
  Prefix string
  Marker string
  MaxKeys int
  IsTruncated bool
  Contents []keyItem
}

type keyItem struct {
  Key string
  LastModified string // TODO: date
  ETag string
  Size int64
  Owner owner
  StorageClass string
}

/* Example output (2011-06-20)
<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
 <Name>apt.abneptis.com</Name>
 <Prefix></Prefix>
 <Marker></Marker>
 <MaxKeys>1000</MaxKeys>
 <IsTruncated>false</IsTruncated>
 <Contents>
  <Key>debian//dists/deli/Release</Key>
  <LastModified>2011-06-19T00:45:37.000Z</LastModified>
  <ETag>&quot;bf185ac183a3cd3ad4c9332beb45fc70&quot;</ETag>
  <Size>14007</Size>
  <Owner>
   <ID>19f48e038756359c402c774f40ea9b193668d906b8836c783823b9fd33b270ef</ID>
   <DisplayName>amazon</DisplayName>
  </Owner>
  <StorageClass>STANDARD</StorageClass>
 </Contents>
*/
