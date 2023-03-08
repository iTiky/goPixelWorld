package world

import (
	"math/rand"

	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/closerange"
	"github.com/itiky/goPixelWorld/world/materials"
	"github.com/itiky/goPixelWorld/world/types"
)

const (
	natureCloudsTimeout     = 60
	natureWindChangeTimeout = 500
)

func (m *Map) initNatureEvents() {
	m.natureCloudsTimeout = natureCloudsTimeout
	m.natureWindChangeTimeout = natureWindChangeTimeout
}

func (m *Map) handleNatureEvents() (inputActions []types.InputAction) {
	m.natureCloudsTimeout--
	m.natureWindChangeTimeout--

	if m.natureCloudsTimeout == 0 {
		for x := 1; x < m.width-1; x++ {
			for y := m.height / 8; y > 1; y-- {
				if !pkg.RollDice(1500) {
					continue
				}

				inputActions = append(inputActions, types.CreateParticlesInputAction{
					X:        x,
					Y:        y,
					Material: materials.AllMaterialsSet[types.MaterialTypeSteam],
				})
			}
		}
		m.natureCloudsTimeout = natureCloudsTimeout
	}
	if m.natureWindChangeTimeout == 0 {
		genWindMag := func() float64 {
			return rand.Float64() / 10.0
		}

		switch rand.Int31n(5) {
		case 0:
			closerange.SetWind(genWindMag(), true)
		case 1:
			closerange.SetWind(genWindMag(), false)
		default:
			closerange.SetWind(0.0, true)
		}

		m.natureWindChangeTimeout = natureWindChangeTimeout
	}

	return inputActions
}
