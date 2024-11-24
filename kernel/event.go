package main

/*
#include <stdlib.h>
typedef enum {
    KIND_TEXT = 0,
    KIND_LINE_FEED,
    KIND_CARRIAGE_RETURN,
    KIND_BACKSPACE,
    KIND_TAB,
    KIND_BELL,
    KIND_VERTICAL_TAB,
    KIND_FORM_FEED,
	KIND_CURSOR_MOVE,
	KIND_ERROR,
} EventKind;
typedef struct KernelEvent {
	char* name;
	void* data;
	EventKind kind;
} KernelEvent;
*/
import "C"
import "unsafe"

var eventQueue = make(chan Event, 256)

type Kind C.EventKind

const (
	EventKindText           Kind = C.KIND_TEXT
	EventKindLineFeed       Kind = C.KIND_LINE_FEED
	EventKindCarriageReturn Kind = C.KIND_CARRIAGE_RETURN
	EventKindBackspace      Kind = C.KIND_BACKSPACE
	EventKindTab            Kind = C.KIND_TAB
	EventKindBell           Kind = C.KIND_BELL
	EventKindVerticalTab    Kind = C.KIND_VERTICAL_TAB
	EventKindFormFeed       Kind = C.KIND_FORM_FEED
	EventKindCursorMove     Kind = C.KIND_CURSOR_MOVE
	EventKindError          Kind = C.KIND_ERROR
)

type Event struct {
	name string
	data any
	kind Kind
}

func PushEvent(e Event) {
	select {
	case eventQueue <- e:
		Log("Parser", "pushed event %s: '%s'\n", e.name, e.data)
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
	cEvent.kind = C.EventKind(inEvent.kind)

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
