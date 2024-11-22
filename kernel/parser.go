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
				name: "LineFeed",
				data: nil,
			}
		} else if b == '\000' {
			return nil, nil
		} else {
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
					}
				}
			}
			return nil, &Event{
				name: "Text",
				data: textBuilder.String(),
			}
		}
	}

	return nil, &Event{
		name: "Error",
		data: bytes,
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
