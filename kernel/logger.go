package main

import "fmt"

func Log(scope string, log string, args ...any) {
	fmt.Printf("[%s] %s", scope, fmt.Sprintf(log, args...))
}
