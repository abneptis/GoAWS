package simpledb

import "http"
import "xml"
import "os"

type Response http.Response


type responseMetadata struct {
  StatusCode string // not all requests generate this.
  RequestId string
  BoxUsage  string // string to avoid float precision issues
}

type listDomainsResult struct {
  DomainName []string
  NextToken string
}

type getAttributesResult struct {
  Attribute []Attribute
}

type Item struct {
  Name string
  Attribute []Attribute
}

type selectResult struct {
  Item []Item
}

type responseError struct {
  Error errorResult
}

type errorResult struct {
  Code string
  Message string
}

type SimpledbResponse struct {
  ListDomainsResult listDomainsResult
  GetAttributesResult getAttributesResult
  ResponseMetadata responseMetadata
  SelectResult selectResult
  Errors responseError
}


func (self Response)ParseResponse()(resp SimpledbResponse, err os.Error){
  err = xml.Unmarshal(self.Body, &resp)
  if err != nil {
    return
  }
  if resp.Errors.Error.Code != "" {
    err = os.NewError(resp.Errors.Error.Code)
  }
  return
}

