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
	elib "github.com/bnaucler/ernst/elib"
)

func cherr(e error) { if e != nil { panic(e) } }

func clines(f *os.File) (lines int) {

	scanner := bufio.NewScanner(f)
	for scanner.Scan() { lines++ }
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

func rval(prompt string, minval, maxval, def int) int {

	var tmp string
	var chint int

	fmt.Printf("%v [%d]: ", prompt, def)
	fmt.Scanln(&tmp)

	if len(tmp) == 0 { return def }

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
	return chint
}

func rtext (prompt string) (resp string) {

	fmt.Printf("%s: ", prompt)
	fmt.Scanln(&resp)
	return
}

func main() {

	var verb bool
	settings := elib.Settings{}

	if len(os.Args) < 3 || len(os.Args) > 4 {
		cherr(fmt.Errorf("Usage: %s <file.txt> <file.db> [v]\n", os.Args[0]))
	}

	sfname := os.Args[1]
	dbname := os.Args[2]
	if len(os.Args) == 4 && os.Args[3] == "v" { verb = true }

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
		err = elib.Wrdb(db, (k+1), []byte(v))
		if verb { fmt.Printf("%d(%d): %v\n", k + 1, pos, v) }
		cherr(err)
	}

	settings.Ircnick = rtext("Bot nick")
	settings.Uname = rtext("Realname")
	settings.Channel = rtext("Channel")
	settings.Server = rtext("server:port (SSL only!)")
	settings.Rate = rval("Rate (0-1000)", 0, elib.Ratemax, 10)
	settings.Kdel = rval("Keystroke delay in ms. (0-1000)", 0, elib.Kdelmax, 100)
	settings.Randel = rval("Random delay in ms. (0-10000)", 0, elib.Randelmax, 700)
	settings.Dnrmem = rval("Avoid repetition (0-20) times", 0, elib.Dnrmemmax, 5)

	s, err:= json.Marshal(settings)
	cherr(err)
	err = elib.Wrdb(db, 0, s)
	cherr(err)
}
