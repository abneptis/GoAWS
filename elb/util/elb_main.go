package elb_util

import (
	"errors"
	"github.com/abneptis/GoAWS"
	"github.com/abneptis/GoAWS/elb"
	. "github.com/abneptis/GoAWS/flags"
	. "github.com/abneptis/GoAWS/util/common"
	"net/url"
)

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

// Safety warning
// These are globals to allow the code to be more readable,
// since the tool is "single-tasked" it has no threading issues.
//
// You are of course encouraged to take a more thread-safe approach
// if you intend to use multiple threads.

var flag_endpoint_url string = ""
var signer *aws.Signer
var service *elb.Service

// Convenience method to clean up calls.
func DefaultELBService() (id *aws.Signer, s *elb.Service, err error) {
	id, err = DefaultSigner()
	if err == nil {
		url_, err := url.Parse(flag_endpoint_url)
		if err == nil {
			s = elb.NewService(url_)
		}
	}
	return
}

func init() {
	AddModule("elb", func() {
		flag.StringVar(&flag_endpoint_url, "elb-endpoint", "https://elasticloadbalancing.amazonaws.com/", "Endpoint to use for S3 calls")
	})
	Modules["elb"].Setup = func() (err error) {
		signer, service, err = DefaultELBService()
		return
	}
	Modules["elb"].Calls["destroy"] = func(args []string) (err error) {
		if len(args) != 1 {
			err = errors.New("Usage: destroy lbname")
		}
		err = service.DeleteLoadBalancer(signer, args[0])
		return
	}
	Modules["elb"].Calls["list"] = func(args []string) (err error) {
		if len(args) != 0 {
			err = errors.New("Usage: list")
		}
		lbd, err := service.DescribeLoadBalancers(signer)
		if err == nil {
			for lb := range lbd {
				fmt.Printf("%s\t%v\t%s\n", lbd[lb].LoadBalancerName, lbd[lb].AvailabilityZones, lbd[lb].DNSName)
			}
		}
		return
	}

	Modules["elb"].Calls["create"] = func(args []string) (err error) {
		if len(args) != 3 {
			return errors.New("Usage: create name zone[,zone2] lbport,iport,proto")
		}
		name := args[0]
		zones := strings.Split(args[1], ",")
		triple := strings.SplitN(args[2], ",", 3)
		if len(triple) != 3 {
			return errors.New("Invalid lbport/iport/proto triple")
		}
		l := elb.Listener{}
		l.InstancePort, err = strconv.Atoi(triple[0])
		if err != nil {
			return
		}
		l.LoadBalancerPort, err = strconv.Atoi(triple[1])
		if err != nil {
			return
		}
		l.Protocol = triple[2]
		err = service.CreateLoadBalancer(signer, name, zones, []elb.Listener{l})

		return
	}

}

func Nil() {}
