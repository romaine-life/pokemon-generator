package gen

import (
	"encoding/hex"
	"image"
	"image/color"
	"image/png"
	"io"
)

// WritePNG draws the sprite scale times its size, each pixel a scale-by-scale block.
func (s Sprite) WritePNG(w io.Writer, scale int) error {
	return png.Encode(w, spriteImage(s, scale))
}

func spriteImage(s Sprite, scale int) *image.NRGBA {
	colors := make([]color.NRGBA, len(s.Palette))
	for i, p := range s.Palette {
		colors[i] = parseHex(p)
	}
	img := image.NewNRGBA(image.Rect(0, 0, s.Width*scale, s.Height*scale))
	for y := 0; y < s.Height*scale; y++ {
		for x := 0; x < s.Width*scale; x++ {
			img.SetNRGBA(x, y, colors[s.Pixels[(y/scale)*s.Width+x/scale]])
		}
	}
	return img
}

func parseHex(s string) color.NRGBA {
	b, err := hex.DecodeString(s[1:])
	if err != nil || len(b) < 3 {
		return color.NRGBA{}
	}
	c := color.NRGBA{b[0], b[1], b[2], 255}
	if len(b) == 4 {
		c.A = b[3]
	}
	return c
}
