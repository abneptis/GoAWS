package aws

type Signer struct {
  AccessKey string
  secretAccessKey string
}

func NewSigner(akid, sak string)(*Signer){
  return &Signer{akid,sak}
}
