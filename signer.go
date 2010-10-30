package goaws

import "com.abneptis.oss/cryptools/signer"

type Signer interface {
  signer.Signer
  PublicIdentity()([]byte)
}

func GetSignerIDString(s Signer)(string){
  return string(s.PublicIdentity())
}
