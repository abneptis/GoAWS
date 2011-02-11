package aws

import "com.abneptis.oss/cryptools"

import "flag"
import "os"

var defaultAccessKey, defaultSecretKey string

func AwsIDFlags(){
  flag.StringVar(&defaultAccessKey, "aws-access-key", os.Getenv("AWS_ACCESS_KEY_ID"), "AWS Access Key")
  flag.StringVar(&defaultSecretKey, "aws-secret-key", os.Getenv("AWS_SECRET_ACCESS_KEY"), "AWS Secret Key")  
}

func DefaultIdentity(mech string)(cryptools.NamedSigner, os.Error){
  flag.Parse()
  return NewIdentity(mech, defaultAccessKey, defaultSecretKey)
}

