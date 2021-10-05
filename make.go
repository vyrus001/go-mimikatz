package main

//go:generate go run packer/packer.go -o stub
//go:generate go build -o stub/stub github.com/vyrus001/go-mimikatz/stub
//go:generate go run fixer/fixer.go -s stub/stub -o go-mimikatz.exe
