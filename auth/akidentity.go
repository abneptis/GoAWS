package auth
/* 
  Copyright (c) 2010, Abneptis LLC.
  See COPYRIGHT and LICENSE for details.
*/

import "com.abneptis.oss/cryptools/hashes"

import "bytes"
import "os"
import "hash"
import "crypto/hmac"

type Identity struct {
  accessKeyID []byte
  secretAccessKey []byte
  sigHasher func()hash.Hash
}

func NewIdentity(mech, ak, sk string)(id Signer, err os.Error){
  akb := bytes.NewBufferString(ak)
  skb := bytes.NewBufferString(sk)
  hf, err := hashes.GetHashFunc(mech)
  if err != nil { return }
  id = &Identity{
    accessKeyID: akb.Bytes(),
    secretAccessKey: skb.Bytes(),
    sigHasher: hf,
  }
  return
}

// We dupe the slice to ensure nobody changes it
// down the internal values.
func (self *Identity)PublicIdentity()(out []byte){
  out = make([]byte, len(self.accessKeyID))
  copy(out, self.accessKeyID)
  return
}

func (self *Identity)Sign(sts []byte)(out []byte, err os.Error){
  hh := hmac.New(self.sigHasher, self.secretAccessKey)
  n, err := hh.Write(sts)
  if err == nil {
    out = hh.Sum()
  }
  if n != len(sts) {
    err = os.NewError("Hash function did not read entire string-to-sign")
  }
  return
}

func (self *Identity)Verify(uvsig, sts []byte)(err os.Error){
  sig, err := self.Sign(sts)
  if err == nil {
    if len(uvsig) == len(sig) {
      for i := range(uvsig) {
        if sig[i] != uvsig[i] {
          err = os.NewError("Signature verification failed")
        }
      }
    } else {
      err = os.NewError("Signature length mismatch")
    }
  }
  return
}
