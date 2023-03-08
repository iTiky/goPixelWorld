package materials

import (
	"image/color"

	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/types"
)

var _ types.Material = Water{}

// Water spreads like water.
// It puts out the Fire and makes the Grass grow faster.
type Water struct {
	base
	surroundingFireDamperStep float64 // surrounding fire damage
}

func NewWater() Water {
	return Water{
		base: newBase(
			types.MaterialTypeWater,
			color.RGBA{R: 0x00, G: 0x6B, B: 0xFF, A: 0xAF},
			withFlags(types.MaterialFlagIsLiquid),
			withCloseRangeType(types.MaterialCloseRangeTypeSurrounding),
			withMass(10.0),
			withSelfHealthReduction(100.0, 0.2),
			withSourceDamping(0.2, 0.0),
		),
		surroundingFireDamperStep: 25.0,
	}
}

func (m Water) ColorAdjusted(health float64) color.Color {
	return m.baseColor
}

func (m Water) ProcessInternal(env types.TileEnvironment) {
	m.commonProcessInternal(env)

	// Reduce health and replace self with steam (is there is some air above)
	env.DampSelfHealth(m.selfHealthDampStep)
	if env.Health() <= 0.0 {
		airTiles, _ := env.SearchNeighbours(
			pkg.ValuePtr(true),
			[]pkg.Direction{pkg.DirectionTopLeft, pkg.DirectionTop, pkg.DirectionTopRight}, true,
			nil, false,
			nil, false,
		)
		if len(airTiles) > 0 {
			env.ReplaceSelf(AllMaterialsSet[types.MaterialTypeSteam])
		} else {
			env.RemoveSelfHealthDamps()
			env.DampSelfHealth(-m.selfHealthInitial)
		}
	}

	env.AddGravity()

	if cnt := env.DampNeighboursHealthByFlag(m.surroundingFireDamperStep, nil, []types.MaterialFlag{types.MaterialFlagIsFire}); cnt > 0 {
		env.DampSelfHealth(float64(cnt) * m.surroundingFireDamperStep)
	}

	if env.Health() <= 0.0 {
		env.RemoveSelfHealthDamps()
		env.ReplaceSelf(AllMaterialsSet[types.MaterialTypeSteam])
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
