package materials

import (
	"image/color"

	"github.com/itiky/goPixelWorld/world/types"
)

var _ types.Material = Rock{}

type Rock struct {
	base
}

func NewRock() Rock {
	return Rock{
		base: newBase(
			color.RGBA{R: 0xA7, G: 0x39, B: 0x00, A: 0xFF},
			withMass(100.0),
			withSelfHealthReduction(100.0, 0.5),
			withSourceDamping(0.9, 0.0),
		),
	}
}

func (m Rock) Type() types.MaterialType {
	return types.MaterialTypeRock
}

func (m Rock) ProcessInternal(env types.TileEnvironment) {
	env.AddGravity()
}

func (m Rock) ProcessCollision(env types.CollisionEnvironment) {
	env.ReflectSourceTargetForces(m.srcForceDamperK)
	env.DampSelfHealthByMassRate(m.selfHealthDampStep)
}
