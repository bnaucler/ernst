/*

		mkdb.go

*/

package main

import (
	"fmt"
	"os"
	"bufio"
	"strconv"
	"github.com/boltdb/bolt"
)

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

func gline(f *os.File, scanner *bufio.Scanner) string {

	f.Seek(1, 1)
	scanner.Scan()
	return scanner.Text()
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

	for k := 0; k < numln; k++ {
		v := gline(f, scanner)
		err = wrdb(db, []byte(strconv.Itoa(k+1)), []byte(v), cbuc)
		fmt.Printf("%d: %v\n", k, v)
		cherr(err)
	}
}
