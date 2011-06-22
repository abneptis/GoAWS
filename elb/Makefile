include $(GOROOT)/src/Make.inc

TARG=aws/elb
GOFILES=\
	consts.go\
	response.go\
	service.go\

DEPS=../

CLEANFILES+=\
	tools/*.[568vq]\
	tools/simpledb

include $(GOROOT)/src/Make.pkg


tools: tools/simpledb

tools/simpledb.$(O): tools/simpledb.go
	$(GC) -o $@ $^

tools/simpledb: tools/simpledb.$(O)
	$(LD) -o $@ $^
