package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"os/exec"
	"unsafe"
)

//export Initialize
func Initialize(shellPath *C.char) {
	goShellPath := C.GoString(shellPath)
	shell := exec.Command(goShellPath)

	stdoutChan := make(chan string)
	stdout, err := shell.StdoutPipe()
	if err != nil {
		fmt.Printf("Failed to get stdout pipe: %v\n", err)
		return
	}

	stderrChan := make(chan string)
	stderr, err := shell.StderrPipe()
	if err != nil {
		fmt.Printf("Failed to get stdout pipe: %v\n", err)
		return
	}

	go Parse(stdout, stdoutChan)
	go Parse(stderr, stderrChan)

	err = shell.Start()
	if err != nil {
		return
	}
	fmt.Printf("running %d:%s\n", shell.Process.Pid, shell.Path)

	for output := range stdoutChan {
		fmt.Printf("%s", output)
	}

	err = shell.Wait()
	fmt.Printf("shell exited: %v\n", err)

}

func main() {
	shellPath := C.CString("/bin/ls")
	defer C.free(unsafe.Pointer(shellPath))

	Initialize(shellPath)
}
