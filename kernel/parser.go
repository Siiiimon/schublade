package main

import (
	"io"
)

func Parse(r io.ReadCloser, outChan chan<- string) {
	Log("Parser", "started\n")
	defer close(outChan)
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		Log("Parser", "read %d bytes\n", n)
		if n > 0 {
			output := string(buf[:n])
			outChan <- output
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
