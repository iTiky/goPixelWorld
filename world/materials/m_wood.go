package materials

import (
	"image/color"

	"github.com/itiky/goPixelWorld/world/types"
)

var _ types.Material = Wood{}

// Wood is a soft and flammable material.
type Wood struct {
	base
}

func NewWood() Wood {
	return Wood{
		base: newBase(
			color.RGBA{R: 0x7A, G: 0x33, B: 0x00, A: 0xFF},
			withFlags(types.MaterialFlagIsUnmovable, types.MaterialFlagIsFlammable),
			withCloseRangeType(types.MaterialCloseRangeTypeNone),
			withMass(1000.0),
			withSelfHealthReduction(100.0, 0.5),
			withSourceDamping(0.5, 0.0),
		),
	}
}

func (m Wood) Type() types.MaterialType {
	return types.MaterialTypeWood
}

func (m Wood) ProcessCollision(env types.CollisionEnvironment) {
	env.ReflectSourceTargetForces(m.srcForceDamperK)
	env.DampSelfHealth(m.selfHealthDampStep)
}
