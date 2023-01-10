package types

import (
	"image/color"
)

type TileI interface {
	MaterialI
	Position() Position
}

type MaterialI interface {
	Color() color.Color
}
