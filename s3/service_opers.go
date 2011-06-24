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

// Initilalize a new Service object with a specific
// S3 endpoint.  If URL is omitted, it defaults to the
// us-east endpoint over HTTPS (https://s3.amazonaws.com/) 
func NewService(url *http.URL) (s *Service) {
	s = &Service{
		URL: url,
	}
	if s.URL == nil {
		s.URL, _ = http.ParseURL("https://" + USEAST_HOST + "/")
	}
	s.conn = aws.NewConn(aws.URLDialer(s.URL, nil))
	return
}

// Returns a new *Bucket with the same URL data as the Service connection.  
// You MUST have already created the bucket in order to make use of the 
// Bucket object.
//
// See CreateBucket to create a new bucket.
func (self *Service) Bucket(name string) *Bucket {
	// We deliberately do NOT re-use our conn here, in order to take advantage
	// of vhosted bucket DNS perks.
	return NewBucket(self.URL, name, nil)
}

func s3Path(bucket, key string) string {
	return path.Join("/", bucket, key)
}

func (self *Service) bucket_url(bucket string) *http.URL {
	return &http.URL{
		Host:   self.URL.Host,
		Path:   path.Join(self.URL.Path, s3Path(bucket, "")),
		Scheme: self.URL.Scheme,
	}
}

// Deletes the named bucket from the S3 service.
func (self *Service) DeleteBucket(id *aws.Signer, name string) (err os.Error) {
	var resp *http.Response
	hreq := aws.NewRequest(self.bucket_url(name), "DELETE", nil, nil)
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

// Creates a new bucket
// TODO: Will (probably) create the bucket in US-east no matter
// what underlying endpoint you've chosen.
func (self *Service) CreateBucket(id *aws.Signer, name string) (err os.Error) {
	var resp *http.Response
	hreq := aws.NewRequest(self.bucket_url(name), "PUT", nil, nil)
	err = id.SignRequestV1(hreq, aws.CanonicalizeS3, 15)

	if err == nil {
		resp, err = self.conn.Request(hreq)
	}
	if err == nil {
		defer resp.Body.Close()
		err = aws.CodeToError(resp.StatusCode)
	}
	return
}


// Returns a list of bucket names known by the endpoint.  Depending on the
// endpoint used, your list may be global or regional in nature.
func (self *Service) ListBuckets(id *aws.Signer) (out []string, err os.Error) {
	var resp *http.Response
	hreq := aws.NewRequest(self.bucket_url(""), "GET", nil, nil)
	err = id.SignRequestV1(hreq, aws.CanonicalizeS3, 15)
	if err == nil {
		resp, err = self.conn.Request(hreq)
	}
	if err == nil {
		err = aws.CodeToError(resp.StatusCode)
	}
	if err == nil {
		defer resp.Body.Close()
		result := listAllMyBucketsResult{}
		err = xml.Unmarshal(resp.Body, &result)
		if err == nil {
			out = result.Buckets
		}
	}
	return
}

type listAllMyBucketsResult struct {
	Owner   owner
	Buckets []string "Buckets>Bucket>Name"
}

// Closes the underlying connection
func (self *Service) Close() (err os.Error) {
	return self.conn.Close()
}
