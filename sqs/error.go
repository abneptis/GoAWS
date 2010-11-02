package sqs

type errorResponse struct {
  Error Error
}

func (self errorResponse)String()(string){
  return self.Error.String()
}

type Error struct {
  Type string
  Code string
  Message string
  Detail string
}

func (self Error)String()(string){
  return "{" + self.Code + "}: " + self.Message
}
