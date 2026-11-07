package kowalski

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHiddenPixels_highlightsBothColoursOfAPair(t *testing.T) {
	first := color.RGBA{R: 0x0a, G: 0x0a, B: 0x0a, A: 255}
	second := color.RGBA{R: 0x0c, G: 0x0c, B: 0x0c, A: 255}

	input := image.NewRGBA(image.Rect(0, 0, 20, 20))
	for y := range 20 {
		for x := range 20 {
			if (x+y)%2 == 0 {
				input.SetRGBA(x, y, first)
			} else {
				input.SetRGBA(x, y, second)
			}
		}
	}

	buf := &bytes.Buffer{}
	require.NoError(t, png.Encode(buf, input))

	reader, err := HiddenPixels(buf)
	require.NoError(t, err)

	output, err := png.Decode(reader)
	require.NoError(t, err)

	black := color.RGBA{A: 255}
	firstHighlight := color.RGBAModel.Convert(output.At(0, 0))
	secondHighlight := color.RGBAModel.Convert(output.At(1, 0))
	assert.NotEqual(t, black, firstHighlight)
	assert.NotEqual(t, black, secondHighlight)
	assert.NotEqual(t, firstHighlight, secondHighlight)
}
