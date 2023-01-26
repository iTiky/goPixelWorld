package materials

import (
	"image/color"

	"github.com/itiky/goPixelWorld/world/types"
)

var _ types.Material = Smoke{}

type Smoke struct {
	base
}

func NewSmoke() Smoke {
	return Smoke{
		base: newBase(
			color.RGBA{R: 0xCD, G: 0xCD, B: 0xCD, A: 0xFF},
			withFlags(types.MaterialFlagIsGas),
			withCloseRangeType(types.MaterialCloseRangeTypeSelfOnly),
			withMass(2.0),
			withSelfHealthReduction(100.0, 0.5),
		),
	}
}

func (m Smoke) Type() types.MaterialType {
	return types.MaterialTypeSmoke
}

func (m Smoke) ColorAdjusted(health float64) color.Color {
	if health < 20.0 {
		return color.RGBA{R: 0x35, G: 0x35, B: 0x35, A: 0x35}
	} else if health < 40.0 {
		return color.RGBA{R: 0x5C, G: 0x5C, B: 0x5C, A: 0x5C}
	} else if health < 60.0 {
		return color.RGBA{R: 0x6F, G: 0x6F, B: 0x6F, A: 0x6F}
	} else if health < 80.0 {
		return color.RGBA{R: 0x9C, G: 0x9C, B: 0x9C, A: 0xFF}
	}

	return m.baseColor
}

func (m Smoke) ProcessInternal(env types.TileEnvironment) {
	env.AddReverseGravity()
	env.DampSelfHealth(m.selfHealthDampStep)
}

func (m Smoke) ProcessCollision(env types.CollisionEnvironment) {
	if env.IsFlagged(types.MaterialFlagIsGas) {
		env.MoveSandSource()
		return
	}
	env.SwapSourceTarget()
}
