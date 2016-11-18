/* 

		dbdump.go
		Dumps the databse on the screen. That's all.

*/

package main

import (
	"fmt"
	"os"
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

func main() {

	settings := Settings{}
	var verb bool

	if len(os.Args) < 3 || len(os.Args) > 4 {
		cherr(fmt.Errorf("Usage: %s <file> <bucket> [v]\n", os.Args[0]))
	}

	dbname := os.Args[1]
	cbuc := []byte(os.Args[2])
	if len(os.Args) == 4 && os.Args[3] == "v" { verb = true }

	db, err := bolt.Open(dbname, 0640, nil)
	cherr(err)
	defer db.Close()

	tmp, err := rdb(db, 0, cbuc)
	cherr(err)
	json.Unmarshal(tmp, &settings)

	for k := 0; k <= settings.Numln; k++ {
		v, err := rdb(db, k, cbuc)
		cherr(err)
		if verb { fmt.Printf("%d: %v\n", settings.Numln, string(v))
		} else { fmt.Printf("%v\n", string(v)) }
	}

	if verb { fmt.Printf("Settings: %+v\n", settings) }
}
