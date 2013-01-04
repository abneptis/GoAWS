package sqs

import (
	"aws"
	"errors"
	"net/url"
)

import (
	"encoding/xml"
	"net/http"
	"strconv"
)

type Queue struct {
	URL  *url.URL
	conn *aws.Conn
}

func NewQueue(url_ *url.URL) *Queue {
	return &Queue{
		URL:  url_,
		conn: aws.NewConn(aws.URLDialer(url_, nil)),
	}
}

func (self *Queue) DeleteQueue(id *aws.Signer) (err error) {
	var resp *http.Response
	parms := url.Values{}
	parms.Set("Action", "DeleteQueue")

	req := aws.NewRequest(self.URL, "GET", nil, parms)
	err = id.SignRequestV2(req, aws.Canonicalize, DEFAULT_VERSION, 15)
	if err == nil {
		resp, err = self.conn.Request(req)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				err = errors.New("Unexpected response")
			}
		}
	}

	return
}

func (self *Queue) Push(id *aws.Signer, body []byte) (err error) {
	var resp *http.Response
	parms := url.Values{}
	parms.Set("Action", "SendMessage")
	parms.Set("MessageBody", string(body))
	req := aws.NewRequest(self.URL, "GET", nil, parms)
	err = id.SignRequestV2(req, aws.Canonicalize, DEFAULT_VERSION, 15)
	if err == nil {
		resp, err = self.conn.Request(req)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				err = errors.New("Unexpected response")
			}
		}
	}
	return
}

// Note: 0 is a valid timeout!!
func (self *Queue) Peek(id *aws.Signer, vt int) (body []byte, msgid string, err error) {
	var resp *http.Response
	parms := url.Values{}
	parms.Set("Action", "ReceiveMessage")
	if vt >= 0 {
		parms.Set("VisibilityTimeout", strconv.Itoa(vt))
	}
	req := aws.NewRequest(self.URL, "GET", nil, parms)
	err = id.SignRequestV2(req, aws.Canonicalize, DEFAULT_VERSION, 15)
	if err == nil {
		resp, err = self.conn.Request(req)
		if err == nil {
			defer resp.Body.Close()
		}
		if err == nil && resp.StatusCode != http.StatusOK {
			err = errors.New("Unexpected response")
		}
		if err == nil {
			msg := message{}
			err = xml.Unmarshal(resp.Body, &msg)
			if err == nil {
				body, msgid = msg.Body, msg.ReceiptHandle
			}
		}
	}
	return
}

// Note: 0 is a valid timeout!!
func (self *Queue) Delete(id *aws.Signer, mid string) (err error) {
	var resp *http.Response
	parms := url.Values{}
	parms.Set("Action", "DeleteMessage")
	parms.Set("ReceiptHandle", mid)
	req := aws.NewRequest(self.URL, "GET", nil, parms)
	err = id.SignRequestV2(req, aws.Canonicalize, DEFAULT_VERSION, 15)
	if err == nil {
		resp, err = self.conn.Request(req)
		if err == nil {
			defer resp.Body.Close()
		}
		if resp.StatusCode != http.StatusOK {
			err = errors.New("Unexpected response")
		}
	}
	return
}

type message struct {
	MessageId     string "ReceiveMessageResult>Message>MessageId"
	ReceiptHandle string "ReceiveMessageResult>Message>ReceiptHandle"
	MD5OfBody     string "ReceiveMessageResult>Message>MD5OfBody"
	Body          []byte "ReceiveMessageResult>Message>Body"
}

// Closes the underlying connection
func (self *Queue) Close() (err error) {
	return self.conn.Close()
}
