package materials

import (
	"image/color"

	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/types"
)

var (
	SandM     = NewSand()
	WaterM    = NewWater()
	WoodM     = NewWood()
	FireM     = NewFire()
	GrassM    = NewGrass()
	SmokeM    = NewSmoke()
	SteamM    = NewSteam()
	MetalM    = NewMetal()
	RockM     = NewRock()
	GravitonM = NewGraviton()
)

type (
	base struct {
		flags     map[types.MaterialFlag]bool
		baseColor color.Color
		//
		mass float64
		//
		selfHealthInitial  float64
		selfHealthDampStep float64
		//
		srcForceDamperK   float64
		srcHealthDampStep float64
		//
		closeRangeType    types.MaterialCloseRangeType
		closeRangeCircleR int
	}

	baseOpt func(*base)
)

func withFlags(flags ...types.MaterialFlag) baseOpt {
	return func(m *base) {
		for _, flag := range flags {
			m.flags[flag] = true
		}
	}
}

func withMass(mass float64) baseOpt {
	return func(m *base) {
		m.mass = mass
	}
}

func withSourceDamping(forceK, healthStep float64) baseOpt {
	return func(m *base) {
		m.srcForceDamperK = forceK
		m.srcHealthDampStep = healthStep
	}
}

func withCloseRangeType(closeRangeType types.MaterialCloseRangeType) baseOpt {
	return func(m *base) {
		m.closeRangeType = closeRangeType
	}
}

func withCloseRangeCircleR(r int) baseOpt {
	return func(m *base) {
		m.closeRangeCircleR = r
	}
}

func withSelfHealthReduction(initial, dampStep float64) baseOpt {
	return func(m *base) {
		m.selfHealthInitial = initial
		m.selfHealthDampStep = dampStep
	}
}

func newBase(baseColor color.Color, opts ...baseOpt) base {
	m := base{
		flags:             make(map[types.MaterialFlag]bool),
		baseColor:         baseColor,
		mass:              100.0,
		selfHealthInitial: 100.0,
		closeRangeType:    types.MaterialCloseRangeTypeNone,
	}
	for _, opt := range opts {
		opt(&m)
	}

	return m
}

func (m base) Color() color.Color {
	return m.baseColor
}

func (m base) ColorAdjusted(health float64) color.Color {
	if health >= 100.0 {
		return m.baseColor
	}

	c := pkg.ColorToNRGBA(m.baseColor)
	c.R = uint8(float64(c.R) * health / 100.0)
	c.G = uint8(float64(c.G) * health / 100.0)
	c.B = uint8(float64(c.B) * health / 100.0)

	return c
}

func (m base) IsFlagged(flags ...types.MaterialFlag) bool {
	for _, flag := range flags {
		if m.flags[flag] {
			return true
		}
	}

	return false
}

func (m base) InitialHealth() float64 {
	return m.selfHealthInitial
}

func (m base) Mass() float64 {
	return m.mass
}

func (m base) CloseRangeCircleRadius() int {
	return m.closeRangeCircleR
}

func (m base) CloseRangeType() types.MaterialCloseRangeType {
	return m.closeRangeType
}

func (m base) ProcessInternal(env types.TileEnvironment) {}
