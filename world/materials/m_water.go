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
			withMass(10.0),
			withForceDamperK(0.2),
		),
		surroundingFireDamperStep: 25.0,
	}
}

func (m Water) Type() types.MaterialType {
	return types.MaterialTypeWater
}

func (m Water) ProcessInternal(env types.TileEnvironment) {
	env.AddGravity()

	if cnt := env.ReduceEnvHealthByFlag(m.surroundingFireDamperStep, types.MaterialFlagIsFire); cnt > 0 {
		env.ReduceHealth(float64(cnt) * m.surroundingFireDamperStep)
	}

	if env.Health() <= 0.0 {
		env.RemoveHealthSelfReductions()
		env.ReplaceSelf(NewSteam())
	}
}

func (m Water) ProcessCollision(env types.CollisionEnvironment) {
	if env.IsFlagged(types.MaterialFlagIsLiquid) {
		if env.MoveLiquidSource() {
			return
		}
	}

	env.SwapSourceTarget()

	env.DampSourceForce(m.forceDamperK)

	return
}
