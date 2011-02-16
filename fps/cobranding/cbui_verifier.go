package cobranding

import "http"
import "log"
import "xml"

type CBUIReturnVerifier struct {
  Success http.Handler
  Failure http.Handler
  GatewayURL *http.URL
  Logger *log.Logger
}

type verifySignatureResult struct {
  VerificationStatus string
}

type responseMetadata struct {
  RequestId string
}

type cbuiResponse struct {
  VerifySignatureResult verifySignatureResult
  ResponseMetadata responseMetadata
}

func (self CBUIReturnVerifier)ServeHTTP(rw http.ResponseWriter, req *http.Request){
  if req.FormValue("signatureVersion") != "2" {
    req.Form["CBUI.Error"] = []string{"Invalid CBUI signature version"}
    self.Failure.ServeHTTP(rw, req)
    self.Logger.Printf("signatureVersion not provided: %v", req.URL.RawQuery)
    return
  }
  myurl := &http.URL {
    Host: req.Host,
    Path: req.URL.Path,
    Scheme: "http",
  }
  if rw.UsingTLS() { myurl.Scheme += "s" }

  vurl := &http.URL{
    Scheme: self.GatewayURL.Scheme,
    Host: self.GatewayURL.Host,
    Path: self.GatewayURL.Path,
    RawQuery: http.EncodeQuery(map[string][]string{
      "UrlEndPoint": []string{myurl.String()},
      "HttpParameters": []string{req.URL.RawQuery},
      "Action": []string{"VerifySignature"},
      "Version": []string{"2008-09-17"},
    }),
  }
  
  self.Logger.Printf("Verifying signature from %s", self.GatewayURL.String())
  
  resp, _, err := http.Get(vurl.String())
  if err != nil {
    self.Logger.Printf("Get Failed: %v", err)
    req.Form["CBUI.Error"] = []string { err.String() }
    self.Failure.ServeHTTP(rw, req)
  } else {
    xresp := cbuiResponse{}
    err = xml.Unmarshal(resp.Body, &xresp)
    if err != nil {
      req.Form["CBUI.Error"] = []string { err.String() }
      self.Failure.ServeHTTP(rw, req)    
      return
    }
    if xresp.VerifySignatureResult.VerificationStatus != "Success" {
      req.Form["CBUI.Error"] = []string{"Amazon refused signature verification"}
      self.Failure.ServeHTTP(rw, req)
      return
    }
    req.Form["CBUI.Ok"] = []string{"true"}
    self.Success.ServeHTTP(rw, req)
  }
}