package materials

import (
	"image/color"

	"github.com/itiky/goPixelWorld/world/types"
)

var _ types.Material = Grass{}

type Grass struct {
	base
	growsRateK                       float64
	waterDrainK                      float64
	surroundingWaterGrowsMultiplierK float64
}

func NewGrass() Grass {
	return Grass{
		base: newBase(
			color.RGBA{R: 0x04, G: 0xDE, B: 0x1E, A: 0xFF},
			withFlags(types.MaterialFlagIsFlammable),
			withMass(7.5),
			withForceDamperK(0.1),
			withInitialHealth(5.0),
		),
		growsRateK:                       0.5, //0.1,
		waterDrainK:                      15.0,
		surroundingWaterGrowsMultiplierK: 3.0,
	}
}

func (m Grass) Type() types.MaterialType {
	return types.MaterialTypeWood
}

func (m Grass) ColorAdjusted(health float64) color.Color {
	if health < 20.0 {
		return color.RGBA{R: 0x7D, G: 0xAC, B: 0x1A, A: 0x35}
	} else if health < 40.0 {
		return color.RGBA{R: 0x69, G: 0xAC, B: 0x2B, A: 0x5C}
	} else if health < 60.0 {
		return color.RGBA{R: 0x54, G: 0xBC, B: 0x21, A: 0x6F}
	} else if health < 80.0 {
		return color.RGBA{R: 0x04, G: 0xBC, B: 0x1E, A: 0xFF}
	}

	return m.baseColor
}

func (m Grass) ProcessInternal(env types.TileEnvironment) {
	if cnt := env.ReduceEnvHealthByType(m.waterDrainK, types.MaterialTypeWater); cnt > 0 {
		env.ReduceHealth(-m.surroundingWaterGrowsMultiplierK * float64(cnt))
	} else {
		env.ReduceHealth(-m.growsRateK)
	}

	if health := env.Health(); health >= 100.0 {
		if env.AddTileGrassStyle(NewGrass()) {
			env.ReduceHealth(health - m.initialHealth)
		}
	}
}

func (m Grass) ProcessCollision(env types.CollisionEnvironment) {
	env.ReflectSourceTargetForces()
}
