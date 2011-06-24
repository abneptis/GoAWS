package ec2

import (
	"aws"
)

import (
	"http"
	"log"
	"os"
	"xml"
)

type Service struct {
	conn *aws.Conn
	URL  *http.URL
}

func NewService(url *http.URL) (s *Service) {
	return &Service{
		URL:  url,
		conn: aws.NewConn(aws.URLDialer(url, nil)),
	}
}

func (self *Service) DescribeInstances(id *aws.Signer, filter http.Values, ic chan Instance) (err os.Error) {
	if filter == nil {
		filter = http.Values{}
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
		err := xml.Unmarshal(resp.Body, &xresp)
		if err == nil {
			log.Printf("XRESP == %+v", xresp)
		} else {
			log.Printf("XERR == %+v", err)
		}
		ob, _ := http.DumpResponse(resp, true)
		os.Stdout.Write(ob)
	}

	return
}

// Closes the underlying connection
func (self *Service) Close() (err os.Error) {
	return self.conn.Close()
}
