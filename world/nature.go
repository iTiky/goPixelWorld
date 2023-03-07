package world

import (
	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/materials"
	"github.com/itiky/goPixelWorld/world/types"
)

const (
	natureCloudsTimeout = 60
)

func (m *Map) initNatureEvents() {
	m.natureCloudsTimeout = natureCloudsTimeout
}

func (m *Map) handleNatureEvents() (inputActions []types.InputAction) {
	m.natureCloudsTimeout--

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

	return inputActions
}
