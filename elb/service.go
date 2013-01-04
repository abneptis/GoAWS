package elb

import (
	aws ".."
	"errors"
	"net/url"
)

import (
	"crypto"
	"encoding/xml"
	. "fmt"
	"net/http/httputil"
	"os"
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

type Listener struct {
	InstancePort     int
	LoadBalancerPort int
	Protocol         string
	SSLCertificateID string
}

func (self Listener) SetValues(v url.Values, i int) {
	v.Set(Sprintf("Listeners.members.%d.LoadBalancerPort", i), strconv.Itoa(self.LoadBalancerPort))
	v.Set(Sprintf("Listeners.members.%d.InstancePort", i), strconv.Itoa(self.InstancePort))
	v.Set(Sprintf("Listeners.members.%d.Protocol", i), self.Protocol)
}

func (self *Service) CreateLoadBalancer(id *aws.Signer, name string, zones []string, listeners []Listener) (err error) {
	parms := url.Values{}
	parms.Set("Action", "CreateLoadBalancer")
	parms.Set("LoadBalancerName", name)
	for zi := range zones {
		parms.Set(Sprintf("AvailabilityZones.members.%d", zi+1), zones[zi])
	}
	for li := range listeners {
		listeners[li].SetValues(parms, li+1)
	}
	req := aws.NewRequest(self.URL, "GET", nil, parms)
	err = id.SignRequestV2(req, aws.Canonicalize, DEFAULT_VERSION, 0)
	if err != nil {
		return
	}
	resp, err := self.conn.Request(req)
	if err == nil {
		defer resp.Body.Close()
		ob, _ := httputil.DumpResponse(resp, true)
		os.Stdout.Write(ob)
	}
	return

}

func (self *Service) DescribeLoadBalancers(id *aws.Signer) (lbs []LoadBalancerDescription, err error) {
	parms := url.Values{}
	parms.Set("Action", "DescribeLoadBalancers")
	req := aws.NewRequest(self.URL, "GET", nil, parms)
	err = id.SignRequestV2(req, aws.Canonicalize, DEFAULT_VERSION, 0)
	if err != nil {
		return
	}
	resp, err := self.conn.Request(req)
	if err == nil {
		qr := LoadBalancerQueryResult{}
		defer resp.Body.Close()
		err = xml.NewDecoder(resp.Body).Decode(&qr)
		if err == nil {
			lbs = qr.LoadBalancerDescription
		}
	}

	return
}

// Users note: amazon will only return an error if the request is bad,
// thus an error will not be raised when deleting a non-existent LB.
func (self *Service) DeleteLoadBalancer(id *aws.Signer, name string) (err error) {
	parms := url.Values{}
	parms.Set("Action", "DeleteLoadBalancer")
	parms.Set("LoadBalancerName", name)
	req := aws.NewRequest(self.URL, "GET", nil, parms)
	err = id.SignRequestV2(req, aws.Canonicalize, DEFAULT_VERSION, 0)
	if err != nil {
		return
	}
	resp, err := self.conn.Request(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = aws.CodeToError(resp.StatusCode)
	}
	qr := LoadBalancerQueryResult{}
	err = xml.NewDecoder(resp.Body).Decode(&qr)
	if err == nil {
		if qr.ErrorCode != "" {
			err = errors.New(qr.ErrorCode)
		}
	}
	return
}

// Closes the underlying connection
func (self *Service) Close() (err error) {
	return self.conn.Close()
}
