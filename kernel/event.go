package main

/*
#include <stdlib.h>
typedef struct KernelEvent {
	char* name;
} KernelEvent;
*/
import "C"
import "unsafe"

var eventQueue = make(chan Event, 256)

type Event struct {
	name string
	data any
}

func PushEvent(e Event) {
	select {
	case eventQueue <- e:
	default:
		Log("Event", "failed to push %s. event queue length is %d\n", e.name, len(eventQueue))
	}
}

func mallocEvent(inEvent Event) *C.KernelEvent {
	cEvent := (*C.KernelEvent)(C.malloc(C.sizeof_KernelEvent))
	if cEvent == nil {
		panic("malloc failed")
	}

	cEvent.name = C.CString(inEvent.name)

	return cEvent
}

//export FreeEvent
func FreeEvent(event *C.KernelEvent) {
	if event != nil {
		if event.name != nil {
			C.free(unsafe.Pointer(event.name))
		}
		C.free(unsafe.Pointer(event))
	}
}

//export PollEvent
func PollEvent() *C.KernelEvent {
	select {
	case event := <-eventQueue:
		return mallocEvent(event)
	default:
		return nil
	}
}
