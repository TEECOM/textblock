package textblock

import (
	"image"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// Alignment indicates how text should be aligned within a text block.
type Alignment int

const (
	AlignmentLeft Alignment = iota
	AlignmentCenter
	AlignmentRight
)

// Options are optional arguments to New.
type Options struct {
	// Spacing is the line spacing for a block of text.
	//
	// A zero value means to use a line spacing of 1.5.
	Spacing float64

	// Alignment specifies how the text should be horizontally aligned
	// within the texblock's bounds.
	//
	// A zero value means to left-align the text.
	Alignment Alignment
}

func (o *Options) spacing() float64 {
	if o != nil && o.Spacing > 0 {
		return o.Spacing
	}
	return 1.5
}

func (o *Options) alignment() Alignment {
	if o != nil {
		return o.Alignment
	}
	return AlignmentLeft
}

type TextBlock interface {
	// BoundsAt returns a bounding rectangle that would encompass the text block
	// if it were positioned at the given point.
	BoundsAt(image.Point) image.Rectangle

	// DrawAt draws the text block positioned at the given point.
	DrawAt(image.Point)
}

// New initializes a new text block.
//
// The returned TextBlock will center the text when drawn. The returned
// text block is not concurrency safe, as the provided Drawer is not.
func New(d *font.Drawer, lines []string, opts *Options) TextBlock {
	spacing := fixed.I(int(float64(d.Face.Metrics().Height.Ceil()) * opts.spacing()))
	spaces := fixed.I(0)
	if len(lines) > 0 {
		spaces = fixed.Int26_6(len(lines) - 1)
	}
	height := d.Face.Metrics().Height * fixed.Int26_6(len(lines)) // heaight of all lines
	height += (spacing - d.Face.Metrics().Height) * spaces        // height of all spaces

	var width fixed.Int26_6
	for _, str := range lines {
		if size := d.MeasureString(str); size > width {
			width = size
		}
	}

	return &textBlock{
		d:       d,
		width:   width / 2,
		height:  height / 2,
		lines:   lines,
		spacing: spacing,
		opts:    opts,
	}
}

type textBlock struct {
	d       *font.Drawer
	width   fixed.Int26_6
	height  fixed.Int26_6
	lines   []string
	spacing fixed.Int26_6
	opts    *Options
}

// DrawAt draws the text block, centered on the given point.
func (tb *textBlock) DrawAt(pt image.Point) {
	pos := fixed.I(pt.Y) - tb.height + (tb.d.Face.Metrics().Height - tb.d.Face.Metrics().Descent)
	for _, str := range tb.lines {
		dx := tb.width
		if tb.opts.alignment() == AlignmentCenter {
			dx = tb.d.MeasureString(str) / 2
		} else if tb.opts.alignment() == AlignmentRight {
			dx = -dx + tb.d.MeasureString(str)
		}
		tb.d.Dot = fixed.Point26_6{
			X: fixed.I(pt.X) - dx,
			Y: pos,
		}
		tb.d.DrawString(str)
		pos += tb.spacing
	}
}

// BoundsAt returns the rectangle that would bound the text block if it were
// centered at the given point.
func (tb *textBlock) BoundsAt(pt image.Point) image.Rectangle {
	return image.Rectangle{
		Min: image.Point{
			X: (fixed.I(pt.X) - tb.width).Ceil(),
			Y: (fixed.I(pt.Y) - tb.height).Ceil(),
		},
		Max: image.Point{
			X: (fixed.I(pt.X) + tb.width).Ceil(),
			Y: (fixed.I(pt.Y) + tb.height).Ceil(),
		},
	}
}
