package types

import (
	"image/color"
)

// TileI defines an interface to export Tile state.
type TileI interface {
	X() int
	Y() int
	Color() color.Color
}

// MaterialI defines an interface to export Material state.
type MaterialI interface {
	Name() string
	Color() color.Color
}
