package materials

import (
	"image/color"

	"github.com/itiky/goPixelWorld/world/types"
)

var _ types.Material = AntiGraviton{}

// AntiGraviton pushes away other Particles in a circle area.
type AntiGraviton struct {
	base
	antiGravityForceMag float64 // anti-gravity force magnitude
}

func NewAntiGraviton() AntiGraviton {
	return AntiGraviton{
		base: newBase(
			types.MaterialTypeAntiGraviton,
			color.RGBA{R: 0x00, G: 0xA7, B: 0x9F, A: 0xFF},
			withCloseRangeType(types.MaterialCloseRangeTypeInCircleRange),
			withMass(1000000.0),
			withFlags(types.MaterialFlagIsUnmovable),
			withCloseRangeCircleR(20),
		),
		antiGravityForceMag: -0.7,
	}
}

func (m AntiGraviton) ProcessInternal(env types.TileEnvironment) {
	m.commonProcessInternal(env)

	env.AddForceInRange(m.antiGravityForceMag, []types.MaterialFlag{types.MaterialFlagIsUnmovable})
}

func (m AntiGraviton) ProcessCollision(env types.CollisionEnvironment) {
	env.ReflectSourceTargetForces(m.srcForceDamperK)
}
