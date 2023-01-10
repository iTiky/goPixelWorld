package materials

import (
	"image/color"

	"github.com/itiky/goPixelWorld/world/types"
)

var _ types.Material = Water{}

type Water struct {
	base
}

func NewWater() Water {
	return Water{
		base: newBase(
			color.RGBA{R: 0x00, G: 0x6B, B: 0xFF, A: 0xAF},
			withForceDamperK(0.2),
			withFlags(types.MaterialFlagIsLiquid),
			withMass(7.5),
		),
	}
}

func (m Water) Type() types.MaterialType {
	return types.MaterialTypeWater
}

func (m Water) ProcessInternal(env types.TileEnvironment) {
	env.AddGravity()
}

func (m Water) ProcessCollision(env types.CollisionEnvironment) {
	env.DampSourceForce(m.forceDamperK)

	if env.IsFlagged(types.MaterialFlagIsLiquid) {
		if env.MoveLiquidSource() {
			return
		}
	}

	env.SwapSourceTarget()

	return
}
