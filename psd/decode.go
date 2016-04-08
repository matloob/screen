package psd

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Decode(r io.Reader) {
	var fileHeader FileHeader

	decodeFileHeader(r, &fileHeader)
	skipColorModeData(r)
	skipImageResources(r)
	skipLayers(r)
	decodeImageData(r)

	fmt.Printf("%+v\n", fileHeader)
	
}

func read(r io.Reader, data interface{}) {
	check(binary.Read(r, binary.BigEndian, data))
}

func decodeFileHeader(r io.Reader, h *FileHeader) {
	read(r, &h.Signature)
	read(r, &h.Version)
	read(r, &h.Reserved)
	read(r, &h.Channels)
	read(r, &h.Height)
	read(r, &h.Width)
	read(r, &h.Depth)
	read(r, &h.ColorMode)
}

func skipColorModeData(r io.Reader) {
	var length int32
	read(r, &length)
	if length != 0 {
		panic("non-empty color mode data")
	}
}

func skipImageResources(r io.Reader) {
	var length int32
	read(r, &length)
	fmt.Println("irlength:", length)
	buf := make([]byte, length)
	n, err := r.Read(buf)
	if err != nil || n < int(length) {
		panic("error in skipImageResources")
	}
	 
}

func skipLayers(r io.Reader) {
	var length int32
	read(r, &length)
	buf := make([]byte, length)
	n, err := r.Read(buf)
	if err != nil || n < int(length) {
		panic("error in skipLayers")
	}	
}

func decodeImageData(r io.Reader) {
	var compression int16
	read(r, &compression)
	fmt.Println("compression:", compression)
	b, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(b)
}

type ColorMode int16

const (
	Bitmap ColorMode = iota
	Grayscale
	Indexed
	RGB
	CMYK
	Multichannel
	Duotone
	Lab
)

type FileHeader struct {
	Signature	[4]byte
	Version uint16
	Reserved [6]byte
	Channels uint16
	Height uint32
	Width uint32
	Depth uint16
	ColorMode ColorMode
}