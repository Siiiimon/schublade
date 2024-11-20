package main

import "C"
import (
	"io"
	"os/exec"
)

type Shell struct {
	Cmd    *exec.Cmd
	Stdout io.ReadCloser
	Stderr io.ReadCloser
}

func InitShell(shellPath string) *Shell {
	cmd := exec.Command(shellPath)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		Log("Shell", "Failed to get stdout pipe: %v\n", err)
		return nil
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		Log("Shell", "Failed to get stderr pipe: %v\n", err)
		return nil
	}

	return &Shell{
		Cmd:    cmd,
		Stdout: stdout,
		Stderr: stderr,
	}
}
