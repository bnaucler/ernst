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

const channel = "#kakapa";
const serverssl = "irc.inet.tele.dk:6697"
const fname = "./skymfer.txt"
const ircnick = "ernst7"
const ircuname = "ErnstHugo"

const rate = 10

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

func sskymf(irccon *irc.Connection, db *bolt.DB, cbuc []byte, rnd *rand.Rand,
	target string, numln int, mindel, maxdel int) bool {

	ln := rnd.Intn(numln)
	skymf, err := rdb(db, ln, cbuc)
	cherr(err)
	time.Sleep(time.Duration(rnd.Intn(maxdel) + mindel) * time.Millisecond)
	resp := fmt.Sprintf("%v: %v", target, skymf)
	irccon.Privmsg(channel, resp)

	return true
}

func main() {

    rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	cbuc := []byte("skymf")
	dbname := "./ernst.db"

	db, err := bolt.Open(dbname, 0640, nil)
	cherr(err)
	defer db.Close()

	mindel := 200
	maxdel := 5000

	addkey := "!skymf "
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

	// irccon.AddCallback("CTCP_VERSION", func(event *irc.Event) {
	// 	irccon.SendRawf("NOTICE %s :\x01VERSION %s\x01", event.Nick, "Skam och skuld")
	// })

	irccon.AddCallback("PRIVMSG", func(event *irc.Event) {
		go func(event *irc.Event) {

			lcnick := strings.ToLower(ircnick)
			lcstr := strings.ToLower(event.Arguments[1])

			if event.Arguments[0] == channel && strings.HasPrefix(event.Arguments[1], addkey) {

				skymf := strings.TrimPrefix(event.Arguments[1], addkey)
				err := wrdb(db, numln, skymf, cbuc)

				if err == nil {
					numln++
					err := wrdb(db, 0, strconv.Itoa(numln), cbuc)
					cherr(err)
					time.Sleep(time.Duration(rnd.Intn(maxdel) + mindel) * time.Millisecond)
					resp := fmt.Sprintf("%v: lade till \"%v\"", event.Nick, skymf)
					irccon.Privmsg(channel, resp)
				}

			} else if event.Arguments[0] == channel &&
				strings.HasPrefix(event.Arguments[1], statkey) {

				resp := fmt.Sprintf("Jag kan %d skymfer.", numln)
				time.Sleep(time.Duration(rnd.Intn(maxdel) + mindel) * time.Millisecond)
				irccon.Privmsg(channel, resp)

			} else if event.Arguments[0] == channel && rnd.Intn(1000) < rate &&
				event.Nick != ircnick {

				sskymf(irccon, db, cbuc, rnd, event.Nick, numln, mindel, maxdel)
			}

			if event.Arguments[0] == channel && strings.Contains(lcstr, lcnick) {
				sskymf(irccon, db, cbuc, rnd, event.Nick, numln, mindel, maxdel)
			}

			if event.Arguments[0] == ircnick {
				target := strings.Split(event.Arguments[1], " ")
				sskymf(irccon, db, cbuc, rnd, target[0], numln, mindel, maxdel)
			}

		}(event)
	});

	err = irccon.Connect(serverssl)
	cherr(err)

	irccon.Loop()
}
