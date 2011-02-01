package main

import "com.abneptis.oss/aws"
import "com.abneptis.oss/aws/elb"

import "flag"
import "fmt"
import "log"
import "strings"


var accessKey = flag.String("access-key", "", "AWS Access Key ID")
var secretKey = flag.String("secret-key", "", "AWS Secret Access Key")
var region    = flag.String("region","us-east-1", "AWS Region to use")
var lbName    = flag.String("name","", "load balancer name")
var zonestr    = flag.String("zones","","availability zones")
var listnersstr = flag.String("listeners","","listeners xport:iport:svc[,...]")

var doCreate = flag.Bool("create", false, "Create an LB")
var doDelete = flag.Bool("delete", false, "Delete an LB")


func main(){
  flag.Parse()
  ep, err := elb.GetEndpoint(*region, nil)
  if err != nil {
    log.Exitf("ERROR: %v", err)
  }
  zones := strings.Split(*zonestr, ",", -1)
  listners := []elb.ELBListener{}
  if *listnersstr != "" {
    listnerstrs := strings.Split(*listnersstr, ",", -1)
    for i := range(listnerstrs) {
      l := elb.ELBListener{}
      _, err = fmt.Sscanf(listnerstrs[i], "%d:%d:%s", &l.LoadBalancerPort, &l.InstancePort, &l.Protocol)
      if err == nil {
        listners = append(listners, l)
      } else {
        log.Exitf("Invalid listener: %v", err)
      }
    }
  }
  ch := aws.NewConnection(ep, "tcp", "", true)
  auth, err := aws.NewIdentity("sha256", *accessKey, *secretKey)
  if err != nil {
    log.Exitf("ERROR: %v", err)
  }

  elbh := elb.NewHandler(ch, auth)
  if *doCreate {
    out, err := elbh.CreateLoadBalancer(zones, listners, *lbName)
    if err != nil {
      log.Exitf("ERROR: %v", err)
    }
    log.Printf("DNSName: %s", out)
  }
  if *doDelete {
    err := elbh.DeleteLoadBalancer(*lbName)
    if err != nil {
      log.Exitf("ERROR: %v", err)
    }
  }
}
