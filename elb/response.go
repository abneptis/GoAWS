package elb

type LoadBalancerDescription struct {
	CanonicalHostedZoneName       string
	CanonicalHostedZoneNameID     string
	CreatedTime                   string
	LoadBalancerName              string
	SourceSecurityGroupOwnerAlias string "sourcesecuritygroup>owneralias"
	SourceSecurityGroupGroupName  string "sourcesecuritygroup>groupname"
	DNSName                       string
	Listeners                     []Listener "listenerdescriptions>member>listener"
	AvailabilityZones             []string   "AvailabilityZones>member"
}

type LoadBalancerQueryResult struct {
	LoadBalancerDescription []LoadBalancerDescription "DescribeLoadBalancersResult>LoadBalancerDescriptions>member"
	RequestId               string                    "requestmetadata>requestid"
	ErrorCode               string                    "error>errorcode"
}
