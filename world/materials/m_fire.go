package materials

import (
	"image/color"

	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/types"
)

var _ types.Material = Fire{}

type Fire struct {
	base
}

func NewFire() Fire {
	return Fire{
		base: newBase(
			color.RGBA{R: 0xFF, G: 0xAD, B: 0x8B, A: 0xFF},
			withHealthDamperStep(1.5),
			withFlags(types.MaterialFlagIsGas),
			withMass(1.0),
		),
	}
}

func (m Fire) Type() types.MaterialType {
	return types.MaterialTypeFire
}

func (m Fire) ColorAdjusted(health float64) color.Color {
	if health < 20.0 {
		return color.RGBA{R: 0x9C, G: 0x00, B: 0x05, A: 0x35}
	} else if health < 40.0 {
		return color.RGBA{R: 0xFF, G: 0x35, B: 0x3B, A: 0x5C}
	} else if health < 60.0 {
		return color.RGBA{R: 0xFF, G: 0x5C, B: 0x4B, A: 0x6F}
	} else if health < 80.0 {
		return color.RGBA{R: 0xFF, G: 0x94, B: 0x59, A: 0xFF}
	}

	return m.baseColor
}

func (m Fire) ProcessInternal(env types.TileEnvironment) {
	env.ReduceHealth(m.healthDamperStep)

	if env.Health() < 10.0 && pkg.FlipCoin() {
		env.ReplaceTile(NewFire(), types.MaterialFlagIsFlammable)
	}

	if pkg.FlipCoin() {
		env.AddTile(NewSmoke())
	}
}

func (m Fire) ProcessCollision(env types.CollisionEnvironment) {
	if env.IsFlagged(types.MaterialFlagIsGas) {
		env.MoveSandSource()
		return
	}

	env.ReflectSourceTargetForces()

	return
}
