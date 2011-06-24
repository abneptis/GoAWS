package s3

import (
	"crypto"
)


// Do S3 endpoints actually play any role?
const (
	USWEST_HOST      = "us-west-1.s3.amazonaws.com"
	USEAST_HOST      = "s3.amazonaws.com"
	APSOUTHEAST_HOST = "ap-southeast-1.s3.amazonaws.com"
	EUWEST_HOST      = "eu-west-1.s3.amazonaws.com"
)

const (
	DEFAULT_HASH = crypto.SHA1
)
