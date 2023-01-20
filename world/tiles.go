package world

import (
	"github.com/itiky/goPixelWorld/world/types"
)

func (m *Map) iterateNonEmptyTiles(fn func(tile *types.Tile)) {
	for _, tile := range m.particles {
		fn(tile)
	}
}

func (m *Map) getTile(x, y int) *types.Tile {
	return m.grid[x][y]
}

func (m *Map) createTile(pos types.Position, material types.Material) {
	tile := types.NewTile(pos, nil)
	if material != nil {
		m.createParticle(tile, material)
	}

	m.grid[tile.Pos.X][tile.Pos.Y] = tile
}

func (m *Map) createParticle(tile *types.Tile, material types.Material) {
	if m.monitor != nil {
		m.monitor.AddParticle()
	}

	particle := types.NewParticle(material)
	tile.Particle = particle

	m.particles[tile.Particle.ID()] = tile
}

func (m *Map) removeParticle(tile *types.Tile) bool {
	if m.monitor != nil {
		m.monitor.RemoveParticle()
	}

	if tile.HasParticle() && tile.Particle.Material().IsFlagged(types.MaterialFlagIsUnremovable) {
		return false
	}

	delete(m.particles, tile.Particle.ID())
	m.grid[tile.Pos.X][tile.Pos.Y].Particle = nil

	return true
}

func (m *Map) moveTile(sourceTile *types.Tile, targetPos types.Position) {
	if m.monitor != nil {
		m.monitor.AddParticleMove()
	}

	targetTile := m.getTile(targetPos.X, targetPos.Y)
	targetTile.Particle = sourceTile.Particle
	sourceTile.Particle = nil

	targetTile.Particle.OnMove()
	m.particles[targetTile.Particle.ID()] = targetTile
}

func (m *Map) swapTiles(tile1, tile2 *types.Tile) {
	if m.monitor != nil {
		m.monitor.AddParticleMove()
	}

	tile1.Particle, tile2.Particle = tile2.Particle, tile1.Particle

	tile1.Particle.OnMove()
	tile2.Particle.OnMove()
	m.particles[tile1.Particle.ID()] = tile1
	m.particles[tile2.Particle.ID()] = tile2
}
