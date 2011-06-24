package sdb_util

import (
	. "aws/flags"
	. "aws/util/common"
	"aws"
	"aws/sdb"
)

import (
	"flag"
	"fmt"
	"http"
	"os"
)

var flag_endpoint_url string
var id *aws.Signer
var service *sdb.Service

func DefaultSDBService() (id *aws.Signer, s *sdb.Service, err os.Error) {
	id, err = DefaultSigner()
	if err == nil {
		url, err := http.ParseURL(flag_endpoint_url)
		if err == nil {
			s = sdb.NewService(url)
		}
	}
	return
}

func init() {
	AddModule("sdb", func() {
		flag.StringVar(&flag_endpoint_url, "sdb-endpoint", "https://sdb.amazonaws.com/", "Endpoint to use for S3 calls")
	})

	Modules["sdb"].Setup = func() (err os.Error) {
		id, service, err = DefaultSDBService()
		return
	}
	Modules["sdb"].Calls["rm"] = func(args []string) (err os.Error) {
		if len(args) < 2 {
			return os.NewError("Usage: rm domain_name item [...]")
		}
		d := service.Domain(args[0])
		args = args[1:]
		for i := range args {
			err = d.DeleteAttribute(id, args[i], nil, nil)
			if err != nil {
				return
			}
		}
		return
	}

	Modules["sdb"].Calls["get"] = func(args []string) (err os.Error) {
		if len(args) < 2 {
			return os.NewError("Usage: get domain_name item [...]")
		}
		d := service.Domain(args[0])
		args = args[1:]
		for i := range args {
			var attrs []sdb.Attribute
			attrs, err = d.GetAttribute(id, args[i], nil, false)
			if err != nil {
				return
			}
			fmt.Printf("Item: %+v", attrs)
		}
		return
	}
	Modules["sdb"].Calls["create"] = func(args []string) (err os.Error) {
		if len(args) != 1 {
			return os.NewError("Usage: create domain_name")
		}
		err = service.CreateDomain(id, args[0])
		return
	}
	Modules["sdb"].Calls["drop"] = func(args []string) (err os.Error) {
		if len(args) != 0 {
			return os.NewError("Usage: drop domain_name")
		}
		err = service.DestroyDomain(id, args[0])
		return
	}
	Modules["sdb"].Calls["select"] = func(args []string) (err os.Error) {
		if len(args) < 2 || len(args) > 3 {
			return os.NewError("Usage: service.lect ('*'|col,col2,...) domain_name [extended expression]")
		}
		colstr := args[0]
		d := service.Domain(args[1])
		expr := ""
		if len(args) == 3 {
			expr = args[2]
		}
		c := make(chan sdb.Item)
		go func() {
			for i := range c {
				fmt.Printf("%s\n", i.Name)
				for ai := range i.Attribute {
					fmt.Printf("\t%s\t%s\n", i.Attribute[ai].Name, i.Attribute[ai].Value)
				}
			}
		}()
		err = d.Select(id, colstr, expr, true, c)
		close(c)

		return
	}
	Modules["sdb"].Calls["domains"] = func(args []string) (err os.Error) {
		doms, err := service.ListDomains(id)
		for i := range doms {
			fmt.Printf("%s\n", doms[i])
		}
		return
	}
}

func Nil() {}
