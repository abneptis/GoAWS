package ec2_util

import (
	"aws"
	"aws/ec2"
	. "aws/flags"
	. "aws/util/common"
	"net/url"
)

import (
	"flag"
	"fmt"
)

// Safety warning
// These are globals to allow the code to be more readable,
// since the tool is "single-tasked" it has no threading issues.
//
// You are of course encouraged to take a more thread-safe approach
// if you intend to use multiple threads.

var flag_endpoint_url string = ""

// Convenience method to clean up calls.
func DefaultEC2Service() (id *aws.Signer, s *ec2.Service, err error) {
	id, err = DefaultSigner()
	if err == nil {
		url_, err := url.Parse(flag_endpoint_url)
		if err == nil {
			s = ec2.NewService(url_)
		}
	}
	return
}

func init() {
	AddModule("ec2", func() {
		flag.StringVar(&flag_endpoint_url, "ec2-endpoint", "https://ec2.amazonaws.com/", "Endpoint to use for EC2 calls")
	})
	Modules["ec2"].Calls["instances"] = func(args []string) (err error) {
		id, s, err := DefaultEC2Service()
		if err != nil {
			return
		}
		c := make(chan ec2.Instance)
		go func() {
			for i := range c {
				fmt.Printf("%+v\n", i)
			}
		}()
		err = s.DescribeInstances(id, nil, c)
		return
	}
}

func Nil() {}
