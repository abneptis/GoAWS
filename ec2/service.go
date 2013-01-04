package ec2

import (
	aws ".."
	"net/url"
)

import (
	"encoding/xml"
	"log"
	"net/http/httputil"
	"os"
)

type Service struct {
	conn *aws.Conn
	URL  *url.URL
}

func NewService(url_ *url.URL) (s *Service) {
	return &Service{
		URL:  url_,
		conn: aws.NewConn(aws.URLDialer(url_, nil)),
	}
}

func (self *Service) DescribeInstances(id *aws.Signer, filter url.Values, ic chan Instance) (err error) {
	if filter == nil {
		filter = url.Values{}
	}
	filter.Set("Action", "DescribeInstances")
	req := aws.NewRequest(self.URL, "GET", nil, filter)
	err = id.SignRequestV2(req, aws.Canonicalize, DEFAULT_VERSION, 15)
	if err != nil {
		return
	}
	resp, err := self.conn.Request(req)
	if err == nil {
		defer resp.Body.Close()
		xresp := describeInstancesResponse{}
		err := xml.NewDecoder(resp.Body).Decode(&xresp)
		if err == nil {
			log.Printf("XRESP == %+v", xresp)
		} else {
			log.Printf("XERR == %+v", err)
		}
		ob, _ := httputil.DumpResponse(resp, true)
		os.Stdout.Write(ob)
	}

	return
}

// Closes the underlying connection
func (self *Service) Close() (err error) {
	return self.conn.Close()
}
