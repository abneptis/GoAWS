package sqs

type Message struct {
  ID string
  Message []byte
  VisibilityTimeout int
}
