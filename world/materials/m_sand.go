package materials

import (
	"image/color"

	"github.com/itiky/goPixelWorld/world/types"
)

var _ types.Material = Sand{}

// Sand spreads like a sand.
type Sand struct {
	base
}

func NewSand() Sand {
	return Sand{
		base: newBase(
			types.MaterialTypeSand,
			color.RGBA{R: 0xFF, G: 0xD5, B: 0x00, A: 0xFF},
			withFlags(types.MaterialFlagIsSand),
			withCloseRangeType(types.MaterialCloseRangeTypeSelfOnly),
			withMass(5.0),
			withSourceDamping(0.9, 0.0),
		),
	}
}

func (m Sand) ProcessInternal(env types.TileEnvironment) {
	m.commonProcessInternal(env)

	env.AddGravity()
}

func (m Sand) ProcessCollision(env types.CollisionEnvironment) {
	if env.IsFlagged(types.MaterialFlagIsSand) || env.IsFlagged(types.MaterialFlagIsLiquid) {
		if env.MoveSandSource() {
			return
		}
	}
	env.ReflectSourceTargetForces(m.srcForceDamperK)
}
