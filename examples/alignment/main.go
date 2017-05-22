package main

import (
	"flag"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/image/font"

	"github.com/TEECOM/textblock"
	"github.com/golang/freetype/truetype"
)

const imgW, imgH = 640, 480

func main() {
	ff := flag.String("f", "font.ttf", "filename of a .ttf font")
	flag.Parse()

	// Read font
	fb, err := ioutil.ReadFile(*ff)
	if err != nil {
		log.Println(err)
	}
	f, err := truetype.Parse(fb)
	if err != nil {
		log.Println(err)
	}

	fg, bg := image.Black, image.White
	rgba := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)

	d := &font.Drawer{
		Dst:  rgba,
		Src:  fg,
		Face: truetype.NewFace(f, &truetype.Options{Size: 50}),
	}

	tb := textblock.New(
		d,
		[]string{"How", "Do", "Youy", "Doy"},
		&textblock.Options{Spacing: 1.5, Alignment: textblock.AlignmentRight},
	)

	pt := image.Point{imgW / 2, imgH / 2}

	tb.DrawAt(pt)

	b := tb.BoundsAt(pt)
	for x := b.Min.X; x < b.Max.X; x++ {
		rgba.Set(x, b.Min.Y, color.RGBA{255, 0, 0, 255})
		rgba.Set(x, b.Max.Y, color.RGBA{255, 0, 0, 255})
	}
	for y := b.Min.Y; y < b.Max.Y; y++ {
		rgba.Set(b.Min.X, y, color.RGBA{255, 0, 0, 255})
		rgba.Set(b.Max.X, y, color.RGBA{255, 0, 0, 255})
	}

	out, err := os.Create("out.png")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	if err := png.Encode(out, rgba); err != nil {
		log.Fatal(err)
	}
}
