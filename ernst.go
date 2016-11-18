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
)

const dbname =		"./ernst.db"

const ratemax = int(1000)
const kdelmax = int(1000)
const randelmax = int(10000)

type Settings struct {
	Numln		int
	Rate		int
	Ircnick		string
	Uname		string
	Channel		string
	Server		string
	// Tword		[]string
	Kdel		int
	Randel		int
}

func cherr(e error) { if e != nil { log.Fatal(e) } }

func rdb(db *bolt.DB, k int, cbuc []byte) ([]byte, error) {

	var v []byte

	err := db.View(func(tx *bolt.Tx) error {
		buc := tx.Bucket(cbuc)
		if buc == nil { return fmt.Errorf("No bucket!") }

		v = buc.Get([]byte(strconv.Itoa(k)))
		return nil
	})
	return v, err
}

func wrdb(db *bolt.DB, k int, v string, cbuc []byte) (err error) {

	err = db.Update(func(tx *bolt.Tx) error {
		buc, err := tx.CreateBucketIfNotExists(cbuc)
		if err != nil { return err }

		err = buc.Put([]byte(strconv.Itoa(k)), []byte(v))
		if err != nil { return err }

		return nil
	})
	return
}

func askymf(db *bolt.DB, irccon *irc.Connection, event *irc.Event,
	rnd *rand.Rand, settings Settings, skymf string, cbuc []byte) int {

	settings.Numln++
	senc, _ := json.Marshal(settings)
	err := wrdb(db, settings.Numln, skymf, cbuc)

	if err == nil {
		err := wrdb(db, 0, string(senc), cbuc)
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
	target string, settings Settings, ln int) bool {

	if ln == 0 { ln = rnd.Intn(settings.Numln) }

	skymf, err := rdb(db, ln, cbuc)
	cherr(err)
	time.Sleep(time.Duration(len(skymf) *
		settings.Kdel + rnd.Intn(settings.Randel)) * time.Millisecond)
	resp := fmt.Sprintf("%v: %v", target, string(skymf))
	irccon.Privmsg(settings.Channel, resp)

	return true
}

func cset(irccon *irc.Connection, db *bolt.DB, cbuc []byte, rnd *rand.Rand,
	event *irc.Event, settings *Settings) bool {

	ssp := strings.Split(event.Arguments[1], " ")

	var (setvar, setval, resp string)

	if len(ssp) > 1 { setvar = strings.ToLower(ssp[1]) }
	if len(ssp) > 2 { setval = ssp[2] }

	if len(setval) == 0 {
		if setvar == "rate" {
			resp = fmt.Sprintf("%v: %v: %d/%d", event.Nick, setvar, settings.Rate, ratemax)
		} else if setvar == "kdel" {
			resp = fmt.Sprintf("%v: %v: %d/%d", event.Nick, setvar, settings.Kdel, kdelmax)
		} else if setvar == "randel" {
			resp = fmt.Sprintf("%v: %v: %d/%d", event.Nick, setvar, settings.Randel, randelmax)
		}
		if len(resp) != 0 { irccon.Privmsg(settings.Channel, resp) }

	} else {

		nval, nerr := strconv.Atoi(setval)

		if setvar == "rate"  && nerr == nil && nval > -1 && nval <= ratemax {
			settings.Rate = nval
			resp = fmt.Sprintf("Nu %d/%d.", settings.Rate, ratemax)

		} else if setvar == "kdel" && nerr == nil && nval > -1 && nval <= kdelmax {
			settings.Kdel = nval
			resp = fmt.Sprintf("Nu %d/%d.", settings.Kdel, kdelmax)

		} else if setvar == "randel" && nerr == nil && nval > -1 && nval <= randelmax {
			settings.Randel = nval
			resp = fmt.Sprintf("Nu %d/%d.", settings.Randel, randelmax)

		}

		if len(resp) != 0 {
			s, err:= json.Marshal(settings)
			cherr(err)
			err = wrdb(db, 0, string(s), cbuc)
			cherr(err)
			irccon.Privmsg(settings.Channel, resp)
		}
	}

	return true
}

func main() {

    rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cbuc = []byte("skymf")
	settings := Settings{}

	db, err := bolt.Open(dbname, 0640, nil)
	cherr(err)
	defer db.Close()

	addkey := "!skymf "
	setkey := "!sset "
	statkey := "!skymfstat"

	tmp, err := rdb(db, 0, cbuc)
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

			} else if event.Arguments[0] == settings.Channel && rnd.Intn(1000) < settings.Rate &&
				event.Nick != settings.Ircnick {
				sskymf(irccon, db, cbuc, rnd, event.Nick, settings, 0)
			}

			if event.Arguments[0] == settings.Channel && strings.Contains(lcstr, lcnick) {
				sskymf(irccon, db, cbuc, rnd, event.Nick, settings, 0)
			}

			if event.Arguments[0] == settings.Ircnick {
				var nval int

				target := strings.Split(event.Arguments[1], " ")
				if len(target) > 1 {
					nval, err = strconv.Atoi(target[1])
					if err == nil && nval > 0 && nval <= settings.Numln {
						sskymf(irccon, db, cbuc, rnd, target[0], settings, nval)
					}

				} else {
					sskymf(irccon, db, cbuc, rnd, target[0], settings, 0)
				}
			}

		}(event)
	});

	err = irccon.Connect(settings.Server)
	cherr(err)

	irccon.Loop()
}
