package aws

import (
  "bytes"
  "crypto"
  "crypto/hmac"
  "encoding/base64"
  "hash"
  "os"
)

type Signer struct {
  AccessKey string
  secretAccessKey []byte
}

func NewSigner(akid, sak string)(*Signer){
  return &Signer{akid,bytes.NewBufferString(sak).Bytes()}
}

func (self Signer)SignBytes(h crypto.Hash, buff []byte)(sig []byte, err os.Error){
  hh := hmac.New(func()(hash.Hash){
     return h.New()
  }, self.secretAccessKey)
  _, err = hh.Write(buff)
  if err == nil {
    sig = hh.Sum()
  }
  return
}

func (self Signer)SignString(h crypto.Hash, s string)(os string, err os.Error){
  ob, err := self.SignBytes(h, bytes.NewBufferString(s).Bytes())
  if err == nil {
    os = string(ob)
  }
  return
}

func (self Signer)SignEncoded(h crypto.Hash, s string, e *base64.Encoding)(out []byte, err os.Error){
  ob, err := self.SignBytes(h, bytes.NewBufferString(s).Bytes())
  if err == nil { 
    out = make([]byte, e.EncodedLen(len(ob)))
    e.Encode(out, ob)
  }
  return
}
