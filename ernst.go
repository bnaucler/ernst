/*

		ernst.go
		A very unfriendly IRC bot

*/

package main

import (
	"github.com/thoj/go-ircevent"
	"crypto/tls"
	"bufio"
	"fmt"
	"time"
	"math/rand"
	"os"
	"strings"
)

const channel = "#ljusdal";
const serverssl = "irc.inet.tele.dk:6697"
const fname = "./skymfer.txt"
const ircnick = "ernst"
const ircuname = "ErnstHugo"

const rate = 10

func cherr(e error) { if e != nil { panic(e) } }

func getskymf(f *os.File, rnd *rand.Rand, numln int) (skymf string) {

	f.Seek(0, 0)
	randln := rnd.Intn(numln)

	scanner := bufio.NewScanner(f)
	for a := 0; a < randln; a++ {
		scanner.Scan()
	}
	skymf = scanner.Text()

	return
}

func wrskymf(f *os.File, rnd *rand.Rand, skymf string, mindel, maxdel int) bool {

	skymf = fmt.Sprintf("%v\n", skymf)
	_, err := f.WriteString(skymf)
	cherr(err)
	time.Sleep(time.Duration(rnd.Intn(maxdel) + mindel) * time.Millisecond)
	return true
}

func clines(f *os.File) (lines int) {

	scanner := bufio.NewScanner(f)
	for scanner.Scan() { lines++ }
	return
}

func sskymf(irccon *irc.Connection, f *os.File, numln int, channel,
	target string, rnd *rand.Rand, mindel, maxdel int) bool {

	skymf := fmt.Sprintf("%v: %v", target, getskymf(f, rnd, numln))
	time.Sleep(time.Duration(rnd.Intn(maxdel) + mindel) * time.Millisecond)
	irccon.Privmsg(channel, skymf)
	return true
}

func main() {

    rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	mindel := 200
	maxdel := 5000

	addkey := "!skymf "
	statkey := "!skymfstat"

	f, err := os.OpenFile(fname, os.O_APPEND|os.O_RDWR, 0644)
	cherr(err)
	defer f.Close()

	numln := clines(f)

	irccon := irc.IRC(ircnick, ircuname)

	irccon.VerboseCallbackHandler = true
	irccon.Debug = true
	irccon.UseTLS = true
	irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	irccon.AddCallback("001", func(e *irc.Event) { irccon.Join(channel) })

	irccon.AddCallback("CTCP_VERSION", func(event *irc.Event) {
		irccon.SendRawf("NOTICE %s :\x01VERSION %s\x01", event.Nick, "Skam och skuld")
	})

	irccon.AddCallback("PRIVMSG", func(event *irc.Event) {
		go func(event *irc.Event) {

			lcnick := strings.ToLower(ircnick)
			lcstr := strings.ToLower(event.Arguments[1])

			if event.Arguments[0] == channel && strings.HasPrefix(event.Arguments[1], addkey) {
				skymf := strings.TrimPrefix(event.Arguments[1], addkey)
				if wrskymf(f, rnd, skymf, mindel, maxdel) {
					numln++
					resp := fmt.Sprintf("%v: lade till \"%v\"", event.Nick, skymf)
					irccon.Privmsg(channel, resp)
				}

			} else if event.Arguments[0] == channel && strings.HasPrefix(event.Arguments[1], statkey) {
				resp := fmt.Sprintf("Jag kan %d skymfer.", numln)
				time.Sleep(time.Duration(rnd.Intn(maxdel) + mindel) * time.Millisecond)
				irccon.Privmsg(channel, resp)

			} else if event.Arguments[0] == channel && rnd.Intn(1000) < rate {
				sskymf(irccon, f, numln, channel, event.Nick, rnd, mindel, maxdel)
			}

			if event.Arguments[0] == channel && strings.Contains(lcstr, lcnick) {
				sskymf(irccon, f, numln, channel, event.Nick, rnd, mindel, maxdel)
			}

			if event.Arguments[0] == ircnick {
				target := strings.Split(event.Arguments[1], " ")
				sskymf(irccon, f, numln, channel, target[0], rnd, mindel, maxdel)
			}

		}(event)
	});

	err = irccon.Connect(serverssl)
	cherr(err)

	irccon.Loop()
}
