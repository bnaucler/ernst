package dbfunc

import (
	"fmt"
	"strconv"
	"github.com/boltdb/bolt"
)

type Settings struct {
	Numln		int
	Rate		int
	Ircnick		string
	Uname		string
	Channel		string
	Server		string
	Kdel		int
	Randel		int
}

func Wrdb(db *bolt.DB, k int, v, cbuc []byte) (err error) {

	err = db.Update(func(tx *bolt.Tx) error {
		buc, err := tx.CreateBucketIfNotExists(cbuc)
		if err != nil { return err }

		err = buc.Put([]byte(strconv.Itoa(k)), v)
		if err != nil { return err }

		return nil
	})
	return
}

func Rdb(db *bolt.DB, k int, cbuc []byte) ([]byte, error) {

	var v []byte

	err := db.View(func(tx *bolt.Tx) error {
		buc := tx.Bucket(cbuc)
		if buc == nil { return fmt.Errorf("No bucket!") }

		v = buc.Get([]byte(strconv.Itoa(k)))
		return nil
	})
	return v, err
}
