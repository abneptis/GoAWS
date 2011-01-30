package main

import "com.abneptis.oss/aws/simpledb"
import "com.abneptis.oss/aws/awsconn"
import "com.abneptis.oss/aws/auth"

import "flag"
import "fmt"
import "http"
import "log"
import "os"
import "strings"

var DoCreate = flag.Bool("create", false, "Create the domain")
var DoDelete = flag.Bool("delete", false, "Delete the domain")
var DoList   = flag.Bool("list", false, "List domains")
var DoGetAttributes = flag.Bool("get-attributes", false, "Get attributes (-item is mandatory")
var DoPutAttributes = flag.Bool("put-attributes", false, "Put attributes (-item is mandatory")

var Item     = flag.String("item","","Item name")
var Domain = flag.String("domain", "", "Domain name")
var DBUrl  = flag.String("DBUrl", "http://sdb.amazonaws.com", "SQS Endpoint")
var AccessKey  = flag.String("access-key-id", "", "Access key")
var SecretAccessKey  = flag.String("secret-access-key", "", "Secret access key")

func FlagAttributes()(out []simpledb.Attribute, err os.Error){
  args := flag.Args()
  for fi := range(args){
    attr := strings.Split(args[fi],  "=",2)
    switch len(attr) {
      case 1: out = append(out, simpledb.Attribute{Name: attr[0]})
      case 2: out = append(out, simpledb.Attribute{Name: attr[0], Value: attr[1]})
      default: err = os.NewError("Unparsable parameter: " + args[fi])
    }
  }
  return
}

func main(){
  flag.Parse()
  url, err := http.ParseURL(*DBUrl)
  if err != nil {
    log.Exitf("DBUrl (%s) invalid: (%v)", *DBUrl, err)
  }
  s, err := auth.NewIdentity("sha256", *AccessKey, *SecretAccessKey)
  if err != nil {
    log.Exitf("Couldn't create identity: %v", err)
  }
  ep := awsconn.NewEndpoint(url, nil)
  dbh := simpledb.NewHandler(simpledb.NewConnection(*ep, "tcp", ""), s)
  if *DoCreate {
    _, err = dbh.CreateDomain(*Domain)
    if err != nil {
      log.Exitf("Couldn't create domain (req): %v", err)
    }
  }
  if *DoList {
    doms, err := dbh.ListDomains("", 100)
    if err != nil {
      log.Exitf("Couldn't list domain (req): %v", err)
    }
    for i := range(doms){
      fmt.Printf("%s\n", doms[i])
    }
  }
  if *DoPutAttributes {
    attrs, err := FlagAttributes()
    if err != nil {
      log.Exitf("Error interpreting attr: %v", err)
    }
    _, err = dbh.PutAttributes(*Domain, *Item, attrs, nil)
  }
  if *DoGetAttributes {
    attrs, err := dbh.GetAttributes(*Domain, *Item, nil, false)
    if err != nil {
      log.Exitf("Couldn't get item")
    } else {
      for attri := range(attrs){
        fmt.Printf("%s\t%s\n", attrs[attri].Name, attrs[attri].Value)
      }
    }
  }
  if *DoDelete {
    _, err = dbh.DeleteDomain(*Domain)
    if err != nil {
      log.Exitf("Couldn't delete domain (req): %v", err)
    }
  }
}
