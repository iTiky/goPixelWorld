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
			withFlags(types.MaterialFlagIsFlammable),
			withMass(1000.0),
			withForceDamperK(0.2),
		),
	}
}

func (m Wood) Type() types.MaterialType {
	return types.MaterialTypeWood
}

func (m Wood) ProcessCollision(env types.CollisionEnvironment) {
	env.ReflectSourceTargetForces()
	env.DampSourceForce(m.forceDamperK)
}
