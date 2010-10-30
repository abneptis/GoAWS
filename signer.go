package goaws
/* 
  Copyright (c) 2010, Abneptis LLC.
  See COPYRIGHT and LICENSE for details.
*/

import "com.abneptis.oss/cryptools/signer"

type Signer interface {
  signer.Signer
  PublicIdentity()([]byte)
}

func GetSignerIDString(s Signer)(string){
  return string(s.PublicIdentity())
}
