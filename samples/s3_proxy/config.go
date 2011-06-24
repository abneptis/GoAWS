package main

import (
	. "aws/flags"
	"aws"
	"aws/s3"
)

import (
	"json"
	"http"
	"os"
)

type conf map[string]*proxyConf

type proxyConf struct {
	Bucket   *s3.Bucket
	Prefix   string
	Identity *aws.Signer
}

func (self *proxyConf) UnmarshalJSON(in []byte) (err os.Error) {
	confmap := map[string]string{}
	err = json.Unmarshal(in, &confmap)
	if err == nil {
		if _, ok := confmap["Bucket"]; !ok {
			return os.NewError("Bucket is required")
		}
		self.Prefix = confmap["Prefix"]
		if self.Prefix == "" {
			self.Prefix = "/"
		}
		self.Bucket = s3.NewBucket(&http.URL{
			Scheme: "http",
			Host:   "s3.amazonaws.com",
			Path:   self.Prefix,
		}, confmap["Bucket"], nil)

		if confmap["AccessKey"] != "" && confmap["SecretKey"] != "" {
			self.Identity = aws.NewSigner(confmap["AccessKey"], confmap["SecretKey"])
		} else {
			self.Identity, err = DefaultSigner()
		}
	}
	return
}
