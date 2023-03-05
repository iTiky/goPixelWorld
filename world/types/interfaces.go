package types

import (
	"image/color"
)

// TileI defines an interface to export Tile state.
type TileI interface {
	MaterialI
	X() int
	Y() int
}

// MaterialI defines an interface to export Material state.
type MaterialI interface {
	Color() color.Color
}
