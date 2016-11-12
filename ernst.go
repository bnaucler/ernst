package main

import (
	"github.com/thoj/go-ircevent"
	"crypto/tls"
	"bufio"
	"fmt"
	"time"
	"math/rand"
	"os"
)

const channel = "#ljusdal";
const serverssl = "irc.inet.tele.dk:6697"
const fname = "./skymfer.txt"
const ircnick = "ernst"
const ircuname = "ErnstHugo"

const rate = 10

func cherr(e error) { if e != nil { panic(e) } }

func getskymf(f *os.File, r1 *rand.Rand, numln int) (skymf string) {

	f.Seek(0, 0)
	randln := r1.Intn(numln)

	scanner := bufio.NewScanner(f)
	for a := 0; a < randln; a++ {
		scanner.Scan()
	}
	skymf = scanner.Text()

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

	irccon := irc.IRC(ircnick, ircuname)

	irccon.VerboseCallbackHandler = true
	irccon.Debug = true
	irccon.UseTLS = true
	irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	irccon.AddCallback("001", func(e *irc.Event) { irccon.Join(channel) })
	// irccon.AddCallback("366", func(e *irc.Event) {  })

	irccon.AddCallback("PRIVMSG", func(event *irc.Event) {
		go func(event *irc.Event) {

			if event.Arguments[0] == channel {
				if r1.Intn(1000) < rate {
					skymf := fmt.Sprintf("%v: %v", event.Nick, getskymf(f, r1, numln))
					time.Sleep(time.Duration(r1.Intn(5000) + 200) * time.Millisecond)
					irccon.Privmsg(channel, skymf)
				}
			}

			// UNTESTED
			if event.Arguments[0] == ircnick {
				skymf := fmt.Sprintf("%v: %v", event.Arguments[1], getskymf(f, r1, numln))
				irccon.Privmsg(channel, skymf)
			}

		}(event)
	});

	err = irccon.Connect(serverssl)
	cherr(err)

	irccon.Loop()
}
