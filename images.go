package kowalski

import (
	"bytes"
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
		scale = max(100 / bounds.Dy(), 100 / bounds.Dx())
	}

	outBounds := image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Min.X + bounds.Dx() * scale, bounds.Min.Y + bounds.Dy() * scale)
	redImage := image.NewRGBA(outBounds)
	greenImage := image.NewRGBA(outBounds)
	blueImage := image.NewRGBA(outBounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			cr, cg, cb, ca := im.At(x, y).RGBA()
			for dx := 0; dx < scale; dx++ {
				for dy := 0; dy < scale; dy++ {
					redImage.SetRGBA(x * scale + dx, y * scale + dy, color.RGBA{R: uint8(cr & 0xff), A: uint8(ca & 0xff)})
					greenImage.SetRGBA(x * scale + dx, y * scale + dy, color.RGBA{G: uint8(cg & 0xff), A: uint8(ca & 0xff)})
					blueImage.SetRGBA(x * scale + dx, y * scale + dy, color.RGBA{B: uint8(cb & 0xff), A: uint8(ca & 0xff)})
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
