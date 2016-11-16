/* 

		dbdump.go
		Dumps a BoltDB on the screen. That's all

*/

package main

import (
	"fmt"
	"github.com/boltdb/bolt"
)

func cherr(e error) {
	if e != nil { panic(e) }
}

func dbdump (db *bolt.DB, cbuc []byte) {

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(cbuc)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("%s=%+v\n", k, string(v))
		}
		return nil
	})
}

func main() {

	cbuc := []byte("skymf")
	dbname := "./ernst.db"

	db, err := bolt.Open(dbname, 0640, nil)
	cherr(err)
	defer db.Close()

	dbdump(db, cbuc)
}
