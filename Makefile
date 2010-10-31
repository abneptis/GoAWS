include $(GOROOT)/src/Make.inc

TARG=com.abneptis.oss/aws
GOFILES=request_map.go\
	connection.go\
	timeformats.go\

include $(GOROOT)/src/Make.pkg

