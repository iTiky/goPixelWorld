package materials

import (
	"image/color"

	"github.com/itiky/goPixelWorld/world/types"
)

var _ types.Material = Sand{}

type Sand struct {
	base
}

func NewSand() Sand {
	return Sand{
		base: newBase(
			color.RGBA{R: 0xFF, G: 0xD5, B: 0x00, A: 0xFF},
			withFlags(types.MaterialFlagIsSand),
			withMass(5.0),
			withForceDamperK(0.9),
		),
	}
}

func (m Sand) Type() types.MaterialType {
	return types.MaterialTypeSand
}

func (m Sand) ProcessInternal(env types.TileEnvironment) {
	env.AddGravity()
}

func (m Sand) ProcessCollision(env types.CollisionEnvironment) {
	if env.IsFlagged(types.MaterialFlagIsSand) || env.IsFlagged(types.MaterialFlagIsLiquid) {
		if env.MoveSandSource() {
			return
		}
	}

	env.ReflectSourceTargetForces()

	env.DampSourceForce(m.forceDamperK)

	return
}
