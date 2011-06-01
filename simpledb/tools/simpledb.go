package main

import "com.abneptis.oss/aws/simpledb"
import "com.abneptis.oss/aws/awsconn"
import "com.abneptis.oss/aws"

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
var DoDeleteAttributes = flag.Bool("del-attributes", false, "Del attributes (-item is mandatory")

var Item     = flag.String("item","","Item name")
var Select   = flag.String("select","","Select query")
var Domain = flag.String("domain", "", "Domain name")
var DBUrl  = flag.String("DBUrl", "http://sdb.amazonaws.com", "SQS Endpoint")

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
  aws.AwsIDFlags()
  flag.Parse()
  url, err := http.ParseURL(*DBUrl)
  if err != nil {
    log.Fatalf("DBUrl (%s) invalid: (%v)", *DBUrl, err)
  }
  s, err := aws.DefaultIdentity("sha256")
  if err != nil {
    log.Fatalf("Couldn't create identity: %v", err)
  }
  ep := awsconn.NewEndpoint(url, nil)
  dbh := simpledb.NewHandler(*ep, s)
  if *DoCreate {
    _, err = dbh.CreateDomain(*Domain)
    if err != nil {
      log.Fatalf("Couldn't create domain (req): %v", err)
    }
  }
  if *DoList {
    doms, err := dbh.ListDomains("", 100)
    if err != nil {
      log.Fatalf("Couldn't list domain (req): %v", err)
    }
    for i := range(doms){
      fmt.Printf("%s\n", doms[i])
    }
  }
  if *DoPutAttributes {
    attrs, err := FlagAttributes()
    if err != nil {
      log.Fatalf("Error interpreting attr: %v", err)
    }
    err = dbh.PutAttributes(*Domain, *Item, attrs, nil)
  }
  if *DoGetAttributes {
    attrs, err := dbh.GetAttributes(*Domain, *Item, nil, false)
    if err != nil {
      log.Fatalf("Couldn't get item")
    } else {
      for attri := range(attrs){
        fmt.Printf("%s\t%s\n", attrs[attri].Name, attrs[attri].Value)
      }
    }
  }
  if *DoDeleteAttributes {
    err = dbh.DeleteAttributes(*Domain, *Item, nil, nil)
    if err != nil {
      log.Fatalf("Couldn't delete item attrs: %v", err)
    }
  }
  if *Select != "" {
    var items []simpledb.Item
    items, err = dbh.Select(*Select,"", false)
    for ii := range(items){
      fmt.Printf("%s\n", items[ii].Name)
      for ai := range(items[ii].Attribute){
        fmt.Printf("\t%s\t%s\n", items[ii].Attribute[ai].Name,
                               items[ii].Attribute[ai].Value)
      }
    }
  }
  if *DoDelete {
    _, err = dbh.DeleteDomain(*Domain)
    if err != nil {
      log.Fatalf("Couldn't delete domain (req): %v", err)
    }
  }
}
