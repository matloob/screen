package psd

import (
	"bytes"
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
	fmt.Println("layers length", length)
	buf := make([]byte, length)
	n, err := r.Read(buf)
	if err != nil || n < int(length) {
		panic("error in skipLayers")
	}
	fmt.Println("layersbuf", buf)

	bb := bytes.NewBuffer(buf)
	skipLayerInfo(bb)
	skipGlobalLayerMaskInfo(bb)	
	fmt.Println("remaining", bb.Bytes())	
}

func skipLayerInfo(r io.Reader) {
	var length int32 
	read(r, &length)
	fmt.Println("layer info length: ", length)
	if length == 0 {
		fmt.Println("X")
		return
	}
	var count int16
	read(r, &count)
	fmt.Println("layer count: ", count)
	buf := make([]byte, length)
	_, err := r.Read(buf)
	if err != nil {
		panic ("error in skipLayerInfo")
	}
}

func skipGlobalLayerMaskInfo(r io.Reader) {
	var length int32
	var overlayColorSpace int16
	var colorComponents [4]int16
	var opacity int16
	var kind byte
	var pad byte

	read(r, &length)
	fmt.Println("global layer mask info length:", length)
	read(r, &overlayColorSpace)
	read(r, &colorComponents)
	read(r, &opacity)
	read(r, &kind)
	read(r, &pad)
	fmt.Println(overlayColorSpace, colorComponents, opacity, kind)
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