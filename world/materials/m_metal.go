package materials

import (
	"image/color"

	"github.com/itiky/goPixelWorld/world/types"
)

var _ types.Material = Metal{}

type Metal struct {
	base
}

func NewMetal() Metal {
	return Metal{
		base: newBase(
			color.RGBA{R: 0xBD, G: 0xC9, B: 0xBE, A: 0xFF},
			withFlags(types.MaterialFlagIsUnmovable),
			withMass(100000.0),
			withSourceDamping(0.7, 0.0),
		),
	}
}

func (m Metal) Type() types.MaterialType {
	return types.MaterialTypeMetal
}

func (m Metal) ProcessCollision(env types.CollisionEnvironment) {
	env.ReflectSourceTargetForces(m.srcForceDamperK)
}
