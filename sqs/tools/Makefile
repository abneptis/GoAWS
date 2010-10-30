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

pop: pop.go common_flags.go
	$(GC) -o .$<.6 $^
	$(LD) -o $@ .$<.6

push: push.go common_flags.go push_flags.go
	$(GC) -o .$<.6 $^
	$(LD) -o $@ .$<.6

clean:
	rm -f create_queue list_queues delete_queue push pop

install: create_queue delete_queue list_queues push pop
