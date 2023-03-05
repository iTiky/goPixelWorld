package pkg

import (
	"image/color"
)

// ColorToNRGBA converts the color to the 8-bit model.
func ColorToNRGBA(c color.Color) color.NRGBA {
	r, g, b, a := c.RGBA()

	return color.NRGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: uint8(a >> 8),
	}
}
