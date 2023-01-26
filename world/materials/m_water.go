package materials

import (
	"image/color"

	"github.com/itiky/goPixelWorld/world/types"
)

var _ types.Material = Water{}

type Water struct {
	base
	surroundingFireDamperStep float64
}

func NewWater() Water {
	return Water{
		base: newBase(
			color.RGBA{R: 0x00, G: 0x6B, B: 0xFF, A: 0xAF},
			withFlags(types.MaterialFlagIsLiquid),
			withCloseRangeType(types.MaterialCloseRangeTypeSurrounding),
			withMass(10.0),
			withSourceDamping(0.2, 0.0),
		),
		surroundingFireDamperStep: 25.0,
	}
}

func (m Water) Type() types.MaterialType {
	return types.MaterialTypeWater
}

func (m Water) ProcessInternal(env types.TileEnvironment) {
	env.AddGravity()

	if cnt := env.DampEnvHealthByFlag(m.surroundingFireDamperStep, types.MaterialFlagIsFire); cnt > 0 {
		env.DampSelfHealth(float64(cnt) * m.surroundingFireDamperStep)
	}

	if env.Health() <= 0.0 {
		env.RemoveSelfHealthDamps()
		env.ReplaceSelf(SteamM)
	}
}

func (m Water) ProcessCollision(env types.CollisionEnvironment) {
	if env.IsFlagged(types.MaterialFlagIsLiquid) {
		if env.MoveLiquidSource() {
			return
		}
	}
	env.SwapSourceTarget()
	env.DampSourceForce(m.srcForceDamperK)
}
