/* 

		dbdump.go
		Dumps the databse on the screen. That's all.

*/

package main

import (
	"fmt"
	"os"
	"encoding/json"
	"github.com/boltdb/bolt"
	elib "github.com/bnaucler/ernst/elib"
)

func cherr(e error) { if e != nil { panic(e) } }

func main() {

	settings := elib.Settings{}
	var verb bool

	if len(os.Args) < 2 || len(os.Args) > 3 {
		cherr(fmt.Errorf("Usage: %s <file> [v]\n", os.Args[0]))
	}

	dbname := os.Args[1]
	if len(os.Args) == 3 && os.Args[2] == "v" { verb = true };

	db, err := bolt.Open(dbname, 0640, nil)
	cherr(err)
	defer db.Close()

	tmp, err := elib.Rdb(db, 0)
	cherr(err)
	json.Unmarshal(tmp, &settings)

	for k := 1; k <= settings.Numln; k++ {
		v, err := elib.Rdb(db, k)
		cherr(err)
		if verb { fmt.Printf("%d: %v\n", settings.Numln, string(v))
		} else { fmt.Printf("%v\n", string(v)) }
	}

	if verb { fmt.Printf("Settings: %+v\n", settings) }
}
