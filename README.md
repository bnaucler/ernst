# ernst v0.2A
A very unfriendly IRC bot.

## Written by
B Naucler (mail@bnaucler.se)

## Usage
1. Add your favorite insults, one per line to <insults.txt> (UTF8 encoding)
2. go build mkdb.go
3. ./mkdb <insults.txt> <ernst.db> skymf
4. go run ernst.go

Insults are added with !skymf <insult> in the channel.

Number of insults are reported with !skymfstat(s).

Options are configured with !sset in the channel. For more configuration, poke around in ernst.go. You can use skymfer.txt instead of your own list if you want to be insulted in Swedish.

To (in reverse) create a text file with insults from a database:
1. go build dbdump.go
2. ./dbdump ernst.db skymf > <insults.txt>

## Thanks to
Thomas Jager (for go-ircevent)  
The BoltDB team  
\#ljusdal @ EFNet

## TODO
* default settings in mkdb
* Multiple channels
* Multiple servers
* Increased ratio on trigger words
* Increased ratio with time
* Search function in privmsg
* Requesting specific insults
	- per number or keyword
* PRIVMSG goroutine as func

## License
MIT:
Do whatever you want
