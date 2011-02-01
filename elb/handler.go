package elb

import "com.abneptis.oss/aws"
import "com.abneptis.oss/cryptools"

import "http"
import "os"
import "strconv"
import "xml"

type ELBHandler struct {
  c  aws.Channel
  signer cryptools.NamedSigner
}

func NewHandler(c aws.Channel, s cryptools.NamedSigner)(*ELBHandler){
  return &ELBHandler{c:c, signer:s }
}

type ELBListener struct {
  LoadBalancerPort int
  InstancePort int
  Protocol string
}

type elbCreateResult struct {
  DNSName string
}

type elbError struct {
  Type string
  Code string
  Message string
}

func (self elbError)String()(string) {
  return self.Code + ": " + self.Message
}

type elbResult struct {
  CreateLoadBalancerResult elbCreateResult
  Error	elbError
  RequestId string
}



func (self *ELBHandler)CreateLoadBalancer(zones []string, listeners []ELBListener, name string)(dnsname string, err os.Error){
  parms := map[string]string {}
  for i := range(zones) {
    parms["AvailabilityZones.members." + strconv.Itoa(1+i)] = zones[i]
  }
  parms["LoadBalancerName"] = name
  for i := range(listeners) {
    parms["Listeners.members." + strconv.Itoa(1+i) + ".InstancePort"] = strconv.Itoa(listeners[i].InstancePort)
    parms["Listeners.members." + strconv.Itoa(1+i) + ".LoadBalancerPort"] = strconv.Itoa(listeners[i].LoadBalancerPort)
    parms["Listeners.members." + strconv.Itoa(1+i) + ".Protocol"] = listeners[i].Protocol
  }
  req, err := newQuery(self.signer, self.c.Endpoint(), "CreateLoadBalancer", parms)
  if err == nil {
    var resp *http.Response
    var eresult elbResult
    resp, err = self.c.WriteRequest(req)
    if err == nil {
      err = xml.Unmarshal(resp.Body, &eresult)
      if err == nil {
        if eresult.Error.Code != "" {
          err = eresult.Error
        } else {
          dnsname = eresult.CreateLoadBalancerResult.DNSName
        }
      }
    }
  }
  return
}

func (self *ELBHandler)DeleteLoadBalancer(name string)(err os.Error){
  parms := map[string]string {}
  parms["LoadBalancerName"] = name
  req, err := newQuery(self.signer, self.c.Endpoint(), "DeleteLoadBalancer", parms)
  if err == nil {
    var resp *http.Response
    var eresult elbResult
    resp, err = self.c.WriteRequest(req)
    if err == nil {
      err = xml.Unmarshal(resp.Body, &eresult)
      if err == nil {
        if eresult.Error.Code != "" {
          err = eresult.Error
        }
      }
    }
  }
  return
}
