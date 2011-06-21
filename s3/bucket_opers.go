package s3

import (
  "aws"
  "http"
  "io"
  "os"
  "path"
  "xml"
)


// Represents a URL and connection to an S3 bucket.
type Bucket struct {
  Name string
  URL  *http.URL
  conn *aws.Conn
}


// URL is the _endpoint_ url (nil will default to https://s3.amazonaws.com/)
// Name should be the bucket name; If you pass this as an empty string, you will regret it.
// conn is OPTIONAL, but allows you to re-use another aws.Conn if you'd like.
func NewBucket(u *http.URL, Name string, conn *aws.Conn)(b *Bucket){
  if u == nil {
    u = &http.URL{Scheme:"https", Host: USEAST_HOST, Path: "/"}
  }
  if conn == nil {
    conn = aws.NewConn(aws.URLDialer(u, nil))
  }
  b = &Bucket {
    URL: u,
    Name: Name,
    conn: conn,
  }
  return
}

func (self *Bucket)key_url(key string)(*http.URL){
  return &http.URL {
    Scheme: self.URL.Scheme,
    Host: self.URL.Host,
    Path: path.Join(self.URL.Path, self.Name, key),
  }
}


// Will open a local file, size it, and upload it to the named key.
// This is a convenience wrapper aroudn PutFile.
func (self *Bucket)PutLocalFile(id *aws.Signer, key, file string)(err os.Error){
  fp, err := os.Open(file)
  if err == nil {
    defer fp.Close()
    err = self.PutFile(id, key,fp)
  }
  return
}

// Will put an open file descriptor to the named key.  Size is determined
// by statting the fd (so a partially read file will not work).
// TODO: ACL's & content-type/headers support
func (self *Bucket)PutFile(id *aws.Signer, key string, fp *os.File)(err os.Error){
  var resp *http.Response
  if fp == nil {
    return os.NewError("invalid file descriptor")
  }
  fi, err := fp.Stat()
  if err == nil {
    fsize := fi.Size
    hdr := http.Header{}
    hreq := aws.NewRequest(self.key_url(key), "PUT", hdr, nil)
    hreq.ContentLength = fsize
    hreq.Body = fp
    err   = id.SignRequestV1(hreq, aws.CanonicalizeS3, 15)
    if err == nil {
      resp, err = self.conn.Request(hreq)
      if err == nil { err = aws.CodeToError(resp.StatusCode) }
    }
  }
  return
}

// Deletes the named key from the bucket.  To delete a bucket, see *Service.DeleteBucket()
func (self *Bucket)Delete(id *aws.Signer, key string)(err os.Error){
  var resp *http.Response
  if key == "" {
    return os.NewError("Key cannot be empty!")
  }
  hreq := aws.NewRequest(self.key_url(key), "DELETE", nil, nil)
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

// Opens the named key and copys it to the named io.Writer.
// Also returns the http headers for convenience.
func (self *Bucket)GetKey(id *aws.Signer, key string, w io.Writer)(hdr http.Header, err os.Error){
  var resp *http.Response
  hreq := aws.NewRequest(self.key_url(key), "GET", nil, nil)
  err   = id.SignRequestV1(hreq, aws.CanonicalizeS3, 15)
  if err == nil {
    resp, err = self.conn.Request(hreq)
  }
  if err == nil {
    err = aws.CodeToError(resp.StatusCode)
    hdr = resp.Header
  } 

  if resp != nil {
    _, err2 := io.Copy(w, resp.Body)
    if err == nil { err = err2 }
  }
  return
}

// Performs a HEAD request on the bucket and returns nil of the key appears
// valid (returns 200).
func (self *Bucket)Exists(id *aws.Signer, key string)(err os.Error){
  var resp *http.Response
  hreq := aws.NewRequest(self.key_url(key), "HEAD", nil, nil)
  err   = id.SignRequestV1(hreq, aws.CanonicalizeS3, 15)
  if err == nil {
    resp, err = self.conn.Request(hreq)
  }
  if err == nil {
    err = aws.CodeToError(resp.StatusCode)
  } 
  return
}

type owner struct {
  ID          string
  DisplayName string
}

// Walks a bucket and writes the resultign strings to the channel.
// * There is currently NO (correct) way to abort a running walk.
func (self *Bucket)ListKeys(id *aws.Signer,
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

    hreq := aws.NewRequest(self.key_url("/"), "GET", nil, form)
    err   = id.SignRequestV1(hreq, aws.CanonicalizeS3, 15)

    if err == nil {
      resp, err = self.conn.Request(hreq)
    }
    if err == nil {
      err = aws.CodeToError(resp.StatusCode)
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
