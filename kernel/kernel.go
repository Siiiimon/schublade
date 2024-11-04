package main

import (
	"C"
)

//export ProcessInput
func ProcessInput(input *C.char) *C.char {
	return input
}

func main() {}
