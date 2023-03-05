package materials

import (
	"image/color"

	"github.com/itiky/goPixelWorld/world/types"
)

var _ types.Material = Graviton{}

// Graviton attracts other Particles in a circle area.
type Graviton struct {
	base
	gravityForceMag float64
}

func NewGraviton() Graviton {
	return Graviton{
		base: newBase(
			color.RGBA{R: 0x8F, G: 0x00, B: 0xA2, A: 0xFF},
			withCloseRangeType(types.MaterialCloseRangeTypeInCircleRange),
			withMass(1000000.0),
			withFlags(types.MaterialFlagIsUnmovable),
			withCloseRangeCircleR(20),
		),
		gravityForceMag: 0.7,
	}
}

func (m Graviton) Type() types.MaterialType {
	return types.MaterialTypeGraviton
}

func (m Graviton) ProcessInternal(env types.TileEnvironment) {
	env.AddForceInRange(m.gravityForceMag, types.MaterialFlagIsUnmovable)
}

func (m Graviton) ProcessCollision(env types.CollisionEnvironment) {
	env.ReflectSourceTargetForces(m.srcForceDamperK)
}
