package main

import (
	"flag"
	"os"

	"github.com/Binject/debug/pe"
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
}
