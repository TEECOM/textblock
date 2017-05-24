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

const imgW, imgH = 600, 1200

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

	// Make white image the size of imgW and imgH.
	rgba := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
	draw.Draw(rgba, rgba.Bounds(), image.White, image.ZP, draw.Src)

	d := &font.Drawer{
		Dst:  rgba,
		Src:  image.Black,
		Face: truetype.NewFace(f, &truetype.Options{Size: 50}),
	}

	o := &textblock.Options{Spacing: 1.2, Alignment: textblock.AlignmentLeft}
	tb := textblock.New(
		d,
		[]string{"Hello", "There", "Loooooooooooooooong", "Word"},
		o,
	)

	// Draw text left aligned
	pt := image.Point{imgW / 2, 200}
	tb.DrawAt(pt)
	drawRect(rgba, tb.BoundsAt(pt), color.RGBA{255, 0, 0, 255})

	// Draw text center aligned
	pt = image.Point{imgW / 2, 600}
	o.Alignment = textblock.AlignmentCenter
	tb.DrawAt(pt)
	drawRect(rgba, tb.BoundsAt(pt), color.RGBA{255, 0, 0, 255})

	// Draw text right aligned
	pt = image.Point{imgW / 2, 1000}
	o.Alignment = textblock.AlignmentRight
	tb.DrawAt(pt)
	drawRect(rgba, tb.BoundsAt(pt), color.RGBA{255, 0, 0, 255})

	out, err := os.Create("out.png")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	if err := png.Encode(out, rgba); err != nil {
		log.Fatal(err)
	}
}

func drawRect(dst draw.Image, rect image.Rectangle, clr color.Color) {
	for x := rect.Min.X; x < rect.Max.X; x++ {
		dst.Set(x, rect.Min.Y, clr)
		dst.Set(x, rect.Max.Y, clr)
	}
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		dst.Set(rect.Min.X, y, clr)
		dst.Set(rect.Max.X, y, clr)
	}
}
