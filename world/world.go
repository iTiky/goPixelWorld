package world

import (
	"fmt"
	"math/rand"

	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/closerange"
	"github.com/itiky/goPixelWorld/world/collision"
	"github.com/itiky/goPixelWorld/world/materials"
	"github.com/itiky/goPixelWorld/world/types"
)

type Map struct {
	width     int
	height    int
	particles map[uint64]*types.Tile
	grid      [][]*types.Tile
	//
	tileEnv      *closerange.Environment
	collisionEnv *collision.Environment
}

func NewMap(width, height int) (*Map, error) {
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("invalid map size: %dx%d", width, height)
	}

	m := Map{
		width:        width,
		height:       height,
		particles:    make(map[uint64]*types.Tile),
		grid:         make([][]*types.Tile, width),
		tileEnv:      closerange.NewEnvironment(nil),
		collisionEnv: collision.NewEnvironment(pkg.DirectionTop, nil, nil),
	}
	for x := 0; x < width; x++ {
		m.grid[x] = make([]*types.Tile, height)
		for y := 0; y < height; y++ {
			var material types.Material
			if x == 0 || y == 0 || x == width-1 || y == height-1 {
				material = materials.NewBorder()
			}

			m.createTile(types.NewPosition(x, y), material)
		}
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

func (m *Map) Update() {
	m.iterateNonEmptyTiles(m.processTile)
}

func (m *Map) CreateParticles(x, y, radius int, materialBz types.MaterialI, randomForce bool) {
	for _, pos := range types.PositionsInCircle(x, y, radius) {
		if !m.isPositionValid(pos.X, pos.Y) {
			continue
		}

		material, ok := materialBz.(types.Material)
		if !ok {
			continue
		}

		tile := m.getTile(pos.X, pos.Y)
		if tile.HasParticle() {
			if material.Type() != types.MaterialTypeFire {
				continue
			}
			if !m.removeParticle(tile) {
				continue
			}
		}
		m.createParticle(tile, material)

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
		if !m.isPositionValid(pos.X, pos.Y) {
			continue
		}

		tile := m.getTile(pos.X, pos.Y)
		if !tile.HasParticle() {
			continue
		}

		m.removeParticle(tile)
	}
}

func (m *Map) FlipGravity() {
	closerange.FlipGravity()
}

func (m *Map) isPositionValid(x, y int) bool {
	return x >= 0 && y >= 0 && x < m.width && y < m.height
}
