package world

import (
	"github.com/itiky/goPixelWorld/world/types"
)

func (m *Map) iterateNonEmptyTiles(fn func(tile *types.Tile)) {
	for _, pos := range m.particles {
		fn(m.getTile(pos))
	}
}

func (m *Map) createTile(pos types.Position, material types.Material) *types.Tile {
	particle := types.NewParticle(material)
	m.particles[particle.ID()] = pos
	m.grid[pos.X][pos.Y] = particle

	return types.NewTile(pos, particle)
}

func (m *Map) removeTile(tile *types.Tile) {
	delete(m.particles, tile.Particle.ID())
	m.grid[tile.Pos.X][tile.Pos.Y] = nil
}

func (m *Map) removeTileAtPos(pos types.Position) {
	tile := m.getTile(pos)
	if tile == nil || !tile.HasParticle() || tile.Particle.Material().IsFlagged(types.MaterialFlagIsUnremovable) {
		return
	}

	m.removeTile(tile)
}

func (m *Map) getTile(pos types.Position) *types.Tile {
	if !m.isPositionValid(pos) {
		return nil
	}

	return types.NewTile(pos, m.grid[pos.X][pos.Y])
}

func (m *Map) moveTile(tile *types.Tile, newPos types.Position) {
	if tile := m.getTile(newPos); tile == nil || tile.HasParticle() {
		panic("FUCK")
	}

	m.grid[newPos.X][newPos.Y] = tile.Particle
	m.grid[tile.Pos.X][tile.Pos.Y] = nil
	m.particles[tile.Particle.ID()] = newPos
	tile.Pos = newPos
}

func (m *Map) swapTiles(tile1, tile2 *types.Tile) {
	pos1, pos2 := tile1.Pos, tile2.Pos
	particle1, particle2 := tile1.Particle, tile2.Particle

	m.particles[tile1.Particle.ID()] = pos2
	m.particles[tile2.Particle.ID()] = pos1
	m.grid[pos1.X][pos1.Y] = particle2
	m.grid[pos2.X][pos2.Y] = particle1
	tile1.Pos, tile2.Pos = pos2, pos1
}
