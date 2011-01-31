include $(GOROOT)/src/Make.inc

TARG=com.abneptis.oss/aws/simpledb
GOFILES=\
	attribute.go\
	consts.go\
	connection.go\
	handler.go\
	query.go\
	response.go\

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
