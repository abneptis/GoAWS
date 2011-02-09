package cobranding

import "com.abneptis.oss/cryptools"
import "com.abneptis.oss/aws/awsconn"

import "http"
import "os"


type PrepaidGateway struct {
  signer cryptools.NamedSigner
  gwep   *awsconn.Endpoint
}

func NewPrepaidGateway(signer cryptools.NamedSigner, epurl *http.URL)(*PrepaidGateway){
  return &PrepaidGateway{
    signer: signer,
    gwep: awsconn.NewEndpoint(epurl, nil),
  }
}

func (self *PrepaidGateway)PrepaidCobrandedURL(params map[string]string)(url *http.URL, err os.Error){
  params["pipelineName"] = "SetupPrepaid"
  if _, ok := params["callerReferenceFunding"]; ! ok { return nil, os.NewError("callerReferenceFunding is mandatory") }
  if _, ok := params["callerReferenceSender"]; ! ok { return nil, os.NewError("callerReferenceSender is mandatory") }
  if _, ok := params["fundingAmount"]; ! ok { return nil, os.NewError("fundingAmount is mandatory") }
  
  ireq, err := calculateCobrandingURL(self.signer, self.gwep, params)
  if err == nil {
    url = ireq.URL
    if len(url.RawQuery) > 0 { url.RawQuery += "&" }
    url.RawQuery += http.EncodeQuery(ireq.Form)
  }
  return
}