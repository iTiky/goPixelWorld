package materials

import (
	"image/color"

	"github.com/itiky/goPixelWorld/world/types"
)

type (
	base struct {
		flags            map[types.MaterialFlag]bool
		baseColor        color.Color
		forceDamperK     float64
		healthDamperStep float64
		mass             float64
	}

	baseOpt func(*base)
)

func withForceDamperK(k float64) baseOpt {
	return func(m *base) {
		m.forceDamperK = k
	}
}

func withHealthDamperStep(step float64) baseOpt {
	return func(m *base) {
		m.healthDamperStep = step
	}
}

func withMass(mass float64) baseOpt {
	return func(m *base) {
		m.mass = mass
	}
}

func withFlags(flags ...types.MaterialFlag) baseOpt {
	return func(m *base) {
		for _, flag := range flags {
			m.flags[flag] = true
		}
	}
}

func newBase(baseColor color.Color, opts ...baseOpt) base {
	m := base{
		flags:     make(map[types.MaterialFlag]bool),
		baseColor: baseColor,
		mass:      100.0,
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
	return m.baseColor
}

func (m base) IsFlagged(flag types.MaterialFlag) bool {
	return m.flags[flag]
}

func (m base) Mass() float64 {
	return m.mass
}

func (m base) ProcessInternal(env types.TileEnvironment) {}
