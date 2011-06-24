package sdb


type listdomainsresponse struct {
	Domains      []string "ListDomainsResult>DomainName"
	NextToken    string   "ListDomainsResult>NextToken"
	RequestId    string   "ResponseMetadata>RequestId"
	BoxUsage     string   "ResponseMetadata>BoxUsage"
	ErrorMessage string   "Errors>Error>Message"
	ErrorCode    string   "Errors>Error>Code"
}


type getattributesresponse struct {
	Attributes   []Attribute "GetAttributesResult>Attribute"
	ErrorMessage string      "Errors>Error>Message"
	ErrorCode    string      "Errors>Error>Code"
}

type selectresponse struct {
	Items        []Item "SelectResult>Item"
	ErrorMessage string "Errors>Error>Message"
	ErrorCode    string "Errors>Error>Code"
	NextToken    string
}
