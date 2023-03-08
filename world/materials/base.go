package materials

import (
	"image/color"

	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/types"
)

// AllMaterialsSet is set by type of all known Materials.
var AllMaterialsSet = map[types.MaterialType]types.Material{
	types.MaterialTypeSand:         NewSand(),
	types.MaterialTypeWater:        NewWater(),
	types.MaterialTypeWood:         NewWood(),
	types.MaterialTypeFire:         NewFire(),
	types.MaterialTypeGrass:        NewGrass(),
	types.MaterialTypeSmoke:        NewSmoke(),
	types.MaterialTypeSteam:        NewSteam(),
	types.MaterialTypeMetal:        NewMetal(),
	types.MaterialTypeRock:         NewRock(),
	types.MaterialTypeGraviton:     NewGraviton(),
	types.MaterialTypeAntiGraviton: NewAntiGraviton(),
	types.MaterialTypeBug:          NewBug(),
}

type (
	// base defines common fields and method for all Materials.
	base struct {
		// Base params
		mType     types.MaterialType
		flags     map[types.MaterialFlag]bool // Material properties set
		baseColor color.Color                 // the main Particle's color
		mass      float64                     // Particle's mass
		// Collision processing
		srcForceDamperK   float64 // source Particle force Vector damper K (the one who has collided to us)
		srcHealthDampStep float64 // source Particle health damper step
		// Self-processing
		selfHealthInitial  float64                      // Particle's initial health
		selfHealthDampStep float64                      // Particle health change step
		closeRangeType     types.MaterialCloseRangeType // the type of CloseRange environment required for self-processing
		closeRangeCircleR  int                          // the circle area radius requirement for the CloseRange environment
	}

	// baseOpt defines a Material constructor option
	baseOpt func(*base)
)

// withFlags sets a Material flags (properties)
func withFlags(flags ...types.MaterialFlag) baseOpt {
	return func(m *base) {
		for _, flag := range flags {
			m.flags[flag] = true
		}
	}
}

// withMass sets a Material mass.
func withMass(mass float64) baseOpt {
	return func(m *base) {
		m.mass = mass
	}
}

// withSourceDamping sets the source Particle damping coefs for collision processing.
func withSourceDamping(forceK, healthStep float64) baseOpt {
	return func(m *base) {
		m.srcForceDamperK = forceK
		m.srcHealthDampStep = healthStep
	}
}

// withCloseRangeType sets the CloseRange environment build required type for self-processing.
func withCloseRangeType(closeRangeType types.MaterialCloseRangeType) baseOpt {
	return func(m *base) {
		m.closeRangeType = closeRangeType
	}
}

// withCloseRangeCircleR sets the circle area radius for self-processing.
func withCloseRangeCircleR(r int) baseOpt {
	return func(m *base) {
		m.closeRangeCircleR = r
	}
}

// withSelfHealthReduction sets the initial and damper step health params for self-processing.
func withSelfHealthReduction(initial, dampStep float64) baseOpt {
	return func(m *base) {
		m.selfHealthInitial = initial
		m.selfHealthDampStep = dampStep
	}
}

// newBase creates a new base Material with defaults.
func newBase(mType types.MaterialType, baseColor color.Color, opts ...baseOpt) base {
	m := base{
		mType:             mType,
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

func (m base) commonProcessInternal(env types.TileEnvironment) {
	env.AddWind()
}

/* The following methods partially implements the types.Material interface */

func (m base) Type() types.MaterialType {
	return m.mType
}

func (m base) Name() string {
	return m.mType.String()
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

func (m base) ProcessInternal(env types.TileEnvironment) {
	m.commonProcessInternal(env)
}
