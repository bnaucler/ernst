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
	"log"
	"time"
	"math/rand"
	"strings"
	"strconv"
	elib "github.com/bnaucler/ernst/elib"
)

const dbname = "./ernst.db"

func cherr(e error) { if e != nil { log.Fatal(e) } }

func askymf(db *bolt.DB, irccon *irc.Connection, event *irc.Event,
	rnd *rand.Rand, settings elib.Settings, skymf string) int {

	settings.Numln++
	senc, _ := json.Marshal(settings)
	err := elib.Wrdb(db, settings.Numln, []byte(skymf))

	if err == nil {
		err := elib.Wrdb(db, 0, senc)
		cherr(err)
		resp := fmt.Sprintf("%v: lade till \"%v\"", event.Nick, skymf)
		time.Sleep(time.Duration(len(resp) * settings.Kdel +
			rnd.Intn(settings.Randel)) * time.Millisecond)
		irccon.Privmsg(settings.Channel, resp)
	} else {
		cherr(err)
	}
	return settings.Numln
}

func sskymf(irccon *irc.Connection, db *bolt.DB, rnd *rand.Rand,
	target string, settings elib.Settings, lastsk []int, ln int) bool {

	if ln == 0 && settings.Dnrmem > 0 {
		inmem := true
		for inmem {
			ln = rnd.Intn(settings.Numln)
			inmem = false
			for a := 0; a < settings.Dnrmem; a++ {
				if ln == lastsk[a] {
					inmem = true
				}
			}
		}
	}

	skymf, err := elib.Rdb(db, ln)
	cherr(err)
	time.Sleep(time.Duration(len(skymf) *
		settings.Kdel + rnd.Intn(settings.Randel)) * time.Millisecond)

	resp := fmt.Sprintf("%v: %v", target, string(skymf))
	irccon.Privmsg(settings.Channel, resp)

	if settings.Dnrmem > 0 {
		for a := settings.Dnrmem - 1; a > 0; a-- { lastsk[a] = lastsk[(a - 1)] }
		lastsk[0] = ln
	}

	return true
}

func fskymf(irccon *irc.Connection, db *bolt.DB, rnd *rand.Rand,
	target string, kw []string, settings elib.Settings, lastsk []int) bool {

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
		cherr(err)
	}

	return sskymf(irccon, db, rnd, target, settings, lastsk, reqnum)
}

func csetlist(event *irc.Event, settings *elib.Settings) (resp string) {

	resp = fmt.Sprintf("%v: rate: %d/%d, kdel: %d/%d, randel: %d/%d dnrmem: %d/%d",
		event.Nick, settings.Rate, elib.Ratemax,
		settings.Kdel, elib.Kdelmax, settings.Randel, elib.Randelmax,
		settings.Dnrmem, elib.Dnrmemmax)

	return
}

func csetshow(event *irc.Event, settings *elib.Settings, setvar string) (resp string) {

	if setvar == "rate" {
		resp = fmt.Sprintf("%v: %v: %d/%d", event.Nick,
			setvar, settings.Rate, elib.Ratemax)
	} else if setvar == "kdel" {
		resp = fmt.Sprintf("%v: %v: %d/%d", event.Nick,
			setvar, settings.Kdel, elib.Kdelmax)
	} else if setvar == "randel" {
		resp = fmt.Sprintf("%v: %v: %d/%d", event.Nick,
			setvar, settings.Randel, elib.Randelmax)
	} else if setvar == "dnrmem" {
		resp = fmt.Sprintf("%v: %v: %d/%d", event.Nick,
			setvar, settings.Dnrmem, elib.Dnrmemmax)
	}

	return
}

func csetset(event *irc.Event, settings *elib.Settings, setvar,
	setval string) (resp string, dbchange bool) {

	nval, nerr := strconv.Atoi(setval)

	if setvar == "rate"  && nerr == nil && nval > -1 && nval <= elib.Ratemax {
		if settings.Rate != nval { dbchange = true }
		settings.Rate = nval
		resp = fmt.Sprintf("%s %d/%d.", setvar, settings.Rate, elib.Ratemax)

	} else if setvar == "kdel" && nerr == nil && nval > -1 && nval <= elib.Kdelmax {
		if settings.Kdel != nval { dbchange = true }
		settings.Kdel = nval
		resp = fmt.Sprintf("%s %d/%d.", setvar, settings.Kdel, elib.Kdelmax)

	} else if setvar == "randel" && nerr == nil && nval > -1 && nval <= elib.Randelmax {
		if settings.Randel != nval { dbchange = true }
		settings.Randel = nval
		resp = fmt.Sprintf("%s %d/%d.", setvar, settings.Randel, elib.Randelmax)

	} else if setvar == "dnrmem" && nerr == nil && nval > -1 && nval <= elib.Dnrmemmax {
		if settings.Dnrmem != nval { dbchange = true }
		settings.Dnrmem = nval
		resp = fmt.Sprintf("%s %d/%d.", setvar, settings.Dnrmem, elib.Dnrmemmax)
	}

	return
}

func cset(irccon *irc.Connection, db *bolt.DB, rnd *rand.Rand,
	event *irc.Event, settings *elib.Settings) bool {

	ssp := strings.Split(event.Arguments[1], " ")

	var (
		setvar, setval, resp string
		dbchange bool
	)

	if len(ssp) > 1 { setvar = strings.ToLower(ssp[1]) }
	if len(ssp) > 2 { setval = ssp[2] }

	if len(ssp) == 1 { resp = csetlist(event, settings)
	} else if len(ssp) == 2 { resp = csetshow(event, settings, setvar)
	} else if len(ssp) == 3 {
		resp, dbchange = csetset (event, settings, setvar, setval)
	}

	if len(resp) != 0 {
		if dbchange {
			s, err:= json.Marshal(settings)
			cherr(err)
			err = elib.Wrdb(db, 0, s)
			cherr(err)
		}
		irccon.Privmsg(settings.Channel, resp)
	}

	return true
}

func main() {

    rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	var incrt = int(0)

	settings := elib.Settings{}

	db, err := bolt.Open(dbname, 0640, nil)
	cherr(err)
	defer db.Close()

	addkey := "!skymf"
	setkey := "!sset"
	statkey := "!skymfstat"

	tmp, err := elib.Rdb(db, 0)
	cherr(err)
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

			if event.Arguments[0] == settings.Channel && strings.HasPrefix(lcstr, addkey) {

				keytr := fmt.Sprintf("%s ", addkey)
				skymf := strings.TrimPrefix(event.Arguments[1], keytr)
				settings.Numln = askymf(db, irccon, event, rnd, settings, skymf)

			} else if event.Arguments[0] == settings.Channel &&
				strings.HasPrefix(lcstr, statkey) {

				resp := fmt.Sprintf("Jag kan %d skymfer.", settings.Numln)
				time.Sleep(time.Duration(len(resp) * settings.Kdel +
					rnd.Intn(settings.Randel)) * time.Millisecond)
				irccon.Privmsg(settings.Channel, resp)

			} else if event.Arguments[0] == settings.Channel &&
				strings.HasPrefix(lcstr, setkey) {

				cset(irccon, db, rnd, event, &settings)

			} else if event.Arguments[0] == settings.Channel &&
				strings.Contains(lcstr, lcnick) {

				kw := strings.Split(event.Arguments[1], " ")

				if strings.Contains(kw[0], lcnick) {
					fskymf(irccon, db, rnd, event.Nick, kw[1:], settings, lastsk)
				} else {
					sskymf(irccon, db, rnd, event.Nick, settings, lastsk, 0)
				}
				incrt = 0

			} else if event.Arguments[0] == settings.Channel &&
				rnd.Intn(elib.Ratemax) < settings.Rate + incrt &&
				event.Nick != settings.Ircnick {

				sskymf(irccon, db, rnd, event.Nick, settings, lastsk, 0)
				incrt = 0
			}

			if event.Arguments[0] == settings.Ircnick {
				var nval int

				target := strings.Split(event.Arguments[1], " ")
				if len(target) > 1 {
					nval, err = strconv.Atoi(target[1])

					if err == nil && nval > 0 && nval <= settings.Numln {
						sskymf(irccon, db, rnd, target[0], settings, lastsk, nval)
						incrt = 0
					} else if nval == 0 {
						fskymf(irccon, db, rnd, target[0], target[1:], settings, lastsk)
						incrt = 0
					}

				} else {

					sskymf(irccon, db, rnd, target[0], settings, lastsk, 0)
					incrt = 0
				}
			}

			incrt++
		}(event)
	});

	err = irccon.Connect(settings.Server)
	cherr(err)

	irccon.Loop()
}
