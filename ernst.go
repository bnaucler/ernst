/*

		ernst.go
		A very unfriendly IRC bot

*/

package main

import (
	"github.com/thoj/go-ircevent"
	"github.com/boltdb/bolt"
	"crypto/tls"
	"fmt"
	"log"
	"time"
	"math/rand"
	"strings"
	"strconv"
)

const channel =		"#kakapa";
const serverssl =	"irc.inet.tele.dk:6697"
const dbname =		"./ernst.db"
const ircnick =		"ernst7"
const ircuname =	"ErnstHugo"
const rate =		10

func cherr(e error) { if e != nil { log.Fatal(e) } }

func rdb(db *bolt.DB, k int, cbuc []byte) (string, error) {

	var v []byte

	err := db.View(func(tx *bolt.Tx) error {
		buc := tx.Bucket(cbuc)
		if buc == nil { return fmt.Errorf("No bucket!") }

		v = buc.Get([]byte(strconv.Itoa(k)))
		return nil
	})
	return string(v), err
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
	rnd *rand.Rand, numln, kdel, randel int, skymf string, cbuc []byte) int {

	err := wrdb(db, numln, skymf, cbuc)

	if err == nil {
		numln++
		err := wrdb(db, 0, strconv.Itoa(numln), cbuc)
		cherr(err)
		resp := fmt.Sprintf("%v: lade till \"%v\"", event.Nick, skymf)
		time.Sleep(time.Duration(len(resp) * kdel +
			rnd.Intn(randel)) * time.Millisecond)
		irccon.Privmsg(channel, resp)
	} else {
		cherr(err)
	}
	return numln
}

func sskymf(irccon *irc.Connection, db *bolt.DB, cbuc []byte, rnd *rand.Rand,
	target string, numln int, kdel, randel int) bool {

	ln := rnd.Intn(numln)
	skymf, err := rdb(db, ln, cbuc)
	cherr(err)
	time.Sleep(time.Duration(len(skymf) * kdel + rnd.Intn(randel)) * time.Millisecond)
	resp := fmt.Sprintf("%v: %v", target, skymf)
	irccon.Privmsg(channel, resp)

	return true
}

func main() {

    rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cbuc = []byte("skymf")

	db, err := bolt.Open(dbname, 0640, nil)
	cherr(err)
	defer db.Close()

	kdel := 100
	randel := 700

	addkey := "!skymf "
	// setkey := "!sset "
	statkey := "!skymfstat"

	tmp, err := rdb(db, 0, cbuc)
	cherr(err)
	numln, err:= strconv.Atoi(tmp)
	cherr(err)

	fmt.Printf("%v, %T\n", numln, numln)
	irccon := irc.IRC(ircnick, ircuname)

	irccon.VerboseCallbackHandler = true
	irccon.Debug = true
	irccon.UseTLS = true
	irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	irccon.AddCallback("001", func(e *irc.Event) { irccon.Join(channel) })

	irccon.AddCallback("PRIVMSG", func(event *irc.Event) {
		go func(event *irc.Event) {

			lcnick := strings.ToLower(ircnick)
			lcstr := strings.ToLower(event.Arguments[1])

			if event.Arguments[0] == channel && strings.HasPrefix(lcstr, addkey) {
				skymf := strings.TrimPrefix(event.Arguments[1], addkey)
				numln = askymf(db, irccon, event, rnd, numln, kdel, randel, skymf, cbuc)

			} else if event.Arguments[0] == channel &&
				strings.HasPrefix(lcstr, statkey) {

				resp := fmt.Sprintf("Jag kan %d skymfer.", numln)
				time.Sleep(time.Duration(len(resp) * kdel +
					rnd.Intn(randel)) * time.Millisecond)
				irccon.Privmsg(channel, resp)

			} else if event.Arguments[0] == channel && rnd.Intn(1000) < rate &&
				event.Nick != ircnick {
				sskymf(irccon, db, cbuc, rnd, event.Nick, numln, kdel, randel)
			}

			if event.Arguments[0] == channel && strings.Contains(lcstr, lcnick) {
				sskymf(irccon, db, cbuc, rnd, event.Nick, numln, kdel, randel)
			}

			if event.Arguments[0] == ircnick {
				target := strings.Split(event.Arguments[1], " ")
				sskymf(irccon, db, cbuc, rnd, target[0], numln, kdel, randel)
			}

		}(event)
	});

	err = irccon.Connect(serverssl)
	cherr(err)

	irccon.Loop()
}
