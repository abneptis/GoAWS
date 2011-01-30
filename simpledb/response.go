package simpledb

import "http"
import "xml"
import "os"
import "log"
import "io"

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

type SimpledbResponse struct {
  ListDomainsResult listDomainsResult
  GetAttributesResult getAttributesResult
  ResponseMetadata responseMetadata
  SelectResult selectResult
}


func (self Response)ParseResponse()(resp SimpledbResponse, err os.Error){
  switch self.StatusCode {
    case http.StatusOK: err = xml.Unmarshal(self.Body, &resp)
    case http.StatusNotFound: err = os.NewError("Not found")
    case http.StatusForbidden: err = os.NewError("Not authorized")
    default:
	err = os.NewError("Unexpected status code")
	log.Printf("Unexpected status code: %d", self.StatusCode)
        io.Copy(os.Stdout, self.Body)
  }
  return
}

