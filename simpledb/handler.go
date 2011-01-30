package simpledb

import "com.abneptis.oss/aws/auth"

import "http"
import "os"
import "strconv"

type Handler struct {
  conn AWSConnection
  signer auth.Signer
}

func NewHandler(c AWSConnection, a auth.Signer)(*Handler){
  return &Handler{ conn: c, signer: a}
}

func (self *Handler)doRequest(req *http.Request)(response SimpledbResponse, err os.Error){
 var resp *http.Response
 resp, err = self.conn.WriteRequest(req)
 if err == nil {
   response, err = ((*Response)(resp)).ParseResponse()
 }
 return
}


func (self *Handler)ListDomains(start string, max int)(domains []string, err os.Error){
  var response SimpledbResponse
  parms := map[string]string{}
  if max != 0 { parms["MaxNumberOfDomains"] = strconv.Itoa(max) }
  if start != "" { parms["NextToken"] = start }

  req, err := newQuery(self.signer, self.conn.Endpoint(), "", "ListDomains", parms)
  if err == nil {
    response, err = self.doRequest(req)
    if err == nil {
      domains = response.ListDomainsResult.DomainName
    }
  }
  return
}

func (self *Handler)CreateDomain(dn string)(response SimpledbResponse, err os.Error){
  req, err := newQuery(self.signer, self.conn.Endpoint(), dn, "CreateDomain", nil)
  if err == nil {
    response, err = self.doRequest(req)
  }
  return
}

func (self *Handler)DeleteDomain(dn string)(response SimpledbResponse, err os.Error){
  req, err := newQuery(self.signer, self.conn.Endpoint(), dn, "DeleteDomain", nil)
  if err == nil {
    response, err = self.doRequest(req)
  }
  return
}


func (self *Handler)PutAttributes(dn string, in string, attrs, expected AttributeList)(response SimpledbResponse, err os.Error){
  parms := map[string]string{
    "ItemName": in,
  }
  for i := range(attrs) {
    itoa := strconv.Itoa(i)
    parms["Attribute." + itoa + ".Name"] = attrs[i].Name
    parms["Attribute." + itoa + ".Value"] = attrs[i].Value
  }
  for i := range(expected) {
    itoa := strconv.Itoa(i)
    parms["Expected." + itoa + ".Name"] = expected[i].Name
    parms["Expected." + itoa + ".Value"] = expected[i].Value
  }
  req, err := newQuery(self.signer, self.conn.Endpoint(), dn, "PutAttributes", parms)
  if err == nil {
    response, err = self.doRequest(req)
  }
  return
}

func (self *Handler)GetAttributes(dn string, in string, attrs []string, consistant bool)(out []Attribute, err os.Error){
  parms := map[string]string{
    "ItemName": in,
  }
  for i := range(attrs) {
   parms["Attribute." + strconv.Itoa(i) + ".Name"] = attrs[i]
  }
  req, err := newQuery(self.signer, self.conn.Endpoint(), dn, "GetAttributes", parms)
  var response SimpledbResponse
  if err == nil {
    response, err = self.doRequest(req)
    if err == nil { out = response.GetAttributesResult.Attribute }
  }
  return
}