include $(GOROOT)/src/Make.inc

create_queue: create_queue.go common_flags.go
	$(GC) -o .$<.6 $^
	$(LD) -o $@ .$<.6

delete_queue: delete_queue.go common_flags.go
	$(GC) -o .$<.6 $^
	$(LD) -o $@ .$<.6

list_queues: list_queues.go common_flags.go
	$(GC) -o .$<.6 $^
	$(LD) -o $@ .$<.6

clean:
	rm -f create_queue list_queues delete_queue

install: create_queue delete_queue list_queues
