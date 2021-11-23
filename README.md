# go-mimikatz
A Go wrapper Mimikatz for the purpose of anti-virus evasion.

# Building
cd into the repo and run `go generate`

### Notes:
* evades windows again (as of 11/23/2021)
* If compiled as position independent code (`-buildmode=pie`) via go 1.15 or newer, this code can be transformed via [donut](https://github.com/Binject/go-donut) and then subsequently injected into another process on the target machine (a hint for those trying to avoid disk writes during deployment)
