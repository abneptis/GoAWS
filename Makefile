include $(GOROOT)/src/Make.inc

TARG=com.abneptis.oss/aws
GOFILES=request_map.go\
	timeformats.go\
	escape.go\

include $(GOROOT)/src/Make.pkg

