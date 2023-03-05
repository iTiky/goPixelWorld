package types

import (
	"image/color"
)

type Pixel struct {
	Ready bool
	X     int
	Y     int
	Color color.Color
}
