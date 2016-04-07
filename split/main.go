package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

var in = flag.String("in", "", "input file")

var out = flag.String("out", "", "output file")

var s1 = flag.String("s1", "", "")
var s2 = flag.String("s2", "", "")

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()
	f, err := os.Open(*in)
	check(err)

	im, err := png.Decode(f)
	check(err)
	check(f.Close())

	outim := image.NewGray(im.Bounds())
	outs1 := image.NewGray(im.Bounds())
	outs2 := image.NewGray(im.Bounds())
	for x := im.Bounds().Min.X; x <= im.Bounds().Max.X; x++ {
		for y := im.Bounds().Min.Y; y <= im.Bounds().Max.Y; y++ {
			outim.Set(x, y, filter(im.At(x, y)))
			c := im.At(x, y)
			outs1.Set(x, y, color.Gray{255})
			outs2.Set(x, y, color.Gray{255})
			if s1f(c) {
				outs1.Set(x, y, color.Gray{0})
			}
			if s2f(c) {
				outs2.Set(x, y, color.Gray{0})
			}
		}
	}

	outf, err := os.Create(*out)
	check(err)
	check(png.Encode(outf, outim))
	check(outf.Close())

	outfs1, err := os.Create(*s1)
	check(err)
	check(png.Encode(outfs1, outs1))
	check(outfs1.Close())

	outfs2, err := os.Create(*s2)
	check(err)
	check(png.Encode(outfs2, outs2))
	check(outfs2.Close())

}

func s1f(c color.Color) bool {
	g := color.GrayModel.Convert(c).(color.Gray)
	return g.Y < 100
}

func s2f(c color.Color) bool {
	g := color.GrayModel.Convert(c).(color.Gray)
	return g.Y >= 100 && g.Y < 196
}

func filter(c color.Color) color.Color {

	if s1f(c) {
		return color.Gray{0}
	}
	if s2f(c) {
		return color.Gray{128}
	}
	return color.Gray{255}
}
