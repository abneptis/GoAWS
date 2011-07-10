package s3

import (
	"aws"
)

import (
	"bytes"
	"http"
	"io"
	"io/ioutil"
	"net"
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

// NewBucket creates a new *Bucket object from an endpoint URL.
//
// URL is the _endpoint_ url (nil will default to https://s3.amazonaws.com/)
// Name should be the bucket name; If you pass this as an empty string, you will regret it.
// conn is OPTIONAL, but allows you to re-use another aws.Conn if you'd like.
//
// If you omit conn, the dialer used will be based off the VHost of your bucket (if possible),
// to ensure best performance (e.g., endpoint associated w/ your bucket, and any
// regional lb's)
func NewBucket(u *http.URL, Name string, conn *aws.Conn) (b *Bucket) {
	if u == nil {
		u = &http.URL{Scheme: "https", Host: USEAST_HOST, Path: "/"}
	}
	if conn == nil {
		vname := VhostName(Name, u)
		addrs, err := net.LookupHost(vname)
		if err == nil && len(addrs) > 0 {
			dial_url := &http.URL{
				Scheme: u.Scheme,
				Host:   vname,
				Path:   u.Path,
			}
			conn = aws.NewConn(aws.URLDialer(dial_url, nil))
		} else {
			conn = aws.NewConn(aws.URLDialer(u, nil))
		}
	}
	b = &Bucket{
		URL:  u,
		Name: Name,
		conn: conn,
	}
	return
}

// Returns the vhost name of the bucket (bucket.s3.amazonaws.com)
func VhostName(b string, ep *http.URL) string {
	return b + "." + ep.Host
}

func (self *Bucket) key_url(key string) *http.URL {
	return &http.URL{
		Scheme: self.URL.Scheme,
		Host:   self.URL.Host,
		Path:   path.Join(self.URL.Path, self.Name, key),
	}
}


// Will open a local file, size it, and upload it to the named key.
// This is a convenience wrapper aroudn PutFile.
func (self *Bucket) PutLocalFile(id *aws.Signer, key, file string) (err os.Error) {
	fp, err := os.Open(file)
	if err == nil {
		defer fp.Close()
		err = self.PutFile(id, key, fp)
	}
	return
}

// Will put an open file descriptor to the named key.  Size is determined
// by statting the fd (so a partially read file will not work).
// TODO: ACL's & content-type/headers support
func (self *Bucket) PutFile(id *aws.Signer, key string, fp *os.File) (err os.Error) {
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
		err = id.SignRequestV1(hreq, aws.CanonicalizeS3, 15)
		if err == nil {
			resp, err = self.conn.Request(hreq)
			if err == nil {
				defer resp.Body.Close()
				err = aws.CodeToError(resp.StatusCode)
			}
		}
	}
	return
}


func (self *Bucket) PutKeyBytes(id *aws.Signer, key string, buff []byte, hdr http.Header) (err os.Error) {
	var resp *http.Response
	hreq := aws.NewRequest(self.key_url(key), "PUT", hdr, nil)
	hreq.ContentLength = int64(len(buff))
	hreq.Body = ioutil.NopCloser(bytes.NewBuffer(buff))
	err = id.SignRequestV1(hreq, aws.CanonicalizeS3, 15)
	if err == nil {
		resp, err = self.conn.Request(hreq)
		if err == nil {
			defer resp.Body.Close()
			err = aws.CodeToError(resp.StatusCode)
		}
	}
	return
}

// NB: Length is required as we do not buffer the reader
// NB(2): We do NOT close your reader (hence the io.Reader), 
// we wrap it with a NopCloser.
func (self *Bucket) PutKeyReader(id *aws.Signer, key string, r io.Reader, l int64, hdr http.Header) (err os.Error) {
	var resp *http.Response
	hreq := aws.NewRequest(self.key_url(key), "PUT", hdr, nil)
	hreq.ContentLength = l
	hreq.Body = ioutil.NopCloser(io.LimitReader(r, l))
	err = id.SignRequestV1(hreq, aws.CanonicalizeS3, 15)
	if err == nil {
		resp, err = self.conn.Request(hreq)
		if err == nil {
			defer resp.Body.Close()
			err = aws.CodeToError(resp.StatusCode)
		}
		if err == aws.ErrorUnexpectedResponse {
			ob, _ := http.DumpResponse(resp, true)
			os.Stdout.Write(ob)
		}
	}
	return
}

// Deletes the named key from the bucket.  To delete a bucket, see *Service.DeleteBucket()
func (self *Bucket) Delete(id *aws.Signer, key string) (err os.Error) {
	var resp *http.Response
	if key == "" {
		return os.NewError("Key cannot be empty!")
	}
	hreq := aws.NewRequest(self.key_url(key), "DELETE", nil, nil)
	err = id.SignRequestV1(hreq, aws.CanonicalizeS3, 15)
	if err == nil {
		resp, err = self.conn.Request(hreq)
	}
	if err == nil {
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNoContent {
			err = aws.CodeToError(resp.StatusCode)
		}
	}

	return
}

// Opens the named key and copys it to the named io.Writer IFF the response.Status is 200.
// Also returns the http headers for convenience (regardless of status code, as long as a  resp is generated).
func (self *Bucket) GetKey(id *aws.Signer, key string, w io.Writer) (hdr http.Header, err os.Error) {
	var resp *http.Response
	hreq := aws.NewRequest(self.key_url(key), "GET", nil, nil)
	err = id.SignRequestV1(hreq, aws.CanonicalizeS3, 15)
	if err == nil {
		resp, err = self.conn.Request(hreq)
	}
	if err == nil {
		defer resp.Body.Close()
		err = aws.CodeToError(resp.StatusCode)
		hdr = resp.Header
		if err == nil {
			_, err2 := io.Copy(w, resp.Body)
			if err == nil {
				err = err2
			}
		}
	}
	return
}

// Performs a HEAD request on the bucket and returns nil of the key appears
// valid (returns 200).
func (self *Bucket) Exists(id *aws.Signer, key string) (err os.Error) {
	_, err = self.HeadKey(id, key)
	return
}

// Performs a HEAD request on the bucket and returns the response object.
// The body is CLOSED, and it is an error to try and read from it.
func (self *Bucket) HeadKey(id *aws.Signer, key string) (resp *http.Response, err os.Error) {
	hreq := aws.NewRequest(self.key_url(key), "HEAD", nil, nil)
	err = id.SignRequestV1(hreq, aws.CanonicalizeS3, 15)
	if err == nil {
		resp, err = self.conn.Request(hreq)
	}
	if err == nil {
		resp.Body.Close()
		err = aws.CodeToError(resp.StatusCode)
	}
	return
}

type owner struct {
	ID          string
	DisplayName string
}

// Walks a bucket and writes the resulting strings to the channel.
// * There is currently NO (correct) way to abort a running walk.
func (self *Bucket) ListKeys(id *aws.Signer,
prefix, delim, marker string, out chan<- string) (err os.Error) {
	var done bool
	var resp *http.Response
	var last string
	form := http.Values{"prefix": []string{prefix},
											"delimeter": []string{delim},
											"marker":[]string{marker}}
	for err == nil && !done {
		result := listBucketResult{}
		result.Prefix = prefix
		result.Marker = marker
		if last != "" {form.Set("marker", last) }

		hreq := aws.NewRequest(self.key_url("/"), "GET", nil, form)
		err = id.SignRequestV1(hreq, aws.CanonicalizeS3, 15)

		if err == nil {
			resp, err = self.conn.Request(hreq)
		}
		if err == nil {
			err = aws.CodeToError(resp.StatusCode)
			if err == nil {
				err = xml.Unmarshal(resp.Body, &result)
				if err == nil {
					for i := range result.Contents {
						out <- result.Contents[i].Key
					}
					if len(result.Contents) > 0 {
						last = result.Contents[len(result.Contents) - 1].Key
					}
					done = !result.IsTruncated
				}
			}
			resp.Body.Close()
		}
	}
	close(out)
	return
}

type listBucketResult struct {
	Name        string
	Prefix      string
	Marker      string
	MaxKeys     int
	IsTruncated bool
	Contents    []keyItem
}

type keyItem struct {
	Key          string
	LastModified string // TODO: date
	ETag         string
	Size         int64
	Owner        owner
	StorageClass string
}

// Closes the underlying connection
func (self *Bucket) Close() (err os.Error) {
	return self.conn.Close()
}
