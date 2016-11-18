/*

		mkdb.go
		Create a ernst database from text file.

*/

package main

import (
	"fmt"
	"os"
	"bufio"
	"unicode"
	"strconv"
	"encoding/json"
	"github.com/boltdb/bolt"
)

type Settings struct {
	Numln		int
	Rate		int
	Ircnick		string
	Uname		string
	Channel		string
	Server		string
	// tword		[]string
	Kdel		int
	Randel		int
}

func cherr(e error) { if e != nil { panic(e) } }

func clines(f *os.File) (lines int) {

	scanner := bufio.NewScanner(f)
	for scanner.Scan() { lines++ }
	return
}

func wrdb(db *bolt.DB, k, v, cbuc []byte) (err error) {

	err = db.Update(func(tx *bolt.Tx) error {
		buc, err := tx.CreateBucketIfNotExists(cbuc)
		if err != nil { return err }

		err = buc.Put(k, v)
		if err != nil { return err }

		return nil
	})
	return
}

func gline(f *os.File, scanner *bufio.Scanner, l int64) (string, int64) {

	_, err := f.Seek(int64(l), 0)
	cherr(err)

	scanner.Scan()

	pos, err := f.Seek(0, 1)
	cherr(err)

	return scanner.Text(), pos
}

func rval(prompt string, minval, maxval int) (chint int) {

	var tmp string
	fmt.Printf("%v: ", prompt)
	fmt.Scanln(&tmp)
	for a := 0; a < len(tmp); a++ {
		r := rune(tmp[a])
		if !unicode.IsDigit(r) {
			panic("Not a number")
		}
	}

	chint, err := strconv.Atoi(tmp)
	if chint < minval || chint > maxval {
		resp := fmt.Sprintf("Number not in %d-%d range", minval, maxval)
		panic(resp)
	}

	cherr(err)
	return
}

func rtext (prompt string) (resp string) {

	fmt.Printf("%s: ", prompt)
	fmt.Scanln(&resp)
	return
}

func main() {

	var verb bool
	settings := Settings{}

	if len(os.Args) < 4 || len(os.Args) > 5 {
		cherr(fmt.Errorf("Usage: %s <file.txt> <file.db> <bucket> [v]\n", os.Args[0]))
	}

	sfname := os.Args[1]
	dbname := os.Args[2]
	cbuc := []byte(os.Args[3])
	if len(os.Args) == 5 && os.Args[4] == "v" { verb = true }

	db, err := bolt.Open(dbname, 0640, nil)
	cherr(err)
	defer db.Close()

	f, err := os.Open(sfname)
	cherr(err)
	defer f.Close()

	settings.Numln = clines(f)

	scanner := bufio.NewScanner(f)
	f.Seek(0, 0)
	var pos = int64(0)
	var v = string("")

	for k := 0; k < settings.Numln; k++ {
		v, pos = gline(f, scanner, pos)
		err = wrdb(db, []byte(strconv.Itoa(k+1)), []byte(v), cbuc)
		if verb { fmt.Printf("%d(%d): %v\n", k + 1, pos, v) }
		cherr(err)
	}

	settings.Ircnick = rtext("Bot nick")
	settings.Channel = rtext("Channel")
	settings.Server = rtext("server:port (SSL only!)")
	settings.Rate = rval("Rate (0-1000)", 0, 1000)
	settings.Kdel = rval("Keystroke delay in ms. (0-1000)", 0, 1000)
	settings.Randel = rval("Random delay in ms. (0-10000)", 0, 10000)

	s, err:= json.Marshal(settings)
	cherr(err)
	err = wrdb(db, []byte("0"), s, cbuc)
	cherr(err)
}
