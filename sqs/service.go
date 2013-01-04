package sqs

import (
	"aws"
	"errors"
	"net/url"
)

import (
	"crypto"
	"encoding/xml"
	"net/http"
	"strconv"
)

const (
	DEFAULT_HASH     = crypto.SHA256
	MAX_MESSAGE_SIZE = 64 * 1024
)

type Service struct {
	URL  *url.URL
	conn *aws.Conn
}

func NewService(url_ *url.URL) *Service {
	return &Service{
		URL:  url_,
		conn: aws.NewConn(aws.URLDialer(url_, nil)),
	}
}

func (self *Service) ListQueues(id *aws.Signer, prefix string) (mq []string, err error) {
	var resp *http.Response
	parms := url.Values{}
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
			err = errors.New("Unexpected response code")
		}
		if err == nil {
			mq = xresp.QueueURL
		}

	}
	return
}

// Create a queue, returning the Queue object.
func (self *Service) CreateQueue(id *aws.Signer, name string, dvtimeout int) (mq *Queue, err error) {
	var resp *http.Response
	parms := url.Values{}
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
					var qrl *url.URL
					qrl, err = url.Parse(xmlresp.QueueURL)
					if err == nil {
						mq = NewQueue(qrl)
					}
				}
			} else {
				err = errors.New("Unexpected response")
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
func (self *Service) Close() (err error) {
	return self.conn.Close()
}
