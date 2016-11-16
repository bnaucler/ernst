/*

		mkdb.go

*/

package main

import (
	"fmt"
	"log"
	"os"
	"bufio"
	"strconv"
	"github.com/boltdb/bolt"
)

func cherr(e error) { if e != nil { log.Fatal(e) } }

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

func main() {

	cbuc := []byte("skymf")
	dbname := "./ernst.db"
	sfname := "./skymfer.txt"

	db, err := bolt.Open(dbname, 0640, nil)
	cherr(err)
	defer db.Close()

	f, err := os.Open(sfname)
	cherr(err)
	defer f.Close()

	numln := clines(f)

	scanner := bufio.NewScanner(f)
	f.Seek(0, 0)
	var pos = int64(0)
	var v = string("")

	for k := 0; k < numln; k++ {
		v, pos = gline(f, scanner, pos)
		err = wrdb(db, []byte(strconv.Itoa(k+1)), []byte(v), cbuc)
		fmt.Printf("%d(%d): %v\n", k + 1, pos, v)
		cherr(err)
	}

	err = wrdb(db, []byte("0"), []byte(strconv.Itoa(numln)), cbuc)
	cherr(err)
}
