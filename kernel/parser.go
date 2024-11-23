package main

import (
	"io"
	"strings"
)

func isText(b byte) bool {
	return b >= 32 && b <= 126
}

// scan may return any bytes not interpreted due to
// reading text and encountering a non-text code,
// multi character escape codes crossing buffer boundaries, etc.
func scan(bytes []byte) ([]byte, *Event) {
	for i, b := range bytes {
		if b == '\r' {
			return nil, &Event{
				name: "CarriageReturn",
				data: nil,
				kind: EventKindCarriageReturn,
			}
		} else if b == '\x00' {
			return nil, nil
		} else if b == '\x07' {
			return nil, &Event{
				name: "Bell",
				data: nil,
				kind: EventKindBell,
			}
		} else if b == '\x08' {
			return nil, &Event{
				name: "Backspace",
				data: nil,
				kind: EventKindBackspace,
			}
		} else if b == '\t' {
			return nil, &Event{
				name: "Tab",
				data: nil,
				kind: EventKindTab,
			}
		} else if b == '\n' {
			return nil, &Event{
				name: "LineFeed",
				data: nil,
				kind: EventKindLineFeed,
			}
		} else if b == '\v' {
			return nil, &Event{
				name: "VerticalTab",
				data: nil,
				kind: EventKindVerticalTab,
			}
		} else if b == '\x0C' {
			return nil, &Event{
				name: "FormFeed",
				data: nil,
				kind: EventKindFormFeed,
			}
		} else if isText(b) {
			// read until non-text byte
			var textBuilder strings.Builder
			var texti int
			for texti = i; texti < len(bytes); texti++ {
				if isText(bytes[texti]) {
					textBuilder.WriteByte(bytes[texti])
				} else {
					return bytes[texti:], &Event{
						name: "Text",
						data: textBuilder.String(),
						kind: EventKindText,
					}
				}
			}
			return nil, &Event{
				name: "Text",
				data: textBuilder.String(),
				kind: EventKindText,
			}
		}
	}

	return nil, &Event{
		name: "Error",
		data: bytes,
		kind: EventKindError,
	}
}

func Parse(r io.ReadCloser, outChan chan<- Event) {
	Log("Parser", "started\n")
	defer close(outChan)
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		Log("Parser", "read %d bytes\n", n)
		if n > 0 {
			remaining, event := scan(buf[:n])
			if event != nil {
				outChan <- *event
			}
			for len(remaining) > 0 {
				remaining, event = scan(remaining)
				outChan <- *event
			}
		}
		if err == io.EOF {
			Log("Parser", "Got EOF\n")
			break
		} else if err != nil {
			Log("Parser", "Parser Error: %v\n", err)
			break
		}
	}
}
