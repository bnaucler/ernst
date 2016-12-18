/*

		ernst.go
		A very unfriendly IRC bot

*/

package main

import (
	"github.com/thoj/go-ircevent"
	"github.com/boltdb/bolt"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"time"
	"math/rand"
	"strings"
	"strconv"
	elib "github.com/bnaucler/ernst/elib"
)

const dbname = "./ernst.db"
const erresp = "Vad fan är det frågan om?"

func askymf(db *bolt.DB, irccon *irc.Connection, event *irc.Event,
	rnd *rand.Rand, settings elib.Settings, skymf string) int {

	settings.Numln++
	senc, _ := json.Marshal(settings)
	err := elib.Wrdb(db, settings.Numln, []byte(skymf))

	if err == nil {
		err := elib.Wrdb(db, 0, senc)
		elib.Cherr(err)
		resp := fmt.Sprintf("%v: lade till \"%v\"", event.Nick, skymf)
		time.Sleep(time.Duration(len(resp) * settings.Kdel +
			rnd.Intn(settings.Randel)) * time.Millisecond)
		irccon.Privmsg(settings.Channel, resp)
	} else {
		elib.Cherr(err)
	}
	return settings.Numln
}

func rskymf(irccon *irc.Connection, db *bolt.DB, event *irc.Event,
	rnd *rand.Rand, settings *elib.Settings, lastsk []int) []int {

	resp := "Kunde ej ta bort skymf"
	dbch := false
	var rsk string

	bar, err := elib.Rdb(db, lastsk[0])
	rsk = string(bar)

	if lastsk[0] != settings.Numln && lastsk[0] != 0 && err == nil {
		tmps, err := elib.Rdb(db, settings.Numln)

		if err == nil {
			err = elib.Wrdb(db, lastsk[0], tmps)
			dbch = true
		}
	} else if lastsk[0] == settings.Numln {
		dbch = true
	}

	if dbch {
		settings.Numln--;
		senc, _ := json.Marshal(settings)
		err = elib.Wrdb(db, 0, senc)
		resp = fmt.Sprintf("Skymf \"%v\" borttagen.", rsk)

		for a := 0; a < elib.Dnrmemmax - 1; a++ { lastsk[a] = lastsk[(a + 1)] }
	}

	time.Sleep(time.Duration(len(resp) * settings.Kdel +
		rnd.Intn(settings.Randel)) * time.Millisecond)
	irccon.Privmsg(settings.Channel, resp)

	return lastsk
}

func sskymf(irccon *irc.Connection, db *bolt.DB, rnd *rand.Rand,
	target string, settings elib.Settings, lastsk []int, ln int) []int {

	if ln == 0 && settings.Dnrmem > 0 {
		inmem := true
		for inmem {
			ln = rnd.Intn(settings.Numln - 1) + 1
			inmem = false
			for a := 0; a < settings.Dnrmem; a++ {
				if ln == lastsk[a] {
					inmem = true
				}
			}
		}
	}

	skymf, err := elib.Rdb(db, ln)
	elib.Cherr(err)
	time.Sleep(time.Duration(len(skymf) *
		settings.Kdel + rnd.Intn(settings.Randel)) * time.Millisecond)

	resp := fmt.Sprintf("%v: %v", target, string(skymf))
	irccon.Privmsg(settings.Channel, resp)

	if settings.Dnrmem > 0 {
		for a := settings.Dnrmem - 1; a > 0; a-- { lastsk[a] = lastsk[(a - 1)] }
		lastsk[0] = ln
	}

	return lastsk
}

func fskymf(irccon *irc.Connection, db *bolt.DB, rnd *rand.Rand,
	target string, kw []string, settings elib.Settings, lastsk []int) []int {

	kwln := len(kw)
	var (reqnum, cqual, tqual int)
	lckw := make([]string, kwln)

	for a := 0; a < kwln; a++ { lckw[a] = strings.ToLower(kw[a]) }

	for k := 0; k <= settings.Numln; k++ {
		v, err := elib.Rdb(db, k)
		cqual = 0
		cstr := string(v)
		for a := 0; a < kwln; a++ {
			if strings.Contains(cstr, lckw[a]) {
				cqual++
				if cqual > tqual {
					tqual = cqual
					reqnum = k
				}
			}
		}
		elib.Cherr(err)
	}

	return sskymf(irccon, db, rnd, target, settings, lastsk, reqnum)
}

func main() {

    rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	settings := elib.Settings{}
	incrt := 0

	db, err := bolt.Open(dbname, 0640, nil)
	elib.Cherr(err)
	defer db.Close()

	tmp, err := elib.Rdb(db, 0)
	elib.Cherr(err)
	json.Unmarshal(tmp, &settings)

	var lastsk = make([]int, elib.Dnrmemmax)

	irccon := irc.IRC(settings.Ircnick, settings.Uname)

	irccon.VerboseCallbackHandler = true
	irccon.Debug = true
	irccon.UseTLS = true
	irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	irccon.AddCallback("001", func(e *irc.Event) { irccon.Join(settings.Channel) })

	irccon.AddCallback("PRIVMSG", func(event *irc.Event) {
		go func(event *irc.Event) {

			lcnick := strings.ToLower(settings.Ircnick)
			lcstr := strings.ToLower(event.Arguments[1])
			keytr := fmt.Sprintf("%s ", elib.Addkey)

			if event.Arguments[0] == settings.Channel && strings.HasPrefix(lcstr, keytr) {

				keytr := fmt.Sprintf("%s ", elib.Addkey)
				skymf := strings.TrimPrefix(event.Arguments[1], keytr)
				settings.Numln = askymf(db, irccon, event, rnd, settings, skymf)

			} else if event.Arguments[0] == settings.Channel &&
				strings.HasPrefix(lcstr, elib.Statkey) {

				resp := fmt.Sprintf("Jag kan %d skymfer.", settings.Numln)
				time.Sleep(time.Duration(len(resp) * settings.Kdel +
					rnd.Intn(settings.Randel)) * time.Millisecond)
				irccon.Privmsg(settings.Channel, resp)

			} else if event.Arguments[0] == settings.Channel &&
				strings.HasPrefix(lcstr, elib.Setkey) {

				cset(irccon, db, rnd, event, &settings)

			} else if event.Arguments[0] == settings.Channel &&
				strings.HasPrefix(lcstr, elib.Rmkey) {

				lastsk = rskymf(irccon, db, event, rnd, &settings, lastsk)

			} else if event.Arguments[0] == settings.Channel &&
				strings.Contains(lcstr, lcnick) {

				kw := strings.Split(event.Arguments[1], " ")

				if strings.Contains(kw[0], lcnick) {
					lastsk = fskymf(irccon, db, rnd, event.Nick, kw[1:], settings, lastsk)
				} else {
					lastsk = sskymf(irccon, db, rnd, event.Nick, settings, lastsk, 0)
				}
				incrt = 0

			} else if event.Arguments[0] == settings.Channel &&
				rnd.Intn(elib.Ratemax) < settings.Rate + incrt &&
				event.Nick != settings.Ircnick {

				lastsk = sskymf(irccon, db, rnd, event.Nick, settings, lastsk, 0)
				incrt = 0

			} else if event.Arguments[0] == settings.Ircnick {
				var nval int

				target := strings.Split(event.Arguments[1], " ")
				if len(target) > 1 {
					nval, err = strconv.Atoi(target[1])

					if err == nil && nval > 0 && nval <= settings.Numln {
						lastsk = sskymf(irccon, db, rnd, target[0], settings, lastsk, nval)
						incrt = 0
					} else if nval == 0 {
						lastsk = fskymf(irccon, db, rnd, target[0], target[1:], settings, lastsk)
						incrt = 0
					}

				} else {

					lastsk = sskymf(irccon, db, rnd, target[0], settings, lastsk, 0)
					incrt = 0
				}
			} else {

				incrt++
			}

		}(event)
	});

	err = irccon.Connect(settings.Server)
	elib.Cherr(err)

	irccon.Loop()
}
