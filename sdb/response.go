package sdb


type listdomainsresponse struct {
  Domains []string "ListDomainsResult>DomainName"
  NextToken string "ListDomainsResult>NextToken"
  RequestId string "ResponseMetadata>RequestId"
  BoxUsage string "ResponseMetadata>BoxUsage"
  ErrorMessage  string "Errors>Error>Message"
  ErrorCode string "Errors>Error>Code"
}

