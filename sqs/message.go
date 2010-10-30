package sqs

import "com.abneptis.oss/uuid"

import "bytes"
import "os"
import "json"

var MaxMessageSize = 64*1024

type rawAttribute struct {
  Name string
  Value string
}

type rawMessage struct {
  MessageId string
  ReceiptHandle string
  MD5OfBody string
  Attribute []rawAttribute
  Body string
}

func (self *rawMessage)Message()(m *Message, err os.Error){
  m = &Message{
    ReceiptHandle: self.ReceiptHandle,
    MD5OfBody: self.MD5OfBody,
    Attributes: map[string]string{},
  }
  m.MessageId, err = uuid.Parse(self.MessageId)
  if err != nil { return }
  for i := range(self.Attribute){
    m.Attributes[ self.Attribute[i].Name ] = self.Attribute[i].Value
  }
  bbuff := bytes.NewBufferString(self.Body)
  m.Body = bbuff.Bytes()
  return
}

type Message struct {
  MessageId *uuid.UUID
  ReceiptHandle string
  MD5OfBody string
  Body []byte
  Attributes map[string]string
}

func (self *Message)MarshalJSON()(out []byte, err os.Error){
  return json.Marshal(map[string]interface{}{
   "MessageId": self.MessageId,
   "ReceiptHandle": self.ReceiptHandle,
   "MD5OfBody": self.MD5OfBody,
   "Body": string(self.Body),
   "Attributes": self.Attributes,
  })
}
