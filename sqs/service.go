package sqs

import (
	"aws"
)

import (
	"crypto"
	"http"
	"os"
	"strconv"
	"xml"
)

const (
	DEFAULT_HASH     = crypto.SHA256
	MAX_MESSAGE_SIZE = 64 * 1024
)

type Service struct {
	URL  *http.URL
	conn *aws.Conn
}

func NewService(url *http.URL) *Service {
	return &Service{
		URL:  url,
		conn: aws.NewConn(aws.URLDialer(url, nil)),
	}
}

func (self *Service) ListQueues(id *aws.Signer, prefix string) (mq []string, err os.Error) {
	var resp *http.Response
	parms := http.Values{}
	parms.Set("Action", "ListQueues")
	if prefix != "" {
		parms.Set("QueueNamePrefix", prefix)
	}

	req := aws.NewRequest(self.URL, "GET", nil, parms)
	err = id.SignRequestV2(req, aws.Canonicalize, DEFAULT_VERSION, 15)
	resp, err = self.conn.Request(req)
	if err == nil {
		defer resp.Body.Close()
		xresp := listQueuesResponse{}
		if resp.StatusCode == http.StatusOK {
			err = xml.Unmarshal(resp.Body, &xresp)
		} else {
			err = os.NewError("Unexpected response code")
		}
		if err == nil {
			mq = xresp.QueueURL
		}

	}
	return
}

// Create a queue, returning the Queue object.
func (self *Service) CreateQueue(id *aws.Signer, name string, dvtimeout int) (mq *Queue, err os.Error) {
	var resp *http.Response
	parms := http.Values{}
	parms.Set("Action", "CreateQueue")
	parms.Set("QueueName", name)
	parms.Set("DefaultVisibilityTimeout", strconv.Itoa(dvtimeout))

	req := aws.NewRequest(self.URL, "GET", nil, parms)
	err = id.SignRequestV2(req, aws.Canonicalize, DEFAULT_VERSION, 15)
	if err == nil {
		resp, err = self.conn.Request(req)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				xmlresp := createQueueResponse{}
				err = xml.Unmarshal(resp.Body, &xmlresp)
				if err == nil {
					var qrl *http.URL
					qrl, err = http.ParseURL(xmlresp.QueueURL)
					if err == nil {
						mq = NewQueue(qrl)
					}
				}
			} else {
				err = os.NewError("Unexpected response")
			}
		}
	}

	return
}

type createQueueResponse struct {
	QueueURL  string "CreateQueueResult>QueueUrl"
	RequestId string "ResponseMetadata>RequestId"
}

type listQueuesResponse struct {
	QueueURL  []string "ListQueuesResult>QueueUrl"
	RequestId string   "ResponseMetadata>RequestId"
}

// Closes the underlying connection
func (self *Service) Close() (err os.Error) {
	return self.conn.Close()
}
