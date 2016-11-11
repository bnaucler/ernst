package main

import (
	"github.com/thoj/go-ircevent"
	"crypto/tls"
	"bufio"
	"time"
	"math/rand"
	"os"
)

const channel = "#test";
const serverssl = "irc.apansson.se:6697"
const fname = "./skymfer.txt"

func cherr(e error) { if e != nil { panic(e) } }

func getskymf(f *os.File, r1 *rand.Rand, numln int) (skymf string) {

	f.Seek(0, 0)
	randln := r1.Intn(numln)

	scanner := bufio.NewScanner(f)
	for a := 0; a < randln; a++ {
		scanner.Scan()
		skymf = scanner.Text()
	}

	return
}

func clines(f *os.File) (lines int) {

	scanner := bufio.NewScanner(f)
	for scanner.Scan() { lines++ }
	return
}

func main() {

	s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)

	f, err := os.Open(fname)
	cherr(err)
	defer f.Close()

	numln := clines(f)

	ircnick := "ernst"
	ircuname := "elak"
	irccon := irc.IRC(ircnick, ircuname)

	irccon.VerboseCallbackHandler = true
	irccon.Debug = true
	irccon.UseTLS = true
	irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	irccon.AddCallback("001", func(e *irc.Event) { irccon.Join(channel) })
	// irccon.AddCallback("366", func(e *irc.Event) {  })
	irccon.AddCallback("PRIVMSG", func(event *irc.Event) {
		skymf := getskymf(f, r1, numln)
		irccon.Privmsg(event.Nick, skymf)
	});
	err = irccon.Connect(serverssl)
	cherr(err)

	irccon.Loop()
}
