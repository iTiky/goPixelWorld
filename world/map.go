package world

import (
	"fmt"
	"math/rand"

	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/materials"
	"github.com/itiky/goPixelWorld/world/types"
)

type Map struct {
	width     int
	height    int
	particles map[uint64]types.Position
	grid      [][]*types.Particle
}

func NewMap(width, height int) (*Map, error) {
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("invalid map size: %dx%d", width, height)
	}

	m := Map{
		width:     width,
		height:    height,
		particles: make(map[uint64]types.Position),
		grid:      make([][]*types.Particle, width),
	}
	for x := 0; x < width; x++ {
		m.grid[x] = make([]*types.Particle, height)
	}

	for x := 0; x < width; x++ {
		posTop, posBottom := types.NewPosition(x, 0), types.NewPosition(x, height-1)
		m.createTile(posTop, materials.NewBorder())
		m.createTile(posBottom, materials.NewBorder())
	}
	for y := 1; y < height-1; y++ {
		posLeft, posRight := types.NewPosition(0, y), types.NewPosition(width-1, y)
		m.createTile(posLeft, materials.NewBorder())
		m.createTile(posRight, materials.NewBorder())
	}

	return &m, nil
}

func (m *Map) Size() (int, int) {
	return m.width, m.height
}

func (m *Map) IterateTiles(fn func(tile types.TileI)) {
	m.iterateNonEmptyTiles(func(tile *types.Tile) {
		fn(tile)
	})
}

func (m *Map) CreateParticles(x, y, radius int, materialBz types.MaterialI, randomForce bool) {
	for _, pos := range types.PositionsInCircle(x, y, radius) {
		if !m.isPositionValid(pos) {
			continue
		}

		material, ok := materialBz.(types.Material)
		if !ok {
			continue
		}
		if material.Type() == types.MaterialTypeFire {
			m.removeTileAtPos(pos)
		}

		if existingTile := m.getTile(pos); existingTile.HasParticle() {
			continue
		}

		tile := m.createTile(pos, material)
		if randomForce {
			forceVec := pkg.NewVector(
				float64(rand.Int31n(4)),
				pkg.RandomAngle(),
			)
			tile.Particle.SetForce(forceVec)
		}
	}
}

func (m *Map) RemoveParticles(x, y, radius int) {
	for _, pos := range types.PositionsInCircle(x, y, radius) {
		if !m.isPositionValid(pos) {
			continue
		}

		if existingTile := m.getTile(pos); !existingTile.HasParticle() {
			continue
		}

		m.removeTileAtPos(pos)
	}
}

func (m *Map) isPositionValid(pos types.Position) bool {
	return pos.X >= 0 && pos.Y >= 0 && pos.X < m.width && pos.Y < m.height
}
