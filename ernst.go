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

const ratemax = int(1000)
const kdelmax = int(1000)
const randelmax = int(10000)

func cherr(e error) { if e != nil { log.Fatal(e) } }

func askymf(db *bolt.DB, irccon *irc.Connection, event *irc.Event,
	rnd *rand.Rand, settings elib.Settings, skymf string, cbuc []byte) int {

	settings.Numln++
	senc, _ := json.Marshal(settings)
	err := elib.Wrdb(db, settings.Numln, []byte(skymf), cbuc)

	if err == nil {
		err := elib.Wrdb(db, 0, senc, cbuc)
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

func sskymf(irccon *irc.Connection, db *bolt.DB, cbuc []byte, rnd *rand.Rand,
	target string, settings elib.Settings, ln int) bool {

	if ln == 0 { ln = rnd.Intn(settings.Numln) }

	skymf, err := elib.Rdb(db, ln, cbuc)
	cherr(err)
	time.Sleep(time.Duration(len(skymf) *
		settings.Kdel + rnd.Intn(settings.Randel)) * time.Millisecond)
	resp := fmt.Sprintf("%v: %v", target, string(skymf))
	irccon.Privmsg(settings.Channel, resp)

	return true
}

func fskymf(irccon *irc.Connection, db *bolt.DB, cbuc []byte, rnd *rand.Rand,
	target string, kw []string, index int, settings elib.Settings) bool {

	kwln := len(kw)
	var (reqnum, cqual, tqual int)

	for k := 0; k <= settings.Numln; k++ {
		v, err := elib.Rdb(db, k, cbuc)
		cqual = 0
		cstr := string(v)
		for a := index; a < kwln; a++ {
			if strings.Contains(cstr, strings.ToLower(kw[a])) {
				cqual++
				if cqual > tqual {
					tqual = cqual
					reqnum = k
				}
			}
		}
		cherr(err)
	}

	return sskymf(irccon, db, cbuc, rnd, target, settings, reqnum)
}

func csetlist(event *irc.Event, settings *elib.Settings) (resp string) {

	resp = fmt.Sprintf("%v: rate: %d/%d, kdel: %d/%d, randel: %d/%d",
		event.Nick, settings.Rate, ratemax,
		settings.Kdel, kdelmax, settings.Randel, randelmax)

	return
}

func csetshow(event *irc.Event, settings *elib.Settings, setvar string) (resp string) {

	if setvar == "rate" {
		resp = fmt.Sprintf("%v: %v: %d/%d", event.Nick, setvar, settings.Rate, ratemax)
	} else if setvar == "kdel" {
		resp = fmt.Sprintf("%v: %v: %d/%d", event.Nick, setvar, settings.Kdel, kdelmax)
	} else if setvar == "randel" {
		resp = fmt.Sprintf("%v: %v: %d/%d", event.Nick, setvar, settings.Randel, randelmax)
	}

	return
}

func csetset(event *irc.Event, settings *elib.Settings, setvar,
	setval string) (resp string, dbchange bool) {

	nval, nerr := strconv.Atoi(setval)

	if setvar == "rate"  && nerr == nil && nval > -1 && nval <= ratemax {
		if settings.Rate != nval { dbchange = true }
		settings.Rate = nval
		resp = fmt.Sprintf("%s %d/%d.", setvar, settings.Rate, ratemax)

	} else if setvar == "kdel" && nerr == nil && nval > -1 && nval <= kdelmax {
		if settings.Kdel != nval { dbchange = true }
		settings.Kdel = nval
		resp = fmt.Sprintf("%s %d/%d.", setvar, settings.Kdel, kdelmax)

	} else if setvar == "randel" && nerr == nil && nval > -1 && nval <= randelmax {
		if settings.Randel != nval { dbchange = true }
		settings.Randel = nval
		resp = fmt.Sprintf("%s %d/%d.", setvar, settings.Randel, randelmax)
	}

	return
}

func cset(irccon *irc.Connection, db *bolt.DB, cbuc []byte, rnd *rand.Rand,
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
			err = elib.Wrdb(db, 0, s, cbuc)
			cherr(err)
		}
		irccon.Privmsg(settings.Channel, resp)
	}

	return true
}

func main() {

    rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cbuc = []byte("skymf")
	var incrt = int(0)

	settings := elib.Settings{}

	db, err := bolt.Open(dbname, 0640, nil)
	cherr(err)
	defer db.Close()

	addkey := "!skymf "
	setkey := "!sset"
	statkey := "!skymfstat"

	tmp, err := elib.Rdb(db, 0, cbuc)
	cherr(err)
	json.Unmarshal(tmp, &settings)

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
				skymf := strings.TrimPrefix(event.Arguments[1], addkey)
				settings.Numln = askymf(db, irccon, event, rnd, settings, skymf, cbuc)

			} else if event.Arguments[0] == settings.Channel &&
				strings.HasPrefix(lcstr, statkey) {

				resp := fmt.Sprintf("Jag kan %d skymfer.", settings.Numln)
				time.Sleep(time.Duration(len(resp) * settings.Kdel +
					rnd.Intn(settings.Randel)) * time.Millisecond)
				irccon.Privmsg(settings.Channel, resp)

			} else if event.Arguments[0] == settings.Channel &&
				strings.HasPrefix(lcstr, setkey) {
				cset(irccon, db, cbuc, rnd, event, &settings)

			} else if event.Arguments[0] == settings.Channel &&
				strings.Contains(lcstr, lcnick) {
				kw := strings.Split(event.Arguments[1], " ")

				if strings.Contains(kw[0], lcnick) {
					fskymf(irccon, db, cbuc, rnd, event.Nick, kw, 1, settings)
				} else {
					sskymf(irccon, db, cbuc, rnd, event.Nick, settings, 0)
				}
				incrt = 0

			} else if event.Arguments[0] == settings.Channel &&
				rnd.Intn(ratemax) < settings.Rate + incrt && event.Nick != settings.Ircnick {
				sskymf(irccon, db, cbuc, rnd, event.Nick, settings, 0)
				incrt = 0
			}

			if event.Arguments[0] == settings.Ircnick {
				var nval int

				target := strings.Split(event.Arguments[1], " ")
				if len(target) > 1 {
					nval, err = strconv.Atoi(target[1])

					debugresp := fmt.Sprintf("%+v, %d, %d", target, nval, len(target))
					irccon.Privmsg(settings.Channel, debugresp)

					if err == nil && nval > 0 && nval <= settings.Numln {
						sskymf(irccon, db, cbuc, rnd, target[0], settings, nval)
						incrt = 0
					} else if nval == 0 {
						fskymf(irccon, db, cbuc, rnd, target[0], target, 1, settings)
						incrt = 0
					}

				} else {
					sskymf(irccon, db, cbuc, rnd, target[0], settings, 0)
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
