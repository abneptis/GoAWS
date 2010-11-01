// Signing and verification interface for AWS functionality.

// The aws/auth package implements a simple interface encouraging
// correct usage of public/private identity pairs (username/password,
// public/private key, etc).
//
// It makes no effort to ensure the data is "secure" (anonymously 
// mmapped, etc), simply that it's procedurally secure.
//
// By implementing the com.abneptis.org/cryptools/signer interface,
// we can quickly get various forms and formats since AWS does not
// use a single canonical signing format, but several similar ones.
package auth
/* 
  Copyright (c) 2010, Abneptis LLC.
  See COPYRIGHT and LICENSE for details.
*/

import "com.abneptis.oss/cryptools/signer"

import "os"

// The signer interface extends the cryptools signer interface by
// requiring a PublicIdentity() function that returns a []byte array
// that identifies the signing party.
type Signer interface {
  signer.Signer
  PublicIdentity()([]byte)
}

// A signable object must implement a canonical string.
type Signable interface {
  Sign(Signer)(os.Error)
}

// Wraps the auth.Signer interface by returning the PublicIdentity
// of a Signer as a string.
func GetSignerIDString(s Signer)(string){
  return string(s.PublicIdentity())
}
