package sdb

import (
	"errors"
	"github.com/abneptis/GoAWS"
	"net/url"
)

import (
	"encoding/xml"
	"net/http"
	"net/http/httputil"
	"os"
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

func (self *Service) Domain(name string) *Domain {
	return &Domain{
		URL: &url.URL{
			Scheme: self.URL.Scheme,
			Host:   self.URL.Host,
			Path:   self.URL.Path,
		},
		conn: self.conn,
		Name: name,
	}
}

func (self *Service) CreateDomain(id *aws.Signer, name string) (err error) {
	var resp *http.Response
	parms := url.Values{}
	parms.Set("DomainName", name)
	parms.Set("Action", "CreateDomain")
	req := aws.NewRequest(self.URL, "GET", nil, parms)
	err = id.SignRequestV2(req, aws.Canonicalize, DEFAULT_API_VERSION, 0)
	if err == nil {
		resp, err = self.conn.Request(req)
	}
	if err == nil {
		if resp.StatusCode != http.StatusOK {
			err = errors.New("Unexpected response")
		}
	}
	return
}

func (self *Service) DestroyDomain(id *aws.Signer, name string) (err error) {
	var resp *http.Response
	parms := url.Values{}
	parms.Set("DomainName", name)
	parms.Set("Action", "DeleteDomain")
	req := aws.NewRequest(self.URL, "GET", nil, parms)

	err = id.SignRequestV2(req, aws.Canonicalize, DEFAULT_API_VERSION, 0)
	if err == nil {
		resp, err = self.conn.Request(req)
	}

	if err == nil {
		resp, err = self.conn.Request(req)
	}
	if err == nil {
		if resp.StatusCode != http.StatusOK {
			err = errors.New("Unexpected response")
		}
	}
	return
}

func (self *Service) ListDomains(id *aws.Signer) (out []string, err error) {
	var resp *http.Response
	parms := url.Values{}
	parms.Set("Action", "ListDomains")
	parms.Set("MaxNumberOfDomains", "100")
	var done bool
	nextToken := ""
	for err == nil && !done {
		xmlresp := listdomainsresponse{}
		if nextToken != "" {
			parms.Set("NextToken", nextToken)
		} else {
			parms.Del("NextToken")
		}
		req := aws.NewRequest(self.URL, "GET", nil, parms)

		err = id.SignRequestV2(req, aws.Canonicalize, DEFAULT_API_VERSION, 0)
		if err == nil {
			resp, err = self.conn.Request(req)
		}

		if err == nil {
			resp, err = self.conn.Request(req)
		}
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				err = errors.New("Unexpected response")
				ob, _ := httputil.DumpResponse(resp, true)
				os.Stdout.Write(ob)
			}
			if err == nil {
				err = xml.NewDecoder(resp.Body).Decode(&xmlresp)
				if err == nil {
					if xmlresp.ErrorCode != "" {
						err = errors.New(xmlresp.ErrorCode)
					}
					if err == nil {
						for d := range xmlresp.Domains {
							out = append(out, xmlresp.Domains[d])
						}
					}
					nextToken = xmlresp.NextToken
				}
			}
		}
		done = (nextToken == "")
	}
	return
}

// Closes the underlying connection
func (self *Service) Close() (err error) {
	return self.conn.Close()
}
