package sqs

type SQSError struct {
	Type    string
	Code    string
	Message string
	// <Detail?>
	RequestId string
}
