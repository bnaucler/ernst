# ernst v0.2A
A very unfriendly IRC bot.

## Written by
B Naucler (mail@bnaucler.se)

## Setup
1. Add your favorite insults, one per line to insults.txt (UTF8 encoding)
2. `go get github.com/bnaucler/ernst`
3. `cd $GOPATH/src/github.com/bnaucler/ernst`
4. `go build -o mkdb/mkdb mkdb/mkdb.go`
5. `mkdb/mkdb insults.txt ernst.db`
6. `go build -o ernst *.go`
7. `./ernst`

You can use `skymfer.txt` instead of your own list if you want to be insulted in Swedish.

## Usage
Insults are added with !skymf \<insult\> in the channel.  
Number of insults are reported with !skymfstat(s).  
Last insult spoken is removed from database with !srm.

## Configuration options
Some options can be configured with !sset in the channel. For more extensive configuration, poke around in ernst.go and elib.go.

* rate (0-1000) chance of random insult triggering
* kdel (0-1000) delay per "keystroke" (in ms). Default: 100
* randel (0-1000) random delay per line (in ms). Default: 700
* dnrmem (0-20) number of insults to remember and avoid repetition. Default: 5

## Dump database to file
1. `go build -o dbdump/dbdump dbdump/dbdump.go`  
2. `dbdump/dbdump ernst.db > insults.txt`

## Thanks to
\#ljusdal @ EFNet  
Thomas Jager  
The BoltDB team  

## TODO
* This code base is a mess and in urgent need of refactoring
* Support for channel keys
* Definition of triggers in settings struct
* Convert settings to map
* Avoid repetition when directly addressed
* Multiple channels
	- As []string
* Multiple servers
	- As server[port]-map
* Support for non-SSL servers
* Auto-reconnect
* Shared delay between threads
* Search function returning number (privmsg)
* Search and remove
* Remove last added
* Increased ratio on trigger words
* PRIVMSG goroutine as func

## License
MIT:
Do whatever you want
