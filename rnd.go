package main

import (
	"fmt"
	"bufio"
	"time"
	"math/rand"
	"os"
)

const fname = "skymfer.txt"

func cherr(e error) {
	if e != nil { panic(e) }
}

func getrandstr() {

}

func main() {

	s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)

	var skymf string

	cfile, err := os.Open(fname)
	cherr(err)
	defer cfile.Close()

	scanner := bufio.NewScanner(cfile)
	numlns := 0

	for scanner.Scan() {
		numlns++
	}

	randln := r1.Intn(numlns)

	cfile.Seek(0, 0)
	scanner = bufio.NewScanner(cfile)
	for a := 0; a < randln; a++ {
		scanner.Scan()
		skymf = scanner.Text()
	}

	fmt.Printf("File obj: %T\n", cfile)
	fmt.Printf("Rand obj: %T\n", r1)
	fmt.Printf("Total: %d\n", numlns)
	fmt.Printf("Random: %d\n", randln)
	fmt.Printf("Skymf: %v\n", skymf)

}
