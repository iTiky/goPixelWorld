package types

import (
	"image/color"
)

var _ TileI = Pixel{}

// Pixel defines the Tile's export state.
type Pixel struct {
	Ready         bool // if false, object is empty
	PosX          int
	PosY          int
	ParticleColor color.Color
}

func (p Pixel) X() int {
	return p.PosX
}

func (p Pixel) Y() int {
	return p.PosY
}

func (p Pixel) Color() color.Color {
	return p.ParticleColor
}
