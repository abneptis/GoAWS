package flags
// The flags module is primarily a convenience module for those
// who want to request the AWS identity on the command line
// It should be noted that this approach is not particularly secure
// against a malicious root, but very little is.

import (
	"aws"
)

import (
	"flag"
	"os"
)


var accessKey *string = flag.String("aws-access-key",
	os.Getenv("AWS_ACCESS_KEY_ID"), "AWS Access Key")

var secretKey *string = flag.String("aws-secret-key",
	os.Getenv("AWS_SECRET_ACCESS_KEY"), "AWS Secret Key")

// Returns the aws.Signer associated with the aws-*-key flags.
func DefaultSigner() (signer *aws.Signer, err os.Error) {
	if accessKey == nil || secretKey == nil {
		flag.Parse()
	}
	if *accessKey == "" || *secretKey == "" {
		err = os.NewError("No default access key provided")
	} else {
		signer = aws.NewSigner(*accessKey, *secretKey)
	}
	return
}

// An empty function to allow easy package usage (for side-effect only imports)
func UseFlags() {}
