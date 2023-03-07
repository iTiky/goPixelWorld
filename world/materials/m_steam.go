package materials

import (
	"image/color"

	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/types"
)

var _ types.Material = Steam{}

// Steam moves up and dissipates with time.
// The dissipation can produce a new Water Particle.
type Steam struct {
	base
}

func NewSteam() Steam {
	return Steam{
		base: newBase(
			types.MaterialTypeSteam,
			color.RGBA{R: 0x05, G: 0x00, B: 0xA7, A: 0xFF},
			withFlags(types.MaterialFlagIsGas),
			withCloseRangeType(types.MaterialCloseRangeTypeSelfOnly),
			withMass(5.0),
			withSelfHealthReduction(100.0, 0.25),
		),
	}
}

func (m Steam) ColorAdjusted(health float64) color.Color {
	if health < 20.0 {
		return color.RGBA{R: 0x14, G: 0x9D, B: 0xC5, A: 0x35}
	} else if health < 40.0 {
		return color.RGBA{R: 0x09, G: 0x75, B: 0xB8, A: 0x5C}
	} else if health < 60.0 {
		return color.RGBA{R: 0x07, G: 0x57, B: 0xBB, A: 0x6F}
	} else if health < 80.0 {
		return color.RGBA{R: 0x02, G: 0x37, B: 0xAB, A: 0xFF}
	}

	return m.baseColor
}

func (m Steam) ProcessInternal(env types.TileEnvironment) {
	env.AddReverseGravity()
	env.DampSelfHealth(m.selfHealthDampStep)

	if env.Health() < 10.0 && pkg.RollDice(3) {
		env.RemoveSelfHealthDamps()
		env.ReplaceSelf(AllMaterialsSet[types.MaterialTypeWater])
	}
}

func (m Steam) ProcessCollision(env types.CollisionEnvironment) {
	if env.IsFlagged(types.MaterialFlagIsGas) {
		env.MoveSandSource()
		return
	}
	env.SwapSourceTarget()
}
