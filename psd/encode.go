package psd

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"io"
)

func Encode(w io.Writer) {
	fileHeader := &FileHeader{
		Version:   1,
		Channels:  1,
		Height:    1800,
		Width:     1800,
		Depth:     1,
		ColorMode: Bitmap,
	}
	encodeFileHeader(w, fileHeader)
	write(w, int32(0)) // no color mode data
	//	write(w, int32(0)) // no image resources
	encodeImageResources(w, dpi300)
	write(w, int32(0)) // no layer & mask info

	// write image data
	write(w, int16(0)) // no compression
	totalSize := (fileHeader.Height*fileHeader.Width + 7) / 8
	imdata := make([]byte, totalSize)
	_, err := rand.Read(imdata)
	check(err)
	write(w, imdata)
	//	write(w, int8(0)) // image data, 1 byte
}

func write(w io.Writer, data interface{}) {
	check(binary.Write(w, binary.BigEndian, data))
}

func encodeFileHeader(w io.Writer, fh *FileHeader) {
	write(w, []byte("8BPS"))
	write(w, fh.Version)
	write(w, fh.Reserved)
	write(w, fh.Channels)
	write(w, fh.Height)
	write(w, fh.Width)
	write(w, fh.Depth)
	write(w, fh.ColorMode)
}

type fixed uint32 // fixed pt number 16 bits of decimal...

type ResolutionInfo struct {
	hRes       fixed
	hResUnit   int16
	WidthUnit  int16
	vRes       fixed
	vResUnit   int16
	heightUnit int16
}

func encodeResolutionInfo(w io.Writer, ri *ResolutionInfo) {
	write(w, ri.hRes)
	write(w, ri.hResUnit)
	write(w, ri.WidthUnit)
	write(w, ri.vRes)
	write(w, ri.vResUnit)
	write(w, ri.heightUnit)
}

func encodeResolutionResource(w io.Writer, ri *ResolutionInfo) {
	write(w, []byte("8BIM"))
	write(w, uint16(0x03ed))
	write(w, [2]byte{}) // empty name
	write(w, int32(16)) // res info is 16 bytes
	encodeResolutionInfo(w, ri)
}

func encodeImageResources(w io.Writer, ri *ResolutionInfo) {
	var buf bytes.Buffer
	encodeResolutionResource(&buf, ri)
	write(w, int32(buf.Len()))
	io.Copy(w, &buf)
}

var dpi300 = &ResolutionInfo{
	fixed(300 << 16),
	1, // RES_UNIT_PPI
	1, // UNIT_INCHES
	fixed(300 << 16),
	1,
	1,
}
