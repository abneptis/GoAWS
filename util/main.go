package main

import (
	"aws/ec2/ec2_util"
	"aws/elb/elb_util"
	. "aws/flags" // AWS ID Flags
	"aws/s3/s3_util"
	"aws/sdb/sdb_util"
	"aws/sqs/sqs_util"
	. "aws/util/common"
	"errors" // AWS ID Flags
)

import (
	"flag"
	"fmt"
	"os"
)

func keys(in map[string]interface{}) (out []string) {
	for k, _ := range in {
		out = append(out, k)
	}
	return
}

func main() {
	flag.Parse()
	if flag.NArg() < 2 {
		fmt.Printf("USAGE: aws [ec2|elb|s3|sdb|sqs] subcommand ...\n")
		os.Exit(1)
	}
	module := flag.Arg(0)
	cmd := flag.Arg(1)
	os.Args = os.Args[2:]
	var err error
	modulenames := []string{}
	for k, _ := range Modules {
		modulenames = append(modulenames, k)
	}
	if m, ok := Modules[module]; ok {
		m.FlagFunc()
		flag.Parse()
		if c, ok := m.Calls[cmd]; ok {
			err = m.Setup()
			if err == nil {
				err = c(flag.Args())
			}
		} else {
			err = errors.New(fmt.Sprintf("Invalid subcommand: %s, expected one of %v",
				cmd, m.Names()))
		}
	} else {
		err = errors.New(fmt.Sprintf("Invalid modulle : %s, expected one of %v",
			flag.Arg(0), modulenames))

	}
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
	UseFlags() // we want the side effects of import...
	ec2_util.Nil()
	elb_util.Nil()
	s3_util.Nil()
	sdb_util.Nil()
	sqs_util.Nil()
}
