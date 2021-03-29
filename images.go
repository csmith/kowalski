package kowalski

import (
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
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
