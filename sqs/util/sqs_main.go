package sqs_util

import (
	aws "../.."
	. "../../flags"
	"../../sqs"
	. "../../util/common"
	"errors"
	"net/url"
)

import (
	"flag"
	"fmt"
	"io"
	"os"
)

var flag_endpoint_url string
var flag_default_timeout int
var flag_pop_timeout int

var signer *aws.Signer
var s *sqs.Service

func DefaultSQSService() (id *aws.Signer, s *sqs.Service, err error) {
	id, err = DefaultSigner()
	if err == nil {
		url_, err := url.Parse(flag_endpoint_url)
		if err == nil {
			s = sqs.NewService(url_)
		}
	}
	return
}

func init() {
	AddModule("sqs", func() {
		flag.StringVar(&flag_endpoint_url, "sqs-endpoint", "https://queue.amazonaws.com/", "Endpoint to use")
		flag.IntVar(&flag_default_timeout, "sqs-queue-timeout", 90, "Queue timeout (create/delete)")
		flag.IntVar(&flag_pop_timeout, "sqs-message-timeout", 90, "Queue timeout (pop/peek)")
	})
	Modules["sqs"].Setup = func() (err error) {
		signer, s, err = DefaultSQSService()
		return
	}
	Modules["sqs"].Calls["create"] = func(args []string) (err error) {
		if len(args) != 1 {
			return errors.New("Usage: create QUEUE")
		}
		Q, err := s.CreateQueue(signer, args[0], flag_default_timeout)
		if err == nil {
			fmt.Printf("%s\n", Q.URL)
		}
		return
	}

	Modules["sqs"].Calls["list"] = func(args []string) (err error) {
		if len(args) != 0 {
			return errors.New("Usage: list")
		}
		qs, err := s.ListQueues(signer, "")
		if err == nil {
			for i := range qs {
				fmt.Printf("%s\n", qs[i])
			}
		}
		return
	}

	Modules["sqs"].Calls["drop"] = func(args []string) (err error) {
		if len(args) != 1 {
			return errors.New("Usage: drop queue")
		}
		Q, err := s.CreateQueue(signer, args[0], flag_default_timeout)
		if err == nil {
			err = Q.DeleteQueue(signer)
		}
		return
	}

	Modules["sqs"].Calls["push"] = func(args []string) (err error) {
		if len(args) != 1 {
			return errors.New("Usage: push queuename")
		}
		Q, err := s.CreateQueue(signer, args[0], flag_default_timeout)
		if err == nil {
			var n int
			lr := io.LimitReader(os.Stdin, sqs.MAX_MESSAGE_SIZE)
			buff := make([]byte, sqs.MAX_MESSAGE_SIZE)
			n, err = io.ReadFull(lr, buff)
			if err == nil || err == io.ErrUnexpectedEOF {
				buff = buff[0:n]
				err = Q.Push(signer, buff)
			}
		}
		return
	}
	Modules["sqs"].Calls["rm"] = func(args []string) (err error) {
		if len(args) != 2 {
			return errors.New("Usage: rm queuename receipthandle")
		}
		Q, err := s.CreateQueue(signer, args[0], flag_default_timeout)
		if err == nil {
			err = Q.Delete(signer, args[1])
		}
		return
	}
	Modules["sqs"].Calls["peek"] = func(args []string) (err error) {
		if len(args) != 1 {
			return errors.New("Usage: peek queuename")
		}
		Q, err := s.CreateQueue(signer, args[0], flag_default_timeout)
		var body []byte
		var id string
		if err == nil {
			body, id, err = Q.Peek(signer, flag_pop_timeout)
		}
		if err == nil {
			fmt.Printf("# MessageId %s\n", id)
			os.Stdout.Write(body)
		}
		return
	}
}

func Nil() {}
