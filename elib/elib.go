package dbfunc

import (
	"fmt"
	"strconv"
	"github.com/boltdb/bolt"
)

var Cbuc = []byte("skymf")
const Ratemax = int(1000)
const Kdelmax = int(1000)
const Randelmax = int(10000)
const Dnrmemmax = int(20)

type Settings struct {
	Numln		int
	Rate		int
	Ircnick		string
	Uname		string
	Channel		string
	Server		string
	Kdel		int
	Randel		int
	Dnrmem		int
}

func Wrdb(db *bolt.DB, k int, v []byte) (err error) {

	err = db.Update(func(tx *bolt.Tx) error {
		buc, err := tx.CreateBucketIfNotExists(Cbuc)
		if err != nil { return err }

		err = buc.Put([]byte(strconv.Itoa(k)), v)
		if err != nil { return err }

		return nil
	})
	return
}

func Rdb(db *bolt.DB, k int) ([]byte, error) {

	var v []byte

	err := db.View(func(tx *bolt.Tx) error {
		buc := tx.Bucket(Cbuc)
		if buc == nil { return fmt.Errorf("No bucket!") }

		v = buc.Get([]byte(strconv.Itoa(k)))
		return nil
	})
	return v, err
}
