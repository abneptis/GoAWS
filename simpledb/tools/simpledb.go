package main

import "com.abneptis.oss/aws/simpledb"
import "com.abneptis.oss/aws/awsconn"
import "com.abneptis.oss/aws/auth"

import "flag"
import "http"
import "log"

var DoCreate = flag.Bool("create", false, "Create the domain")
var DoDelete = flag.Bool("delete", false, "Delete the domain")
var Domain = flag.String("domain", "", "Domain name")
var DBUrl  = flag.String("DBUrl", "http://sdb.amazonaws.com", "SQS Endpoint")
var AccessKey  = flag.String("access-key-id", "", "Access key")
var SecretAccessKey  = flag.String("secret-access-key", "", "Secret access key")


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
    resp, err := dbh.CreateDomain(*Domain)
    if err != nil {
      log.Exitf("Couldn't create domain (req): %v", err)
    }
    log.Printf("Response: %v", resp)
  }
  if *DoDelete {
    resp, err := dbh.DeleteDomain(*Domain)
    if err != nil {
      log.Exitf("Couldn't delete domain (req): %v", err)
    }
    log.Printf("Response: %v", resp)
  }
}
