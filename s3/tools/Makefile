include $(GOROOT)/src/Make.inc

all: list_buckets create_bucket delete_bucket get_object list_objects delete_object put_object

create_bucket: create_bucket.go common_flags.go 
	$(GC) -o .$<.6 $^
	$(LD) -o $@ .$<.6

delete_object: delete_object.go common_flags.go 
	$(GC) -o .$<.6 $^
	$(LD) -o $@ .$<.6

delete_bucket: delete_bucket.go common_flags.go 
	$(GC) -o .$<.6 $^
	$(LD) -o $@ .$<.6

list_buckets: list_buckets.go common_flags.go 
	$(GC) -o .$<.6 $^
	$(LD) -o $@ .$<.6

list_objects: list_objects.go common_flags.go 
	$(GC) -o .$<.6 $^
	$(LD) -o $@ .$<.6

get_object: get_object.go common_flags.go 
	$(GC) -o .$<.6 $^
	$(LD) -o $@ .$<.6

put_object: put_object.go common_flags.go put_flags.go
	$(GC) -o .$<.6 $^
	$(LD) -o $@ .$<.6

clean:
	rm -f list_buckets create_bucket delete_bucket get_object list_objects delete_object put_object

install: all
