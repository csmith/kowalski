package kowalski

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"sort"

	_ "golang.org/x/image/webp"
)

type ColourCount struct {
	Colour color.Color
	Count  int
}

// ExtractColours returns a sorted slice containing each individual colour used in the image, and the total number of
// pixels that have that colour. Colours are sorted from most-used to least-used.
func ExtractColours(reader io.Reader) ([]ColourCount, error) {
	im, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	count := make(map[color.Color]int)
	b := im.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			count[im.At(x, y)]++
		}
	}

	var res []ColourCount
	for c, n := range count {
		res = append(res, ColourCount{c, n})
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Count > res[j].Count
	})

	return res, nil
}

// HiddenPixels attempts to find patterns of hidden pixels (those that are consistently used near very similar colours).
func HiddenPixels(reader io.Reader) (io.Reader, error) {
	im, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	colourAt := func(x, y int) uint32 {
		r, g, b, _ := im.At(x, y).RGBA()
		return (r&0xff)<<16 | (g&0xff)<<8 | (b & 0xff)
	}

	diff := func(c1, c2 uint32) uint32 {
		if c1 < c2 {
			return c2 - c1
		}
		return c1 - c2
	}

	similar := func(c1, c2 uint32) bool {
		return c1 != c2 && diff(c1&0xff, c2&0xff) <= 10 &&
			diff((c1>>8)&0xff, (c2>>8)&0xff) <= 10 &&
			diff((c1>>16)&0xff, (c2>>16)&0xff) <= 10
	}

	pair := func(c1, c2 uint32) uint32 {
		if c1 < c2 {
			return (c1 << 24) | c2
		} else {
			return (c2 << 24) | c1
		}
	}

	pairs := make(map[uint32]int)
	b := im.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			myPos := colourAt(x, y)
			if x > 0 {
				otherPos := colourAt(x-1, y)
				if similar(myPos, otherPos) {
					pairs[pair(myPos, otherPos)]++
				}
			}

			if y > 0 {
				otherPos := colourAt(x, y-1)
				if similar(myPos, otherPos) {
					pairs[pair(myPos, otherPos)]++
				}
			}
		}
	}

	colours := []color.RGBA{
		{R: 255, A: 255},
		{G: 255, A: 255},
		{B: 255, A: 255},
		{R: 255, G: 255, A: 255},
		{R: 255, B: 255, A: 255},
		{B: 255, G: 255, A: 255},
	}
	replacements := make(map[uint32]color.RGBA)
	nextColour := 0
	for i := range pairs {
		if pairs[i] > 100 {
			replacements[i&0xffffff] = colours[nextColour]
			nextColour = (nextColour + 1) % len(colours)
			replacements[(i>>24)&0xffffff] = colours[nextColour]
			nextColour = (nextColour + 1) % len(colours)
		}
	}

	if len(replacements) == 0 {
		return nil, errors.New("no obvious hidden pixels found")
	}

	output := image.NewRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if c, ok := replacements[colourAt(x, y)]; ok {
				output.SetRGBA(x, y, c)
			} else {
				output.SetRGBA(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 255})
			}
		}
	}

	outputBuffer := &bytes.Buffer{}
	if err := png.Encode(outputBuffer, output); err != nil {
		return nil, err
	}
	return outputBuffer, nil
}

// SplitRGB reads an image, and returns three readers containing PNG encoded images containing just the red,
// green and blue channels respectively. If the input image is excessively small, the outputs will be scaled
// up.
func SplitRGB(reader io.Reader) (r, g, b io.Reader, err error) {
	im, _, err := image.Decode(reader)
	if err != nil {
		return nil, nil, nil, err
	}

	bounds := im.Bounds()
	scale := 1
	if bounds.Dy() < 100 || bounds.Dx() < 100 {
		scale = max(100/bounds.Dy(), 100/bounds.Dx())
	}

	outBounds := image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Min.X+bounds.Dx()*scale, bounds.Min.Y+bounds.Dy()*scale)
	redImage := image.NewRGBA(outBounds)
	greenImage := image.NewRGBA(outBounds)
	blueImage := image.NewRGBA(outBounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			cr, cg, cb, ca := im.At(x, y).RGBA()
			for dx := 0; dx < scale; dx++ {
				for dy := 0; dy < scale; dy++ {
					redImage.SetRGBA(x*scale+dx, y*scale+dy, color.RGBA{R: uint8(cr & 0xff), A: uint8(ca & 0xff)})
					greenImage.SetRGBA(x*scale+dx, y*scale+dy, color.RGBA{G: uint8(cg & 0xff), A: uint8(ca & 0xff)})
					blueImage.SetRGBA(x*scale+dx, y*scale+dy, color.RGBA{B: uint8(cb & 0xff), A: uint8(ca & 0xff)})
				}
			}
		}
	}

	redBuf := &bytes.Buffer{}
	greenBuf := &bytes.Buffer{}
	blueBuf := &bytes.Buffer{}

	if err := png.Encode(redBuf, redImage); err != nil {
		return nil, nil, nil, err
	}

	if err := png.Encode(greenBuf, greenImage); err != nil {
		return nil, nil, nil, err
	}

	if err := png.Encode(blueBuf, blueImage); err != nil {
		return nil, nil, nil, err
	}

	return redBuf, greenBuf, blueBuf, nil
}
