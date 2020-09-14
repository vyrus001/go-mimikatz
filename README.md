# go-mimikatz
A Go wrapper around a recompiled version of the Mimikatz executable for the purpose of anti-virus evasion.

This application will download the current version of mimikatz as a binary, extract the version appropriate for the architecture the application is being run on from the download zip file, transform that executable into encrypted shell code, then call the shell code in its own thread and transfer the current run time context to it.
### Notes:
* While this binary as of the date of its Github commit does not set off any defender flags, running it does, due to behavioral flagging. If you want to evade Windows memory protection, you have to add your own special sauce to the payload ;)
* If compiled as position independent code (`-buildmode=pie`) via go 1.15 or newer, this code can be transformed via [donut](https://github.com/Binject/go-donut) and then subsequently injected into another process on the target machine (a hint for those trying to avoid disk writes during deployment)
* There are a lot of forks of this code, some of them providing casual code - convenience fixes and then subsequently asking for BTC. If this is you, stop doing that shit, just submit a PR and don't be a dick 