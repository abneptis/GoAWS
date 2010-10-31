include $(GOROOT)/src/Make.inc

list_buckets: list_buckets.go common_flags.go 
	$(GC) -o .$<.6 $^
	$(LD) -o $@ .$<.6

clean:
	rm -f list_buckets

install: list_buckets
