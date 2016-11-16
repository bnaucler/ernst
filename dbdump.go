/* 

		dbdump.go
		Dumps a BoltDB on the screen. That's all

*/

package main

import (
	"fmt"
	"os"
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

func rdb(db *bolt.DB, k, cbuc []byte) (v []byte, err error) {

	err = db.View(func(tx *bolt.Tx) error {
		buc := tx.Bucket(cbuc)
		if buc == nil { return fmt.Errorf("No bucket!") }

		v = buc.Get(k)
		return nil
	})
	return
}

func main() {


	if len(os.Args) != 3 {
		cherr(fmt.Errorf("Usage: %s <file> <bucket>\n", os.Args[0]))
	}

	dbname := os.Args[1]
	cbuc := []byte(os.Args[2])

	db, err := bolt.Open(dbname, 0640, nil)
	cherr(err)
	defer db.Close()

	dbdump(db, cbuc)
	// for a := 1; a < 1063; a++ {
	// 	v, err := rdb(db, []byte(strconv.Itoa(a)), cbuc)
	// 	cherr(err)
	// 	fmt.Printf("%d: %v\n", a, string(v))
	// }
}
