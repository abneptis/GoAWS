include $(GOROOT)/src/Make.inc

TARG=aws
GOFILES=\
	consts.go\
	dialer.go\
	http_dialer.go\
	signer.go\
	escape.go\
	timeformats.go\

include $(GOROOT)/src/Make.pkg

module.%: %
	make -C $*

module_install.%: module.%
	make -C $* install

module_clean.%:
	make -C $* clean

modules: module.sqs module.s3 module.sdb
modules_install: module_install.sqs module_install.s3 module_install.sdb
modules_clean:  module_clean.sqs module_clean.s3 module_clean.sdb
