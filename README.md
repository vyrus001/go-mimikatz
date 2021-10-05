# go-mimikatz
A Go wrapper Mimikatz for the purpose of anti-virus evasion.

# Building
cd into the repo and run `go generate`

### Notes:
* Now with 100% windows defender evasion! (as of 10/4/2021)
* If compiled as position independent code (`-buildmode=pie`) via go 1.15 or newer, this code can be transformed via [donut](https://github.com/Binject/go-donut) and then subsequently injected into another process on the target machine (a hint for those trying to avoid disk writes during deployment)
* There are a lot of forks of this code, some of them providing casual code - convenience fixes and then subsequently asking for BTC. If this is you, stop doing that shit, just submit a PR and don't be a dick 
