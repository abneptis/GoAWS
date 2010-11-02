package sqs

// This file contains nominal structures to make XML parsing
// easier.  They should almost never be used directly.

// All responses will contain a responseMetadata

type responseMetadata struct {
  RequestId string
}

/* Outer response types */

type createQueueResponse struct {
  CreateQueueResult createQueueResult
  ResponseMetadata *responseMetadata
}

type sendMessageResponse struct {
  SendMessageResult  sendMessageResult
  ResponseMetadata *responseMetadata
}

type listQueuesResponse struct {
  ListQueuesResult listQueuesResult
  ResponseMetadata *responseMetadata
}

type receiveMessageResponse struct {
  ReceiveMessageResult receiveMessageResult
  ResponseMetadata *responseMetadata
}

type deleteMessageResponse struct {
  ResponseMetadata responseMetadata
}

type deleteQueueResponse struct {
  ResponseMetadata responseMetadata
}

/* Inner response types */

type createQueueResult struct {
  QueueUrl string
}


type listQueuesResult struct {
  QueueUrl []string
}


type sendMessageResult struct {
  MD5OfMessageBody string
  MessageId string
}

type receiveMessageResult struct {
  Message []*rawMessage
}


