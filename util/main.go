package main

import (
	"errors" // AWS ID Flags
	ec2_util "github.com/abneptis/GoAWS/ec2/util"
	elb_util "github.com/abneptis/GoAWS/elb/util"
	. "github.com/abneptis/GoAWS/flags" // AWS ID Flags
	s3_util "github.com/abneptis/GoAWS/s3/util"
	sdb_util "github.com/abneptis/GoAWS/sdb/util"
	sqs_util "github.com/abneptis/GoAWS/sqs/util"
	. "github.com/abneptis/GoAWS/util/common"
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
