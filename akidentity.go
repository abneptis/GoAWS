package aws
/* 
  Copyright (c) 2010, Abneptis LLC.
  See COPYRIGHT and LICENSE for details.
*/

import "com.abneptis.oss/cryptools/hashes"
import "com.abneptis.oss/cryptools"

import "bytes"
import "os"
import "hash"
import "crypto/hmac"

// An AWS identity (AccessKey/SecretKey)
type Identity struct {
  accessKeyID []byte
  secretAccessKey []byte
  sigHasher func()hash.Hash
}

// Constructs a new signer object based off of the
// an ak/sk string pair.
func NewIdentity(mech, ak, sk string)(id cryptools.NamedSigner, err os.Error){
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
func (self *Identity)PublicIdentity()(out []byte){
  out = make([]byte, len(self.accessKeyID))
  copy(out, self.accessKeyID)
  return
}

// (cryptools/signer/Sign()) - Implements the Sign() interface,
// returns a raw byte signature based off of a raw byte string-to-sign.
// Errors can only be returned on bad/short writes to the hash function.
func (self *Identity)Sign(s cryptools.Signable)(sig cryptools.Signature, err os.Error){
  hh := hmac.New(self.sigHasher, self.secretAccessKey)
  sb, err := s.SignableBytes()
  if err != nil { return }
  n, err := hh.Write(sb)
  if err == nil {
    sig = cryptools.NewSignature(hh.Sum())
  }
  if n != len(sb) {
    err = os.NewError("Hash function did not read entire string-to-sign")
  }
  return
}

// (cryptools/signer/Verify()) - Returns an error if the signature
// cannot be validated.
//
// NB: If the signing function returns an empty signature, AND the
// verification signature is empty, it is considered a pass.
func (self *Identity)VerifySignature(sig cryptools.Signature, o cryptools.Signable)(err os.Error){
  esig, err := self.Sign(o)
  if err == nil {
    eb := esig.SignatureBytes()
    sb := sig.SignatureBytes()
    if len(eb) == len(sb) {
      for i := range(eb) {
        if eb[i] != sb[i] {
          err = os.NewError("Signature verification failed")
        }
      }
    } else {
      err = os.NewError("Signature length mismatch")
    }
  }
  return
}

func (self *Identity)SignerName()(string){
  return string(self.accessKeyID)
}
