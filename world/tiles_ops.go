package world

import (
	"github.com/itiky/goPixelWorld/world/materials"
	"github.com/itiky/goPixelWorld/world/types"
)

// initGrid inits / reinits the Map grid creating surrounding borders.
// Method also inits the procOutput buffer with empty Pixels.
func (m *Map) initGrid(width, height int) {
	m.width = width
	m.height = height

	// Cleanup the grid for reinit case
	for pID := range m.particles {
		delete(m.particles, pID)
	}

	m.grid = make([][]*types.Tile, m.width)
	m.procOutput = make([]types.Pixel, 0, m.width*m.height)
	for x := 0; x < m.width; x++ {
		m.grid[x] = make([]*types.Tile, m.height)
		for y := 0; y < m.height; y++ {
			var material types.Material
			if x == 0 || y == 0 || x == m.width-1 || y == m.height-1 {
				material = materials.NewBorder()
			}

			m.createTile(types.NewPosition(x, y), material)
			m.procOutput = append(m.procOutput, types.Pixel{})
		}
	}
}

// iterateNonEmptyTiles iterates over non-empty grid Tiles.
func (m *Map) iterateNonEmptyTiles(fn func(tile *types.Tile)) {
	for _, tile := range m.particles {
		fn(tile)
	}
}

// getTile returns a Tile by its Position.
func (m *Map) getTile(x, y int) *types.Tile {
	return m.grid[x][y]
}

// createTile creates a new Tile.
// If {material} is not nil, also creates a corresponding Particle.
func (m *Map) createTile(pos types.Position, material types.Material) {
	tile := types.NewTile(pos, nil)
	if material != nil {
		m.createParticle(tile, material)
	}

	m.grid[tile.Pos.X][tile.Pos.Y] = tile
}

// createParticle creates a new Particle on the specified Tile.
func (m *Map) createParticle(tile *types.Tile, material types.Material) {
	if m.monitor != nil {
		m.monitor.AddParticle()
	}

	particle := types.NewParticle(material)
	tile.Particle = particle

	m.particles[tile.Particle.ID()] = tile
}

// removeParticle removes a single Tile's Particle.
// Skips the operations if a Particle is not removable based on Material properties.
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

// moveTile moves a Particle from the source to the destination Tile.
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

// swapTiles swaps Particles between two Tiles.
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
