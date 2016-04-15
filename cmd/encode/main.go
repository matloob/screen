package main

import (
	"os"

	"matloob.io/screen/psd"
)

func main() {
	psd.Encode(os.Stdout)
}
