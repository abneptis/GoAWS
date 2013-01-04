package main

import (
	aws "../.."
	. "../../flags"
	"../../s3"
	"errors"
	"net/url"
)

import "encoding/json"

type conf map[string]*proxyConf

type proxyConf struct {
	Bucket   *s3.Bucket
	Prefix   string
	Identity *aws.Signer
}

func (self *proxyConf) UnmarshalJSON(in []byte) (err error) {
	confmap := map[string]string{}
	err = json.Unmarshal(in, &confmap)
	if err == nil {
		if _, ok := confmap["Bucket"]; !ok {
			return errors.New("Bucket is required")
		}
		self.Prefix = confmap["Prefix"]
		if self.Prefix == "" {
			self.Prefix = "/"
		}
		self.Bucket = s3.NewBucket(&url.URL{
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
