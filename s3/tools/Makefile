include $(GOROOT)/src/Make.inc

create_bucket: create_bucket.go common_flags.go 
	$(GC) -o .$<.6 $^
	$(LD) -o $@ .$<.6

delete_bucket: delete_bucket.go common_flags.go 
	$(GC) -o .$<.6 $^
	$(LD) -o $@ .$<.6

list_buckets: list_buckets.go common_flags.go 
	$(GC) -o .$<.6 $^
	$(LD) -o $@ .$<.6

clean:
	rm -f list_buckets create_bucket delete_bucket

install: list_buckets create_bucket delete_bucket
