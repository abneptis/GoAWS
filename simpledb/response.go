package simpledb

import "http"
import "xml"
import "os"
import "log"
import "io"

type Response http.Response


type ResponseMetadata struct {
  StatusCode string // not all requests generate this.
  RequestId string
  BoxUsage  string // string to avoid float precision issues
}

type ListDomainsResult struct {
  DomainName []string
  NextToken string
}

type GetAttributesResult struct {
  Attribute []Attribute
}

type Item struct {
  Name string
  Attribute []Attribute
}

type SelectResult struct {
  Item []Item
}

type SimpledbResponse struct {
  ListDomainsResult ListDomainsResult
  GetAttributesResult GetAttributesResult
  ResponseMetadata ResponseMetadata
  SelectResult SelectResult
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

