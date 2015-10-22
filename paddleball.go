package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// check args
	if len(os.Args) < 2 {
		fmt.Println("usage: " + os.Args[0] + " <file>")
		os.Exit(1)
	}

	// create pad files
	pad0, err := os.Create(os.Args[1] + ".0.pad")
	check(err)
	defer pad0.Close()
	pad1, err := os.Create(os.Args[1] + ".1.pad")
	check(err)
	defer pad1.Close()

	// open file
	file, err := os.Open(os.Args[1])
	check(err)
	defer file.Close()

	// loop through file
	fileByte := make([]byte, 1)
	for {
		_, err := file.Read(fileByte)
		if err == io.EOF {
			break
		} else {
			check(err)
		}

		// get random byte
		randByte := byte(rand.Intn(255))

		// xor random byte against pad byte
		padByte := fileByte[0] ^ randByte

		// coin flip
		if rand.Intn(1) == 1 { // write pads
			_, err := pad0.Write([]byte{randByte})
			check(err)

			_, err = pad1.Write([]byte{padByte})
			check(err)
		} else {
			_, err := pad1.Write([]byte{randByte})
			check(err)
			_, err = pad0.Write([]byte{padByte})
			check(err)
		}
	}
}
