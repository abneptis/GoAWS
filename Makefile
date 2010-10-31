include $(GOROOT)/src/Make.inc

TARG=com.abneptis.oss/aws
GOFILES=request_map.go\
	connection.go\
	timeformats.go\
	endpoint/endpoint.go\

include $(GOROOT)/src/Make.pkg

