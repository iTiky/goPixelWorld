package materials

import (
	"image/color"

	"github.com/itiky/goPixelWorld/world/types"
)

var _ types.Material = Wood{}

type Wood struct {
	base
}

func NewWood() Wood {
	return Wood{
		base: newBase(
			color.RGBA{R: 0x7A, G: 0x33, B: 0x00, A: 0xFF},
			withForceDamperK(0.05),
			withFlags(types.MaterialFlagIsFlammable),
		),
	}
}

func (m Wood) Type() types.MaterialType {
	return types.MaterialTypeWood
}

func (m Wood) ProcessCollision(env types.CollisionEnvironment) {
	env.DampSourceForce(m.forceDamperK)
	env.ReflectSourceForce()
}
