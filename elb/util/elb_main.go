package elb_util

import (
	"aws/elb"
	. "aws/flags"
	. "aws/util/common"
	"aws"
)

import (
	"flag"
	"fmt"
	"http"
	"os"
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
func DefaultELBService() (id *aws.Signer, s *elb.Service, err os.Error) {
	id, err = DefaultSigner()
	if err == nil {
		url, err := http.ParseURL(flag_endpoint_url)
		if err == nil {
			s = elb.NewService(url)
		}
	}
	return
}


func init() {
	AddModule("elb", func() {
		flag.StringVar(&flag_endpoint_url, "elb-endpoint", "https://elasticloadbalancing.amazonaws.com/", "Endpoint to use for S3 calls")
	})
	Modules["elb"].Setup = func() (err os.Error) {
		signer, service, err = DefaultELBService()
		return
	}
	Modules["elb"].Calls["destroy"] = func(args []string) (err os.Error) {
		if len(args) != 1 {
			err = os.NewError("Usage: destroy lbname")
		}
		err = service.DeleteLoadBalancer(signer, args[0])
		return
	}
	Modules["elb"].Calls["list"] = func(args []string) (err os.Error) {
		if len(args) != 0 {
			err = os.NewError("Usage: list")
		}
		lbd, err := service.DescribeLoadBalancers(signer)
		if err == nil {
			for lb := range lbd {
				fmt.Printf("%s\t%v\t%s\n", lbd[lb].LoadBalancerName, lbd[lb].AvailabilityZones, lbd[lb].DNSName)
			}
		}
		return
	}

	Modules["elb"].Calls["create"] = func(args []string) (err os.Error) {
		if len(args) != 3 {
			return os.NewError("Usage: create name zone[,zone2] lbport,iport,proto")
		}
		name := args[0]
		zones := strings.Split(args[1], ",", -1)
		triple := strings.Split(args[2], ",", 3)
		if len(triple) != 3 {
			return os.NewError("Invalid lbport/iport/proto triple")
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
