package main

import (
	"flag"
	"log"
	"os"

	"matloob.io/screen/psd"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()
	f, err := os.Open(flag.Arg(0))
	check(err)
	psd.Decode(f)
}