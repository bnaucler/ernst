# ernst v0.2A
A very unfriendly IRC bot.

## Written by
B Naucler (mail@bnaucler.se)

## Setup
1. Add your favorite insults, one per line to insults.txt (UTF8 encoding)
2. `go build mkdb.go`
3. `./mkdb insults.txt ernst.db skymf`
4. `go run ernst.go`

You can use skymfer.txt instead of your own list if you want to be insulted in Swedish.

## Usage
Insults are added with !skymf \<insult\> in the channel.  
Number of insults are reported with !skymfstat(s).

## Configuration options
Some options can be configured with !sset in the channel. For more extensive configuration, poke around in ernst.go.

* rate (0-1000) chance of random insult triggering
* kdel (0-1000) delay per "keystroke" (in ms). Default: 100
* randel (0-1000) random delay per line (in ms). Default: 700

## Dump database to file
1. `go build dbdump.go`  
2. `./dbdump ernst.db skymf > insults.txt`

## Thanks to
\#ljusdal @ EFNet  
Thomas Jager  
The BoltDB team  

## TODO
* Support for channel keys
* Avoid unnecessary type casting
* Convert settings to map
* Multiple channels
* Multiple servers
* Increased ratio on trigger words
* Increased ratio with time
* Requesting specific insults per keyword
	- DB search function
* PRIVMSG goroutine as func

## License
MIT:
Do whatever you want
