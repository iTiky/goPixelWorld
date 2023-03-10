package materials

import (
	"image/color"

	"github.com/itiky/goPixelWorld/world/types"
)

var _ types.Material = Metal{}

// Metal is very strong and reflects other Particles with a very low force damping.
type Metal struct {
	base
}

func NewMetal() Metal {
	return Metal{
		base: newBase(
			types.MaterialTypeMetal,
			color.RGBA{R: 0xBD, G: 0xC9, B: 0xBE, A: 0xFF},
			withFlags(types.MaterialFlagIsUnmovable),
			withMass(100000.0),
			withSourceDamping(0.7, 0.0),
		),
	}
}

func (m Metal) ProcessCollision(env types.CollisionEnvironment) {
	env.ReflectSourceTargetForces(m.srcForceDamperK)
}
