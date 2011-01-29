package simpledb

import "http"
import "xml"
import "os"
import "log"
import "io"

type Response http.Response


type ResponseMetadata struct {
  RequestId string
  BoxUsage  string // string to avoid float precision issues
}

type SimpledbResponse struct {
  ResponseMetadata ResponseMetadata
}


func (self Response)ParseResponse()(resp SimpledbResponse, err os.Error){
  switch self.StatusCode {
    case http.StatusOK: err = xml.Unmarshal(self.Body, &resp)
    case http.StatusNotFound: err = os.NewError("Not found")
    default:
	err = os.NewError("Unexpected status code")
	log.Printf("Unexpected status code: %d", self.StatusCode)
        io.Copy(os.Stdout, self.Body)
  }
  return
}

