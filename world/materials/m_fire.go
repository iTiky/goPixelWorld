package materials

import (
	"image/color"

	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/types"
)

var _ types.Material = Fire{}

// Fire burns itself and surrounding burnable neighbours.
// When Fire health is low it replaces itself with the Smoke.
type Fire struct {
	base
	fireDamageDampStep float64 // flammable surrounding damage
}

func NewFire() Fire {
	return Fire{
		base: newBase(
			types.MaterialTypeFire,
			color.RGBA{R: 0xFF, G: 0xAD, B: 0x8B, A: 0xFF},
			withFlags(
				types.MaterialFlagIsGas,
				types.MaterialFlagIsFire,
			),
			withCloseRangeType(types.MaterialCloseRangeTypeSurrounding),
			withMass(1.0),
			withSelfHealthReduction(100.0, 1.5),
			withSourceDamping(0.0, 5.0),
		),
		fireDamageDampStep: 2.0,
	}
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
	m.commonProcessInternal(env)

	env.DampSelfHealth(m.selfHealthDampStep)
	env.DampNeighboursHealthByFlag(m.fireDamageDampStep, nil, []types.MaterialFlag{types.MaterialFlagIsFlammable})

	health := env.Health()
	if health < 30.0 && pkg.RollDice(3) {
		env.ReplaceNeighbourTile(AllMaterialsSet[types.MaterialTypeFire], []types.MaterialFlag{types.MaterialFlagIsFlammable})
	}

	if pkg.RollDice(3) {
		env.AddNewNeighbourTile(AllMaterialsSet[types.MaterialTypeSmoke], nil)
	}
}

func (m Fire) ProcessCollision(env types.CollisionEnvironment) {
	env.DampSourceHealth(m.srcHealthDampStep, types.MaterialFlagIsFlammable)

	if env.IsFlagged(types.MaterialFlagIsGas) {
		env.MoveSandSource()
		return
	}
	env.ReflectSourceTargetForces(m.srcForceDamperK)
}
