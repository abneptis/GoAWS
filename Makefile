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
	gomake -C $*

module_install.%: module.%
	gomake -C $* install

module_clean.%:
	gomake -C $* clean

modules: module.sqs module.s3 module.sdb module.ec2 module.elb
modules_install: module_install.sqs module_install.s3 module_install.sdb module_install.ec2 module_install.elb
modules_clean:  module_clean.sqs module_clean.s3 module_clean.sdb module_clean.elb module_clean.ec2
