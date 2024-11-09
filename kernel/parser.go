package main

import (
	"fmt"
	"io"
)

func Parse(r io.ReadCloser, outChan chan<- string) {
	defer close(outChan)
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			output := string(buf[:n])
			outChan <- output
		}
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error reading output:", err)
			break
		}
	}
}
