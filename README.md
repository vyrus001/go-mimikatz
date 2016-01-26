# go-mimikatz
A Go wrapper around a pre-compiled version of the Mimikatz executable for the purpose of anti-virus evasion.
### Requirements:
	go-bindata => https://github.com/jteeuwen/go-bindata
	MemoryModule => https://github.com/fancycode/MemoryModule

This application utilizes 3 segmented components to provide a Go wrapper for the Mimikatz application that is not considered malicious by most anti-virus software without additional packing, and can be dynamically built utilizing a repeatable build recipie. This is done by deviding the mimikatz executible into 2 randomly generated pads that are then stored as strings within the compiled Go binary and combined, and subsiquently loaded from within the existing process memory space at run time.

### Build Process:
1. Build or aquire Mimikatz 32 bit and 64 bit executibles
2. Use paddleball.go to devide the executibles into "pad" files

	Example: go run paddleball.go mimikatz32.exe (will output mimikatz32.exe.0.pad and mimikatz32.exe.1.pad)
3. Store the pad files within the main package of the go-mimikatz.go application

	Example: go-bindata mimikatz32.exe.0.pad and mimikatz32.exe.1.pad mimikatz64.exe.0.pad and mimikatz64.exe.1.pad (will
output bindata.go)
4. Build the MemoryModule library with MinGW

	(thanks Kost for figuring out how to do this using multiple build utilities AND operating systems!!! https://github.com/kost/MemoryModule)
	note, MemoryModule must be in the same path as mimikatz.go
5. run go build
