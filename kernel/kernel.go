package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"github.com/creack/pty"
	"os"
	"unsafe"
)

//export Initialize
func Initialize(shellPath *C.char) {
	goShellPath := C.GoString(shellPath)
	shell := InitShell(goShellPath)

	stdoutData := make(chan string)
	//stderrChan := make(chan string)

	ptmx, err := pty.Start(shell.Cmd)
	if err != nil {
		Log("Shell", "Failed to start PTY")
		return
	}

	defer func(ptmx *os.File) {
		err := ptmx.Close()
		if err != nil {
			Log("Shell", "Failed to close PTY")
		}
	}(ptmx)

	go Parse(ptmx, stdoutData)

	Log("Kernel", "Running %d:%s\n", shell.Cmd.Process.Pid, shell.Cmd.Path)

	for {
		output, ok := <-stdoutData
		if ok {
			Log("Kernel", "stdout got: %s\n", output)
		} else {
			break
		}
	}

	Log("Kernel", "stdout/stderr channels closed\n")

	Log("Kernel", "waiting for shell proc to exit...")
	err = shell.Cmd.Wait()
	if err != nil {
		fmt.Printf(" failed with error %v\n", err)
	} else {
		fmt.Printf(" DONE\n")
	}
}

func main() {
	shellPath := C.CString("/bin/zsh")
	defer C.free(unsafe.Pointer(shellPath))

	Initialize(shellPath)
}
