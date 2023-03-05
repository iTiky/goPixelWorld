package world

import (
	"math/rand"

	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/closerange"
	"github.com/itiky/goPixelWorld/world/types"
)

// PushInputAction pushes a new input action to the queue.
// An actual handling is done in the ExportState.
func (m *Map) PushInputAction(action types.InputAction) {
	if action == nil {
		return
	}

	m.inputActions = append(m.inputActions, action)
}

// handleCreateParticlesInput handles the CreateParticlesInputAction input action.
func (m *Map) handleCreateParticlesInput(input types.CreateParticlesInputAction) {
	for _, pos := range types.PositionsInCircle(input.X, input.Y, input.Radius, true) {
		// Skip out-of-range Positions
		if !m.isPositionValid(pos.X, pos.Y) {
			continue
		}

		// Skip an unknown Material
		material, ok := input.Material.(types.Material)
		if !ok {
			continue
		}

		// If a target Tile has a particle, remove it (if it is removable)
		tile := m.getTile(pos.X, pos.Y)
		if tile.HasParticle() {
			if mType := material.Type(); mType != types.MaterialTypeFire && mType != types.MaterialTypeAntiGraviton {
				continue
			}
			if !m.removeParticle(tile) {
				continue
			}
		}
		m.createParticle(tile, material)

		// Apply an initial random force
		if input.ApplyForce {
			forceVec := pkg.NewVector(
				float64(rand.Int31n(5)),
				pkg.RandomAngle(),
			)
			tile.Particle.SetForce(forceVec)
		}
	}
}

// handleRemoveParticlesInput handles the DeleteParticlesInputAction input action.
func (m *Map) handleRemoveParticlesInput(input types.DeleteParticlesInputAction) {
	for _, pos := range types.PositionsInCircle(input.X, input.Y, input.Radius, true) {
		// Skip out-of-range Positions
		if !m.isPositionValid(pos.X, pos.Y) {
			continue
		}

		// Skip empty Tiles
		tile := m.getTile(pos.X, pos.Y)
		if !tile.HasParticle() {
			continue
		}

		m.removeParticle(tile)
	}
}

// handleFlipGravityInput handles the FlipGravityInputAction input action.
func (m *Map) handleFlipGravityInput() {
	closerange.FlipGravity()
}
