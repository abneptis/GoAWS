include $(GOROOT)/src/Make.inc

TARG=com.abneptis.oss/aws
GOFILES=\
	channel.go\
	escape.go\
	request_map.go\
	timeformats.go\
	akidentity.go\

DEPS=awsconn

include $(GOROOT)/src/Make.pkg

module.%: %
	make -C $*

module_install.%: module.%
	make -C $* install

module_clean.%:
	make -C $* clean

modules: module.sqs module.s3 module.simpledb module.awsconn
modules_install: module_install.sqs module_install.s3 module_install.simpledb module_install.awsconn
modules_clean:  module_clean.sqs module_clean.s3 module_clean.simpledb module_clean.awsconn
